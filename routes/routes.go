package routes

import (
    "github.com/chandan/trading-dashboard/auth"
    "github.com/chandan/trading-dashboard/controllers"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173",
            "http://trading-dashboard-chandan.s3-website.eu-north-1.amazonaws.com",
            "https://doui0bo89ibz6.cloudfront.net",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: false,
    }))

    // Public endpoints
    r.GET("/prices", controllers.GetPrices)
    r.GET("/ws", controllers.WSHandler)

    // Auth
    authGroup := r.Group("/auth")
    authGroup.POST("/login", auth.LoginHandler)

    // Orders (protected)
    orders := r.Group("/orders")
    orders.Use(auth.AuthMiddleware())
    orders.GET("", controllers.GetOrders)
    orders.POST("", controllers.CreateOrder)

    return r
}
