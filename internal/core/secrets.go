package core

import (
	"context"
	"encoding/json"
)

func (s *secretsManager) Set(ctx context.Context, appName string, val string) error {
	if !isValidJSON(val) {
		return ErrInvalidJSON
	}

	return s.db.Set(ctx, appName, val)
}

func (s *secretsManager) Get(ctx context.Context, appName string) (string, error) {
	return s.db.Get(ctx, appName)
}

func isValidJSON(str string) bool {
	err := json.Unmarshal([]byte(str), &map[string]interface{}{})

	return err == nil
}
