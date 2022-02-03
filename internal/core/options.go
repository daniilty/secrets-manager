package core

import "github.com/daniilty/secrets-manager/internal/db"

// SecretsManagerOption - DI configuration option.
type SecretsManagerOption func(*secretsManager)

// WithDB - set db.
func WithDB(db db.DB) SecretsManagerOption {
	return func(s *secretsManager) {
		s.db = db
	}
}
