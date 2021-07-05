package dandelion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// New creates a new Engine instance.
//
// It requires a token, provided by @BotFather on Telegram.
func New(token string) (*Engine, error) {
	return NewEngine(token, APIEndpoint, &http.Client{})
}

// NewEngine creates a new Engine instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewEngine(token, apiEndpoint string, client HTTPClient) (*Engine, error) {
	bot := &Engine{
		Token:       token,
		Client:      client,
		apiEndpoint: apiEndpoint,
	}
	bot.ctx, bot.cancelFunc = context.WithCancel(context.Background())
	bot.pool.New = func() interface{} {
		return &Context{
			Engine: bot,
		}
	}
	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// SetAPIEndpoint changes the Telegram Bot API endpoint used by the instance.
func (bot *Engine) SetAPIEndpoint(apiEndpoint string) {
	bot.apiEndpoint = apiEndpoint
}
func (bot *Engine) buildMethod(endpoint string) string {
	return fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *Engine) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v\n", endpoint, params)
	}

	resp, err := bot.Client.PostForm(bot.buildMethod(endpoint), params.Build())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err = apiResp.Decode(resp.Body); err != nil {
		return &apiResp, err
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// UploadFiles makes a request to the API with files.
func (bot *Engine) UploadFiles(endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	// This code modified from the very helpful @HirbodBehnam
	// https://github.com/go-telegram-bot-api/telegram-bot-api/issues/354#issuecomment-663856473
	go func() {
		defer w.Close()
		defer m.Close()

		var err error

		for field, value := range params {
			if err = m.WriteField(field, value); err != nil {
				w.CloseWithError(err)
				return
			}
		}

		for _, file := range files {
			switch f := file.File.(type) {
			case string:
				var fileHandle *os.File
				var part io.Writer

				if fileHandle, err = os.Open(f); err == nil {
					if part, err = m.CreateFormFile(file.Name, fileHandle.Name()); err == nil {
						_, err = io.Copy(part, fileHandle)
					}
					_ = fileHandle.Close()
				}
			case FileBytes:
				var part io.Writer
				if part, err = m.CreateFormFile(file.Name, f.Name); err == nil {
					_, err = io.Copy(part, bytes.NewBuffer(f.Bytes))
				}
			case FileReader:
				var part io.Writer
				if part, err = m.CreateFormFile(file.Name, f.Name); err == nil {
					_, err = io.Copy(part, f.Reader)
				}
			case FileURL:
				err = m.WriteField(file.Name, string(f))
			case FileID:
				err = m.WriteField(file.Name, string(f))
			default:
				err = errors.New(ErrBadFileType)
			}

			if err != nil {
				_ = w.CloseWithError(err)
				return
			}
		}
	}()

	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v, with %d files\n", endpoint, params, len(files))
	}

	req, err := http.NewRequest("POST", bot.buildMethod(endpoint), r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err = apiResp.Decode(resp.Body); err != nil {
		return &apiResp, err
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// GetFileDirectURL returns direct URL to file
//
// It requires the FileID.
func (bot *Engine) GetFileDirectURL(fileID string) (string, error) {
	file, err := bot.GetFile(FileConfig{FileID: fileID})

	if err != nil {
		return "", err
	}

	return file.Link(bot.Token), nil
}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from Engine.Self without the need for
// another request.
func (bot *Engine) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}

// IsMessageToMe returns true if message directed to this bot.
//
// It requires the Message.
func (bot *Engine) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+bot.Self.UserName)
}

func hasFilesNeedingUpload(files []RequestFile) bool {
	for _, file := range files {
		switch file.File.(type) {
		case string, FileBytes, FileReader:
			return true
		}
	}

	return false
}

// Request sends a Chattable to Telegram, and returns the APIResponse.
func (bot *Engine) Request(c Chattable) (*APIResponse, error) {
	params, err := c.Params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		files := t.Files()

		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFilesNeedingUpload(files) {
			return bot.UploadFiles(t.Method(), params, files)
		}

		// However, if there are no files to be uploaded, there's likely things
		// that need to be turned into params instead.
		for _, file := range files {
			var s string

			switch f := file.File.(type) {
			case string:
				s = f
			case FileID:
				s = string(f)
			case FileURL:
				s = string(f)
			default:
				return nil, errors.New(ErrBadFileType)
			}

			params[file.Name] = s
		}
	}

	return bot.MakeRequest(c.Method(), params)
}

// Send will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *Engine) Send(c Chattable) (Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return Message{}, err
	}

	var msg Message
	err = json.Unmarshal(resp.Result, &msg)

	return msg, err
}

// SendSet will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *Engine) SendSet(c Chattable) (bool, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return false, err
	}

	var msg bool
	err = json.Unmarshal(resp.Result, &msg)

	return msg, err
}

// SendMediaGroup sends a media group and returns the resulting messages.
func (bot *Engine) SendMediaGroup(config MediaGroupConfig) ([]Message, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return nil, err
	}

	var messages []Message
	err = json.Unmarshal(resp.Result, &messages)

	return messages, err
}

// GetUserProfilePhotos gets a user's profile photos.
//
// It requires UserID.
// Offset and Limit are optional.
func (bot *Engine) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	err = json.Unmarshal(resp.Result, &profilePhotos)

	return profilePhotos, err
}

