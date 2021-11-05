package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var _instance *Config

func init() {
	// //TODO: 環境変数による上書き処理
	_instance = LoadConfig()
}

type Config struct {
	CleanupTickEvery     time.Duration //TODO: Durationを直接持たないで secoundをもつ
	Namespace            string
	ManifestDirs         []string
	KubeServer           string
	Slack                SlackConfig
	OverwriteRole        bool
	OverwriteRoleBinding bool
}

func GetConfig() *Config {
	return _instance
}

func LoadConfig() *Config {
	viper.SetConfigName("cinderella") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name

	viper.SetDefault("Namespace", "default")
	viper.SetDefault("CleanupTickEvery", 10*time.Second)
	viper.SetDefault("ManifestDirs", []string{"/etc/cinderella"})
	viper.SetDefault("OverwriteRole", false)
	viper.SetDefault("OverwriteRoleBinding", false)

	viper.AddConfigPath("/etc/cinderella/") // path to look for the config file in
	viper.AddConfigPath(".")                // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v , %w", config, err))
	}

	return &config
}

func (c *Config) WriteConfig() {
	// viper.SetConfigName("cinderella") // name of config file (without extension)
	// viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name

	// viper.SafeWriteConfig()
	// viper.WriteConfigAs("~/cinderella.config")
}
