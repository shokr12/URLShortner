# Shortify Frontend

Modern, high-performance dashboard for URL shortening.

## Setup Instructions

1. **Navigate to the frontend directory**:
   ```bash
   cd frontend
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Start the development server**:
   ```bash
   npm run dev
   ```

## Design System
- **Colors**: Action Blue (#0058be), Slate Surface (#f9f9ff)
- **Typography**: Inter (Extreme scale variance)
- **Aesthetics**: Glassmorphism, tonal layering, no-line borders.

## Backend Compatibility
The Go backend at `http://localhost:8080` has been updated with:
- `GET /links`: List all shortened URLs.
- `GET /stats`: Dashboard summary stats.
- `DELETE /links/:shortKey`: Remove digital assets.
- `GET /analytics/:shortKey`: Performance tracking metrics.
- CORS support for frontend communication.
