.PHONY: dev-backend dev-frontend install build build-backend lint-backend lint-frontend clean

## Install all dependencies
install:
	cd backend && go mod tidy
	cd frontend && npm install

## Run backend dev server
dev-backend:
	cd backend && go run main.go

## Run frontend dev server
dev-frontend:
	cd frontend && npm run dev

## Build frontend for production
build:
	cd frontend && npm run build

## Build backend binary
build-backend:
	cd backend && go build -o bin/server main.go

## Remove build artifacts and DB
clean:
	rm -f backend/incidents.db backend/bin/server
	rm -rf frontend/dist

## Run go vet + check for common issues
lint-backend:
	cd backend && go vet ./...

# Run ESLint for frontend
lint-frontend:
	cd frontend && npm run lint