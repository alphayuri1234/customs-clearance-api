package middleware

import (
	"net/http"
	"strings"

	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

const ClaimsContextKey = "jwt_claims"

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("authorization header wajib diisi", nil))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("format authorization harus Bearer token", nil))
			return
		}

		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("token tidak valid atau kedaluwarsa", err.Error()))
			return
		}

		ctx.Set(ClaimsContextKey, claims)
		ctx.Next()
	}
}

func OfficerOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawClaims, exists := ctx.Get(ClaimsContextKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("user belum terautentikasi", nil))
			return
		}

		claims, ok := rawClaims.(*services.JWTClaims)
		if !ok || claims.Role != models.RoleOfficer {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse("akses hanya untuk officer", nil))
			return
		}

		ctx.Next()
	}
}

func CurrentUserClaims(ctx *gin.Context) (*services.JWTClaims, bool) {
	rawClaims, exists := ctx.Get(ClaimsContextKey)
	if !exists {
		return nil, false
	}

	claims, ok := rawClaims.(*services.JWTClaims)
	return claims, ok
}
