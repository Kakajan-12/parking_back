package routes

import (
	"backend/contrib/middleware"
	"backend/contrib/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/contrib/controllers"
	"backend/core"
)

func CameraRoute(app fiber.Router, db *gorm.DB) {
	cameraGroup := app.Group("/camera")

	// Wrap middleware with db
	cameraGroup.Get("/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveCamerasApi(c, db, commons)
	}).Name("camera-list")
	cameraGroup.Get("/count/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		return controllers.CountCamerasApi(c, db)
	}).Name("camera-count")
	cameraGroup.Post("/create/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		return controllers.CreateCameraApi(c, db)
	}).Name("camera-create")
	cameraGroup.Get("/:id/detail/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		return controllers.CarDetailApi(c, db)
	}).Name("camera-detail")
	cameraGroup.Patch("/:id/update/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		return controllers.UpdateCameraApi(c, db)
	}).Name("camera-update")
	cameraGroup.Delete("/:id/delete/", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true), func(c *fiber.Ctx) error {
		return controllers.CameraDeleteApi(c, db)
	})

}
