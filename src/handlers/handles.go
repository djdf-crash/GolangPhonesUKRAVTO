package handlers

import (
	"github.com/gin-gonic/gin"
)

func AllPhonesHandler(ctx *gin.Context) {
	GetAllPhonesFunc(ctx)
}

func GetPhonesLastUpdateHandler(ctx *gin.Context) {
	GetPhonesLastUpdateFunc(ctx)
}

func AddUsersHandler(ctx *gin.Context) {
	AddUsersFunc(ctx)
}

func TokenIsExistHandler(ctx *gin.Context) {
	TokenIsExistFunc(ctx)
}

func UpdateUsersHandler(ctx *gin.Context) {
	UpdateUsersFunc(ctx)
}

func CheckEmailHandler(ctx *gin.Context) {
	CheckEmailFunc(ctx)
}

func PhonesByOrganizationNameHandler(ctx *gin.Context) {
	GetPhonesByOrganizationNameFunc(ctx)
}

func PhonesByOrganizationIDHandler(ctx *gin.Context) {
	GetPhonesByOrganizationIDFunc(ctx)
}

func AllOrganizationHandler(ctx *gin.Context) {
	GetAllOrganizationFunc(ctx)
}

func GetLastUpdateAPKHandler(ctx *gin.Context) {
	GetLastUpdateAPKFunc(ctx)
}

func DownloadLastUpdateAPKHandler(ctx *gin.Context) {
	DownloadLastUpdateAPKFunc(ctx)
}
