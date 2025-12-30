package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/baltop/commet/internal/config"
	"github.com/baltop/commet/internal/database"
	"github.com/baltop/commet/internal/handlers"
	"github.com/baltop/commet/internal/middleware"
	"github.com/baltop/commet/internal/repository"
	"github.com/baltop/commet/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 설정 로드
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Gin 모드 설정
	gin.SetMode(cfg.Server.Mode)

	// 데이터베이스 연결
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 마이그레이션 실행
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 샘플 데이터 시드
	if err := database.SeedSampleData(); err != nil {
		log.Printf("Warning: Failed to seed sample data: %v", err)
	}

	// Repository 초기화
	userRepo := repository.NewUserRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Service 초기화
	authService := services.NewAuthService(userRepo, cfg.JWT)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// Handler 초기화
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	healthHandler := handlers.NewHealthHandler()

	// Gin 라우터 생성
	r := gin.Default()

	// 템플릿 로드
	r.SetFuncMap(template.FuncMap{
		"dict": func(values ...interface{}) map[string]interface{} {
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key, _ := values[i].(string)
				dict[key] = values[i+1]
			}
			return dict
		},
		"slice": func(s string, start, end int) string {
			runes := []rune(s)
			if start < 0 {
				start = 0
			}
			if end > len(runes) {
				end = len(runes)
			}
			return string(runes[start:end])
		},
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
	})

	// 템플릿 파일 로드
	loadTemplates(r)

	// 정적 파일 제공
	r.Static("/static", "./web/static")

	// Health check
	r.GET("/api/health", healthHandler.Health)

	// 홈페이지 - 로그인 페이지로 리다이렉트
	r.GET("/", func(c *gin.Context) {
		// 이미 로그인된 경우 대시보드로
		_, err := c.Cookie(middleware.CookieName)
		if err == nil {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
		c.Redirect(http.StatusFound, "/auth/login")
	})

	// 인증 라우트 (Guest only)
	auth := r.Group("/auth")
	auth.Use(middleware.GuestMiddleware())
	{
		auth.GET("/login", authHandler.LoginPage)
		auth.POST("/login", authHandler.Login)
		auth.GET("/register", authHandler.RegisterPage)
		auth.POST("/register", authHandler.Register)
	}

	// 로그아웃은 인증된 사용자만
	r.POST("/auth/logout", middleware.AuthMiddleware(authService), authHandler.Logout)

	// 대시보드 라우트 (Auth required)
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(authService))
	{
		dashboard.GET("", dashboardHandler.Index)
		dashboard.GET("/charts/line", dashboardHandler.LineChart)
		dashboard.GET("/charts/bar", dashboardHandler.BarChart)
		dashboard.GET("/charts/pie", dashboardHandler.PieChart)
	}

	// 서버 시작
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadTemplates(r *gin.Engine) {
	tmpl := template.New("").Funcs(r.FuncMap)

	// 템플릿 디렉토리 순회
	err := filepath.Walk("web/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// 템플릿 이름 생성 (web/templates/ 이후 경로)
		name := strings.TrimPrefix(path, "web/templates/")

		// 파일 내용 읽기
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Warning: Could not read template %s: %v", path, err)
			return nil
		}

		// 템플릿 파싱
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			log.Printf("Warning: Could not parse template %s: %v", path, err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Warning: Error walking templates: %v", err)
	}

	r.SetHTMLTemplate(tmpl)
}
