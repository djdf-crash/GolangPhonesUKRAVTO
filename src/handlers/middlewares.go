package handlers

import "github.com/gin-gonic/gin"

func CheckAuthenticationMiddleware(ctx *gin.Context) {

	CheckAuthenticationFunc(ctx)

}
