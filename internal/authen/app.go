package authen

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserName string `form:"user_name" binding:"required,name_validator"`
	Passwd   string `form:"pass_word" binding:"required"`
	Phone    string `form:"phone" binding:"required,phone_validator"`
	Email    string `form:"email"`
}

func FirstPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tyler/index.tmpl", gin.H{
		"title": " Tyler!"})
}

func registerFunc(c *gin.Context) {
	user := UserInfo{}
	// This c.ShouldBind consumes c.Request.Body and it cannot be reused.
	if errA := c.ShouldBind(&user); errA != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User validation failed!",
			"error":   errA.Error()})
	}
	wirteUserInfoToPg(&user)
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
