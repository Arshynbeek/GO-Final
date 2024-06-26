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
	router.LoadHTMLGlob("/root/templates/*.html")

	router.GET("/", OptionalAuthMiddleware(), HomePage)

	router.GET("/signin/", SignInPage)
	router.GET("/signup/", SignUpPage)

	router.GET("/profile/:id", AuthMiddleware(), ProfilePage)
	router.GET("/edit-profile/:id", AuthMiddleware(), EditProfilePage)

	router.GET("/product/:id", OptionalAuthMiddleware(), ProductPage)
	router.GET("/edit-product/:id", OptionalAuthMiddleware(), EditProductPage)
	router.GET("/new-product/", AuthMiddleware(), NewProductPage)

	router.GET("/report/", AuthMiddleware(), AdminMiddleware(), ReportPage)

	router.GET("/about/", AboutPage)

	router.Static("/static", "/root/public")

	api.SetupAPIRoutes(router)

	router.Run(":2024")
	return router
}

func HomePage(c *gin.Context) {
	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var categories []structs.Category
	if result := server.DB.Find(&categories); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var orders []structs.Order
	if result := server.DB.Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
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

func ProfilePage(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists || fmt.Sprintf("%d", userID) != id {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"status": http.StatusUnauthorized,
			"error":  "Unauthorized access to this profile",
		})
		return
	}

	var user structs.User
	if result := server.DB.First(&user, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var orders []structs.Order
	if result := server.DB.Where("user_id = ?", id).Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	if user.Admin {
		var users []structs.User
		if result := server.DB.Find(&users); result.Error != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"status": http.StatusInternalServerError,
				"error":  result.Error.Error(),
			})
			return
		}

		var orders []structs.Order
		if result := server.DB.Find(&orders); result.Error != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"status": http.StatusInternalServerError,
				"error":  result.Error.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "admin.html", gin.H{
			"admin":  user,
			"users":  users,
			"orders": orders,
			"foods":  foods,
		})
	} else {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"user":   user,
			"orders": orders,
			"foods":  foods,
		})
	}
}

func EditProfilePage(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists || fmt.Sprintf("%d", userID) != id {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"status": http.StatusUnauthorized,
			"error":  "Unauthorized access to this profile",
		})
		return
	}

	var user structs.User
	if result := server.DB.First(&user, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "edit-profile.html", gin.H{
		"user": user,
	})
}

func NewProductPage(c *gin.Context) {
	var categories []structs.Category
	if result := server.DB.Find(&categories); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "new-product.html", gin.H{
		"title":      "New Product",
		"categories": categories,
	})
}

func ProductPage(c *gin.Context) {
	id := c.Param("id")
	var food structs.Food
	if result := server.DB.First(&food, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"status": http.StatusNotFound,
			"error":  result.Error.Error(),
		})
		return
	}

	var foods []structs.Food
	if result := server.DB.Find(&foods); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var orders []structs.Order
	if result := server.DB.Find(&orders); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var feedbacks []structs.Feedback
	if result := server.DB.Find(&feedbacks, "food_id = ?", id); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	var users []structs.User
	if result := server.DB.Find(&users, "admin = ?", false); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	user, exists := c.Get("userID")
	var data structs.User
	if exists {
		server.DB.First(&data, user)
	}

	c.HTML(http.StatusOK, "product.html", gin.H{
		"food":      food,
		"foods":     foods,
		"orders":    orders,
		"feedbacks": feedbacks,
		"users":     users,
		"user":      data,
	})
}

func EditProductPage(c *gin.Context) {
	id := c.Param("id")
	var food structs.Food
	if result := server.DB.First(&food, "id = ?", id); result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"status": http.StatusNotFound,
			"error":  result.Error.Error(),
		})
		return
	}

	var categories []structs.Category
	if result := server.DB.Find(&categories); result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  result.Error.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "edit-product.html", gin.H{
		"food":       food,
		"categories": categories,
	})
}

func ReportPage(c *gin.Context) {
	sales, err1 := api.SaleReports()
	inventory, err2 := api.InventoryReports()
	preferences, err3 := api.PreferenceReports()
	revenue, err4 := api.RevenueReports()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Failed to fetch reports: \n" + err1.Error() + "\n" + err2.Error() + "\n" + err3.Error() + "\n" + err4.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "report.html", gin.H{
		"title":       "Reports",
		"sales":       sales,
		"inventory":   inventory,
		"preferences": preferences,
		"revenue":     revenue,
	})
}

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "About US",
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
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"status": http.StatusUnauthorized,
				"error":  "You need to be signed in to access this page",
			})
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
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var user structs.User
		if err := server.DB.First(&user, userID).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"status": http.StatusUnauthorized,
				"error":  "User not found",
			})
			c.Abort()
			return
		}

		if !user.Admin {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"status": http.StatusForbidden,
				"error":  "Access denied - You do not have permission to access this page",
			})
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
