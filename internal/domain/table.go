package domain

type TableStatus string

const (
	TableBusy    TableStatus = "busy"
	TableReserve TableStatus = "reserve"
	TableFree    TableStatus = "free"
)

type Table struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Status TableStatus `json:"status"`
}
