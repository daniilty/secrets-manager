package core

import (
	"context"

	"github.com/daniilty/secrets-manager/internal/db"
)

// SecretsManager - core service logic.
type SecretsManager interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
}

type secretsManager struct {
	db db.DB
}

func NewSecretsManager(opts ...SecretsManagerOption) SecretsManager {
	s := &secretsManager{}

	for i := range opts {
		opts[i](s)
	}

	return s
}
