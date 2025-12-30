package services

import (
	"github.com/baltop/commet/internal/models"
	"github.com/baltop/commet/internal/repository"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

type ChartData struct {
	Labels []string  `json:"labels"`
	Values []float64 `json:"values"`
}

func (s *DashboardService) GetSalesData() (*ChartData, error) {
	data, err := s.dashboardRepo.GetDataByCategory("sales")
	if err != nil {
		return nil, err
	}
	return toChartData(data), nil
}

func (s *DashboardService) GetProductsData() (*ChartData, error) {
	data, err := s.dashboardRepo.GetDataByCategory("products")
	if err != nil {
		return nil, err
	}
	return toChartData(data), nil
}

func (s *DashboardService) GetTrafficData() (*ChartData, error) {
	data, err := s.dashboardRepo.GetDataByCategory("traffic")
	if err != nil {
		return nil, err
	}
	return toChartData(data), nil
}

func (s *DashboardService) GetSummaryStats() map[string]interface{} {
	// 샘플 요약 통계
	return map[string]interface{}{
		"totalUsers":    1234,
		"totalRevenue":  45678.90,
		"totalOrders":   567,
		"conversionRate": 3.2,
	}
}

func toChartData(data []models.DashboardData) *ChartData {
	chartData := &ChartData{
		Labels: make([]string, len(data)),
		Values: make([]float64, len(data)),
	}
	for i, d := range data {
		chartData.Labels[i] = d.Label
		chartData.Values[i] = d.Value
	}
	return chartData
}
