package middleware

import (
	"net/http"

	"github.com/baltop/commet/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	CookieName      = "auth_token"
	UserContextKey  = "user"
	ClaimsContextKey = "claims"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// HTTP-Only Cookie에서 토큰 가져오기
		tokenString, err := c.Cookie(CookieName)
		if err != nil {
			// 인증이 필요한 페이지인 경우 로그인 페이지로 리다이렉트
			if isHTMXRequest(c) {
				c.Header("HX-Redirect", "/auth/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 토큰 검증
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			// 토큰이 유효하지 않으면 쿠키 삭제 후 로그인 페이지로 리다이렉트
			c.SetCookie(CookieName, "", -1, "/", "", false, true)
			if isHTMXRequest(c) {
				c.Header("HX-Redirect", "/auth/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 사용자 정보를 컨텍스트에 저장
		c.Set(ClaimsContextKey, claims)
		c.Next()
	}
}

func GuestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 이미 로그인된 사용자는 대시보드로 리다이렉트
		_, err := c.Cookie(CookieName)
		if err == nil {
			c.Redirect(http.StatusFound, "/dashboard")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isHTMXRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

// GetCurrentUser 컨텍스트에서 현재 사용자 정보 가져오기
func GetCurrentUser(c *gin.Context) *services.Claims {
	claims, exists := c.Get(ClaimsContextKey)
	if !exists {
		return nil
	}
	return claims.(*services.Claims)
}
