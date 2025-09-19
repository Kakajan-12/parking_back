package zreport

//import (
//	"fmt"
//
//	"github.com/gofiber/fiber/v2"
//)
//
//type ZReport struct {
//	TotalPayment int    `json:"totalPayment"`
//	Username     string `json:"username"`
//	ParkNumber   string `json:"parkNumber"`
//}
//
//// GetZData godoc
//// @Summary Create a New Tarif
//// @Description Creates a new Report and saves it to the database.
//// @Tags ZReport
//// @Accept json
//// @Produce json
//// @Param tariff body ZReport true "Tariff details to be created"
//// @Success 201 {object} ZReport "Successfully created"
//// @Failure 400 {object} resmodel.ErrorResponse "Invalid request data"
//// @Failure 500 {object} resmodel.ErrorResponse "Failed to save data to the database"
//// @Router /z-report [post]
//func GetZData(c *fiber.Ctx) error {
//	var data ZReport
//	if err := c.BodyParser(&data); err != nil {
//		return c.Status(400).JSON("Press err")
//	}
//	fmt.Println("payment:", data.TotalPayment)
//	fmt.Println("Username:", data.Username)
//	fmt.Println("ParkNumber:", data.ParkNumber)
//	return c.Status(200).JSON(data)
//}
