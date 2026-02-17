# Incident Report App

A full-stack incident reporting system built with **Go** (backend) and **React + TypeScript** (frontend), using **SQLite** for data persistence.

## Tech Stack
| Layer    | Tech                       |
|----------|----------------------------|
| Frontend | React 18, TypeScript, Vite |
| Backend  | Go 1.21, Chi router        |
| Database | SQLite (pure-Go, no CGO)   |

## Getting Started

### Backend
cd backend && go mod tidy && go run main.go

### Frontend
cd frontend && npm install && npm run dev