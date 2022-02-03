package healthcheck

import "github.com/daniilty/secrets-manager/internal/db"

// CheckerOption - di option.
type CheckerOption func(*checker)

// WithMongoDBPinger - set mongo db pinger.
func WithMongoDBPinger(pinger db.Pinger) func(*checker) {
	return func(c *checker) {
		c.mongoPinger = pinger
	}
}
