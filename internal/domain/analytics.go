// internal/core/domain/analytics.go
package domain

import "time"

// DashboardData contains all data for the dashboard
type DashboardData struct {
	Summary       SalesSummary   `json:"summary"`
	CategorySales []CategorySale `json:"category_sales"`
	PopularDishes []PopularDish  `json:"popular_dishes"`
}

// SalesSummary contains key metrics for sales
type SalesSummary struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalOrders       int     `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
	RevenueChange     float64 `json:"revenue_change"`   // % change from previous period
	OrdersChange      float64 `json:"orders_change"`    // % change from previous period
	AvgValueChange    float64 `json:"avg_value_change"` // % change from previous period
}

// CategorySale represents sales data for a category
type CategorySale struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Revenue      float64 `json:"revenue"`
	Percentage   float64 `json:"percentage"` // percentage of total revenue
}

// PopularDish represents a popular dish with sales data
type PopularDish struct {
	DishID   int     `json:"dish_id"`
	DishName string  `json:"dish_name"`
	QtySold  int     `json:"qty_sold"`
	Revenue  float64 `json:"revenue"`
}

// WaiterPerformance represents performance metrics for a waiter
type WaiterPerformance struct {
	WaiterID   int     `json:"waiter_id"`
	WaiterName string  `json:"waiter_name"`
	OrderCount int     `json:"order_count"`
	Revenue    float64 `json:"revenue"`
	AvgCheck   float64 `json:"avg_check"`
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders     int     `json:"total_orders"`
	CompletedOrders int     `json:"completed_orders"`
	AverageTime     float64 `json:"average_time"` // in minutes
	PeakHour        int     `json:"peak_hour"`    // hour of day (0-23)
}

// IngredientTurnover represents ingredient usage statistics
type IngredientTurnover struct {
	IngredientID   int     `json:"ingredient_id"`
	IngredientName string  `json:"ingredient_name"`
	Used           float64 `json:"used"`
	Unit           string  `json:"unit"`
	CurrentStock   float64 `json:"current_stock"`
}

// TableUtilization represents table usage statistics
type TableUtilization struct {
	TableID         int     `json:"table_id"`
	TableName       string  `json:"table_name"`
	TimesUsed       int     `json:"times_used"`
	UtilizationRate float64 `json:"utilization_rate"` // percentage
}

// HourlyRevenue represents revenue by hour
type HourlyRevenue struct {
	Hour    int     `json:"hour"` // 0-23
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

// PeriodType represents the type of period for analytics
type PeriodType string

const (
	PeriodYesterday    PeriodType = "yesterday"
	PeriodToday        PeriodType = "today"
	PeriodCurrentMonth PeriodType = "current_month"
	PeriodCustom       PeriodType = "custom"
)

// DateRange represents a date range for analytics
type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type DishAvailability struct {
	DishID       int     `json:"dish_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	PortionsLeft int     `json:"portions_left"`
}
