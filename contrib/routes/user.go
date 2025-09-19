package routes

import (
	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"
	"backend/core"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoute(app fiber.Router, db *gorm.DB) {
	userGroup := app.Group("/user", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true))

	// Wrap middleware with db
	userGroup.Get("/", func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveUsersApi(c, db, commons)
	}).Name("user-list")
	userGroup.Get("/count/", func(c *fiber.Ctx) error {
		return controllers.CountUsersApi(c, db)
	}).Name("user-count")
	userGroup.Post("/create/", func(c *fiber.Ctx) error {
		return controllers.CreateUserApi(c, db)
	}).Name("user-create")
	userGroup.Get("/:id/detail/", func(c *fiber.Ctx) error {
		return controllers.UserDetailApi(c, db)
	}).Name("user-detail")
	userGroup.Patch("/:id/update/", func(c *fiber.Ctx) error {
		return controllers.UpdateUserApi(c, db)
	}).Name("user-update")
	userGroup.Delete("/:id/delete/", func(c *fiber.Ctx) error {
		return controllers.UserDeleteApi(c, db)
	})
}
