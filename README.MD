# Food Order Backend (Golang)

A monolith backend for food ordering system (like ShopeeFood), written in Golang + GORM + MySQL + Redis + Mailhog.

## Features
- Event Sourcing architecture
- REST API + real-time WebSocket updates
- Notification (email + Vue toast)
- Dev setup with Docker Compose

## Dev Start
```bash
make deps
make run
```