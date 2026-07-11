package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iskelk/ecommerce-yt/apiresponse"
	"github.com/iskelk/ecommerce-yt/tokens"
)

func extractToken(c *gin.Context) string {
	if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		if token, ok := strings.CutPrefix(authHeader, "Bearer "); ok && token != "" {
			return token
		}
	}
	return c.Request.Header.Get("token")
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := extractToken(c)
		if clientToken == "" {
			apiresponse.Unauthorized(c, apiresponse.MsgNoToken)
			c.Abort()
			return
		}

		claims, err := tokens.ValidateToken(clientToken)
		if err != nil {
			switch {
			case errors.Is(err, tokens.ErrTokenExpired):
				apiresponse.Unauthorized(c, apiresponse.MsgTokenExpired)
			case errors.Is(err, tokens.ErrTokenInvalid):
				apiresponse.Unauthorized(c, apiresponse.MsgInvalidToken)
			default:
				apiresponse.Internal(c, err)
			}
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Set("is_admin", claims.Is_Admin)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, ok := c.Get("is_admin")
		if !ok || isAdmin != true {
			apiresponse.Forbidden(c, apiresponse.MsgForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
