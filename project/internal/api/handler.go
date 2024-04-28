package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"project/mod/internal/structs"
	"project/mod/pkg/server"
	"project/mod/pkg/utils"
)

func SetupAPIRoutes(router *gin.Engine) {
	router.POST("/api/v1/signin/", SignIn)
	router.POST("/api/v1/signup/", SignUp)
}

func SignUp(c *gin.Context) {
	var newUser structs.User
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)

	if result := server.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/signin/")
}

func SignIn(c *gin.Context) {
	var Credentials struct {
		Username string `form:"Username"`
		Password string `form:"Password"`
	}

	if err := c.ShouldBind(&Credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var user structs.User
	if result := server.DB.Where("username = ?", Credentials.Username).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
	c.Redirect(http.StatusFound, "/")
}
