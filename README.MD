# Food Order Backend (Golang)

A monolith backend for food ordering system (like ShopeeFood), written in Golang + GORM + MySQL + Redis + Mailhog.

## Features
- Event Sourcing architecture
- REST API + real-time WebSocket updates
- Notification (email + Vue toast)
- Dev setup with Docker Compose
- **Data Consistency Testing** - Unit tests verify aggregate state consistency

## Dev Start
```bash
make deps
make run
```

## Testing Data Consistency
make test-consistency