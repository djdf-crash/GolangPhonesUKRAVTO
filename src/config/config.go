package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type server struct {
	ModeStart      string `json:"mode_start,required"`
	Port           string `json:"port,required"`
	SecretKeyToken string `json:"secret_key_token"`
}

type dataBase struct {
	NameDriver string `json:"name_driver,required"`
	Path       string `json:"path,required"`
	LogMode    bool   `json:"log_mode"`
}

type SettingUpdateFileAPK struct {
	PathFile string `json:"path_file"`
	Path     string `json:"path"`
}

type SettingsParseFile struct {
	PathFile                string `json:"path_file"`
	NumberColumnPost        int    `json:"number_column_post"`
	NumberColumnFullName    int    `json:"number_column_full_name"`
	NumberColumnEmail       int    `json:"number_column_email"`
	NumberColumnContactInfo int    `json:"number_column_contact_info"`
	NumberColumnPhone       int    `json:"number_column_phone"`
}

type config struct {
	Server                     *server               `json:"server,required"`
	DataBase                   *dataBase             `json:"data_base,required"`
	SettingsParseFile          *SettingsParseFile    `json:"settings_parse_file"`
	SettingsParseUpdateAPKFile *SettingUpdateFileAPK `json:"settings_parse_update_apk_file"`
	RootDirPath                string
}

var AppConfig *config

func InitConfig(pathConfigFile string) error {

	configFile, err := os.Open(pathConfigFile)
	if err != nil {
		return err
	}
	defer configFile.Close()

	dec := json.NewDecoder(configFile)
	err = dec.Decode(&AppConfig)
	if err != nil {
		return err
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	AppConfig.RootDirPath = dir + string(os.PathSeparator)

	return nil

}
