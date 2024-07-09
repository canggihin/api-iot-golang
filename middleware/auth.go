package middleware

import (
	"encoding/json"
	"log"
	"os"

	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(roleParams ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("x-authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"message": "Token is required"})
			c.Abort()
			return
		}

		pubKeyHex := os.Getenv("PUBLIC_KEY_PASSETO")
		log.Println(pubKeyHex)
		if pubKeyHex == "" {
			c.JSON(500, gin.H{"message": "Public key not set in environment"})
			c.Abort()
			return
		}

		pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(pubKeyHex)
		if err != nil {
			c.JSON(500, gin.H{"message": "Invalid public key"})
			c.Abort()
			return
		}

		parser := paseto.NewParser()
		token, err := parser.ParseV4Public(pubKey, tokenString, nil)
		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		var claims map[string]interface{}
		if err := json.Unmarshal(token.ClaimsJSON(), &claims); err != nil {
			c.JSON(500, gin.H{"message": "Failed to parse token claims"})
			c.Abort()
			return
		}

		userData, ok := claims["data"].(map[string]interface{})
		if !ok {
			c.JSON(500, gin.H{"message": "Invalid token structure"})
			c.Abort()
			return
		}

		userSubs, ok := userData["sub"].(bool)
		if !ok {
			c.JSON(500, gin.H{"message": "Role not found in token"})
			c.Abort()
			return
		}

		userStatus, ok := userData["status"].(string)
		if !ok {
			c.JSON(500, gin.H{"message": "Role not found in token"})
			c.Abort()
			return
		}

		userRole, ok := userData["role"].(string)
		if !ok {
			c.JSON(500, gin.H{"message": "Role not found in token"})
			c.Abort()
			return
		}

		for _, role := range roleParams {
			if role == userRole {
				c.Next()
				return
			}
		}
		if userSubs && userStatus == "active" {
			c.Next()
			return
		}

		c.JSON(403, gin.H{"message": "Forbidden: Insufficient role"})
		c.Abort()
	}
}
