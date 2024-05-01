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
	router.GET("/api/v1/signout/", SignOut)
	router.POST("/api/v1/signin/", SignIn)
	router.POST("/api/v1/signup/", SignUp)
	router.POST("/api/v1/edit-profile", EditProfile)

	router.POST("/api/v1/new/", NewProduct)
	router.POST("/api/v1/buy/", BuyProduct)
	router.POST("/api/v1/add/", AddProduct)
	router.POST("/api/v1/remove/", RemoveProduct)

	router.POST("/api/v1/feedback/", Feedback)
}

func SignUp(c *gin.Context) {
	var Credentials struct {
		Name     string `form:"Name"`
		Username string `form:"Username"`
		Email    string `form:"Email"`
		Password string `form:"Password"`
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}

	var newUser structs.User
	newUser.Name = Credentials.Name
	newUser.Username = Credentials.Username
	newUser.Email = Credentials.Email
	newUser.Password = Credentials.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)

	file, err := c.FormFile("Picture")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		path := filepath.Join("../../frontend/public/images/pfp/", newFileName)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file" + err.Error()})
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

func EditProfile(c *gin.Context) {
	var Changes struct {
		UserID      uint   `form:"UserID"`
		Name        string `form:"Name"`
		Username    string `form:"Username"`
		Email       string `form:"Email"`
		OldPassword string `form:"OldPassword"`
		NewPassword string `form:"NewPassword"`
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form: " + err.Error()})
		return
	}

	if err := c.ShouldBind(&Changes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind data: " + err.Error()})
		return
	}

	var data structs.User
	if result := server.DB.Where("id = ?", Changes.UserID).First(&data); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "User Not Found"})
			return
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(Changes.OldPassword)); err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Invalid User Old Password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Changes.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var user structs.User
	if result := server.DB.Where("id = ?", Changes.UserID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Name = Changes.Name
	user.Username = Changes.Username
	user.Email = Changes.Email
	user.Password = string(hashedPassword)

	file, err := c.FormFile("Picture")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		path := filepath.Join("../../frontend/public/images/pfp/", newFileName)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file: " + err.Error()})
			return
		}

		webPath := "/static/images/pfp/" + newFileName
		webPath = strings.Replace(webPath, "\\", "/", -1)
		user.Picture = webPath
	}

	if result := server.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user: " + result.Error.Error()})
		return
	}

	redirect := fmt.Sprintf("/profile/%d", Changes.UserID)
	c.Redirect(http.StatusFound, redirect)
}

func NewProduct(c *gin.Context) {
	var Product struct {
		CategoryID  uint   `form:"Category"`
		NewCategory string `form:"NewCategory"`
		Name        string `form:"Name"`
		Description string `form:"Description"`
		Quantity    uint   `form:"Quantity"`
		Price       uint   `form:"Price"`
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form: " + err.Error()})
		return
	}

	if err := c.ShouldBind(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect form data: " + err.Error()})
		return
	}

	var category structs.Category
	if Product.CategoryID == 0 {
		if Product.NewCategory != "" {
			category.Name = Product.NewCategory
			if err := server.DB.Create(&category).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new Category: " + err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New category name required"})
			return
		}
	} else {
		if err := server.DB.First(&category, "id = ?", Product.CategoryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Category not found"})
			return
		}
	}

	newProduct := structs.Food{
		Name:        Product.Name,
		Description: Product.Description,
		Quantity:    Product.Quantity,
		Price:       Product.Price,
		CategoryID:  category.ID,
		Pictures:    []string{},
	}

	for i := 1; i <= 4; i++ {
		file, err := c.FormFile(fmt.Sprintf("Picture%d", i))
		if err == nil {
			extension := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
			path := filepath.Join("../../frontend/public/images/products/", newFileName)

			if err := c.SaveUploadedFile(file, path); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
				return
			}

			webPath := "/static/images/products/" + newFileName
			webPath = strings.Replace(webPath, "\\", "/", -1)
			newProduct.Pictures = append(newProduct.Pictures, webPath)
		} else {
			newProduct.Pictures = append(newProduct.Pictures, "/static/images/image.jpg")
		}
	}

	if err := server.DB.Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new product: " + err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func BuyProduct(c *gin.Context) {}

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
		UserID uint `form:"UserID"`
		FoodID uint `form:"FoodID"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := structs.Order{
		UserID: Values.UserID,
		FoodID: Values.FoodID,
	}

	if err := server.DB.Where("user_id = ? AND food_id = ?", Values.UserID, Values.FoodID).Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from cart"})
		return
	}

	redirect := fmt.Sprintf("/profile/%d", Values.UserID)
	c.Redirect(http.StatusFound, redirect)
}

func Feedback(c *gin.Context) {
	var Values struct {
		UserID  uint   `form:"UserID"`
		FoodID  uint   `form:"FoodID"`
		Rating  uint   `form:"Rating"`
		Comment string `form:"Comment"`
	}

	if err := c.ShouldBind(&Values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback := structs.Feedback{
		UserID:  Values.UserID,
		FoodID:  Values.FoodID,
		Rating:  Values.Rating,
		Comment: Values.Comment,
	}

	if err := server.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record feedback to product"})
		return
	}

	redirect := fmt.Sprintf("/product/%d", Values.FoodID)
	c.Redirect(http.StatusFound, redirect)
}
