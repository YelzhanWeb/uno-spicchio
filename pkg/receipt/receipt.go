package receipt

import (
	"fmt"
	"os"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jung-kurt/gofpdf"
)

// статичный баркод для красоты
var barcodeValue = "8937261273610"

func GenerateOrderReceiptPDF(order *domain.Order, waiterName, folder string) (string, error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/order_%d.pdf", folder, order.ID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 15, 25)
	pdf.AddPage()

	// --- HEADER ---
	pdf.SetFont("Helvetica", "B", 24)
	pdf.CellFormat(0, 12, "UNO Spicchio", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(0, 10, "Rakhmet!", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(0, 8, "Sizdin esepshotynyz daiyn!", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 20)
	pdf.CellFormat(0, 8, "Qurmetpen IS-45!", "", 1, "C", false, 0, "")

	pdf.Ln(5)
	pdf.Line(25, pdf.GetY(), 185, pdf.GetY())
	pdf.Ln(5)

	// --- TICKET INFO ---
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(90, 5, "TICKET ID", "", 0, "L", false, 0, "")
	pdf.CellFormat(90, 5, "Amount", "", 1, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(90, 6, fmt.Sprintf("%012d", order.ID), "", 0, "L", false, 0, "")
	pdf.CellFormat(90, 6, fmt.Sprintf("%.2f KZT", order.Total), "", 1, "R", false, 0, "")

	pdf.Ln(3)

	// Date / Table-Info
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(90, 5, "DATE & TIME", "", 0, "L", false, 0, "")
	pdf.CellFormat(90, 5, "Table / Waiter", "", 1, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(90, 6, order.CreatedAt.Format("02 Jan 2006 at 15:04"), "", 0, "L", false, 0, "")
	pdf.CellFormat(90, 6,
		fmt.Sprintf("Table #%d  |  %s", order.TableNumber, waiterName),
		"", 1, "R", false, 0, "",
	)

	pdf.Ln(4)
	pdf.Line(25, pdf.GetY(), 185, pdf.GetY())
	pdf.Ln(6)

	// --- ITEMS ---
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 6, "Order items", "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 11)

	for _, item := range order.Items {
		itemName := fmt.Sprintf("x%d  %s", item.Qty, item.Dish.Name)

		// LEFT: name
		pdf.CellFormat(120, 6, itemName, "", 0, "L", false, 0, "")
		// RIGHT: price
		pdf.CellFormat(50, 6,
			fmt.Sprintf("%.2f", item.Price*float64(item.Qty)),
			"", 1, "R", false, 0, "")
	}

	// line before TOTAL
	pdf.Ln(3)
	pdf.Line(25, pdf.GetY(), 185, pdf.GetY())
	pdf.Ln(2)

	// --- TOTAL ---
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(120, 8, "Total", "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 8, fmt.Sprintf("%.2f KZT", order.Total), "", 1, "R", false, 0, "")

	pdf.Ln(10)

	// --- FOOTER TEXT ---
	pdf.SetFont("Helvetica", "", 10)
	footerText := fmt.Sprintf("Order #%d   /   %s", order.ID, order.CreatedAt.Format("02.01.2006 15:04"))
	pdf.CellFormat(0, 6, footerText, "", 1, "C", false, 0, "")

	pdf.Ln(4)

	// --- BARCODE ---
	pdf.SetFont("Courier", "", 36)
	pdf.CellFormat(0, 20, barcodeValue, "", 1, "C", false, 0, "")

	// Save
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}
