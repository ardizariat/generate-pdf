package data

// Indonesian reference data for dummy documents. Every slice is package-level
// and must not be modified at runtime.

var indonesianData = dataset{
	firstNames: []string{
		"Budi", "Siti", "Agus", "Dewi", "Rizky", "Andi",
		"Maya", "Hendra", "Lestari", "Fajar", "Nur", "Yusuf",
	},
	lastNames: []string{
		"Santoso", "Wijaya", "Pratama", "Sari", "Hidayat", "Kusuma",
		"Nugroho", "Rahmawati", "Saputra", "Wibowo", "Purnama", "Handayani",
	},
	cities: []string{
		"Jakarta", "Bandung", "Surabaya", "Medan", "Semarang", "Makassar",
		"Yogyakarta", "Denpasar", "Palembang", "Balikpapan", "Malang", "Bogor",
	},
	streets: []string{
		"Jl. Merdeka", "Jl. Sudirman", "Jl. Diponegoro", "Jl. Gatot Subroto",
		"Jl. Ahmad Yani", "Jl. Pemuda", "Jl. Kartini", "Jl. Pahlawan",
		"Jl. Melati", "Jl. Kenanga", "Jl. Cempaka", "Jl. Anggrek",
	},
	companies: []string{
		"PT Sinar Nusantara", "CV Karya Mandiri", "PT Bumi Sejahtera",
		"PT Cahaya Abadi", "CV Rizki Jaya", "PT Digital Prima",
		"PT Karya Utama", "CV Sumber Rejeki", "PT Anugerah Sentosa",
		"PT Mitra Bersama", "CV Tunas Harapan", "PT Global Data",
	},
	goods: []string{
		"Kertas A4", "Tinta Printer", "Pulpen Gel", "Map Plastik",
		"Buku Agenda", "Kalkulator", "Stapler", "Spidol Whiteboard",
		"Flashdisk 32GB", "Kabel HDMI", "Kursi Kerja", "Meja Lipat",
	},
	units: []string{
		"pcs", "box", "rim", "unit", "set", "lusin",
		"karung", "meter", "liter", "kg", "pak", "lembar",
	},
	positions: []string{
		"Direktur", "Manajer Keuangan", "Kepala Bagian Umum", "Staf Administrasi",
		"Sekretaris", "Bendahara", "Manajer Operasional", "Kepala Gudang",
		"Supervisor", "Koordinator", "Administrasi", "Pemilik",
	},
	purposes: []string{
		"Pembayaran jasa konsultasi", "Pembelian perlengkapan kantor",
		"Pembayaran sewa ruangan", "Jasa pemeliharaan jaringan",
		"Pembelian perangkat komputer", "Biaya pelatihan karyawan",
		"Pembayaran jasa desain", "Pengadaan alat tulis",
		"Biaya perjalanan dinas", "Pembayaran tagihan listrik",
		"Jasa pencetakan dokumen", "Pembelian bahan baku",
	},
	deliveryNotes: []string{
		"Baik", "Sesuai pesanan", "Tanpa cacat", "Dikemas rapi",
		"Fragile", "Kirim hari ini", "Sudah diperiksa", "Segel utuh",
		"Jumlah sesuai", "Dus penyok", "Prioritas", "Simpan di tempat kering",
	},
	prose: []string{
		"Informasi pada dokumen ini disusun secara ringkas dan mudah dibaca.",
		"Setiap bagian disusun agar dapat dipakai sebagai contoh format resmi.",
		"Nomor dan tanggal dokumen dihasilkan secara konsisten dari nilai seed.",
		"Isi dokumen dapat digunakan untuk pengujian tampilan dan tata letak.",
		"Penggunaan nama dan alamat di sini hanya untuk ilustrasi.",
		"Struktur dokumen mengikuti kebiasaan surat-menyurat di Indonesia.",
		"Bagian penutup dilengkapi ruang tanda tangan sebagaimana lazimnya.",
		"Format nominal mengikuti gaya penulisan rupiah yang berlaku.",
	},
	intro:              "Dokumen %s ini dihasilkan secara otomatis untuk keperluan contoh dan demonstrasi.",
	fictional:          "Seluruh data fiktif pada dokumen ini tidak merujuk pada orang, perusahaan, atau transaksi nyata.",
	phonePattern:       "08##########",
	addressNumberFirst: false,
	price:              moneyRange{min: 10_000, max: 500_000, div: 1},
	receipt:            moneyRange{min: 100_000, max: 50_000_000, div: 1},
	report:             moneyRange{min: 100_000, max: 10_000_000, div: 1},
}
