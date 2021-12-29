package apiclient

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hestingames/hg-hebe-bot/internal/environment"
	"github.com/hestingames/hg-hebe-bot/internal/logs"
)

var (
	logger *logs.Logger
	ctx    context.Context
)

func InitializeApiClient() {
	ctx = context.Background()
	logger = logs.FromContext(ctx)
}

func DisableCertificateCheck() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if !environment.IsLocal() {
		logger.Warn("SSL certificate check is disabled")
	}
}
