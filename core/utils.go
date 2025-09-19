package core

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetCommonsFromContext(c *fiber.Ctx) CommonsModel {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return CommonsModel{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func ParseUUID(value string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, fmt.Errorf("empty string cannot be parsed as UUID")
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return id, nil
}
