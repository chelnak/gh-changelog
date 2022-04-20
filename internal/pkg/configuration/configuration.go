package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	home, _ := os.UserHomeDir()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	cfgPath := filepath.Join(home, ".config", "gh-changelog")
	viper.AddConfigPath(cfgPath)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgPath, 0750); err != nil {
			return fmt.Errorf("failed to create config directory: %s", err)
		}
	}

	SetDefaults()
	err := viper.SafeWriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write config: %s", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %s", err)
	}

	return nil
}

func SetDefaults() {
	viper.SetDefault("file_name", "CHANGELOG.md")
	viper.SetDefault("excluded_labels", []string{"maintenance"})

	sections := make(map[string][]string)
	sections["Changed"] = []string{"backwards-incompatible"}
	sections["Added"] = []string{"feature", "enhancement"}
	sections["Fixed"] = []string{"bug", "bugfix", "documentation"}

	viper.SetDefault("sections", sections)

	viper.SetDefault("skip_entries_without_label", false)
}
