package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"
	"path/filepath"

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

	router.POST("/api/v1/feedback/", Feedback)
}

func SignUp(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}

	var newUser structs.User
	newUser.Name = c.PostForm("Name")
	newUser.Username = c.PostForm("Username")
	newUser.Email = c.PostForm("Email")
	newUser.Password = c.PostForm("Password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)

	file, err := c.FormFile("Picture")
	if err == nil {
		fileExt := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
		path := filepath.Join("../../frontend/public/images/pfp/", newFileName)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}

		webPath := "/static/images/pfp/" + newFileName
		webPath = strings.Replace(webPath, "\\", "/", -1)
		newUser.Picture = webPath
	} else {
		newUser.Picture = "/static/images/icons/camera_off_icon.svg"
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.Redirect(http.StatusFound, "/")
}

func RemoveProduct(c *gin.Context) {
	var Values struct {
		UserID string `form:"UserID"`
		FoodID string `form:"FoodID"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	redirect := fmt.Sprintf("/profile/%d", userID)
	c.Redirect(http.StatusFound, redirect)
}

func BuyProdct(c *gin.Context) {}

func Feedback(c *gin.Context) {
	var Values struct {
		UserID  string `form:"UserID"`
		FoodID  string `form:"FoodID"`
		Rating  string `form:"Rating"`
		Comment string `form:"Comment"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	rating, err := strconv.Atoi(Values.Rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in rating data": rating})
		return
	}

	feedback := structs.Feedback{
		UserID:  uint(userID),
		FoodID:  uint(foodID),
		Rating:  uint(rating),
		Comment: Values.Comment,
	}

	if err := server.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record feedback to product"})
		return
	}

	redirect := fmt.Sprintf("/product/%s", Values.FoodID)
	c.Redirect(http.StatusFound, redirect)
}
