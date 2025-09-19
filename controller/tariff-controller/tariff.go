package tarifcontrol

//import (
//	"encoding/json"
//	"math"
//	resmodel "park/controller/getdata/resModel"
//	"park/database"
//	"park/models/tariff"
//	"park/util"
//	"strconv"
//	"time"
//
//	"github.com/gofiber/fiber/v2"
//)
//
//var TimeFormat = "2006-01-02 15:04:05"
//
//type Tariff struct {
//	Id        int       `json:"id"`
//	Plate     string    `json:"plate"`
//	Name      string    `json:"name"`
//	StartTime time.Time `json:"startTime"`
//	EndTime   time.Time `json:"endTime"`
//	Price     int       `json:"price"`
//}
//
//func (t *Tariff) UnmarshalJSON(data []byte) error {
//	type Alias Tariff
//	aux := &struct {
//		StartTime string `json:"startTime"`
//		EndTime   string `json:"endTime"`
//		*Alias
//	}{
//		Alias: (*Alias)(t),
//	}
//
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//
//	startTime, err := time.Parse(TimeFormat, aux.StartTime)
//	if err != nil {
//		return err
//	}
//	t.StartTime = startTime
//
//	endTime, err := time.Parse(TimeFormat, aux.EndTime)
//	if err != nil {
//		return err
//	}
//	t.EndTime = endTime
//
//	return nil
//}
//
//// CreateTariff godoc
//// @Summary Create a New Tarif
//// @Description Creates a new tarif and saves it to the database.
//// @Tags Tarif
//// @Accept json
//// @Produce json
//// @Param tarif body tarif.Tarif true "Tarif details to be created"
//// @Success 201 {object} tarif.Tarif "Successfully created"
//// @Failure 400 {object} resmodel.ErrorResponse "Invalid request data"
//// @Failure 500 {object} resmodel.ErrorResponse "Failed to save data to the database"
//// @Router /api/v1/accountant/tariff [post]
//func CreateTariff(c *fiber.Ctx) error {
//	var tariffDbj Tariff
//
//	if err := c.BodyParser(&tariffDbj); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(resmodel.ErrorResponse{
//			Error:   "Failed to parse request body",
//			Details: err.Error(),
//		})
//	}
//	if err := database.DB.Create(&tariffDbj).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
//			Error:   "Failed to save data to the database",
//			Details: err.Error(),
//		})
//	}
//	util.LoadVIPPlates()
//	return c.Status(201).JSON(tariffDbj)
//}
//
//// DeleteTariff godoc
//// @Summary Delete Tarif
//// @Description Deletes a tarif by its ID.
//// @Tags Tarif
//// @Param id path int true "ID of the tarif to delete"
//// @Success 200 {string} string "Tarif successfully deleted"
//// @Failure 400 {object} resmodel.ErrorResponse "Invalid ID format"
//// @Failure 404 {object} resmodel.ErrorResponse "Tarif not found"
//// @Failure 500 {object} resmodel.ErrorResponse "Database error"
//// @Router /api/v1/accountant/tariff/{id} [delete]
//func DeleteTariff(c *fiber.Ctx) error {
//	id, err := strconv.Atoi(c.Params("id"))
//	if err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(resmodel.ErrorResponse{
//			Error:   "Invalid ID format",
//			Details: "ID must be a number",
//		})
//	}
//
//	var tariffDbj tariff.Tariff
//	if err := database.DB.First(&tariffDbj, id).Error; err != nil {
//		return c.Status(fiber.StatusNotFound).JSON(resmodel.ErrorResponse{
//			Error:   "Tariff not found",
//			Details: err.Error(),
//		})
//	}
//
//	if err := database.DB.Delete(&tariffDbj).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
//			Error:   "Failed to delete tariff",
//			Details: err.Error(),
//		})
//	}
//	util.LoadVIPPlates()
//
//	return c.Status(200).JSON(fiber.Map{"message": "Tariff successfully deleted"})
//}
//
//type PaginatedResponse struct {
//	Data       interface{} `json:"data"`
//	Page       int         `json:"page"`
//	Limit      int         `json:"limit"`
//	TotalPages int         `json:"totalPages"`
//	HasNext    bool        `json:"hasNext"`
//	HasPrev    bool        `json:"hasPrev"`
//	TotalPrice int         `json:"totalPrice"`
//}
//
//// GetAllTariff godoc
//// @Summary Get all Tariffs with pagination
//// @Description Retrieves all tariffs from the database with pagination support.
//// @Tags Tarif
//// @Accept json
//// @Produce json
//// @Param page query int false "Page number" default(1)
//// @Param limit query int false "Number of items per page" default(10)
//// @Success 200 {object} PaginatedResponse "List of tariffs"
//// @Failure 500 {object} resmodel.ErrorResponse "Database error"
//// @Router /api/v1/accountant/tariff [get]
//func GetAllTariff(c *fiber.Ctx) error {
//	page, err := strconv.Atoi(c.Query("page", "1"))
//	if err != nil || page <= 0 {
//		page = 1
//	}
//
//	limit, err := strconv.Atoi(c.Query("limit", "10"))
//	if err != nil || limit <= 0 {
//		limit = 10
//	}
//
//	var tariffs []tariff.Tariff
//	var totalCount int64
//	var totalPrice int
//
//	if err := database.DB.Model(&tariff.Tariff{}).Count(&totalCount).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
//			Error:   "Failed to get total count of tariffs",
//			Details: err.Error(),
//		})
//	}
//
//	if err := database.DB.Offset((page - 1) * limit).Limit(limit).Find(&tariffs).Error; err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
//			Error:   "Failed to retrieve tariffs",
//			Details: err.Error(),
//		})
//	}
//	for _, tariffDbj := range tariffs {
//		totalPrice += tariffDbj.Price
//	}
//	totalPages := int(totalCount / int64(limit))
//	if totalCount%int64(limit) != 0 {
//		totalPages++
//	}
//	hasNext := page < totalPages
//	hasPrev := page > 1
//
//	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
//		Data:       tariffs,
//		Page:       page,
//		Limit:      limit,
//		TotalPages: totalPages,
//		HasNext:    hasNext,
//		HasPrev:    hasPrev,
//		TotalPrice: totalPrice,
//	})
//}
//
//// SearchCar godoc
//// @Summary Search for cars by plate number
//// @Description Retrieves cars from the database that match the given plate number with pagination.
//// @Tags Tarif
//// @Accept json
//// @Produce json
//// @Param car_number query string false "Car plate number to search for"
//// @Param page query int false "Page number" default(1)
//// @Param limit query int false "Number of items per page" default(5)
//// @Success 200 {object} PaginatedResponse "List of matching cars with pagination"
//// @Router /api/v1/accountant/search-car [get]
//func SearchCar(c *fiber.Ctx) error {
//	var plate []Tariff
//	var totalCount int64
//
//	carNumber := c.Query("car_number")
//	pageStr := c.Query("page", "1")
//	limitStr := c.Query("limit", "5")
//
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page < 1 {
//		return c.Status(400).JSON(fiber.Map{"message": "Invalid page number"})
//	}
//
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil || limit < 1 {
//		return c.Status(400).JSON(fiber.Map{"message": "Invalid limit number"})
//	}
//
//	query := database.DB.Model(&Tariff{})
//
//	if carNumber != "" {
//		query = query.Where("plate LIKE ?", "%"+carNumber+"%")
//	}
//
//	if err := query.Count(&totalCount).Error; err != nil {
//		return c.Status(400).JSON(fiber.Map{"message": "Error counting cars", "error": err.Error()})
//	}
//
//	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
//	hasNext := page < totalPages
//	hasPrev := page > 1
//	offset := (page - 1) * limit
//
//	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&plate).Error; err != nil {
//		return c.Status(400).JSON(fiber.Map{"message": "Error retrieving cars", "error": err.Error()})
//	}
//
//	return c.Status(200).JSON(PaginatedResponse{
//		Data:       plate,
//		Page:       page,
//		Limit:      limit,
//		TotalPages: totalPages,
//		HasNext:    hasNext,
//		HasPrev:    hasPrev,
//		TotalPrice: 0, // Optional: Modify to sum up total prices if needed
//	})
//}
