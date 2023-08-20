package http

import (
	orderHttp "goshop/internal/order/port/http"
	productHttp "goshop/internal/product/port/http"
	userHttp "goshop/internal/user/port/http"
)

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	userHttp.Routes(v1, s.db, s.validator)
	productHttp.Routes(v1, s.db, s.validator, s.cache)
	orderHttp.Routes(v1, s.db, s.validator)
	return nil
}
