package http

import (
	"net/http"

	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

// func Routes(db *sql.DB) *http.ServeMux {
// 	mux := http.NewServeMux()

// 	repos := repository.NewRepositoryWithDB(db)
// 	service := service.NewService(repos)

// 	handlerInv := handler.NewInventoryHandler(service.InventoryService)
// 	handlerMenu := handler.NewMenuHandler(service.MenuService)
// 	handlerOrder := handler.NewOrderHandler(service.OrderService)
// 	handlerReport := handler.NewReportHandler(service.ReportService)

// 	mux.HandleFunc("/orders", handlerOrder.OrdersHandler)
// 	mux.HandleFunc("/orders/", handlerOrder.OrderItemHandler)
// 	mux.HandleFunc("/orders/batch-process", handlerOrder.BatchProcessOrders)

// 	mux.HandleFunc("/menu", handlerMenu.MenuHandler)
// 	mux.HandleFunc("/menu/", handlerMenu.MenuItemHandler)

// 	mux.HandleFunc("/inventory", handlerInv.InventoryHandler)
// 	mux.HandleFunc("/inventory/", handlerInv.InventoryItemHandler)

// 	mux.HandleFunc("/reports/total-sales", handlerReport.HandleGetTotalSales)
// 	mux.HandleFunc("/reports/popular-items", handlerReport.HandleGetPopularItems)
// 	mux.HandleFunc("/reports/search", handlerReport.HandleSearch)
// 	mux.HandleFunc("/reports/orderedItemsByPeriod", handlerOrder.GetOrderedItemsByPeriod)
// 	return mux
// }

func Router(service ports.Service) http.Handler {
	// r := chi.NewRouter()

	// handler := NewHandler(service)

	// r.Post("/user")

	return http.NewServeMux()
}
