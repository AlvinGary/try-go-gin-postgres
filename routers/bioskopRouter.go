package routers

import (
	"database/sql"
	"try-go-gin-postgres/controllers"

	"github.com/gin-gonic/gin"
)

func SetupBioskopRoutes(router *gin.Engine, db *sql.DB) {
	bioskopGroup := router.Group("/bioskop")
	{
		bioskopGroup.POST("", func(c *gin.Context) {
			controllers.CreateBioskop(c, db)
		})
		bioskopGroup.GET("", func(c *gin.Context) {
			controllers.GetBioskops(c, db)
		})
		bioskopGroup.GET("/:ID", func(c *gin.Context) {
			controllers.GetBioskopByID(c, db)
		})
		bioskopGroup.PUT("/:ID", func(c *gin.Context) {
			controllers.UpdateBioskop(c, db)
		})
		bioskopGroup.DELETE("/:ID", func(c *gin.Context) {
			controllers.DeleteBioskop(c, db)
		})
	}
}