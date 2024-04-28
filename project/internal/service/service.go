package service

import (
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"project/mod/internal/api"
	"project/mod/internal/structs"
	"project/mod/pkg/server"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("../../frontend/templates/*.html")

	router.GET("/", HomePage)
	router.GET("/signin/", SignInPage)
	router.GET("/signup/", SignUpPage)
	router.GET("/product/:id", ProductPage)
	router.GET("/profile/:id", ProfilePage)

	router.Static("/static", "./static")

	api.SetupAPIRoutes(router)

	router.Run(":2024")
	return router
}

func HomePage(c *gin.Context) {
	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error retrieving food items from database.",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "home",
		"foods": foods,
	})
}

func ProductPage(c *gin.Context) {
	id := c.Param("id")

	var food structs.Food
	result := server.DB.First(&food, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Product not found"})
		return
	} else if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "product.html", gin.H{"food": food})
}

func ProfilePage(c *gin.Context) {
	id := c.Param("id")

	var user structs.User
	result := server.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	} else if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}
	c.HTML(http.StatusOK, "profile.html", user)
}

func SignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{
		"title": "Sign Up",
	})
}

func SignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{
		"title": "Sign In",
	})
}
