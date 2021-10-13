package middlewares

import (
	"net/http"

	"github.com/geforce6t/go-server/utils"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := utils.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error User seems to be unauthenticated!",
			})
			c.Abort()
			return
		}

		if _, fetchErr := utils.FetchAuth(ad); fetchErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error User seems to be unauthenticated!",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
