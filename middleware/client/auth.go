package client

import (
	"context"

	"github.com/acd19ml/TalentRank/middleware/server"
)

func NewAuthentication(ak, sk string) *Authentication {
	return &Authentication{
		clientId:     ak,
		clientSecret: sk,
	}
}

type Authentication struct {
	clientId     string
	clientSecret string
}

func (a *Authentication) build() map[string]string {
	return map[string]string{
		server.ClientHeaderAccessKey: a.clientId,
		server.ClientHeaderSecretKey: a.clientSecret,
	}
}

func (a *Authentication) GetRequestMetadata(
	ctx context.Context, uri ...string) (
	map[string]string, error) {
	return a.build(), nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
