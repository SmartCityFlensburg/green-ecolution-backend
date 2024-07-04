package http

import (
	"github.com/SmartCityFlensburg/green-space-management/internal/server/http/handler/info"
	"github.com/SmartCityFlensburg/green-space-management/internal/server/http/handler/tree"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) router() *fiber.App {
	app := fiber.New()

	app.Mount("/info", info.RegisterRoutes(s.services.InfoService))
  app.Mount("/tree", tree.RegisterRoutes(s.services.TreeService))

	return app
}
