# URL Shortener

**Developed a highly scalable and performant URL Shortener service, capable of handling high-throughput redirection with sub-millisecond latency, by implementing a cache-aside architecture using Redis, persistent storage with PostgreSQL, and offloading click analytics to asynchronous background Go workers.**

## 🚀 Features

- **Blazing Fast Redirects**: Uses Redis as a caching layer to achieve sub-millisecond redirect times.
- **Persistent Storage**: Uses PostgreSQL for robust and reliable storage of URL mappings and analytics.
- **Asynchronous Analytics**: Click tracking is offloaded to a background worker to ensure redirection latency is not impacted.
- **Custom Keys**: Users can specify custom short URLs.
- **Graceful Shutdown**: Ensures all in-flight requests and background jobs complete before shutting down.

---

## 🛠 Tech Stack

- **Go (Golang)**: Core programming language.
- **Gin**: High-performance HTTP web framework.
- **PostgreSQL**: Relational database (`jackc/pgx/v5` driver).
- **Redis**: In-memory data store for caching and queues (`redis/go-redis/v9`).
- **godotenv**: Environment variable management.

---

## 📦 Project Structure

```text
├── cmd/
│   ├── main.go            # Application entry point
│   └── .env               # Environment configuration
├── config/
│   ├── config.go          # Configuration loading and DB/Redis initialization
│   └── migrate.go         # Auto-migration scripts
├── internal/
│   ├── handler/           # HTTP controllers and routing (Gin)
│   ├── model/             # Data structures and entities
│   ├── repo/              # Database and Redis interaction layers
│   ├── service/           # Business logic
│   └── worker/            # Background worker for async tasks
```

---

## 🚦 Getting Started

### Prerequisites

Ensure you have the following installed:
- [Go](https://go.dev/dl/) (v1.20+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/url-shortener.git
   cd url-shortener
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Configure Environment Variables:**
   Create a `.env` file in the `cmd/` directory (or use the one provided) and configure your database and Redis connections:
   ```env
   DB_USER=postgres
   DB_PASS=yourpassword
   DB_NAME=UrlShortner
   DB_HOST=localhost
   DB_PORT=5432

   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASS=
   REDIS_DB=0

   SERVER_ADDR=:8080
   BASE_URL=http://localhost:8080
   ```

4. **Run the application:**
   ```bash
   cd cmd
   go run main.go
   ```

*(Note: The application will automatically create the required `url` table in PostgreSQL on startup).*

---

## 📡 API Endpoints

### 1. Shorten a URL
Creates a new short URL.

- **URL**: `/shorten`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
      "original_url": "https://www.verylongurl.com/example/path/here",
      "custom_key": "my-link" // Optional
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
      "short_key": "my-link",
      "short_url": "http://localhost:8080/my-link"
  }
  ```

### 2. Redirect
Redirects to the original URL and increments the click count asynchronously.

- **URL**: `/:shortKey`
- **Method**: `GET`
- **Success Response**: `302 Found` (Redirects to original URL)

---

## ⚙️ Architecture Deep Dive

### Read Path Optimization (Cache-Aside Pattern)
When a user visits a short link, the `repo.GetURLByShortKey` method first checks **Redis**. 
- **Cache Hit**: Returns the URL immediately.
- **Cache Miss**: Queries **PostgreSQL**, retrives the URL, and updates the Redis cache for subsequent requests.

### Write Path & Asynchronous Analytics
When a short link is visited, updating the `click_count` synchronously would slow down the user's redirect. Instead, the handler pushes the `shortKey` to a Go channel. A background `ClickWorker` processes this channel, incrementing the counters in both Postgres and Redis without blocking the HTTP response.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
