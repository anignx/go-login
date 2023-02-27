package manager

import (
	"context"
	"net/http"

	"my.service/go-login/conf"
	buzzhttp "my.service/go-login/package/daenerys/http"
)

// Manager represents middleware component
// such as, kafka, http client or rpc client, etc.
type Manager struct {
	c          *conf.Config
	baseClient *http.Client
}

func New(conf *conf.Config) *Manager {
	return &Manager{
		c:          conf,
		baseClient: buzzhttp.InitHttp("buzz.app.go-login2"),
	}
}

func (m *Manager) Ping(ctx context.Context) error {
	return nil
}

func (m *Manager) Close() error {
	return nil
}
