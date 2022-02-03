package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTP) setRoutes(r *mux.Router) {
	const secretsPath = "/secrets"

	r.HandleFunc(secretsPath,
		h.getSecretHandler,
	).Methods(http.MethodGet)

	r.HandleFunc(secretsPath,
		h.setSecretHandler,
	).Methods(http.MethodPost)
}
