package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Party adalah pihak yang terlibat dalam dokumen.
type Party struct {
	Name     string
	Company  string
	Address  string
	City     string
	Email    string
	Phone    string
	Position string
}

// InvoiceItem adalah satu baris barang pada invoice.
type InvoiceItem struct {
	Name     string
	Qty      int
	Price    decimal.Decimal
	Subtotal decimal.Decimal
}

// InvoiceTotals berisi ringkasan nominal invoice; Tax adalah PPN 11%.
type InvoiceTotals struct {
	Subtotal decimal.Decimal
	Tax      decimal.Decimal
	Total    decimal.Decimal
}

// DeliveryItem adalah satu baris barang pada surat jalan.
type DeliveryItem struct {
	Name string
	Qty  int
	Unit string
	Note string
}

// ReportRow adalah satu baris tabel laporan.
type ReportRow struct {
	No    int
	Name  string
	City  string
	Email string
	Value decimal.Decimal
	Date  time.Time
}

// Report adalah payload lengkap yang dibutuhkan layout PDF.
type Report struct {
	Type          string
	Title         string
	DocID         string
	Number        string
	Date          time.Time
	Seed          uint64
	Rows          int
	From          Party
	To            Party
	Purpose       string
	InvoiceItems  []InvoiceItem
	InvoiceTotals InvoiceTotals
	DeliveryItems []DeliveryItem
	ReportRows    []ReportRow
	Amount        decimal.Decimal
	Paragraphs    []string
}

// typeTitles read-only: judul bawaan per jenis dokumen.
var typeTitles = map[string]string{
	"invoice":     "Invoice",
	"surat-resmi": "Surat Resmi",
	"laporan":     "Laporan",
	"kuitansi":    "Kuitansi",
	"surat-jalan": "Surat Jalan",
	"sertifikat":  "Sertifikat",
}

// proseSentences read-only: kalimat pelengkap paragraf.
var proseSentences = [...]string{
	"Informasi pada dokumen ini disusun secara ringkas dan mudah dibaca.",
	"Setiap bagian disusun agar dapat dipakai sebagai contoh format resmi.",
	"Nomor dan tanggal dokumen dihasilkan secara konsisten dari nilai seed.",
	"Isi dokumen dapat digunakan untuk pengujian tampilan dan tata letak.",
	"Penggunaan nama dan alamat di sini hanya untuk ilustrasi.",
	"Struktur dokumen mengikuti kebiasaan surat-menyurat di Indonesia.",
	"Bagian penutup dilengkapi ruang tanda tangan sebagaimana lazimnya.",
	"Format nominal mengikuti gaya penulisan rupiah yang berlaku.",
}

// Generate membentuk payload dokumen secara deterministik dari jenis, jumlah baris, dan seed.
func Generate(docType string, rows int, seed uint64) (Report, error) {
	title, ok := typeTitles[docType]
	if !ok {
		return Report{}, fmt.Errorf("jenis dokumen tidak dikenal: %q", docType)
	}
	if rows < 0 {
		rows = 0
	}

	faker := gofakeit.New(seed)
	date := time.Unix(int64(seed), 0).UTC()

	report := Report{
		Type:       docType,
		Title:      title,
		DocID:      uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s-%d", docType, seed))).String(),
		Number:     fmt.Sprintf("%03d/GENPDF/%d", faker.IntRange(1, 999), date.Year()),
		Date:       date,
		Seed:       seed,
		Rows:       rows,
		From:       newParty(faker),
		To:         newParty(faker),
		Purpose:    faker.RandomString(purposes[:]),
		Paragraphs: newParagraphs(faker, title),
	}

	switch docType {
	case "invoice":
		report.InvoiceItems, report.InvoiceTotals = newInvoice(faker, rows)
	case "surat-jalan":
		report.DeliveryItems = newDeliveryItems(faker, rows)
	case "laporan":
		report.ReportRows = newReportRows(faker, rows, date)
	case "kuitansi":
		report.Amount = decimal.NewFromInt(int64(faker.IntRange(100, 50_000) * 1000))
	}

	return report, nil
}

func newParty(faker *gofakeit.Faker) Party {
	name := faker.RandomString(firstNames[:]) + " " + faker.RandomString(lastNames[:])
	city := faker.RandomString(cities[:])

	return Party{
		Name:     name,
		Company:  faker.RandomString(companies[:]),
		Address:  faker.RandomString(streets[:]) + " No. " + strconv.Itoa(faker.IntRange(1, 200)),
		City:     city,
		Email:    strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@example.com",
		Phone:    faker.Numerify("08##########"),
		Position: faker.RandomString(positions[:]),
	}
}

func newInvoice(faker *gofakeit.Faker, rows int) ([]InvoiceItem, InvoiceTotals) {
	items := make([]InvoiceItem, 0, rows)
	subtotal := decimal.Zero

	for i := 0; i < rows; i++ {
		qty := faker.IntRange(1, 10)
		price := decimal.NewFromInt(int64(faker.IntRange(10_000, 500_000)))
		line := price.Mul(decimal.NewFromInt(int64(qty)))

		items = append(items, InvoiceItem{
			Name:     faker.RandomString(goods[:]),
			Qty:      qty,
			Price:    price,
			Subtotal: line,
		})
		subtotal = subtotal.Add(line)
	}

	tax := subtotal.Mul(decimal.NewFromInt(11)).Div(decimal.NewFromInt(100)).Round(2)
	return items, InvoiceTotals{Subtotal: subtotal, Tax: tax, Total: subtotal.Add(tax)}
}

func newDeliveryItems(faker *gofakeit.Faker, rows int) []DeliveryItem {
	items := make([]DeliveryItem, 0, rows)
	for i := 0; i < rows; i++ {
		items = append(items, DeliveryItem{
			Name: faker.RandomString(goods[:]),
			Qty:  faker.IntRange(1, 20),
			Unit: faker.RandomString(units[:]),
			Note: faker.RandomString(deliveryNotes[:]),
		})
	}
	return items
}

func newReportRows(faker *gofakeit.Faker, rows int, date time.Time) []ReportRow {
	result := make([]ReportRow, 0, rows)
	for i := 0; i < rows; i++ {
		name := faker.RandomString(firstNames[:]) + " " + faker.RandomString(lastNames[:])
		result = append(result, ReportRow{
			No:    i + 1,
			Name:  name,
			City:  faker.RandomString(cities[:]),
			Email: strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@example.com",
			Value: decimal.NewFromInt(int64(faker.IntRange(100, 10_000) * 1000)),
			Date:  date.AddDate(0, 0, faker.IntRange(0, 365)),
		})
	}
	return result
}

func newParagraphs(faker *gofakeit.Faker, title string) []string {
	first := fmt.Sprintf("Dokumen %s ini dihasilkan secara otomatis untuk keperluan contoh dan demonstrasi.", title)
	last := "Seluruh data fiktif pada dokumen ini tidak merujuk pada orang, perusahaan, atau transaksi nyata."

	return []string{
		first + " " + faker.RandomString(proseSentences[:]) + " " + faker.RandomString(proseSentences[:]),
		faker.RandomString(proseSentences[:]) + " " + faker.RandomString(proseSentences[:]) + " " + faker.RandomString(proseSentences[:]),
		last + " " + faker.RandomString(proseSentences[:]) + " " + faker.RandomString(proseSentences[:]),
	}
}
