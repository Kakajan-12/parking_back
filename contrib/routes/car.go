package routes

import (
	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"
	"backend/core"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CarRoute(app fiber.Router, db *gorm.DB) {
	carGroup := app.Group("/car", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true))

	// Wrap middleware with db
	carGroup.Get("/", func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveCarsApi(c, db, commons)
	}).Name("car-list")
	carGroup.Get("/count/", func(c *fiber.Ctx) error {
		return controllers.CountCarsApi(c, db)
	}).Name("car-count")
	carGroup.Post("/create/", func(c *fiber.Ctx) error {
		return controllers.CreateCarApi(c, db)
	}).Name("car-create")
	carGroup.Get("/:id/detail/", func(c *fiber.Ctx) error {
		return controllers.CarDetailApi(c, db)
	}).Name("car-detail")
	carGroup.Patch("/:id/update/", func(c *fiber.Ctx) error {
		return controllers.UpdateCarApi(c, db)
	}).Name("car-update")

	
}
