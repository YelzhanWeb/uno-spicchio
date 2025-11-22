// internal/adapters/http/handlers/analytics_handler.go
package handlers

import (
	"net/http"

	"strconv"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
)

type AnalyticsHandler struct {
	analyticsService ports.AnalyticsService
}

func NewAnalyticsHandler(analyticsService ports.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboard returns complete dashboard data
// Query params: period (yesterday|today|current_month|custom), from (for custom), to (for custom)
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	periodStr := r.URL.Query().Get("period")
	if periodStr == "" {
		periodStr = "today" // default
	}

	period := domain.PeriodType(periodStr)

	var from, to time.Time
	var err error

	// Parse custom dates if provided
	if period == domain.PeriodCustom {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		if fromStr == "" || toStr == "" {
			response.BadRequest(w, "custom period requires 'from' and 'to' parameters")
			return
		}

		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			response.BadRequest(w, "invalid 'from' date format, use YYYY-MM-DD")
			return
		}

		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			response.BadRequest(w, "invalid 'to' date format, use YYYY-MM-DD")
			return
		}

		// Add one day to 'to' to include the entire day
		to = to.Add(24 * time.Hour)
	}

	dashboard, err := h.analyticsService.GetDashboard(r.Context(), period, from, to)
	if err != nil {
		response.InternalError(w, "failed to get dashboard data")
		return
	}

	response.Success(w, dashboard)
}

// GetSalesSummary returns sales summary with comparison
func (h *AnalyticsHandler) GetSalesSummary(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	summary, err := h.analyticsService.GetSalesSummary(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get sales summary")
		return
	}

	response.Success(w, summary)
}

// GetSalesByCategory returns sales grouped by category
func (h *AnalyticsHandler) GetSalesByCategory(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	sales, err := h.analyticsService.GetSalesByCategory(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get sales by category")
		return
	}

	response.Success(w, sales)
}

// GetPopularDishes returns top selling dishes
func (h *AnalyticsHandler) GetPopularDishes(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Parse limit parameter
	limit := 10 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			response.BadRequest(w, "invalid limit parameter")
			return
		}
	}

	dishes, err := h.analyticsService.GetPopularDishes(r.Context(), from, to, limit)
	if err != nil {
		response.InternalError(w, "failed to get popular dishes")
		return
	}

	response.Success(w, dishes)
}

// GetWaiterPerformance returns performance metrics for waiters
func (h *AnalyticsHandler) GetWaiterPerformance(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	performance, err := h.analyticsService.GetWaiterPerformance(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get waiter performance")
		return
	}

	response.Success(w, performance)
}

// GetOrderStats returns order statistics
func (h *AnalyticsHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	stats, err := h.analyticsService.GetOrderStats(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get order stats")
		return
	}

	response.Success(w, stats)
}

// GetIngredientTurnover returns ingredient usage statistics
func (h *AnalyticsHandler) GetIngredientTurnover(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	turnover, err := h.analyticsService.GetIngredientTurnover(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get ingredient turnover")
		return
	}

	response.Success(w, turnover)
}

// GetTableUtilization returns table usage statistics
func (h *AnalyticsHandler) GetTableUtilization(w http.ResponseWriter, r *http.Request) {
	from, to, err := h.parseDateRange(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	utilization, err := h.analyticsService.GetTableUtilization(r.Context(), from, to)
	if err != nil {
		response.InternalError(w, "failed to get table utilization")
		return
	}

	response.Success(w, utilization)
}

// GetHourlyRevenue returns revenue by hour for a specific date
func (h *AnalyticsHandler) GetHourlyRevenue(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(w, "invalid date format, use YYYY-MM-DD")
		return
	}

	hourlyData, err := h.analyticsService.GetHourlyRevenue(r.Context(), date)
	if err != nil {
		response.InternalError(w, "failed to get hourly revenue")
		return
	}

	response.Success(w, hourlyData)
}

// parseDateRange is a helper function to parse from/to query parameters
func (h *AnalyticsHandler) parseDateRange(r *http.Request) (from, to time.Time, err error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Default to today if not provided
	now := time.Now()
	if fromStr == "" {
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toStr == "" {
		to = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		// Add one day to include the entire day
		to = to.Add(24 * time.Hour)
	}

	return from, to, nil
}

func (h *AnalyticsHandler) GetDishAvailability(w http.ResponseWriter, r *http.Request) {
	data, err := h.analyticsService.GetDishAvailability(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get dish availability")
		return
	}

	response.Success(w, data)
}
