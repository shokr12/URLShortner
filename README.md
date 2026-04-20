# Shortify | Modern URL Shortener

**Shortify** is a high-performance, full-stack URL shortening service. It features a scalable Go backend with a cache-aside architecture and an architectural, editorial-style Next.js dashboard for digital curators.

---

## 🚀 Key Features

### Backend (Performance & Scalability)
- **Sub-Millisecond Redirects**: Cache-aside architecture using Redis for near-instant link resolution.
- **Persistent Reliability**: PostgreSQL storage for long-term data integrity.
- **Asynchronous Analytics**: Click tracking offloaded to background Go workers to ensure minimal redirect latency.
- **Custom Keys**: Support for personalized short keys.
- **Graceful Shutdown**: Reliable cleanup of in-flight requests and background jobs.

### Frontend (Editorial Experience)
- **Precision Dashboard**: A "Digital Architect" inspired interface with tonal layering and glassmorphism.
- **Link Management**: Advanced data grid for viewing, copying, and deleting assets.
- **Performance Tracking**: Real-time analytics charts showing engagement trends.
- **Modern Tech**: Built with Next.js 14, Tailwind CSS, and Framer Motion for premium animations.

---

## 🛠 Tech Stack

### Backend
- **Go (Golang)**: Core logic and high-concurrency handling.
- **Gin**: High-performance HTTP web framework.
- **PostgreSQL**: Primary transactional database.
- **Redis**: In-memory caching and real-time click tracking.

### Frontend
- **Next.js 14**: React framework with App Router.
- **Tailwind CSS**: Architectural styling and layout.
- **Lucide React**: Premium iconography.
- **Recharts**: Data visualization for performance metrics.

---

## 📦 Project Structure

```text
├── cmd/               # Go Application Entry (main.go, .env)
├── config/            # DB/Redis Initialization & Config
├── internal/          # Domain Logic (Repo, Service, Handler, Worker)
├── frontend/          # Next.js Dashboard
│   ├── src/app/       # Routing and UI Components
│   └── tailwind.config.ts # Design System Configuration
└── README.md          # Project Documentation
```

---

## 🚦 Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) (v1.20+)
- [Node.js](https://nodejs.org/) (v18+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download) (or WSL/Docker)

### 1. Backend Setup
```bash
# 1. Install dependencies
go mod download

# 2. Configure .env in /cmd
# Default: localhost:8080, DB User: postgres, DB Name: UrlShortner

# 3. Start the server
cd cmd
go run main.go
```

### 2. Frontend Setup
```bash
# 1. Enter directory
cd frontend

# 2. Install dependencies
npm install

# 3. Start development server
npm run dev
```
Visit `http://localhost:3000` to access the dashboard.

---

## 📡 API Reference

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/shorten` | Create a new short URL |
| `GET` | `/:shortKey` | Redirect to original URL (Async analytics) |
| `GET` | `/links` | List all digital assets |
| `GET` | `/stats` | Global dashboard summary |
| `DELETE`| `/links/:key` | Remove an asset |
| `GET` | `/analytics/:key` | Detailed performance data |

---

## ⚙️ Architecture: Cache-Aside & Async Workers

Shortify follows the **Cache-Aside Pattern**:
1. Check **Redis** for the short link.
2. If **Cache Hit**, redirect immediately.
3. If **Cache Miss**, query **PostgreSQL**, update the cache, and redirect.

**Asynchronous Workflows**:
To maintain high performance, the redirect handler doesn't wait for database updates. It pushes click data to a Go channel, which is then processed by a background **ClickWorker** that increments counters in both Postgres and Redis without blocking the user.

---

## 🤝 Contributing
Contributions are welcome! Please feel free to open issues or pull requests.
1. Fork the repo.
2. Create your branch (`feature/AmazingFeature`).
3. Commit your changes.
4. Open a PR.
