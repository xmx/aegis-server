//go:build !windows

package restapi

import "github.com/labstack/echo/v5"

type SystemService struct{}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (ps *SystemService) RegisterRoute(r *echo.Group) error {
	return nil
}
