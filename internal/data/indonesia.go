package data

// Indonesian reference data for dummy documents. Every slice is an array
// so it cannot be changed at runtime.

var firstNames = [...]string{
	"Budi", "Siti", "Agus", "Dewi", "Rizky", "Andi",
	"Maya", "Hendra", "Lestari", "Fajar", "Nur", "Yusuf",
}

var lastNames = [...]string{
	"Santoso", "Wijaya", "Pratama", "Sari", "Hidayat", "Kusuma",
	"Nugroho", "Rahmawati", "Saputra", "Wibowo", "Purnama", "Handayani",
}

var cities = [...]string{
	"Jakarta", "Bandung", "Surabaya", "Medan", "Semarang", "Makassar",
	"Yogyakarta", "Denpasar", "Palembang", "Balikpapan", "Malang", "Bogor",
}

var streets = [...]string{
	"Jl. Merdeka", "Jl. Sudirman", "Jl. Diponegoro", "Jl. Gatot Subroto",
	"Jl. Ahmad Yani", "Jl. Pemuda", "Jl. Kartini", "Jl. Pahlawan",
	"Jl. Melati", "Jl. Kenanga", "Jl. Cempaka", "Jl. Anggrek",
}

var companies = [...]string{
	"PT Sinar Nusantara", "CV Karya Mandiri", "PT Bumi Sejahtera",
	"PT Cahaya Abadi", "CV Rizki Jaya", "PT Digital Prima",
	"PT Karya Utama", "CV Sumber Rejeki", "PT Anugerah Sentosa",
	"PT Mitra Bersama", "CV Tunas Harapan", "PT Global Data",
}

var goods = [...]string{
	"Kertas A4", "Tinta Printer", "Pulpen Gel", "Map Plastik",
	"Buku Agenda", "Kalkulator", "Stapler", "Spidol Whiteboard",
	"Flashdisk 32GB", "Kabel HDMI", "Kursi Kerja", "Meja Lipat",
}

var units = [...]string{
	"pcs", "box", "rim", "unit", "set", "lusin",
	"karung", "meter", "liter", "kg", "pak", "lembar",
}

var positions = [...]string{
	"Direktur", "Manajer Keuangan", "Kepala Bagian Umum", "Staf Administrasi",
	"Sekretaris", "Bendahara", "Manajer Operasional", "Kepala Gudang",
	"Supervisor", "Koordinator", "Administrasi", "Pemilik",
}

var purposes = [...]string{
	"Pembayaran jasa konsultasi", "Pembelian perlengkapan kantor",
	"Pembayaran sewa ruangan", "Jasa pemeliharaan jaringan",
	"Pembelian perangkat komputer", "Biaya pelatihan karyawan",
	"Pembayaran jasa desain", "Pengadaan alat tulis",
	"Biaya perjalanan dinas", "Pembayaran tagihan listrik",
	"Jasa pencetakan dokumen", "Pembelian bahan baku",
}

var deliveryNotes = [...]string{
	"Baik", "Sesuai pesanan", "Tanpa cacat", "Dikemas rapi",
	"Fragile", "Kirim hari ini", "Sudah diperiksa", "Segel utuh",
	"Jumlah sesuai", "Dus penyok", "Prioritas", "Simpan di tempat kering",
}
