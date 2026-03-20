# 🤖 WhatsApp Bot — Go + whatsmeow

Bot WhatsApp berbasis Go menggunakan library [whatsmeow](https://github.com/tulir/whatsmeow) (implementasi native Go dari protokol WhatsApp Web). Mendukung koneksi via **QR Code** maupun **Pairing Code**, sistem plugin yang modular, dan penanganan error yang robust.

---

## ✨ Fitur

| Fitur | Keterangan |
|---|---|
| 🔌 Plugin System | Tambah fitur baru cukup dengan membuat 1 file |
| 📱 QR & Pairing Code | Pilih metode login yang kamu suka |
| 🔄 Auto Reconnect | Reconnect otomatis dengan exponential backoff |
| 🛡️ Panic Recovery | Plugin yang crash tidak mematikan bot |
| 💾 Session Persistence | Sesi disimpan di SQLite, tidak perlu login ulang |
| 📋 Dynamic Help | `!help` otomatis menampilkan semua plugin terdaftar |

---

## 📁 Struktur Proyek

```
whatsapp-bot/
├── main.go                     # Entry point
├── go.mod / go.sum
├── .env.example                # Template konfigurasi
├── Makefile                    # Shortcut commands
├── Dockerfile
│
├── config/
│   └── config.go               # Load konfigurasi dari env
│
├── logger/
│   └── logger.go               # Wrapper slog
│
├── bot/
│   ├── bot.go                  # Core: connect, event loop, dispatch
│   └── plugins.go              # Daftar plugin yang aktif
│
└── plugin/
    ├── plugin.go               # Interface Plugin + Registry
    └── plugins/
        ├── ping.go             # !ping
        ├── help.go             # !help / !menu
        ├── echo.go             # !echo
        ├── info.go             # !info
        └── _template.go        # Template untuk plugin baru
```

---

## 🚀 Cara Pakai

### 1. Prasyarat

- Go 1.22+
- GCC (untuk SQLite CGO): `apt install gcc` / `brew install gcc`

### 2. Clone & Install

```bash
git clone <repo>
cd whatsapp-bot
cp .env.example .env
go mod tidy
```

### 3. Konfigurasi `.env`

```env
# Pilih: "qr" atau "pairing"
CONNECT_METHOD=qr

# Wajib jika CONNECT_METHOD=pairing (tanpa + atau spasi)
PHONE_NUMBER=6281234567890

BOT_PREFIX=!
BOT_NAME=GoBot
AUTO_RECONNECT=true
```

### 4. Jalankan

```bash
# Via QR Code (default)
make run

# Via Pairing Code
CONNECT_METHOD=pairing PHONE_NUMBER=6281234567890 make run

# Build binary
make build
./bot
```

### 5. Login

**QR Code:** Scan QR yang muncul di terminal dengan WhatsApp kamu.

**Pairing Code:** Buka WhatsApp → *Linked Devices* → *Link a Device* → masukkan kode 8 digit yang muncul.

---

## 🔌 Membuat Plugin Baru

1. Salin `plugin/plugins/_template.go` → `plugin/plugins/nama_plugin.go`
2. Isi 4 method wajib:

```go
func (p *MyPlugin) Name() string        { return "myplugin" }
func (p *MyPlugin) Description() string { return "Deskripsi plugin" }
func (p *MyPlugin) Commands() []string  { return []string{"cmd", "alias"} }
func (p *MyPlugin) Execute(ctx context.Context, msg *plugin.MessageContext) error {
    return msg.Reply("Hello!")
}
```

3. Daftarkan di `bot/plugins.go`:

```go
r.Register(plugins.NewMyPlugin())
```

### API `MessageContext`

```go
msg.Args        // []string — [0]=command, [1:]=argumen
msg.Text        // string — pesan lengkap
msg.Sender      // types.JID — pengirim
msg.Chat        // types.JID — chat/grup
msg.IsGroup     // bool
msg.Client      // *whatsmeow.Client — akses penuh ke WA

msg.Reply("teks")  // balas dengan quote
msg.Send("teks")   // kirim tanpa quote
```

---

## 🐳 Docker

```bash
make docker-build
make docker-run

# Atau dengan docker-compose:
docker-compose up -d
```

---

## ⚙️ Environment Variables

| Variable | Default | Keterangan |
|---|---|---|
| `CONNECT_METHOD` | `qr` | `qr` atau `pairing` |
| `PHONE_NUMBER` | — | Nomor HP untuk pairing code |
| `DB_PATH` | `./data/sessions.db` | Path database sesi |
| `BOT_PREFIX` | `!` | Prefix perintah |
| `BOT_NAME` | `GoBot` | Nama bot |
| `OWNER_JID` | — | JID owner (opsional) |
| `AUTO_RECONNECT` | `true` | Reconnect otomatis |

---

## 🛡️ Error Handling

- **Panic recovery** — setiap plugin dibungkus `defer recover()`, panic tidak mematikan bot
- **Exponential backoff** — reconnect dimulai dari 5 detik, max 60 detik
- **User notification** — error plugin otomatis dikirim ke user (`❌ Error: ...`)
- **Graceful shutdown** — `Ctrl+C` disconnect dengan bersih

---

## 📦 Dependensi Utama

| Library | Fungsi |
|---|---|
| `go.mau.fi/whatsmeow` | WhatsApp Web protocol |
| `github.com/mattn/go-sqlite3` | Session storage |
| `github.com/mdp/qrterminal` | Render QR di terminal |
