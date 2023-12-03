package authentication

import (
	"fmt"
	"math"
	"net/http"
	"regexp"

	pb "example.com/go_backend/internal/grpc/authen"
	"example.com/go_backend/internal/grpc/authen/app"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Name     string `form:"name" binding:"required_without=Phone,max=50,name_validator" time_format:"2023-11-11"`
	Phone    uint   `form:"phone" binding:"required_without=Name,phone_validator"`
	Email    string `form:"email" binding:"email"`
	PassWord string `form:"pass_word" binding:"required,min=3"`
}

var nameValidator validator.Func = func(f validator.FieldLevel) bool {
	val := f.Field().String()
	re := regexp.MustCompile(`(?i)^[a-z]+[a-z0-9]*$`)
	return re.Match([]byte(val))
}

var phoneValidator validator.Func = func(f validator.FieldLevel) bool {
	val := float64(f.Field().Uint())
	math.Pow10(10)
	return val > math.Pow10(10) && val < math.Pow10(11)
}

var translateUser = func(user *User) *pb.User {
	rpcUser := pb.User{}
	rpcUser.Name = user.Name
	rpcUser.Phone = uint64(user.Phone)
	rpcUser.Email = user.Email
	rpcUser.PassWord = user.PassWord
	return &rpcUser
}

func firstPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tyler/index.tmpl", nil)
}

func registerFunc(c *gin.Context) {
	user := User{}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("name_validator", nameValidator)
		v.RegisterValidation("phone_validator", phoneValidator)
	}
	// This c.ShouldBind consumes c.Request.Body and it cannot be reused.
	if errA := c.ShouldBind(&user); errA != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User validation failed!",
			"error":   errA.Error()})
	}
	// gRPC cal

	id, token := app.AuthenRegisterRPC(translateUser(&user))
	fmt.Println(id, token)
}

func loginFunc(c *gin.Context) {
	user := User{}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("name_validator", nameValidator)
		v.RegisterValidation("phone_validator", phoneValidator)
	}
	fmt.Println(user)
	//gRPC call
}

// func serveFile(c *gin.Context) {
// 	fmt.Println(c.Request.URL.Path)
// 	reader, err := os.Open("./assets/html_templates" + c.Request.URL.Path)
// 	if err != nil {
// 		c.Status(http.StatusServiceUnavailable)
// 		return
// 	}
// 	defer reader.Close()
// 	header := make([]byte, 512)
// 	reader.Read(header)
// 	contentType := http.DetectContentType(header)
// 	fileStatus, _ := reader.Stat()
// 	contentLength := int64(fileStatus.Size())
// 	extraHeaders := map[string]string{"Content-Disposition": `attachment; filename="style.css"`}
// 	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

// }
