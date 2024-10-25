package service

import (
	"context"
	"net/http"
)

// ServiceStatus ???
func (s *service) ServiceStatus(_ context.Context) int {
	return http.StatusOK
}
