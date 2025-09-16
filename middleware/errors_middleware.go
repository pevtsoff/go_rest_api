package errors_middleware

import "github.com/gin-gonic/gin"


type ErrorResponse struct {
	Error string `json:"error"`
}


func JSONErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.AbortWithStatusJSON(c.Writer.Status(), ErrorResponse{
				Error: c.Errors[0].Error(),
			})
		}
	}
}