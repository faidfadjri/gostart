package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func GetURLParamInt(r *http.Request, key string) (int, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, fmt.Errorf("parameter %s is required", key)
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be a valid integer", key)
	}

	return value, nil
}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryParamInt(r *http.Request, key string) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return 0, fmt.Errorf("query parameter %s is required", key)
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("query parameter %s must be a valid integer", key)
	}

	return value, nil
}

func GetQueryParamWithDefault(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}