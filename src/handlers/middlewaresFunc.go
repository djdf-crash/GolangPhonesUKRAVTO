package handlers

import (
	"db"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func CheckAuthenticationFunc(ctx *gin.Context) {
	var resp ResponseModelToken
	var user db.User

	token := ctx.GetHeader("token")
	if len(token) == 0 {

		resp = ResponseModelToken{
			Result: false,
			Error:  "No found token!",
			Token:  "",
		}
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}
	userDB := db.GetUserByToken(token)

	if reflect.DeepEqual(user, userDB) {
		resp = ResponseModelToken{
			Result: false,
			Error:  "No found user with your token!",
			Token:  "",
		}
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}

	ctx.Set("user", userDB)
	ctx.Set("name_organization", ctx.Query("nameorganization"))
	ctx.Set("id_organization", ctx.Query("idorganization"))

}
