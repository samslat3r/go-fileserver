package config

import (
	"time"
	"github.com/spf13/viper"
)

type App struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	StorageDir string        `mapstructure:"storage_dir"`
	Verbose   bool          `mapstructure:"verbose"`
	SessionKey string      `mapstructure:"session_key"`
	Adminuser string      `mapstructure:"admin_user"`
	BcryptHash string      `mapstructure:"bcrypt_hash"`
	TokenTTL   time.Duration `mapstructure:"token_ttl"`
}

func Load(cfgPath string, ovverrides map[string]any) (*App, error) {
	v := viper.New()

	if cfgPath != "" {
		v.SetConfigFile(cfgPath)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/appname/")
	}
	v.SetConfigType("yaml")

	// Default values
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 8080)
	v.SetDefault("token_ttl", "24h")

	// Apply overrides
	for k, v2 := range ovverrides {
		v.Set(k, v2)
	}

	var c App
	if err := v.Unmarshal(&c); err != nil {
		return nil, err 
	}
	return &c, nil
}

