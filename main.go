package main

import (
	"backend/contrib/routes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/gofiber/websocket/v2"
	
	"backend/config"
	"backend/core"
	"backend/database"
	_ "backend/docs"

	"github.com/gofiber/swagger"

	"backend/contrib/controllers"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Detail string `json:"detail"`
	}
)

// @title Openapi spec
// @version 0.1.0
// @description This is a sample Fiber API with OAuth2 password flow
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/auth/login
// @scope.read Grants read access
// @scope.write Grants write access
func main() {
	// Load configuration
	config.LoadConfig()

	log.SetLevel(log.LevelDebug)
	// Load database
	db, err := database.ConnectDB(config.AppConfig.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to databse: %v", err)
	}

	// Todo enable vip plates load
	//util.LoadVIPPlates()

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Log the error (optional, but recommended)
			fmt.Printf("Error: %v\n", err)

			// Handle ValidationErrorResponse
			var valErr *core.ValidationErrorResponse
			if errors.As(err, &valErr) {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(valErr)
			}

			// Handle Fiber internal errors (like fiber.ErrNotFound)
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				return c.Status(fiberErr.Code).JSON(GlobalErrorHandlerResp{
					Detail: fiberErr.Message,
				})
			}

			// Fallback to 500 Internal Server Error
			return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
				Detail: "Internal Server Error",
			})
		},
		StrictRouting: false, // Enable strict routing for trailing slashes
	})
	app.Use(logger.New())
	app.Use(cors.New(
		cors.Config{
			//AllowOrigins:     strings.Join(config.AppConfig.CORSAllowOrigins, ", "),
			AllowOrigins:     "*",
			AllowCredentials: false,
			AllowHeaders:     "Origin, Content-Type, Accept",

			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodHead,
				fiber.MethodPut,
				fiber.MethodDelete,
				fiber.MethodPatch,
				fiber.MethodOptions,
			}, ","),
		},
	))
	go controllers.HandleMessages()

	app.Static("/swagger/swagger.yaml", "./docs/swagger.yaml")

	baseUrl := "http://" + config.AppConfig.BackendHost + ":" + strconv.Itoa(config.AppConfig.BackendPort)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/swagger.yaml",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: os.Getenv("SECRET_KEY_JWT"),
		},
		OAuth2RedirectUrl: baseUrl + "/swagger/oauth2-redirect.html",
	}))
	routes.ApiV1Route(app, db)
	app.Get("/ws/car-session/", websocket.New(controllers.CarSessionWebsocket))
	//go operator.HandleMessages()
	//go imagetoplate.WatchDirectory("image", database.DB)
	//routes.InitAdminRoute(app)
	//routes.CameraRoutes(app)
	//routes.AccountantRoutes(app)
	//routes.InitZReport(app)
	//routes.InitRealtime(app)
	//routes.FixRoute(app)
	//routes.Init(app)

	// Initialize default config (Assign the middleware to /metrics)
	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":" + strconv.Itoa(config.AppConfig.BackendPort)))

}
