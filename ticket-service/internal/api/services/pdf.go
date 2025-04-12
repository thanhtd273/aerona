package services

import (
	"bytes"
	"html/template"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type Passenger struct {
	Title          string
	FullName       string
	Type           string
	Route          string
	HandLuggage    string
	CheckedLuggage string
	TotalWeight    string
}

type TicketData struct {
	DateTime         string
	DepartureTime    string
	DepartureCity    string
	DepartureCode    string
	DepartureAirport string
	Duration         string
	ArrivalDate      string
	ArrivalTime      string
	ArrivalCity      string
	ArrivalCode      string
	ArrivalAirport   string
	BookingID        string
	PNR              string
	Refundable       string
	Passengers       []Passenger
}

type ServiceItem struct {
	ServiceName string
	BookingCode string
	Quantity    int
	UnitPrice   string
	TotalAmount string
}

type ReceiptData struct {
	ReceiptNumber   string
	TransactionDate string
	BookingID       string
	PaymentMethod   string
	PaymentStatus   string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	PassengerName   string
	ServiceItems    []ServiceItem
	Total           string
	ServiceFee      string
	RescheduleFee   string
	TotalPaid       string
}

type PDFService struct {
	templatePath string
}

func NewPDFService(templatePath string) *PDFService {
	return &PDFService{templatePath: templatePath}
}

func (s *PDFService) GenerateTicketPDF(data TicketData) ([]byte, error) {
	templ, err := template.ParseFiles(s.templatePath)
	if err != nil {
		return nil, err
	}

	var htmlBuffer bytes.Buffer
	err = templ.Execute(&htmlBuffer, data)
	if err != nil {
		return nil, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes())))

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}
	pdfBytes := pdfg.Bytes()
	return pdfBytes, nil
}

func (s *PDFService) Add(a, b int) int {
	return a + b
}

func (s *PDFService) GenerateReceiptPDF(data ReceiptData, outputPath string) error {
	tmpl, err := template.New("receipt.html").Funcs(template.FuncMap{
		"add": s.Add,
	}).ParseFiles(s.templatePath)
	if err != nil {
		return err
	}

	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, data)
	if err != nil {
		return err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes())))

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	err = pdfg.WriteFile(outputPath)
	if err != nil {
		return err
	}

	log.Printf("PDF generated and saved to %s", outputPath)
	return nil
}

// func (s *PDFService) GenerateSampleTicket() error {
// 	data := TicketData{
// 		DateTime:         "Monday, 20 January 2025 / Thứ Hai, 20 tháng 1 2025",
// 		DepartureTime:    "23:59",
// 		DepartureCity:    "Ho Chi Minh City",
// 		DepartureCode:    "SGN",
// 		DepartureAirport: "Tansonnhat Intl\nSân bay Tân Sơn Nhất",
// 		Duration:         "01:25",
// 		ArrivalDate:      "21 tháng 1",
// 		ArrivalTime:      "01:25",
// 		ArrivalCity:      "Da Nang",
// 		ArrivalCode:      "DAD",
// 		ArrivalAirport:   "Da Nang Airport\nSân bay Đà Nẵng",
// 		BookingID:        "1210070435",
// 		PNR:              "ORkJ3C",
// 		Refundable:       "Không hoàn vé",
// 		Passengers: []Passenger{
// 			{
// 				Title:          "Ông",
// 				FullName:       "TRINH DINH THANH",
// 				Type:           "Người lớn",
// 				Route:          "SGN - DAD",
// 				HandLuggage:    "7 KG Hành lý xách tay",
// 				CheckedLuggage: "Không được nhận hành lý ký gửi",
// 				TotalWeight:    "23 KG Hành lý",
// 			},
// 		},
// 	}

// 	return s.GenerateTicketPDF(data, "temp/ticket.pdf")
// }
