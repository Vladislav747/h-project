package openapi

import (
	"net/http"
)

/*
 * if /healthz path returns a success code, the kubelet considers the container to be alive and healthy
 * if the handler returns a failure code, the kubelet kills the container and restarts it
 */
func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
