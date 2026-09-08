package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/septalfauzan/saku-api/app/server/domain"
	"github.com/septalfauzan/saku-api/app/server/datasources/remote"
	"github.com/septalfauzan/saku-api/app/server/services"
)

func OCRReceiptImage(service services.ReceiptService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input remote.OCRRequest

		if err := c.BodyParser(&input); err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"invalid request body: "+err.Error(),
			)
		}

		if input.ImageBase64 == "" {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"image is required",
			)
		}

		receipt, err := service.ExtractReceipt(c.Context(), input)
		if err != nil {
			log.Printf("ERROR extract receipt: %v", err)
			return sendError(
				c,
				fiber.StatusInternalServerError,
				"failed to extract receipt",
			)
		}

		if len(receipt.Items) == 0 {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"image does not contain receipt items",
			)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"data": receipt,
		})
	}
}

func sendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(domain.Error{
		Text: message,
	})
}
