package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Vars = struct {
	Host              string `mapstructure:"host"`
	PixivRefreshToken string `mapstructure:"pixiv_refresh_token"`
	PixivToken        string `mapstructure:"pixiv_token"`
	ThreadMax         int    `mapstructure:"thread_max"`
	VersionName       string `mapstructure:"version_name"`
}{}
var Vipers = viper.New()

func LoadVars() bool {
	if err := Vipers.ReadInConfig(); err != nil {
		fmt.Printf("Read config-settings.json Error:[%v]\n", err)
		return false
	}
	if err := Vipers.WriteConfig(); err != nil {
		fmt.Println("Update config file failed,please check the permission.")
	}
	return true
}

func SaveVars() {
	if LoadVars() {
		if err := Vipers.SafeWriteConfig(); err != nil {
			fmt.Println("Safe write config file failed,please check the permission or create config-settings.json file manually.")
			fmt.Println("Detailed error message as follows:", err)
		} else {
			fmt.Println("safe write config file success!")
		}
	}
}

func VarsConfigInit() {
	Vipers.SetConfigName("config-settings")
	Vipers.SetConfigType("json")
	Vipers.AddConfigPath(".")
	// path to look for the config file in
	Vipers.SetDefault("HOST", "https://app-api.pixiv.net")
	Vipers.SetDefault("thread_max", 16)
	Vipers.SetDefault("version_name", "1.0.9")
	Vipers.SetDefault("PIXIV_TOKEN", "")
	Vipers.SetDefault("PIXIV_REFRESH_TOKEN", "")
	SaveVars()
	if err := Vipers.Unmarshal(&Vars); err != nil {
		panic(err)
	}
}
