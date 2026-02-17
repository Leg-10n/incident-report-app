# Incident Report App

A full-stack incident reporting system built with **Go** (backend) and **React + TypeScript** (frontend), using **SQLite** for data persistence.

## Tech Stack

| Layer    | Tech                       |
| -------- | -------------------------- |
| Frontend | React 18, TypeScript, Vite |
| Backend  | Go 1.21, Chi router        |
| Database | SQLite (pure-Go, no CGO)   |

---

## 🚀 Quick Start

### 1. Install all dependencies

```bash
make install
```

### 2. Run backend server

```bash
make dev-backend
```

Backend runs on: http://localhost:8080

### 3. Run frontend dev server (new terminal)

```bash
make dev-frontend
```

Frontend runs on: http://localhost:5173

---

## 🧪 Linting

Run backend lint:

```bash
make lint-backend
```

Run frontend lint:

```bash
make lint-frontend
```

---

## 🏗 Build (optional)

Build frontend for production:

```bash
make build
```

Build backend binary:

```bash
make build-backend
```

---

## 🧹 Clean project

Remove database and build files:

```bash
make clean
```

---

## Environment Variables (optional)

You can configure:

```
PORT
FRONTEND_ORIGIN
```

Defaults will be used if not provided.
