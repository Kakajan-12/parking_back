package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"
	"backend/core"
)

func UserSessionRoute(app fiber.Router, db *gorm.DB) {
	userSessionGroup := app.Group("/user-session", middleware.AuthToken(db, []models.RoleType{models.AdminRole}, true))

	// Wrap middleware with db
	userSessionGroup.Get("/", func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveUserSessionsApi(c, db, commons)
	}).Name("user-session-list")
	userSessionGroup.Get("/:id/detail/", func(c *fiber.Ctx) error {
		return controllers.UserSessionDetailApi(c, db)
	}).Name("user-session-detail")
	userSessionGroup.Get("/:id/revoke/", func(c *fiber.Ctx) error {
		return controllers.UserSessionRevokeApi(c, db)
	}).Name("user-session-revoke")
}
