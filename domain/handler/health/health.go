package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HealthHandler struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func NewHealthHandler(logger *zap.SugaredLogger, db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		db:     db,
	}
}

func (h *HealthHandler) Health(c echo.Context) error {
	response := HealthResponse{
		Status:   "up",
		Database: "up",
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		response.Database = "down"
		response.Status = "degraded"
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	if err := sqlDB.Ping(); err != nil {
		response.Database = "down"
		response.Status = "degraded"
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) Ready(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}

func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "alive"})
}

func RegisterRoutes(echoInstance *echo.Echo, handler *HealthHandler) {
	echoInstance.GET("/health", handler.Health)
	echoInstance.GET("/ready", handler.Ready)
	echoInstance.GET("/live", handler.Live)
}

var Module = fx.Options(
	fx.Provide(NewHealthHandler),
	fx.Invoke(RegisterRoutes),
)
