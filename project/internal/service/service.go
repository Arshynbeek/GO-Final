package service

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"project/mod/internal/api"
	"project/mod/internal/structs"
	"project/mod/pkg/server"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("../../frontend/templates/*.html")

	router.GET("/", OptionalAuthMiddleware(), HomePage)
	router.GET("/signin/", SignInPage)
	router.GET("/signup/", SignUpPage)
	router.GET("/product/:id", OptionalAuthMiddleware(), ProductPage)
	router.GET("/profile/:id", AuthMiddleware(), ProfilePage)

	router.Static("/static", "../../frontend/public")

	api.SetupAPIRoutes(router)

	router.Run(":2024")
	return router
}

func HomePage(c *gin.Context) {
	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var categories []structs.Category
	if result := server.DB.Find(&categories); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var orders []structs.Order
	if result := server.DB.Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	user, exists := c.Get("userID")
	var data structs.User
	if exists {
		server.DB.First(&data, user)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Home",
		"foods":      foods,
		"orders":     orders,
		"categories": categories,
		"user":       data,
	})
}

func ProductPage(c *gin.Context) {
	id := c.Param("id")
	var food structs.Food
	if result := server.DB.First(&food, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var categories []structs.Category
	if result := server.DB.Find(&categories); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var orders []structs.Order
	if result := server.DB.Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var feedbacks []structs.Feedback
	if result := server.DB.Find(&feedbacks); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	user, exists := c.Get("userID")
	var data structs.User
	if exists {
		server.DB.First(&data, user)
	}

	c.HTML(http.StatusOK, "product.html", gin.H{
		"food":       food,
		"foods":      foods,
		"categories": categories,
		"orders":     orders,
		"feedbacks":  feedbacks,
		"user":       data,
	})
}

func ProfilePage(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists || fmt.Sprintf("%d", userID) != id {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Unauthorized access to this profile"})
		return
	}

	var user structs.User
	if result := server.DB.First(&user, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var orders []structs.Order
	if result := server.DB.Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user":   user,
		"orders": orders,
		"foods":  foods,
	})
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to be logged in to access this page"})
			c.Abort()
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("arch"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))
			c.Set("userID", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err == nil {
			token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte("arch"), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := uint(claims["user_id"].(float64))
				c.Set("userID", userID)
			}
		}
		c.Next()
	}
}
