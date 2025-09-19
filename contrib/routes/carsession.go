package routes

import (
	"backend/contrib/controllers"
	"backend/contrib/middleware"
	"backend/contrib/models"
	"backend/core"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CarSessionRoute(app fiber.Router, db *gorm.DB) {
	group := app.Group("/session-car")

	// Wrap middleware with db
	group.Get("/", 
	// middleware.AuthToken(
	// 	db,
	// 	[]models.RoleType{models.AdminRole, models.AccountantRole, models.OperatorRole},
	// 	true,
	// ),
	 func(c *fiber.Ctx) error {
		commons := core.GetCommonsFromContext(c)
		return controllers.RetrieveCarSessionsApi(c, db, commons)
	}).Name("car-session-list")

	group.Post(
		"/event/entry/",
		func(c *fiber.Ctx) error {

			return controllers.CarSessionEventEntryApi(c, db)
		},
	).Name("car-session-event-entry")
	
	group.Post(
		"/event/exit/",
		func(c *fiber.Ctx) error {
			return controllers.CarSessionEventExitApi(c, db)
		},
	).Name("car-session-event-exit")
	
	group.Get(
		"/count/",
		middleware.AuthToken(
			db,
			[]models.RoleType{models.AdminRole,
				models.AccountantRole, models.OperatorRole}, true,
		),
		func(c *fiber.Ctx) error {
			return controllers.CountCarSessionsApi(c, db)
		}).Name("car-session-count")
}
