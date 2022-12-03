package middleware

import (
	"regexp"
	"log"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	
}

func NewMiddleware() *Middleware {
	return &Middleware{

	}
}
const miniappRegex = `^https:\/\/prod-app[0-9]*-[a-z 0-9]*\.pages-ac\.vk-apps\.com$`

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		isAllowed := false
		
		isAllowed, _ = regexp.MatchString(miniappRegex,origin)

		if !isAllowed {
			log.Print("CORS not allowed origin = ", origin)
			return
		}

        c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}