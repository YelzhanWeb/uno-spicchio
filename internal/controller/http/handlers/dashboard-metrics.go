package handlers

import (
	"net/http"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
)

// GetTodayMetrics — отдельный эндпоинт с метриками только за сегодня.
// Его ты уже использовала для "Today's Overview".
func (h *AnalyticsHandler) GetTodayMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	now := time.Now()
	loc := now.Location()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 1) Сводка продаж за сегодня
	summary, err := h.analyticsService.GetSalesSummary(ctx, startOfDay, endOfDay)
	if err != nil {
		response.InternalError(w, "failed to get sales summary for today")
		return
	}

	// 2) Статистика заказов за сегодня (всего + оплаченные)
	orderStats, err := h.analyticsService.GetOrderStats(ctx, startOfDay, endOfDay)
	if err != nil {
		response.InternalError(w, "failed to get order stats for today")
		return
	}

	metrics := map[string]interface{}{
		"today_revenue":    summary.TotalRevenue,
		"revenue_change":   summary.RevenueChange,
		"total_orders":     orderStats.TotalOrders,
		"orders_change":    summary.OrdersChange,
		"avg_order_value":  summary.AverageOrderValue,
		"avg_value_change": summary.AvgValueChange,
		"completed_orders": orderStats.CompletedOrders,
	}

	response.Success(w, metrics)
}

// GetDashboardMetrics — главный эндпоинт для дашборда:
// /api/analytics/dashboard?period=today|yesterday|current_month
func (h *AnalyticsHandler) GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "today"
	}

	now := time.Now()
	loc := now.Location()

	var from, to time.Time
	switch period {
	case "yesterday":
		y := now.AddDate(0, 0, -1)
		from = time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, loc)
		to = from.Add(24 * time.Hour)
	case "current_month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		to = from.AddDate(0, 1, 0)
	default: // today
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		to = from.Add(24 * time.Hour)
	}

	// 1) Сводка по оплаченных заказам (revenue, avg check и т.п.)
	summary, err := h.analyticsService.GetSalesSummary(ctx, from, to)
	if err != nil {
		response.InternalError(w, "failed to get sales summary")
		return
	}

	// 2) Статистика заказов: ВСЕ и оплаченные
	orderStats, err := h.analyticsService.GetOrderStats(ctx, from, to)
	if err != nil {
		response.InternalError(w, "failed to get order stats")
		return
	}

	// 3) Продажи по категориям и популярные блюда
	categorySales, err := h.analyticsService.GetSalesByCategory(ctx, from, to)
	if err != nil {
		response.InternalError(w, "failed to get category sales")
		return
	}

	popularDishes, err := h.analyticsService.GetPopularDishes(ctx, from, to, 5)
	if err != nil {
		response.InternalError(w, "failed to get popular dishes")
		return
	}

	// ⚠️ ВАЖНО: total_orders берём из orderStats, а не из summary
	summaryMap := map[string]interface{}{
		"total_revenue":       summary.TotalRevenue,
		"total_orders":        orderStats.TotalOrders, // 👈 ВСЕ заказы
		"average_order_value": summary.AverageOrderValue,
		"revenue_change":      summary.RevenueChange,
		"orders_change":       summary.OrdersChange,
		"avg_value_change":    summary.AvgValueChange,
		"completed_orders":    orderStats.CompletedOrders, // только paid
	}

	data := map[string]interface{}{
		"summary":        summaryMap,
		"category_sales": categorySales,
		"popular_dishes": popularDishes,
	}

	response.Success(w, data)
}
