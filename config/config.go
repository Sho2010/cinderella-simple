package config

import (
	"time"
)

var instance *Config

func init() {

	// //TODO: 環境変数による上書き処理
	instance = &Config{
		CleanupTickEvery: 10 * time.Second,
		Namespace:        "",
		ManifestDirs:     []string{},
		KubeServer:       "",
	}
}

type Config struct {
	CleanupTickEvery time.Duration
	Namespace        string
	ManifestDirs     []string
	KubeServer       string
}

func GetConfig() *Config {
	return instance
}
