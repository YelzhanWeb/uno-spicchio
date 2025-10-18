// internal/core/services/analytics_service.go
package usecase

import (
	"context"

	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type AnalyticsService struct {
	analyticsRepo ports.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo ports.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: analyticsRepo}
}

// GetDashboard returns complete dashboard data for the specified period
func (s *AnalyticsService) GetDashboard(ctx context.Context, period domain.PeriodType, from, to time.Time) (*domain.DashboardData, error) {
	// Calculate date range based on period type
	fromDate, toDate := s.calculatePeriod(period, from, to)

	// Get sales summary with comparison
	summary, err := s.GetSalesSummary(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get sales by category
	categorySales, err := s.analyticsRepo.GetSalesByCategory(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get popular dishes (top 5)
	popularDishes, err := s.analyticsRepo.GetPopularDishes(ctx, fromDate, toDate, 5)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardData{
		Summary:       *summary,
		CategorySales: categorySales,
		PopularDishes: popularDishes,
	}, nil
}

// GetSalesSummary returns sales summary with comparison to previous period
func (s *AnalyticsService) GetSalesSummary(ctx context.Context, from, to time.Time) (*domain.SalesSummary, error) {
	// Get current period summary
	currentSummary, err := s.analyticsRepo.GetSalesSummary(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Get previous period summary for comparison
	previousSummary, err := s.analyticsRepo.GetPreviousPeriodSummary(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Calculate percentage changes
	currentSummary.RevenueChange = calculatePercentageChange(currentSummary.TotalRevenue, previousSummary.TotalRevenue)
	currentSummary.OrdersChange = calculatePercentageChange(float64(currentSummary.TotalOrders), float64(previousSummary.TotalOrders))
	currentSummary.AvgValueChange = calculatePercentageChange(currentSummary.AverageOrderValue, previousSummary.AverageOrderValue)

	return currentSummary, nil
}

// GetSalesByCategory returns sales grouped by category
func (s *AnalyticsService) GetSalesByCategory(ctx context.Context, from, to time.Time) ([]domain.CategorySale, error) {
	return s.analyticsRepo.GetSalesByCategory(ctx, from, to)
}

// GetPopularDishes returns top selling dishes
func (s *AnalyticsService) GetPopularDishes(ctx context.Context, from, to time.Time, limit int) ([]domain.PopularDish, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	return s.analyticsRepo.GetPopularDishes(ctx, from, to, limit)
}

// GetWaiterPerformance returns performance metrics for waiters
func (s *AnalyticsService) GetWaiterPerformance(ctx context.Context, from, to time.Time) ([]domain.WaiterPerformance, error) {
	return s.analyticsRepo.GetWaiterPerformance(ctx, from, to)
}

// GetOrderStats returns order statistics
func (s *AnalyticsService) GetOrderStats(ctx context.Context, from, to time.Time) (*domain.OrderStats, error) {
	return s.analyticsRepo.GetOrderStats(ctx, from, to)
}

// GetIngredientTurnover returns ingredient usage statistics
func (s *AnalyticsService) GetIngredientTurnover(ctx context.Context, from, to time.Time) ([]domain.IngredientTurnover, error) {
	return s.analyticsRepo.GetIngredientTurnover(ctx, from, to)
}

// GetTableUtilization returns table usage statistics
func (s *AnalyticsService) GetTableUtilization(ctx context.Context, from, to time.Time) ([]domain.TableUtilization, error) {
	return s.analyticsRepo.GetTableUtilization(ctx, from, to)
}

// GetHourlyRevenue returns revenue by hour for a specific date
func (s *AnalyticsService) GetHourlyRevenue(ctx context.Context, date time.Time) ([]domain.HourlyRevenue, error) {
	return s.analyticsRepo.GetHourlyRevenue(ctx, date)
}

// calculatePeriod calculates the from/to dates based on period type
func (s *AnalyticsService) calculatePeriod(period domain.PeriodType, customFrom, customTo time.Time) (from, to time.Time) {
	now := time.Now()
	location := now.Location()

	switch period {
	case domain.PeriodYesterday:
		from = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, location)
		to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	case domain.PeriodToday:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
		to = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)

	case domain.PeriodCurrentMonth:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
		to = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, location)

	case domain.PeriodCustom:
		from = customFrom
		to = customTo

	default:
		// Default to today
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
		to = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	}

	return from, to
}

// calculatePercentageChange calculates the percentage change between two values
func calculatePercentageChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // 100% increase from 0
	}
	return ((current - previous) / previous) * 100
}

// Helper function to round float to 2 decimal places
func roundFloat(val float64) float64 {
	return float64(int(val*100)) / 100
}
