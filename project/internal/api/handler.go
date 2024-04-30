package api

import (
	"net/http"
	"strconv"

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
	router.GET("/api/v1/signout/", SignOut)

	router.POST("/api/v1/buy/", BuyProdct)
	router.POST("/api/v1/add/", AddProduct)
	router.POST("/api/v1/remove/", RemoveProduct)
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
	if result := server.DB.Where("Username = ?", Credentials.Username).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "User Not Found"})
			return
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Credentials.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("jwt", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func SignOut(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func AddProduct(c *gin.Context) {
	var Values struct {
		UserID string `form:"UserID"`
		FoodID string `form:"FoodID"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, err := strconv.Atoi(Values.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in user id": userID})
		return
	}

	foodID, err := strconv.Atoi(Values.FoodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in food id": foodID})
		return
	}

	order := structs.Order{
		UserID: uint(userID),
		FoodID: uint(foodID),
	}

	if err := server.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
		return
	}
}

func RemoveProduct(c *gin.Context) {
	var Values struct {
		UserID string `form:"UserID"`
		FoodID string `form:"FoodID"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, err := strconv.Atoi(Values.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in user id": userID})
		return
	}

	foodID, err := strconv.Atoi(Values.FoodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in food id": foodID})
		return
	}

	order := structs.Order{
		UserID: uint(userID),
		FoodID: uint(foodID),
	}

	if err := server.DB.Where("user_id = ? AND food_id = ?", userID, foodID).Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from cart"})
		return
	}
}

func BuyProdct(c *gin.Context) {}
