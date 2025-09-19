package routes

import (
	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(api fiber.Router, db *gorm.DB) {
	authGroup := api.Group("/auth")

	authGroup.Post("/login/", func(c *fiber.Ctx) error {
		return controllers.LoginApi(c, db)
	}).Name("login")

	// Wrap middleware with db
	authGroup.Get("/logout/", middleware.AuthToken(db, []models.RoleType{}, false), func(c *fiber.Ctx) error {
		return controllers.LogoutApi(c, db)
	}).Name("logout")
	authGroup.Get("/me/", middleware.AuthToken(db, []models.RoleType{}, false), func(c *fiber.Ctx) error {
		return controllers.MeApi(c)
	}).Name("me")

}
