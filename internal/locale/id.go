package locale

// indonesian is read-only after initialization.
var indonesian = Labels{
	TypeTitles: map[string]string{
		"invoice":       "Invoice",
		"letter":        "Surat Resmi",
		"report":        "Laporan",
		"receipt":       "Kuitansi",
		"delivery-note": "Surat Jalan",
		"certificate":   "Sertifikat",
	},
	TypeDescriptions: map[string]string{
		"invoice":       "Invoice penjualan dengan tabel barang dan PPN",
		"letter":        "Surat resmi dengan kop dan tanda tangan",
		"report":        "Laporan tabel yang mengalir multi-halaman",
		"receipt":       "Kuitansi pembayaran dengan nominal terbilang",
		"delivery-note": "Surat jalan pengiriman dengan dua tanda tangan",
		"certificate":   "Sertifikat landscape dengan QR verifikasi",
	},

	InvoiceNo: "No. Invoice",
	Date:      "Tanggal",
	Seller:    "Penjual",
	Buyer:     "Pembeli",
	Item:      "Barang",
	Qty:       "Qty",
	Price:     "Harga",
	Subtotal:  "Subtotal",
	Tax:       "PPN 11%",
	Total:     "Total",

	Ref:       "Nomor",
	Subject:   "Perihal",
	Attn:      "Kepada Yth.",
	Dear:      "Dengan hormat,",
	Sincerely: "Hormat kami,",

	PagePattern: "Halaman {current} dari {total}",

	No:    "No",
	Name:  "Nama",
	City:  "Kota",
	Email: "Email",
	Value: "Nilai",

	Receipt:      "KUITANSI",
	ReceiptNo:    "No",
	ReceivedFrom: "Telah terima dari",
	Amount:       "Uang sejumlah",
	InWords:      "Terbilang",
	ForPayment:   "Untuk pembayaran",
	ReceivedBy:   "Penerima",
	CurrencyName: "rupiah",

	DeliveryNote: "SURAT JALAN",
	DeliveryNo:   "No. Surat Jalan",
	To:           "Kepada:",
	Unit:         "Satuan",
	Note:         "Keterangan",
	Sender:       "Pengirim",
	Receiver:     "Penerima",

	Certificate:         "SERTIFIKAT",
	CertificateFallback: "Dokumen ini dihasilkan secara otomatis dan seluruh datanya fiktif.",

	PromptFileCount:  "Jumlah file",
	PromptTypes:      "Jenis file (bisa lebih dari satu, dipakai berurutan)",
	PromptPages:      "Jumlah halaman per file",
	PromptInteger:    "Masukkan angka bulat",
	PromptMinOneFile: "Minimal 1 file",
	PromptMinOnePage: "Minimal 1 halaman",
	PromptMinOneType: "Pilih minimal satu jenis",
	ProgressRows:     "[%d/%d] %s · target %d baris",
	ProgressPages:    "[%d/%d] %s · target %d halaman",
	BatchFailed:      "gagal  %s: %v\n",
	BatchOK:          "ok     %s  %d halaman  %d bytes\n",
	BatchSummary:     "ringkasan: %d sukses, %d gagal\n",
}
