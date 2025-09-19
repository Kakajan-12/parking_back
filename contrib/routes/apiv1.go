package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApiV1Route(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api/v1") // /api

	AuthRoute(api, db)
	CameraRoute(api, db)
	CarRoute(api, db)
	CarSessionRoute(api, db)
	TariffRoute(api, db)
	UserRoute(api, db)
	UserSessionRoute(api, db)
}