// GetFile returns a File which can download a file from Telegram.
//
// Requires FileID.
func (bot *Engine) GetFile(config FileConfig) (File, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return File{}, err
	}

	var file File
	err = json.Unmarshal(resp.Result, &file)

	return file, err
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, Timeout, and AllowedUpdates are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests so you can get updates
// instantly instead of having to wait between requests.
func (bot *Engine) GetUpdates(config UpdateConfig) ([]Update, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)

	return updates, err
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *Engine) GetWebhookInfo() (WebhookInfo, error) {
	resp, err := bot.MakeRequest("getWebhookInfo", nil)
	if err != nil {
		return WebhookInfo{}, err
	}

	var info WebhookInfo
	err = json.Unmarshal(resp.Result, &info)

	return info, err
}

// Run starts and returns a channel for getting updates.
func (bot *Engine) Run() {
	go func() {
		for {
			select {
			case <-bot.ctx.Done():
				return
			default:
			}

			updates, err := bot.GetUpdates(bot.UpdateConfig)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= bot.UpdateConfig.Offset {
					bot.UpdateConfig.Offset = update.UpdateID + 1
					bot.serve(update)
				}
			}
		}
	}()
}

// Stop stops the go routine which receives updates
func (bot *Engine) Stop() {
	bot.cancelFunc()
}

// ListenForWebhook return a http func
func (bot *Engine) ListenForWebhook(w http.ResponseWriter, r *http.Request) {
	select {
	case <-bot.ctx.Done():
		return
	default:
	}
	update, err := bot.HandleUpdate(r)
	if err != nil {
		errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(errMsg)
		return
	}

	bot.serve(*update)
}

// HandleUpdate parses and returns update received via webhook
func (bot *Engine) HandleUpdate(r *http.Request) (*Update, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("wrong HTTP method required POST")
	}

	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

// WriteToHTTPResponse writes the request to the HTTP ResponseWriter.
//
// It doesn't support uploading files.
//
// See https://core.telegram.org/bots/api#making-requests-when-getting-updates
// for details.
func WriteToHTTPResponse(w http.ResponseWriter, c Chattable) error {
	params, err := c.Params()
	if err != nil {
		return err
	}

	if t, ok := c.(Fileable); ok {
		if hasFilesNeedingUpload(t.Files()) {
			return errors.New("unable to use http response to upload files")
		}
	}

	values := params.Build()
	values.Set("method", c.Method())

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = w.Write([]byte(values.Encode()))
	return err
}

// GetChat gets information about a chat.
func (bot *Engine) GetChat(config ChatInfoConfig) (Chat, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return Chat{}, err
	}

	var chat Chat
	err = json.Unmarshal(resp.Result, &chat)

	return chat, err
}

// GetChatAdministrators gets a list of administrators in the chat.
//
// If none have been appointed, only the creator will be returned.
// Bots are not shown, even if they are an administrator.
func (bot *Engine) GetChatAdministrators(config ChatAdministratorsConfig) ([]ChatMember, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []ChatMember{}, err
	}

	var members []ChatMember
	err = json.Unmarshal(resp.Result, &members)

	return members, err
}

// GetChatMembersCount gets the number of users in a chat.
func (bot *Engine) GetChatMembersCount(config ChatMemberCountConfig) (int, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return -1, err
	}

	var count int
	err = json.Unmarshal(resp.Result, &count)

	return count, err
}

// GetChatMember gets a specific chat member.
func (bot *Engine) GetChatMember(config GetChatMemberConfig) (ChatMember, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return ChatMember{}, err
	}

	var member ChatMember
	err = json.Unmarshal(resp.Result, &member)

	return member, err
}

// GetGameHighScores allows you to get the high scores for a game.
func (bot *Engine) GetGameHighScores(config GetGameHighScoresConfig) ([]GameHighScore, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []GameHighScore{}, err
	}

	var highScores []GameHighScore
	err = json.Unmarshal(resp.Result, &highScores)

	return highScores, err
}

// GetInviteLink get InviteLink for a chat
func (bot *Engine) GetInviteLink(config ChatInviteLinkConfig) (string, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return "", err
	}

	var inviteLink string
	err = json.Unmarshal(resp.Result, &inviteLink)

	return inviteLink, err
}

// GetStickerSet returns a StickerSet.
func (bot *Engine) GetStickerSet(config GetStickerSetConfig) (StickerSet, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return StickerSet{}, err
	}

	var stickers StickerSet
	err = json.Unmarshal(resp.Result, &stickers)

	return stickers, err
}

// StopPoll stops a poll and returns the result.
func (bot *Engine) StopPoll(config StopPollConfig) (Poll, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return Poll{}, err
	}

	var poll Poll
	err = json.Unmarshal(resp.Result, &poll)

	return poll, err
}

// GetMyCommands gets the currently registered commands.
func (bot *Engine) GetMyCommands() ([]BotCommand, error) {
	config := GetMyCommandsConfig{}

	resp, err := bot.Request(config)
	if err != nil {
		return nil, err
	}

	var commands []BotCommand
	err = json.Unmarshal(resp.Result, &commands)

	return commands, err
}

// CopyMessage copy messages of any kind. The method is analogous to the method
// forwardMessage, but the copied message doesn't have a link to the original
// message. Returns the MessageID of the sent message on success.
func (bot *Engine) CopyMessage(config CopyMessageConfig) (MessageID, error) {
	params, err := config.Params()
	if err != nil {
		return MessageID{}, err
	}

	resp, err := bot.MakeRequest(config.Method(), params)
	if err != nil {
		return MessageID{}, err
	}

	var messageID MessageID
	err = json.Unmarshal(resp.Result, &messageID)

	return messageID, err
}
