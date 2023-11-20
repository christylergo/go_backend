package authen

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func RunMain() {
	router := gin.Default()
	router.Delims("{[{", "}]}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLFiles("/home/tyler/go/src/go_backend/cmd/testdata/raw.tmpl")

	router.GET("/raw", func(c *gin.Context) {
		fmt.Println(c.Request.URL, "***url")
		c.HTML(http.StatusOK, "raw.tmpl", gin.H{
			"now": time.Date(2030, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		})
	})

	router.Run(":8080")
}
