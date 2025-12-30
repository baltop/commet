package handlers

import (
	"net/http"

	"github.com/baltop/commet/internal/middleware"
	"github.com/baltop/commet/internal/models"
	"github.com/baltop/commet/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// GET /auth/login - 로그인 페이지
func (h *AuthHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"title": "로그인",
	})
}

// POST /auth/login - 로그인 처리
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Form 데이터 바인딩
	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")

	if req.Email == "" || req.Password == "" {
		renderAuthError(c, "auth/login.html", "이메일과 비밀번호를 입력해주세요.", req.Email)
		return
	}

	user, token, err := h.authService.Login(&req)
	if err != nil {
		renderAuthError(c, "auth/login.html", "이메일 또는 비밀번호가 올바르지 않습니다.", req.Email)
		return
	}

	// HTTP-Only Cookie 설정
	c.SetCookie(
		middleware.CookieName,
		token,
		60*60*24, // 24시간
		"/",
		"",
		false, // Secure (프로덕션에서는 true)
		true,  // HttpOnly
	)

	// HTMX 요청인 경우 리다이렉트 헤더 설정
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/dashboard")
		c.Status(http.StatusOK)
		return
	}

	_ = user // user 변수 사용
	c.Redirect(http.StatusFound, "/dashboard")
}

// GET /auth/register - 회원가입 페이지
func (h *AuthHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"title": "회원가입",
	})
}

// POST /auth/register - 회원가입 처리
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Form 데이터 바인딩
	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")
	req.Name = c.PostForm("name")
	confirmPassword := c.PostForm("confirm_password")

	// 유효성 검사
	if req.Email == "" || req.Password == "" || req.Name == "" {
		renderRegisterError(c, "모든 필드를 입력해주세요.", req.Email, req.Name)
		return
	}

	if len(req.Password) < 6 {
		renderRegisterError(c, "비밀번호는 6자 이상이어야 합니다.", req.Email, req.Name)
		return
	}

	if req.Password != confirmPassword {
		renderRegisterError(c, "비밀번호가 일치하지 않습니다.", req.Email, req.Name)
		return
	}

	_, err := h.authService.Register(&req)
	if err != nil {
		if err == services.ErrUserExists {
			renderRegisterError(c, "이미 사용 중인 이메일입니다.", req.Email, req.Name)
			return
		}
		renderRegisterError(c, "회원가입 중 오류가 발생했습니다.", req.Email, req.Name)
		return
	}

	// HTMX 요청인 경우
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/auth/login?registered=true")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, "/auth/login?registered=true")
}

// POST /auth/logout - 로그아웃 처리
func (h *AuthHandler) Logout(c *gin.Context) {
	// 쿠키 삭제
	c.SetCookie(
		middleware.CookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	// HTMX 요청인 경우
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/auth/login")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, "/auth/login")
}

func renderAuthError(c *gin.Context, template, errMsg, email string) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "components/alert.html", gin.H{
			"type":    "error",
			"message": errMsg,
		})
		return
	}
	c.HTML(http.StatusOK, template, gin.H{
		"title": "로그인",
		"error": errMsg,
		"email": email,
	})
}

func renderRegisterError(c *gin.Context, errMsg, email, name string) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "components/alert.html", gin.H{
			"type":    "error",
			"message": errMsg,
		})
		return
	}
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"title": "회원가입",
		"error": errMsg,
		"email": email,
		"name":  name,
	})
}
