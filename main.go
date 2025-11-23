package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gin-contrib/cors"
	"gorm.io/gorm"
)

func main() {
	InitDB()

	r := gin.Default()

	// FIX: Add CORS support
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.POST("/users", createUser)
	r.GET("/users", listUsers)
	r.POST("/users/login", loginUser)

	r.POST("/items", createItem)
	r.GET("/items", listItems)

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/carts", addToCart)
		auth.GET("/carts", listCarts)

		auth.POST("/orders", createOrder)
		auth.GET("/orders", listOrders)
	}

	r.Run(":8080")
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := User{
		Username: req.Username,
		Password: req.Password,
	}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func listUsers(c *gin.Context) {
	var users []User
	DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func loginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/password"})
		return
	}

	token := uuid.NewString()
	user.Token = token
	DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type CreateItemRequest struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status"`
}

func createItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status == "" {
		req.Status = "active"
	}
	item := Item{Name: req.Name, Status: req.Status}
	if err := DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func listItems(c *gin.Context) {
	var items []Item
	DB.Find(&items)
	c.JSON(http.StatusOK, items)
}

type AddToCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

func getUserFromContext(c *gin.Context) User {
	u, _ := c.Get("user")
	return u.(User)
}

func addToCart(c *gin.Context) {
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := getUserFromContext(c)

	var cart Cart
	err := DB.Where("user_id = ? AND status = ?", user.ID, "open").First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = Cart{
			UserID: user.ID,
			Name:   "Cart for " + user.Username,
			Status: "open",
		}
		if err := DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cartItem := CartItem{
		CartID: cart.ID,
		ItemID: req.ItemID,
	}
	if err := DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.CartID = &cart.ID
	DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message":   "item added to cart",
		"cart_id":   cart.ID,
		"item_id":   req.ItemID,
		"cart_item": cartItem,
	})
}

func listCarts(c *gin.Context) {
	var carts []Cart
	DB.Preload("Items").Find(&carts)
	c.JSON(http.StatusOK, carts)
}

type CreateOrderRequest struct {
	CartID uint `json:"cart_id" binding:"required"`
}

func createOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := getUserFromContext(c)

	var cart Cart
	if err := DB.Where("id = ? AND user_id = ? AND status = ?", req.CartID, user.ID, "open").First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart not found or not open"})
		return
	}

	order := Order{
		CartID: cart.ID,
		UserID: user.ID,
	}
	if err := DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart.Status = "checked_out"
	DB.Save(&cart)

	c.JSON(http.StatusOK, gin.H{
		"message": "order created",
		"order":   order,
	})
}

func listOrders(c *gin.Context) {
	user := getUserFromContext(c)
	var orders []Order
	DB.Where("user_id = ?", user.ID).Find(&orders)
	c.JSON(http.StatusOK, orders)
}