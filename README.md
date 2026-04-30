# Relay Take-Home

## Local Setup (<5min)

1. Start Postgres:
```bash
docker compose up -d db
```

2. Run migrations:
```bash
docker exec -i relay-db-1 psql -U postgres -d relay < migrations/001_initial_schema.sql
```

3. Backend:
```bash
cd api
go mod tidy
go run main.go
```

4. Frontend:
```bash
cd web
npm install
npm run dev
```

API: http://localhost:8080/health
Web: http://localhost:3000/chat

## Demo Flow
1. /chat - Send message as customer → AI replies citing KB → session ends with sentiment/ticket
2. POST /events complaint → rule creates ticket
3. /admin/events - See firings

See DESIGN.md AI_LOG.md

