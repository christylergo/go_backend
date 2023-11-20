package authen

import (
	"github.com/gin-gonic/gin"
)

func AddAuthenticationRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLFiles("./assets/html_templates/index.tmpl")
	// g.StaticFile("/style.css", "./assets/html_templates/style.css")
	g.Static("/style.css", "./assets/html_templates")
	// g.GET("/new.txt", serveFile) // quivalent to codes above
	g.GET("/", FirstPage)
	g.POST("/register", registerFunc)
	// g.Post("/login", loginFunc)
	// g.Get("/refresh_token", refreshFunc)
	// g.Get("/logout")
	// g.Use(parseToken)
	return g
}
