package model

import (
	"encoding/base64"
	"errors"
	"time"

	"go.x2ox.com/tea"
)

type Token struct {
	Type   Type // 预览草稿箱，阅读文章，分享
	NoteID uint64
	Time   time.Time
}

func (t Token) Valid() bool {
	switch t.Type {
	case Preview:
		return Conf.Token.Preview == 0 ||
			t.Time.Add(time.Duration(Conf.Token.Preview)*time.Minute).After(time.Now())
	case View:
		return Conf.Token.View == 0 ||
			t.Time.Add(time.Duration(Conf.Token.View)*time.Minute).After(time.Now())
	case Share:
		return Conf.Token.Share == 0 ||
			t.Time.Add(time.Duration(Conf.Token.Share)*time.Minute).After(time.Now())
	}
	return false
}

type Type uint8

const (
	None Type = iota
	Preview
	View
	Share
)

func NewToken(t Type, noteID uint64) Token {
	return Token{
		Type:   t,
		NoteID: noteID,
		Time:   time.Now(),
	}
}

func ParseToken(s string) *Token {
	t, err := decode(s)
	if err != nil {
		return nil
	}
	return t
}

func (t Token) String() string {
	return encode(t)
}

func decode(s string) (*Token, error) {
	var (
		key = GetKey()
		arr []byte
		err error
		t   *tea.TinyEncryptionAlgorithm
	)

	if arr, err = base64.URLEncoding.DecodeString(s); err != nil {
		return nil, err
	}
	if t, err = tea.NewTEA(key[:]); err != nil {
		return nil, err
	}

	t.Decrypt(arr[:8], arr[:8])
	t.Decrypt(arr[8:16], arr[8:16])
	t.Decrypt(arr[16:], arr[16:])

	if arr[0] != key[1] || arr[1] != key[3] || arr[2] != key[4] ||
		arr[3] != key[5] || arr[4] != key[2] || arr[5] != key[0] ||
		arr[6] != key[9] || arr[7] == byte(None) {
		return nil, errors.New("parse failure")
	}

	return &Token{
		Type: Type(arr[7]),
		NoteID: uint64(arr[8])<<56 | uint64(arr[9])<<48 | uint64(arr[10])<<40 | uint64(arr[11])<<32 |
			uint64(arr[12])<<24 | uint64(arr[13])<<16 | uint64(arr[14])<<8 | uint64(arr[15]),
		Time: time.Unix(int64(arr[16])<<56|int64(arr[17])<<48|int64(arr[18])<<40|int64(arr[19])<<32|
			int64(arr[20])<<24|int64(arr[21])<<16|int64(arr[22])<<8|int64(arr[23]), 0),
	}, nil
}

func encode(token Token) string {
	var (
		key  = GetKey()
		arr  = make([]byte, 24)
		t, _ = tea.NewTEA(key[:])
		ts   = token.Time.Unix()
	)

	t.Encrypt(arr[:], []byte{key[1], key[3], key[4], key[5], key[2], key[0], key[9], byte(token.Type)})
	t.Encrypt(arr[8:], []byte{
		uint8(token.NoteID >> 56), uint8(token.NoteID >> 48), uint8(token.NoteID >> 40), uint8(token.NoteID >> 32),
		uint8(token.NoteID >> 24), uint8(token.NoteID >> 16), uint8(token.NoteID >> 8), uint8(token.NoteID),
	})
	t.Encrypt(arr[16:], []byte{
		uint8(ts >> 56), uint8(ts >> 48), uint8(ts >> 40), uint8(ts >> 32),
		uint8(ts >> 24), uint8(ts >> 16), uint8(ts >> 8), uint8(ts),
	})

	return base64.URLEncoding.EncodeToString(arr)
}
