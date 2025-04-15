package config

import (
	"github.com/spf13/viper"

	"vpn-ads-router/pkg/logger"
)

// -- config file names --
type Config struct {
	Proxy       ProxyConfig       `mapstructure:"proxy"`
	PLC         PLCConfig         `mapstructure:"plc"`
	Fingerprint FingerprintConfig `mapstructure:"fingerprint"`
}

type ProxyConfig struct {
	EthernetInterface string `mapstructure:"ethernetInterface"`
	StaticNetidSuffix string `mapstructure:"staticNetidSuffix"`
}

type PLCConfig struct {
	Credentials Credentials `mapstructure:"credentials"`
}

type FingerprintConfig struct {
	Subnets        []string         `mapstructure:"subnets"`
	PlcFingerprint []PlcFingerprint `mapstructure:"ports"`
}

// -- sub names --
type Credentials struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type PlcFingerprint struct {
	Port     int    `mapstructure:"port"`
	Label    string `mapstructure:"label"`
	Required bool   `mapstructure:"required"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/")

	viper.AutomaticEnv() // Optional: allow ENV vars to override config file

	if err := viper.ReadInConfig(); err != nil {
		logger.GlobalLogger.Error(logger.ComponentService, "Error reading config file: %v", err)
		return err
	} else {
		logger.GlobalLogger.Info(logger.ComponentService, "Config file found: %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.GlobalLogger.Error(logger.ComponentService, "Error unmarshalling config file: %v", err)
		return err
	} else {
		logger.GlobalLogger.Info(logger.ComponentService, "Config file unmarshalled successfully")
	}

	logger.GlobalLogger.Info(logger.ComponentService, "Config loaded successfully")
	return nil
}
