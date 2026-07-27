package service

import (
	"fmt"
	"net/url"
	"xray-subscription-go/internal/config"
)

type VlessLinkBuilder struct {
	config *config.Config
}

func NewVlessLinkBuilder(config *config.Config) *VlessLinkBuilder {
	return &VlessLinkBuilder{config: config}
}

func (b *VlessLinkBuilder) BuildVlesLink(userUUID, userEmail string) string {
	u := &url.URL{
		Scheme:   "vless",
		User:     url.User(userUUID),
		Host:     fmt.Sprintf("%s:%d", b.config.XrayServerIP, b.config.XrayServerPort),
		Fragment: userEmail,
	}

	query := u.Query()
	query.Set("type", "tcp")
	query.Set("security", "reality")
	query.Set("pbk", b.config.XrayRealityPublicKey)
	query.Set("fp", b.config.XrayRealityFingerprint)
	query.Set("sni", b.config.XrayRealitySNI)
	query.Set("sid", b.config.XrayRealityShortID)
	query.Set("flow", "xtls-rprx-vision")

	u.RawQuery = query.Encode()

	return u.String()
}
