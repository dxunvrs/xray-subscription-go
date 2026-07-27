package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerPort int `env:"SERVER_PORT" env-default:"12258"`
	Xray       struct {
		ServerIP           string `env:"XRAY_SERVER_IP" env-default:"127.0.0.1"`
		ServerPort         int    `env:"XRAY_SERVER_PORT" env-default:"443"`
		RealityPublicKey   string `env:"XRAY_REALITY_PUBLIC_KEY"`
		RealityShortID     string `env:"XRAY_REALITY_SHORT_ID"`
		RealitySNI         string `env:"XRAY_REALITY_SNI"`
		RealityFingerprint string `env:"XRAY_REALITY_FINGERPRINT"`
		GrpcHost           string `env:"XRAY_GRPC_HOST"`
		GrpcPort           int    `env:"XRAY_GRPC_PORT" env-default:"10058"`
	}

	Admin struct {
		Username string `env:"ADMIN_USERNAME"`
		Password string `env:"ADMIN_PASSWORD"`
	}
}

func Load() *Config {
	var config Config

	err := cleanenv.ReadConfig(".env", &config)
	if err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}

	return &config
}
