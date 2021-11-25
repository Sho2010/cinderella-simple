package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var _instance *Config

const TemplateExtension = ".yaml.tmpl"

func init() {
	// //TODO: 環境変数による上書き処理
	_instance = &Config{}
}

type Config struct {
	CleanupTickEverySeconds int `validate:"gte=10,lte=600"`
	Namespace               string
	ManifestDirs            []string
	KubeServer              string
	Slack                   SlackConfig
	OverwriteRole           bool
	OverwriteRoleBinding    bool
}

func GetConfig() *Config {
	return _instance
}

func LoadConfig() *Config {
	viper.SetConfigName("cinderella") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name

	viper.SetDefault("Namespace", "default")
	viper.SetDefault("CleanupTickEverySeconds", 10)
	viper.SetDefault("ManifestDirs", []string{"/etc/cinderella"})
	viper.SetDefault("OverwriteRole", false)
	viper.SetDefault("OverwriteRoleBinding", false)

	viper.AddConfigPath("/etc/cinderella/") // path to look for the config file in
	viper.AddConfigPath(".")                // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if err := viper.Unmarshal(_instance); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v , %w", _instance, err))
	}

	for i, dir := range _instance.ManifestDirs {
		if strings.HasPrefix(dir, "~/") {
			homedir, _ := os.UserHomeDir()
			dir = homedir + dir[1:]
			_instance.ManifestDirs[i] = dir
		}
	}

	return _instance
}

// SearchManifest is search manifest template from manifest directorys
// name argument is file name without extension
func SearchManifest(name string) string {
	return searchManifest(name, GetConfig().ManifestDirs)
}

// SearchManifest is search manifest template from manifest directorys
// name argument is file name without extension
func searchManifest(name string, manifestDirs []string) string {
	for _, dir := range manifestDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("%s is not existing\n", dir)
			continue
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Failed to read dir %s", dir)
		}

		for _, file := range files {
			tmpl := name + TemplateExtension
			if file.Name() == tmpl {
				return fmt.Sprintf("%s/%s", dir, file.Name())
			}
		}
	}

	fmt.Printf("%s not found\n", name)
	return ""
}
