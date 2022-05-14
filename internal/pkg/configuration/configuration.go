//Package configuration contains a number of methods that are used
//to provide configuration to the wider application. It uses viper
//to pull config from either the environment or a config file then
//unmarhsals the config into the configuration struct. The configuration struct
//is made available to the application via a package level variable
//called Config.
package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Config configuration

type configuration struct {
	FileName                string              `mapstructure:"file_name"`
	ExcludedLabels          []string            `mapstructure:"excluded_labels"`
	Sections                map[string][]string `mapstructure:"sections"`
	SkipEntriesWithoutLabel bool                `mapstructure:"skip_entries_without_label"`
	ShowUnreleased          bool                `mapstructure:"show_unreleased"`
	CheckForUpdates         bool                `mapstructure:"check_for_updates"`
}

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

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		err := viper.SafeWriteConfig()
		if err != nil {
			return fmt.Errorf("failed to write config: %s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("changelog")

	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("file_name", "CHANGELOG.md")
	viper.SetDefault("excluded_labels", []string{"maintenance"})

	sections := make(map[string][]string)
	sections["changed"] = []string{"backwards-incompatible"}
	sections["added"] = []string{"feature", "enhancement"}
	sections["fixed"] = []string{"bug", "bugfix", "documentation"}

	viper.SetDefault("sections", sections)

	viper.SetDefault("skip_entries_without_label", false)

	viper.SetDefault("show_unreleased", true)

	viper.SetDefault("check_for_updates", true)
}
