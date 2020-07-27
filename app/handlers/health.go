package handlers

import (
	"net/http"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/netx/httpx"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJson(w, http.StatusOK, http.StatusText(http.StatusOK))
}
