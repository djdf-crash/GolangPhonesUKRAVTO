package handlers

import (
	"config"
	"db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
	"utils"
)

func GetAllPhonesFunc(ctx *gin.Context) {
	arrEmployee := db.GetAllEmployee()

	resp := ResponseModelEmployee{
		Result:    true,
		Error:     "",
		Employees: arrEmployee,
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func GetPhonesByOrganizationNameFunc(ctx *gin.Context) {

	var arrEmployee []db.Employee
	var resp ResponseModelEmployee
	name := ctx.GetString("name_organization")

	user, _ := ctx.Get("user")
	userDB := user.(db.User)

	org := db.GetOrganizationByNameAndLastUpdate(name, userDB.LastUpdate)

	arrEmployee = db.GetEmployeesByOrganizationIDLastUpdate(org.ID, userDB.LastUpdate)
	resp = ResponseModelEmployee{
		Result:    true,
		Error:     "",
		Employees: arrEmployee,
	}

	//userDB.LastUpdate = time.Now()
	//db.UpdateUser(&userDB)

	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func GetPhonesByOrganizationIDLastUpdateFunc(ctx *gin.Context) {
	var arrEmployee []db.Employee
	var resp ResponseModelEmployee

	idStr := ctx.GetString("id_organization")

	u64ISStr, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp = ResponseModelEmployee{
			Result:    false,
			Error:     err.Error(),
			Employees: arrEmployee,
		}
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}

	user, _ := ctx.Get("user")
	userDB := user.(db.User)

	org := db.GetOrganizationByIDAndLastUpdate(uint(u64ISStr), userDB.LastUpdate)

	arrEmployee = db.GetEmployeesByOrganizationIDLastUpdate(org.ID, userDB.LastUpdate)
	resp = ResponseModelEmployee{
		Result:    true,
		Error:     "",
		Employees: arrEmployee,
	}

	//userDB.LastUpdate = time.Now()
	//db.UpdateUser(&userDB)

	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func GetPhonesByOrganizationIDFunc(ctx *gin.Context) {

	var arrEmployee []db.Employee
	var resp ResponseModelEmployee

	idStr := ctx.Param("id")

	u64ISStr, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp = ResponseModelEmployee{
			Result:    false,
			Error:     err.Error(),
			Employees: arrEmployee,
		}
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}

	org := db.GetOrganizationByID(uint(u64ISStr))

	arrEmployee = db.GetEmployeesByOrganizationID(org.ID)
	resp = ResponseModelEmployee{
		Result:    true,
		Error:     "",
		Employees: arrEmployee,
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)

}

func GetAllOrganizationFunc(ctx *gin.Context) {

	arrOrg := db.GetAllOrganizations()

	resp := ResponseModelOrganization{
		Result:       true,
		Error:        "",
		Organization: arrOrg,
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)

}

func AddUsersFunc(ctx *gin.Context) {

	var user db.User
	var resp ResponseModelToken

	email := ctx.GetString("user_email")
	deviceID := ctx.GetString("user_device_id")

	userDB := db.GetUserByEmailAndDeviceID(email, deviceID)

	if reflect.DeepEqual(user, userDB) {
		userDB.Email = email
		userDB.DeviceID = deviceID
		userDB.Token = utils.ComputeHmac256(userDB.Email+userDB.DeviceID, config.AppConfig.Server.SecretKeyToken)
		db.AddUser(&userDB)

		resp = ResponseModelToken{
			Result: true,
			Error:  "",
			Token:  userDB.Token,
		}

	} else {
		resp = ResponseModelToken{
			Result: true,
			Error:  "",
			Token:  userDB.Token,
		}
		db.UpdateUser(&userDB, time.Time{})
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func TokenIsExistFunc(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userDB := user.(db.User)
	resp := ResponseModelToken{
		Result: true,
		Error:  "",
		Token:  userDB.Token,
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func UpdateUsersFunc(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userDB := user.(db.User)

	db.UpdateUser(&userDB, time.Now())

	resp := ResponseModelToken{
		Result: true,
		Error:  "",
		Token:  "",
	}
	RespondWithMessage(http.StatusOK, 0, resp, ctx)
}

func GetPhonesLastUpdateFunc(ctx *gin.Context) {

	var empsDB []db.Employee

	user, _ := ctx.Get("user")
	userDB := user.(db.User)

	nameOrg := ctx.GetString("name_organization")
	idOrg := ctx.GetString("id_organization")

	if len(nameOrg) == 0 && len(idOrg) == 0 {
		empsDB = db.GetEmployeesByLastUpdate(userDB.LastUpdate)

		resp := ResponseModelEmployee{
			Result:    true,
			Error:     "",
			Employees: empsDB,
		}

		RespondWithMessage(http.StatusOK, 0, resp, ctx)

		//userDB.LastUpdate = time.Now()
		//db.UpdateUser(&userDB)
	} else if len(nameOrg) != 0 {
		GetPhonesByOrganizationNameFunc(ctx)
	} else if len(idOrg) != 0 {
		GetPhonesByOrganizationIDLastUpdateFunc(ctx)
	}
}

func CheckEmailFunc(ctx *gin.Context) {

	var emp db.Employee

	email := ctx.Query("email")
	deviceID := ctx.Query("deviceid")
	if len(email) == 0 || len(deviceID) == 0 {

		resp := ResponseModelToken{
			Result: false,
			Error:  "No found parameter email or device id!",
			Token:  "",
		}

		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return

	}

	empDb := db.GetEmployeesByEmail(email)

	if reflect.DeepEqual(emp, empDb) {
		resp := ResponseModelToken{
			Result: false,
			Error:  "No found your corporate email!",
			Token:  "",
		}

		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}
	ctx.Set("user_email", email)
	ctx.Set("user_device_id", deviceID)
	ctx.Next()
}

func GetLastUpdateAPKFunc(ctx *gin.Context) {

	var autoGenerate []AutoGenerated

	resp := ResponseModelUpdateAPK{
		Result:            false,
		NewVersionNameAPK: "",
	}

	configFile, err := os.Open(config.AppConfig.RootDirPath +
		config.AppConfig.SettingsParseUpdateAPKFile.Path +
		string(os.PathSeparator) +
		config.AppConfig.SettingsParseUpdateAPKFile.PathFile)

	if err != nil {
		resp.Error = err.Error()
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}
	defer configFile.Close()

	dec := json.NewDecoder(configFile)
	err = dec.Decode(&autoGenerate)
	if err != nil {
		resp.Error = err.Error()
		RespondWithMessage(http.StatusOK, 0, resp, ctx)
		return
	}
	currentVersionName := ctx.GetString("version_name_apk")
	if len(autoGenerate) > 0 {
		if currentVersionName != autoGenerate[0].ApkInfo.VersionName {
			resp.Result = true
			resp.NewVersionNameAPK = autoGenerate[0].ApkInfo.VersionName
			RespondWithMessage(http.StatusOK, 0, resp, ctx)
			return
		}
	}

	RespondWithMessage(http.StatusOK, 0, resp, ctx)

}

func DownloadLastUpdateAPKFunc(ctx *gin.Context) {
	var autoGenerate []AutoGenerated

	configFile, err := os.Open(config.AppConfig.RootDirPath +
		config.AppConfig.SettingsParseUpdateAPKFile.Path +
		string(os.PathSeparator) +
		config.AppConfig.SettingsParseUpdateAPKFile.PathFile)

	if err != nil {
		return
	}
	defer configFile.Close()

	dec := json.NewDecoder(configFile)
	err = dec.Decode(&autoGenerate)
	if err != nil {
		return
	}

	fullNamePath := config.AppConfig.RootDirPath + config.AppConfig.SettingsParseUpdateAPKFile.Path + string(os.PathSeparator) + autoGenerate[0].ApkInfo.OutputFile

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+autoGenerate[0].ApkInfo.OutputFile)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(fullNamePath)

}

func RespondWithMessage(codeResponse int, codeError int, message interface{}, ctx *gin.Context) {

	if codeResponse != http.StatusOK {
		response := map[string]interface{}{
			"code":    codeError,
			"message": message,
		}

		ctx.JSON(codeResponse, &response)
	} else {
		ctx.JSON(codeResponse, &message)
	}

	if ctx.Request.Method == http.MethodGet {
		ctx.Header("Status", string(codeResponse))
	}

	ctx.Abort()
}
