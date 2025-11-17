package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler accepts username/password and returns JWT
func LoginHandler(c *gin.Context) {
	var p loginPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}

	// Check against hardcoded user (demo)
	if p.Username != HardcodedUser.Username || p.Password != HardcodedUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := GenerateToken(HardcodedUser.ID, HardcodedUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": HardcodedUser.ID, "username": HardcodedUser.Username},
	})
}
