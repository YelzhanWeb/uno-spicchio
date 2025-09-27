package app

import (
	"database/sql"
	"net/http"
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

func Routes(db *sql.DB) *http.ServeMux {
	return http.NewServeMux()
}
