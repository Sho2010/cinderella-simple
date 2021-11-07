package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	viper.AddConfigPath("./testdata/") // optionally look for config in the working directory
	LoadConfig()

	code := m.Run()
	os.Exit(code)
}

func Testチルダを展開する(t *testing.T) {
	result := GetConfig().ManifestDirs[0]
	homedir, _ := os.UserHomeDir()
	homedir = homedir + "/_example_example_manifest_dir"
	assert.Equal(t, result, homedir)
}

func TestSearchManifest(t *testing.T) {
	result := SearchManifest("default-role")
	assert.Equal(t, "testdata/manifest01/default-role.yaml.tmpl", result)

	result = SearchManifest("only-02-dir-file")
	assert.Equal(t, "testdata/manifest02/only-02-dir-file.yaml.tmpl", result)

	result = SearchManifest("both-exists-file")
	assert.Equal(t, "testdata/manifest01/both-exists-file.yaml.tmpl", result, "複数のディレクトリに存在した場合、最初に見つかったファイルを返す")

}
