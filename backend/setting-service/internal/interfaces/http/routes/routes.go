package routes

import (
	"nurseshift/setting-service/internal/interfaces/http/handlers"
	mw "nurseshift/setting-service/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handlers.SettingHandler) {
	api := app.Group("/api/v1")
	settings := api.Group("/settings")
	settings.Use(mw.AuthMiddleware())
	settings.Get("/", h.GetSettings)
	settings.Put("/", h.UpdateSettings)
	settings.Post("/shifts", h.CreateShift)
	settings.Put("/shifts/:id", h.UpdateShift)
	settings.Patch("/shifts/:id/toggle", h.ToggleShift)
	settings.Delete("/shifts/:id", h.DeleteShift)
	settings.Post("/holidays", h.CreateHoliday)
	settings.Delete("/holidays/:id", h.DeleteHoliday)
	settings.Put("/holidays/:id", h.UpdateHoliday)
}
