package validator

import (
	"backend/contrib/decimal"
	"backend/contrib/models"
	"backend/core"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validate Global validator instance
var Validate *validator.Validate

func EnumValidator(allowedValues []string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()

		// Handle pointers
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return true // allow nil for "omitempty"
			}
			field = field.Elem()
		}

		// Convert any string-like type to string
		value := fmt.Sprint(field.Interface())

		for _, v := range allowedValues {
			if value == v {
				return true
			}
		}
		return false
	}
}

func ValidateCarParkType(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(models.CarParkType)
	if !ok {
		return false
	}

	switch value {
	case models.Park3, models.Park4:
		return true
	default:
		return false
	}
}

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	Validate.RegisterValidation(
		"cameratype",
		EnumValidator([]string{string(models.InsideCamera), string(models.OutsideCamera)}),
	)
	Validate.RegisterValidation(
		"carparktype",
		EnumValidator([]string{string(models.Park3), string(models.Park4)}),
	)
	// RoleType
	Validate.RegisterValidation(
		"roletype",
		EnumValidator([]string{string(models.OperatorRole), string(models.AdminRole), string(models.AccountantRole)}),
	)

	Validate.RegisterValidation("decimal_gt0", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(decimal.Decimal)
		if !ok {
			return false
		}
		return val.GreaterThan(decimal.Zero)
	})
}

var validationMessages = map[string]string{
	"required":    "This field is required",
	"carparktype": "Invalid car park type, must be 'park-3' or 'park-4'",
	"roletype":    "Invalid role type, must be one of 'operator', 'admin', 'accountant'",
	"cameratype":  "Invalid camera type",
}

// ParseBody parses the JSON request body into the provided struct pointer
// and validates it using the global validator. Returns structured error if invalid.
func ParseBody[T any](c *fiber.Ctx) (*T, error) {
	var input T
	raw := c.Body()
	contentType := c.Get("Content-Type")

	fmt.Printf("Content-Type: %s\n", contentType)
	fmt.Printf("Request body: %s\n", raw)

	if err := c.BodyParser(&input); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := Validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMap := make(map[string]string, len(ve))
			for _, e := range ve {
				// Map struct field name to JSON name
				field := getJSONFieldName(input, e.Field())

				// Use custom message if available, otherwise fallback
				msg, ok := validationMessages[e.Tag()]
				if !ok {
					msg = "Failed on " + e.Tag()
				}
				errMap[field] = msg
			}

			// Return ValidationErrorResponse as error
			return nil, &core.ValidationErrorResponse{
				Detail: "Validation error",
				Errors: errMap,
			}
		}
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return &input, nil
}

// ParseFormBody parses multipart/form-data into the provided struct pointer
// and validates it. Returns structured ValidationErrorResponse if invalid.
func ParseFormBody[T any](c *fiber.Ctx, fileFields ...string) (*T, map[string]*multipart.FileHeader, error) {
	var input T
	files := make(map[string]*multipart.FileHeader)

	form, err := c.MultipartForm()
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid form-data")
	}
	mapFormToStruct(&input, form.Value)
	// Map files
	for _, field := range fileFields {
		if uploadedFiles, ok := form.File[field]; ok && len(uploadedFiles) > 0 {
			files[field] = uploadedFiles[0]
		}
	}

	// TODO: manually map form values to struct fields here
	if err := Validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMap := make(map[string]string, len(ve))
			for _, e := range ve {

				// Map struct field name to JSON name
				field := getJSONFieldName(input, e.Field())

				// Use custom message if available, otherwise fallback
				msg, ok := validationMessages[e.Tag()]
				if !ok {
					msg = "Failed on " + e.Tag()
				}
				errMap[field] = msg
			}

			return nil, files, &core.ValidationErrorResponse{
				Detail: "Validation error",
				Errors: errMap,
			}
		}
		return nil, files, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return &input, files, nil
}

func mapFormToStruct[T any](input *T, form map[string][]string) {
	v := reflect.ValueOf(input).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		if values, ok := form[tag]; ok && len(values) > 0 {
			f := v.Field(i)
			if f.CanSet() {
				switch f.Kind() {
				case reflect.String:
					f.SetString(values[0])
				case reflect.Bool:
					f.SetBool(values[0] == "true")
					// Add other kinds as needed: int, float64, etc.
				}
			}
		}
	}
}

func getJSONFieldName(structType interface{}, fieldName string) string {
	t := reflect.TypeOf(structType)
	// Dereference pointer if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if f, ok := t.FieldByName(fieldName); ok {
		tag := f.Tag.Get("json")
		// tag can be "carPark,omitempty", split by comma
		return strings.Split(tag, ",")[0]
	}
	return fieldName
}
