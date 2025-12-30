package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/baltop/commet/internal/middleware"
	"github.com/baltop/commet/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GET /dashboard - 대시보드 메인 페이지
func (h *DashboardHandler) Index(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)

	stats := h.dashboardService.GetSummaryStats()

	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"title":          "대시보드",
		"user":           claims,
		"totalUsers":     stats["totalUsers"],
		"totalRevenue":   stats["totalRevenue"],
		"totalOrders":    stats["totalOrders"],
		"conversionRate": stats["conversionRate"],
	})
}

// GET /dashboard/charts/line - 라인 차트 데이터 (HTMX partial)
func (h *DashboardHandler) LineChart(c *gin.Context) {
	data, err := h.dashboardService.GetSalesData()
	if err != nil {
		c.HTML(http.StatusOK, "components/alert.html", gin.H{
			"type":    "error",
			"message": "데이터를 불러오는데 실패했습니다.",
		})
		return
	}

	labelsJSON, _ := json.Marshal(data.Labels)
	valuesJSON, _ := json.Marshal(data.Values)

	c.HTML(http.StatusOK, "dashboard/partials/chart_line.html", gin.H{
		"labels": string(labelsJSON),
		"values": string(valuesJSON),
		"title":  "월별 매출 추이",
	})
}

// GET /dashboard/charts/bar - 바 차트 데이터 (HTMX partial)
func (h *DashboardHandler) BarChart(c *gin.Context) {
	data, err := h.dashboardService.GetProductsData()
	if err != nil {
		c.HTML(http.StatusOK, "components/alert.html", gin.H{
			"type":    "error",
			"message": "데이터를 불러오는데 실패했습니다.",
		})
		return
	}

	labelsJSON, _ := json.Marshal(data.Labels)
	valuesJSON, _ := json.Marshal(data.Values)

	c.HTML(http.StatusOK, "dashboard/partials/chart_bar.html", gin.H{
		"labels": string(labelsJSON),
		"values": string(valuesJSON),
		"title":  "제품별 판매량",
	})
}

// GET /dashboard/charts/pie - 파이 차트 데이터 (HTMX partial)
func (h *DashboardHandler) PieChart(c *gin.Context) {
	data, err := h.dashboardService.GetTrafficData()
	if err != nil {
		c.HTML(http.StatusOK, "components/alert.html", gin.H{
			"type":    "error",
			"message": "데이터를 불러오는데 실패했습니다.",
		})
		return
	}

	labelsJSON, _ := json.Marshal(data.Labels)
	valuesJSON, _ := json.Marshal(data.Values)

	c.HTML(http.StatusOK, "dashboard/partials/chart_pie.html", gin.H{
		"labels": string(labelsJSON),
		"values": string(valuesJSON),
		"title":  "트래픽 소스",
	})
}
