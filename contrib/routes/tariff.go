package routes

import (
	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"
	"backend/core"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TariffRoute(app fiber.Router, db *gorm.DB) {
	tariffGroup := app.Group("/tariff", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true))

	// Wrap middleware with db
	tariffGroup.Get("/", func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveTariffsApi(c, db, commons)
	}).Name("tariff-list")
	tariffGroup.Get("/count/", func(c *fiber.Ctx) error {
		return controllers.CountTariffsApi(c, db)
	}).Name("tariff-count")
	tariffGroup.Post("/create/", func(c *fiber.Ctx) error {
		return controllers.CreateTariffApi(c, db)
	}).Name("tariff-create")
	tariffGroup.Get("/:id/detail/", func(c *fiber.Ctx) error {
		return controllers.TariffDetailApi(c, db)
	}).Name("tariff-detail")
	tariffGroup.Patch("/:id/update/", func(c *fiber.Ctx) error {
		return controllers.UpdateTariffApi(c, db)
	}).Name("tariff-update")
}
