package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/daniilty/secrets-manager/internal/core"
	"github.com/daniilty/secrets-manager/internal/metrics"
)

type secretResp struct {
	jsonSecret string
}

type secretReq struct {
	AppName string `json:"app_name"`
	Secret  string `json:"secret"`
}

func (s *secretReq) validate() error {
	if s.AppName == "" {
		return fmt.Errorf("app_name: cannot be empty")
	}

	if s.Secret == "" {
		return fmt.Errorf("secret: cannot be empty")
	}

	return nil
}

func (s *secretResp) writeJSON(w http.ResponseWriter) error {
	if s.jsonSecret == "" {
		w.WriteHeader(http.StatusNoContent)

		return nil
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(s.jsonSecret))

	return err
}

func (h *HTTP) getSecretHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.getSecretResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getSecretResponse(r *http.Request) response {
	const appNameParamName = "app_name"

	timeNow := time.Now()

	appName := r.FormValue(appNameParamName)
	if appName == "" {
		return getBadRequestWithMsgResponse(appNameParamName + ": cannot be empty")
	}

	secret, err := h.service.Get(r.Context(), appName)
	if err != nil {
		h.logger.Errorw("Get secret.", "app_name", appName, "err", err)
		metrics.TotalErrReqs.With(map[string]string{
			"method":     r.Method,
			"error_code": strconv.Itoa(http.StatusInternalServerError),
		}).Inc()
		metrics.ReqDurationSeconds.With(map[string]string{
			"method":   r.Method,
			"is_error": "true",
		}).Observe(time.Since(timeNow).Seconds())

		return getInternalServerErrorResponse()
	}

	metrics.TotalReqs.With(map[string]string{
		"method": r.Method,
	}).Inc()
	metrics.ReqDurationSeconds.With(map[string]string{
		"method":   r.Method,
		"is_error": "false",
	}).Observe(time.Since(timeNow).Seconds())

	return &secretResp{
		jsonSecret: secret,
	}
}

func (h *HTTP) setSecretHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.setSecretResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) setSecretResponse(r *http.Request) response {
	req := &secretReq{}

	timeNow := time.Now()

	err := unmarshalReader(r.Body, req)
	if err != nil {
		metrics.TotalErrReqs.With(map[string]string{
			"method":     r.Method,
			"error_code": strconv.Itoa(http.StatusBadRequest),
		}).Inc()
		metrics.ReqDurationSeconds.With(map[string]string{
			"method":   r.Method,
			"is_error": "true",
		}).Observe(time.Since(timeNow).Seconds())

		return getBadRequestWithMsgResponse(err.Error())
	}

	err = req.validate()
	if err != nil {
		metrics.TotalErrReqs.With(map[string]string{
			"method":     r.Method,
			"error_code": strconv.Itoa(http.StatusBadRequest),
		}).Inc()
		metrics.ReqDurationSeconds.With(map[string]string{
			"method":   r.Method,
			"is_error": "true",
		}).Observe(time.Since(timeNow).Seconds())

		return getBadRequestWithMsgResponse(err.Error())
	}

	err = h.service.Set(r.Context(), req.AppName, req.Secret)
	if err != nil {
		if errors.Is(err, core.ErrInvalidJSON) {
			metrics.TotalErrReqs.With(map[string]string{
				"method":     r.Method,
				"error_code": strconv.Itoa(http.StatusBadRequest),
			}).Inc()
			metrics.ReqDurationSeconds.With(map[string]string{
				"method":   r.Method,
				"is_error": "true",
			}).Observe(time.Since(timeNow).Seconds())

			return getBadRequestWithMsgResponse(err.Error())
		}

		metrics.TotalErrReqs.With(map[string]string{
			"method":     r.Method,
			"error_code": strconv.Itoa(http.StatusInternalServerError),
		}).Inc()
		metrics.ReqDurationSeconds.With(map[string]string{
			"method":   r.Method,
			"is_error": "true",
		}).Observe(time.Since(timeNow).Seconds())

		h.logger.Errorw("Get secret.", "req", req, "err", err)

		return getInternalServerErrorResponse()
	}

	return &secretResp{
		jsonSecret: req.Secret,
	}
}
