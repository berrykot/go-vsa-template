package health

import (
	"go-vsa-template/internal/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type HandlerProtected struct {
	logger zerolog.Logger
}

func NewHandlerProtected(l zerolog.Logger) *HandlerProtected {
	return &HandlerProtected{logger: l}
}

// Register регистрирует защищённые роуты (группа уже с userInfo)
func (h *HandlerProtected) Register(g *gin.RouterGroup) {
	g.POST("/api/my-protected", h.Generate)
}

func (h *HandlerProtected) Generate(c *gin.Context) {
	var req MyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	userInfo, ok := auth.GetAuthUser(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, MyResponse{
		Data: req.Data + userInfo.Name,
	})
}
