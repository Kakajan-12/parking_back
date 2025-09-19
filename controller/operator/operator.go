package operator

//import (
//	"fmt"
//	"math"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/gofiber/fiber/v2"
//	"gorm.io/gorm"
//
//	"park/database"
//	modelscar "park/models/modelsCar"
//)
//
//const statusExited = "Exited"
//
//type GetCarsResponse struct {
//	Results    []modelscar.CarModel `json:"results"`
//	Page       int                  `json:"page"`
//	Limit      int                  `json:"limit"`
//	TotalPages int                  `json:"totalPages"`
//	HasNext    bool                 `json:"hasNext"`
//	HasPrev    bool                 `json:"hasPrev"`
//}
//
//// GetCars godoc
//// @Summary Get list of cars
//// @Description Get list of cars with pagination
//// @Tags cars
//// @Accept  json
//// @Produce  json
//// @Param page query int false "Page number" default(1)
//// @Param limit query int false "Number of items per page" default(5)
//// @Success 200 {object} GetCarsResponse
//// @Failure 400 {object} ErrorResponse
//// @Router /api/v1/get-all-cars [get]
//func GetCars(c *fiber.Ctx) error {
//	var cars []modelscar.CarModel
//	var totalCount int64
//	pageStr := c.Query("page", "1")
//	limitStr := c.Query("limit", "5")
//	parkNumber := c.Locals("parkNumber")
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page < 1 {
//		return c.Status(400).JSON(fiber.Map{
//			"message": "Invalid page number",
//		})
//	}
//
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil || limit < 1 {
//		return c.Status(400).JSON(fiber.Map{
//			"message": "Invalid limit number",
//		})
//	}
//
//	query := database.DB.Model(&modelscar.CarModel{})
//	if parkNumber != "" {
//		query = query.Where("park_no = ?", parkNumber)
//	}
//
//	query.Count(&totalCount)
//	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
//	hasNext := page < totalPages
//	hasPrev := page > 1
//
//	offset := (page - 1) * limit
//	query.Order("id desc").Limit(limit).Offset(offset).Find(&cars)
//
//	if len(cars) == 0 {
//		cars = []modelscar.CarModel{}
//	}
//	ip := os.Getenv("HOST")
//	port := os.Getenv("PORT")
//
//	for i := range cars {
//		cars[i].ImageUrl = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, cars[i].ImageUrl)
//	}
//	return c.Status(200).JSON(fiber.Map{
//		"results":    cars,
//		"page":       page,
//		"limit":      limit,
//		"totalPages": totalPages,
//		"hasNext":    hasNext,
//		"hasPrev":    hasPrev,
//	})
//}
//
//// GetCar godoc
//// @Summary Get a car by ID
//// @Description Get a car by ID
//// @Tags cars
//// @Accept  json
//// @Produce  json
//// @Param id path int true "Car ID"
//// @Success 200 {object} modelscar.CarModel
//// @Failure 404 {object} ErrorResponse
//// @Router /api/v1/get-car/{id} [get]
//func GetCar(c *fiber.Ctx) error {
//	id := c.Params("id")
//	var car modelscar.CarModel
//	database.DB.Where("id = ?", id).First(&car)
//	if car.ID == 0 {
//		return c.Status(404).JSON(fiber.Map{
//			"message": "Car not found",
//		})
//	}
//	ip := os.Getenv("HOST")
//	port := os.Getenv("PORT")
//
//	car.ImageUrl = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, car.ImageUrl)
//
//	c.Status(200)
//	return c.JSON(car)
//}
//
//type UpdateCarResponse struct {
//	Message string             `json:"message"`
//	Data    modelscar.CarModel `json:"data"`
//}
//
//// UpdateCar godoc
//// @Summary Update a car by plate number
//// @Description Updates a car's status and calculates payment and duration based on start and end times.
//// @Tags cars
//// @Accept  json
//// @Produce  json
//// @Param plate path string true "Car plate number"
//// @Param car body modelscar.CarUpdate true "Car details to update"
//// @Success 200 {object} UpdateCarResponse "Updated car details"
//// @Failure 400 {object} ErrorResponse "Car already exited or invalid request"
//// @Failure 404 {object} ErrorResponse "Car not found"
//// @Failure 500 {object} ErrorResponse "Error parsing time"
//// @Router /api/v1/camera/update-car/{plate} [put]
//func UpdateCar(c *fiber.Ctx) error {
//	plate := c.Params("plate")
//	userIDVal := c.Locals("username")
//
//	var car modelscar.CarModel
//	if err := database.DB.Order("id desc").Where("car_number = ?", plate).First(&car).Error; err != nil {
//		return c.Status(404).JSON(fiber.Map{"message": "Car not found", "error": err.Error()})
//	}
//	if car.Status == statusExited {
//		return c.Status(400).JSON("Car already Exited")
//	}
//
//	var updatedCar modelscar.CarModel
//	if err := c.BodyParser(&updatedCar); err != nil {
//		return c.Status(400).JSON(fiber.Map{"message": "Invalid request", "error": err.Error()})
//	}
//
//	updatedCar.Status = statusExited
//
//	if updatedCar.Reason == "" {
//		updatedCar.Reason = "Payment succeeded"
//		updatedCar.TotalPayment = car.TotalPayment
//	} else {
//		updatedCar.TotalPayment = 0
//	}
//
//	car.Reason = updatedCar.Reason
//	car.TotalPayment = updatedCar.TotalPayment
//	updatedCar.EndTime = car.EndTime
//
//	userId, ok := userIDVal.(string)
//	if !ok {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"message": "Internal Server Error - Invalid user ID type",
//		})
//	}
//	updatedCar.UserID = userId
//
//	updates := map[string]interface{}{
//		"reason":        updatedCar.Reason,
//		"total_payment": updatedCar.TotalPayment,
//		"user_id":       updatedCar.UserID,
//		"status":        updatedCar.Status,
//		"end_time":      updatedCar.EndTime,
//	}
//
//	if err := database.DB.Model(&car).Updates(updates).Error; err != nil {
//		return c.Status(500).JSON(fiber.Map{"message": "Database update failed", "error": err.Error()})
//	}
//
//	updatedCar.ID = car.ID
//	updatedCar.CarNumber = car.CarNumber
//	updatedCar.StartTime = car.StartTime
//	updatedCar.ParkNumber = car.ParkNumber
//	updatedCar.EndTime = car.EndTime
//	updatedCar.ImageUrl = car.ImageUrl
//
//	return c.Status(200).JSON(fiber.Map{
//		"message": "Car updated successfully",
//		"data":    updatedCar,
//	})
//}
//
//// SearchCar godoc
//// @Summary Search for cars
//// @Description Retrieve a paginated list of cars with optional filtering by car number, enter time range, end time range, park number, and status.
//// @Tags cars
//// @Accept json
//// @Produce json
//// @Param carNumber query string false "Filter by car plate number (partial match allowed)"
//// @Param enterTime query string false "Start of enter time range (YYYY-MM-DD)"
//// @Param endTime query string false "End of end time range (YYYY-MM-DD)"
//// @Param parkNumber query string false "Filter by parking spot number"
//// @Param status query string false "Filter by car status (Inside, Exited)"
//// @Param page query int false "Page number" default(1)
//// @Param limit query int false "Number of items per page" default(5)
//// @Success 200 {object} GetCarsResponse
//// @Failure 400 {object} ErrorResponse
//// @Router /api/v1/search-car [get]
//func SearchCar(c *fiber.Ctx) error {
//	var cars []modelscar.CarModel
//	var totalCount int64
//
//	carNumber := c.Query("carNumber")
//	enterTime := c.Query("enterTime")
//	endTime := c.Query("endTime")
//	parkNumber, _ := c.Locals("parkNumber").(string)
//	status := c.Query("status")
//	pageStr := c.Query("page", "1")
//	limitStr := c.Query("limit", "5")
//
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page < 1 {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid page number"})
//	}
//
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil || limit < 1 {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid limit number"})
//	}
//
//	baseQuery := database.DB.Model(&modelscar.CarModel{}).Debug()
//
//	if carNumber != "" {
//		baseQuery = baseQuery.Where("car_number LIKE ?", "%"+strings.TrimSpace(carNumber)+"%")
//	}
//
//	if enterTime != "" {
//		if _, err := time.Parse("2006-01-02", enterTime); err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid enter_time format. Use YYYY-MM-DD."})
//		}
//		baseQuery = baseQuery.Where("start_time >= ?", enterTime+" 00:00:00")
//	}
//
//	if endTime != "" {
//		if _, err := time.Parse("2006-01-02", endTime); err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid end_time format. Use YYYY-MM-DD."})
//		}
//		baseQuery = baseQuery.Where("end_time <= ?", endTime+" 23:59:59")
//	}
//
//	if parkNumber != "" {
//		baseQuery = baseQuery.Where("park_number = ?", parkNumber)
//	}
//
//	if status != "" {
//		validStatuses := map[string]bool{"Inside": true, "Exited": true}
//		if !validStatuses[status] {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid status. Use Inside or Exited."})
//		}
//		baseQuery = baseQuery.Where("status = ?", status)
//	}
//
//	if err := baseQuery.Count(&totalCount).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"message": "Error counting cars",
//			"error":   err.Error(),
//		})
//	}
//
//	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
//	hasNext := page < totalPages
//	hasPrev := page > 1
//	offset := (page - 1) * limit
//
//	query := baseQuery.Session(&gorm.Session{})
//	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"message": "Error retrieving cars",
//			"error":   err.Error(),
//		})
//	}
//
//	ip := os.Getenv("HOST")
//	port := os.Getenv("PORT")
//	if ip == "" || port == "" {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"message": "Server configuration error: HOST or PORT not set",
//		})
//	}
//	for i := range cars {
//		cars[i].ImageUrl = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, cars[i].ImageUrl)
//	}
//
//	return c.Status(fiber.StatusOK).JSON(GetCarsResponse{
//		Results:    cars,
//		Page:       page,
//		Limit:      limit,
//		TotalPages: totalPages,
//		HasNext:    hasNext,
//		HasPrev:    hasPrev,
//	})
//}
//
//type ErrorResponse struct {
//	Message string `json:"message"`
//	Error   string `json:"error,omitempty"`
//}
