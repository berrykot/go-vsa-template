package server

import (
	"go-vsa-template/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

// Router — Engine и группы для публичных и защищённых роутов.
// Регистрация в Public — без JWT, в Protected — с обязательным Bearer JWT (Supabase).
type Router struct {
	Engine    *gin.Engine
	Public    *gin.RouterGroup
	Protected *gin.RouterGroup
}

func New() *Router {
	r := gin.Default()

	public := r.Group("")
	protected := r.Group("")
	protected.Use(auth.Middleware())

	return &Router{
		Engine:    r,
		Public:    public,
		Protected: protected,
	}
}
