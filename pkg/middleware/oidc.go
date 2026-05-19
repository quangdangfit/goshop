package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goshop/pkg/oidc"
)

// OIDCAuth validates the Bearer ID token issued by Authentik and injects user
// info into the Gin context. The keys it sets (`userId`, `role`) match those
// produced by the legacy JWT middleware so downstream handlers stay unchanged.
func OIDCAuth(validator *oidc.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		ctx := c.Request.Context()
		idToken, err := validator.Verify(ctx, authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var claims struct {
			Sub    string   `json:"sub"`
			Email  string   `json:"email"`
			Name   string   `json:"name"`
			Groups []string `json:"groups,omitempty"`
		}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		role := mapAuthentikGroupsToRole(claims.Groups)

		c.Set("userId", claims.Sub)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("role", role)
		c.Next()
	}
}

// mapAuthentikGroupsToRole maps Authentik group names to GoShop roles.
// Defaults to "customer"; only the groups listed below promote to "admin".
func mapAuthentikGroupsToRole(groups []string) string {
	for _, g := range groups {
		if g == "goshop-admins" || g == "admin" {
			return "admin"
		}
	}
	return "customer"
}
