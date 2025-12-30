package database

import (
	"log"

	"github.com/baltop/commet/internal/config"
	"github.com/baltop/commet/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return DB, nil
}

func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.DashboardData{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

func SeedSampleData() error {
	log.Println("Seeding sample dashboard data...")

	// 샘플 데이터가 이미 있는지 확인
	var count int64
	DB.Model(&models.DashboardData{}).Count(&count)
	if count > 0 {
		log.Println("Sample data already exists, skipping...")
		return nil
	}

	// 샘플 데이터 생성
	sampleData := []models.DashboardData{
		// 월별 매출 데이터 (라인 차트용)
		{Category: "sales", Label: "1월", Value: 1200},
		{Category: "sales", Label: "2월", Value: 1900},
		{Category: "sales", Label: "3월", Value: 3000},
		{Category: "sales", Label: "4월", Value: 2500},
		{Category: "sales", Label: "5월", Value: 2800},
		{Category: "sales", Label: "6월", Value: 3200},

		// 제품별 판매량 (바 차트용)
		{Category: "products", Label: "제품 A", Value: 450},
		{Category: "products", Label: "제품 B", Value: 320},
		{Category: "products", Label: "제품 C", Value: 280},
		{Category: "products", Label: "제품 D", Value: 190},
		{Category: "products", Label: "제품 E", Value: 520},

		// 트래픽 소스 (파이 차트용)
		{Category: "traffic", Label: "직접 방문", Value: 35},
		{Category: "traffic", Label: "검색 엔진", Value: 40},
		{Category: "traffic", Label: "소셜 미디어", Value: 15},
		{Category: "traffic", Label: "추천", Value: 10},
	}

	for _, data := range sampleData {
		if err := DB.Create(&data).Error; err != nil {
			return err
		}
	}

	log.Println("Sample data seeded successfully")
	return nil
}
