package api

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/tpl"
)

func previewAction(c *gin.Context) {
	ar := model.GetKey()
	fmt.Println("key:", hex.EncodeToString(ar[:]))
	tokenStr := c.Param("token")
	tk := model.ParseToken(tokenStr)
	if tk == nil || !tk.Valid() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Token", tokenStr, 0, "", "", true, true)

	if tk.NoteID == 0 {
		c.HTML(http.StatusOK, "tpl.html", tpl.ToHTML(db.Input.String()))
		return
	}

	note := db.Note.GetWithID(tk.NoteID)
	if note == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "tpl.html", note.HTML())
}
