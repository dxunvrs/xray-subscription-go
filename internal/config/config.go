package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	XrayServerIP           string `env:"XRAY_SERVER_IP" env-default:"127.0.0.1"`
	XrayServerPort         int    `env:"XRAY_SERVER_PORT" env-default:"443"`
	XrayRealityPublicKey   string `env:"XRAY_REALITY_PUBLIC_KEY"`
	XrayRealityShortID     string `env:"XRAY_REALITY_SHORT_ID"`
	XrayRealitySNI         string `env:"XRAY_REALITY_SNI"`
	XrayRealityFingerprint string `env:"XRAY_REALITY_FINGERPRINT"`
	XrayGrpcHost           string `env:"XRAY_GRPC_HOST"`
	XrayGrpcPort           int    `env:"XRAY_GRPC_PORT" env-default:"10058"`

	AdminUsername string `env:"ADMIN_USERNAME"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
}

func Load() *Config {
	var config Config

	err := cleanenv.ReadConfig(".env", &config)
	if err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}

	return &config
}
