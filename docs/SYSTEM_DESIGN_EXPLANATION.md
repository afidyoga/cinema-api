# Penjelasan Solusi Sistem — Point A

## 1. Solusi Pemilihan Tempat Duduk (Seat Locking)

### Permasalahan
Ketika ratusan user membuka halaman pemilihan kursi secara bersamaan, ada risiko dua user memilih kursi yang sama dan keduanya berhasil membayar, padahal kursi hanya bisa diisi satu orang.

### Solusi: Optimistic UI + Redis Distributed Lock

**Alur teknis:**

1. User membuka halaman pemilihan kursi → sistem membaca status kursi dari **Redis cache** (bukan langsung ke PostgreSQL) untuk performa tinggi.
2. User klik kursi → API menjalankan operasi **SET NX** (Set if Not Exists) di Redis dengan key:
   ```
   seat_lock:{schedule_id}:{seat_id} = {user_id}
   TTL = 600 detik (10 menit)
   ```
3. Jika berhasil: kursi berhasil dikunci untuk user tersebut, user lanjut ke halaman bayar.
4. Jika gagal (key sudah ada): kursi sudah dikunci oleh user lain → tampilkan pesan "Kursi sedang dipesan orang lain, silakan pilih kursi lain."
5. Jika user tidak menyelesaikan pembayaran dalam 10 menit, Redis otomatis hapus key (TTL expired), kursi kembali tersedia.
6. Setelah pembayaran berhasil, status kursi dipindah ke tabel `tickets` di PostgreSQL dengan status `active` , lock di Redis dihapus.

**Kenapa Redis?**
- Operasi `SET NX` di Redis bersifat **atomic**, tidak ada race condition meski ribuan request datang bersamaan.
- Latensi Redis sangat rendah (<1ms), sehingga tampilan denah kursi real-time tidak memberatkan database.
- TTL otomatis menangani skenario user yang menutup browser sebelum bayar.

**Skalabilitas:**
- Redis dijalankan dalam mode **Cluster** untuk high availability.
- API server bisa di-scale horizontal (multiple instance) karena lock state ada di Redis, bukan memory lokal.

---

## 2. Pencatatan dan Restok Tiket

### Alur pencatatan saat tiket terjual:

```
Pembayaran berhasil
  → INSERT INTO transactions (status = PAID)
  → INSERT INTO tickets (status = ACTIVE, ticket_code = UUID unik)
  → UPDATE seat di Redis: hapus lock, tandai SOLD
  → Kurangi available_count di cache schedule
  → Kirim e-ticket ke email/WhatsApp customer
```

### Restok otomatis (kursi kembali tersedia):

Restok terjadi dalam tiga kondisi:

| Kondisi | Trigger | Aksi sistem |
|---------|---------|-------------|
| TTL expired | Redis hapus lock otomatis | Kursi kembali tersedia tanpa intervensi |
| Pembayaran gagal/timeout | Webhook dari payment gateway | API hapus lock Redis, kursi tersedia kembali |
| Pembatalan dari bioskop | Admin set schedule = CANCELLED | Batch job restok semua kursi, proses refund |

### Cara sistem mengetahui kursi tersedia:

Denah kursi dibangun dari kombinasi:
- Tabel `seats` (semua kursi di studio tersebut)
- Tabel `tickets` WHERE schedule_id = X AND status = 'active' → kursi sudah terjual
- Redis keys `seat_lock:{schedule_id}:*` kursi sedang dikunci (dalam proses bayar)
- Sisanya = **tersedia**

---

## 3. Refund dan Pembatalan dari Pihak Bioskop

### Alur refund saat bioskop membatalkan jadwal:

```
Admin set schedules.status = 'cancelled'
  → Background worker mendeteksi perubahan status
  → Query semua tickets WHERE schedule_id = X AND status = 'active'
  → Untuk setiap tiket:
      a. UPDATE tickets SET status = 'cancelled'
      b. Query transaction terkait
      c. Kirim request refund ke payment gateway (Midtrans/Xendit)
      d. INSERT INTO refunds (status = 'pending', amount = total_amount)
      e. UPDATE transactions SET payment_status = 'refunded'
  → Update denah kursi di Redis: hapus semua lock & sold mark
  → Kirim notifikasi ke semua customer (email + push notif)
     → "Jadwal tayang [film] pada [waktu] dibatalkan. Refund akan masuk dalam 1-3 hari kerja."
```

### Jaminan konsistensi:

- Proses refund berjalan dalam **database transaction** — jika satu langkah gagal, semua di-rollback.
- Background worker menggunakan **idempotency key** sehingga jika worker restart, tidak ada double refund.
- Status refund dipantau lewat webhook dari payment gateway, update `refunds.status` ke `processed` atau `failed`.
- Jika refund gagal dari payment gateway, masuk antrian retry dengan exponential backoff.

### Kebijakan refund:

- Pembatalan dari bioskop, refund 100% otomatis, tidak perlu request dari customer.
- Customer tidak bisa request refund mandiri melalui API ini (scope di luar test).

---

## 4. Keputusan Desain Lain

**Mengapa tidak pakai database lock (SELECT FOR UPDATE)?**
Karena `SELECT FOR UPDATE` di PostgreSQL akan memblokir baris selama transaksi berjalan. Jika user membuka halaman pemilihan kursi 10 menit, maka baris kursi di-lock 10 menit di level DB, tidak scalable untuk ribuan concurrent user.

**Mengapa TTL 10 menit?**
Cukup waktu untuk menyelesaikan pembayaran, tapi tidak terlalu lama sehingga kursi tidak "tertahan" terlalu lama jika user meninggalkan halaman.

**Read replica PostgreSQL:**
Query baca (denah kursi, daftar jadwal) diarahkan ke **replica** sehingga primary database tidak terbebani oleh read traffic yang tinggi.
