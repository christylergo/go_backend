package authen

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserInfo struct {
	Id   int    `form:"id"`
	Name string `form:"name" binding:"required,name_validator"`
	Age  int    `form:"age" binding:"required,lt=200"`
}

func (u *UserInfo) TableName() string {
	return "userinfo"
}

func FirstPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tyler/index.tmpl", gin.H{
		"title": " Tyler!"})
}

func registerFunc(c *gin.Context) {
	user := UserInfo{}
	var customValidator validator.Func = func(f validator.FieldLevel) bool {
		val := f.Field().String()
		if strings.Contains(val, " ") {
			return false
		}
		return true
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("name_validator", customValidator)
	}
	// This c.ShouldBind consumes c.Request.Body and it cannot be reused.
	if errA := c.ShouldBind(&user); errA != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User validation failed!",
			"error":   errA.Error()})
	}
	fmt.Println(user)
	wirteUserInfoToPg(&user)
	//		curl -X POST localhost:8080/register
	//	   -H "Content-Type: application/x-www-form-urlencoded"
	//	   -d "param1=value1&param2=value2"
}

func loginFunc(c *gin.Context) {
	//
}

func serveFile(c *gin.Context) {
	fmt.Println(c.Request.URL.Path)
	reader, err := os.Open("./assets/html_templates" + c.Request.URL.Path)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	defer reader.Close()
	header := make([]byte, 512)
	reader.Read(header)
	contentType := http.DetectContentType(header)
	fileStatus, _ := reader.Stat()
	contentLength := int64(fileStatus.Size())
	extraHeaders := map[string]string{"Content-Disposition": `attachment; filename="style.css"`}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}
