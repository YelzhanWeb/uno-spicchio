package postgre

// internal/adapters/postgres/analytics_repository.go

import (
	"context"
	"database/sql"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetSalesSummary(ctx context.Context, from, to time.Time) (*domain.SalesSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(total), 0) as total_revenue,
			COUNT(*) as total_orders,
			COALESCE(AVG(total), 0) as average_order_value
		FROM orders
		WHERE status = 'paid' AND created_at >= $1 AND created_at < $2`

	summary := &domain.SalesSummary{}
	err := r.db.QueryRowContext(ctx, query, from, to).Scan(
		&summary.TotalRevenue,
		&summary.TotalOrders,
		&summary.AverageOrderValue,
	)

	return summary, err
}

func (r *AnalyticsRepository) GetPreviousPeriodSummary(ctx context.Context, from, to time.Time) (*domain.SalesSummary, error) {
	duration := to.Sub(from)
	prevFrom := from.Add(-duration)
	prevTo := from

	return r.GetSalesSummary(ctx, prevFrom, prevTo)
}

func (r *AnalyticsRepository) GetSalesByCategory(ctx context.Context, from, to time.Time) ([]domain.CategorySale, error) {
	query := `
		SELECT 
			c.id,
			c.name,
			COALESCE(SUM(oi.price * oi.qty), 0) as revenue
		FROM categories c
		LEFT JOIN dishes d ON d.category_id = c.id
		LEFT JOIN order_items oi ON oi.dish_id = d.id
		LEFT JOIN orders o ON o.id = oi.order_id
		WHERE o.status = 'paid' AND o.created_at >= $1 AND o.created_at < $2
		GROUP BY c.id, c.name
		ORDER BY revenue DESC`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.CategorySale
	var totalRevenue float64

	// First pass: collect all sales and calculate total
	for rows.Next() {
		var sale domain.CategorySale
		if err := rows.Scan(&sale.CategoryID, &sale.CategoryName, &sale.Revenue); err != nil {
			return nil, err
		}
		totalRevenue += sale.Revenue
		sales = append(sales, sale)
	}

	// Second pass: calculate percentages
	for i := range sales {
		if totalRevenue > 0 {
			sales[i].Percentage = (sales[i].Revenue / totalRevenue) * 100
		}
	}

	return sales, rows.Err()
}

func (r *AnalyticsRepository) GetPopularDishes(ctx context.Context, from, to time.Time, limit int) ([]domain.PopularDish, error) {
	query := `
		SELECT 
			d.id,
			d.name,
			COALESCE(SUM(oi.qty), 0) as qty_sold,
			COALESCE(SUM(oi.price * oi.qty), 0) as revenue
		FROM dishes d
		JOIN order_items oi ON oi.dish_id = d.id
		JOIN orders o ON o.id = oi.order_id
		WHERE o.status = 'paid' AND o.created_at >= $1 AND o.created_at < $2
		GROUP BY d.id, d.name
		ORDER BY revenue DESC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []domain.PopularDish
	for rows.Next() {
		var dish domain.PopularDish
		if err := rows.Scan(&dish.DishID, &dish.DishName, &dish.QtySold, &dish.Revenue); err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	return dishes, rows.Err()
}

func (r *AnalyticsRepository) GetWaiterPerformance(ctx context.Context, from, to time.Time) ([]domain.WaiterPerformance, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			COUNT(o.id) as order_count,
			COALESCE(SUM(o.total), 0) as revenue,
			COALESCE(AVG(o.total), 0) as avg_check
		FROM users u
		LEFT JOIN orders o ON o.waiter_id = u.id AND o.status = 'paid' 
			AND o.created_at >= $1 AND o.created_at < $2
		WHERE u.role = 'waiter' AND u.is_active = true
		GROUP BY u.id, u.username
		ORDER BY revenue DESC`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performance []domain.WaiterPerformance
	for rows.Next() {
		var perf domain.WaiterPerformance
		if err := rows.Scan(&perf.WaiterID, &perf.WaiterName, &perf.OrderCount, &perf.Revenue, &perf.AvgCheck); err != nil {
			return nil, err
		}
		performance = append(performance, perf)
	}

	return performance, rows.Err()
}

func (r *AnalyticsRepository) GetOrderStats(
	ctx context.Context,
	from, to time.Time,
) (*domain.OrderStats, error) {

	// Считаем ВСЕ заказы за период + сколько из них paid
	query := `
        SELECT 
            COUNT(*) AS total_orders,
            COALESCE(SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END), 0) AS completed_orders
        FROM orders
        WHERE created_at >= $1 AND created_at < $2;
    `

	stats := &domain.OrderStats{}

	err := r.db.QueryRowContext(ctx, query, from, to).Scan(
		&stats.TotalOrders,     // все заказы (любой статус)
		&stats.CompletedOrders, // только paid
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Если нужны PeakHour и т.п. — можно потом отдельным запросом посчитать.
	// Пока просто 0, чтобы структура была валидной.
	stats.PeakHour = 0

	return stats, nil
}

func (r *AnalyticsRepository) GetIngredientTurnover(ctx context.Context, from, to time.Time) ([]domain.IngredientTurnover, error) {
	query := `
		SELECT 
			i.id,
			i.name,
			i.unit,
			i.qty as current_stock,
			COALESCE(SUM(di.qty_per_dish * oi.qty), 0) as used
		FROM ingredients i
		LEFT JOIN dish_ingredients di ON di.ingredient_id = i.id
		LEFT JOIN order_items oi ON oi.dish_id = di.dish_id
		LEFT JOIN orders o ON o.id = oi.order_id AND o.status = 'paid'
			AND o.created_at >= $1 AND o.created_at < $2
		GROUP BY i.id, i.name, i.unit, i.qty
		ORDER BY used DESC`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turnover []domain.IngredientTurnover
	for rows.Next() {
		var turn domain.IngredientTurnover
		if err := rows.Scan(&turn.IngredientID, &turn.IngredientName, &turn.Unit, &turn.CurrentStock, &turn.Used); err != nil {
			return nil, err
		}
		turnover = append(turnover, turn)
	}

	return turnover, rows.Err()
}

func (r *AnalyticsRepository) GetTableUtilization(ctx context.Context, from, to time.Time) ([]domain.TableUtilization, error) {
	// Calculate total hours in period for utilization rate
	totalHours := to.Sub(from).Hours()

	query := `
		SELECT 
			t.id,
			t.name,
			COUNT(o.id) as times_used,
			COALESCE(SUM(EXTRACT(EPOCH FROM (o.updated_at - o.created_at)) / 3600), 0) as hours_used
		FROM tables t
		LEFT JOIN orders o ON o.table_number = t.id AND o.status = 'paid'
			AND o.created_at >= $1 AND o.created_at < $2
		GROUP BY t.id, t.name
		ORDER BY times_used DESC`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var utilization []domain.TableUtilization
	for rows.Next() {
		var util domain.TableUtilization
		var hoursUsed float64
		if err := rows.Scan(&util.TableID, &util.TableName, &util.TimesUsed, &hoursUsed); err != nil {
			return nil, err
		}

		if totalHours > 0 {
			util.UtilizationRate = (hoursUsed / totalHours) * 100
		}

		utilization = append(utilization, util)
	}

	return utilization, rows.Err()
}

func (r *AnalyticsRepository) GetHourlyRevenue(ctx context.Context, date time.Time) ([]domain.HourlyRevenue, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT 
			EXTRACT(HOUR FROM created_at) as hour,
			COALESCE(SUM(total), 0) as revenue,
			COUNT(*) as orders
		FROM orders
		WHERE status = 'paid' AND created_at >= $1 AND created_at < $2
		GROUP BY EXTRACT(HOUR FROM created_at)
		ORDER BY hour`

	rows, err := r.db.QueryContext(ctx, query, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hourlyData []domain.HourlyRevenue
	for rows.Next() {
		var data domain.HourlyRevenue
		if err := rows.Scan(&data.Hour, &data.Revenue, &data.Orders); err != nil {
			return nil, err
		}
		hourlyData = append(hourlyData, data)
	}

	return hourlyData, rows.Err()
}

func (r *AnalyticsRepository) GetDishAvailability(ctx context.Context) ([]domain.DishAvailability, error) {
	query := `
		SELECT 
			d.id,
			d.name,
			d.price,
			COALESCE(
				FLOOR(
					MIN(
						CASE 
							WHEN di.qty_per_dish > 0 
							THEN (i.qty / di.qty_per_dish)
							ELSE NULL
						END
					)
				),
				0
			) AS portions_left
		FROM dishes d
		JOIN dish_ingredients di ON di.dish_id = d.id
		JOIN ingredients i ON i.id = di.ingredient_id
		WHERE d.is_active = TRUE
		GROUP BY d.id, d.name, d.price
		ORDER BY d.id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DishAvailability
	for rows.Next() {
		var d domain.DishAvailability
		if err := rows.Scan(&d.DishID, &d.Name, &d.Price, &d.PortionsLeft); err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, rows.Err()
}
