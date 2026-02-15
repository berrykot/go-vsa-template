package health

import (
	"go-vsa-template/internal/infrastructure/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger zerolog.Logger
	db     *database.Client
}

func NewHandler(l zerolog.Logger, db *database.Client) *Handler {
	return &Handler{
		logger: l,
		db:     db,
	}
}

// Register регистрирует публичные роуты этого слайса (без JWT)
func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	g.GET("/api/health", h.Check)
}

func (h *Handler) Check(c *gin.Context) {
	err := h.db.OpenConnection()

	if err != nil {
		h.logger.Error().Err(err).Msg("database health check failed")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"db":     "unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"db":     "connected",
	})
}
