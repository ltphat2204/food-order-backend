# Food Order Backend (Golang)

A monolith backend for food ordering system (like ShopeeFood), written in Golang + GORM + PostgreSQL + Redis + Mailhog.

## Features
- Event Sourcing architecture
- **Event-Driven Architecture** with Redis Streams
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

## Event-Driven Architecture

This backend now uses Redis Streams for event-driven communication, replacing the previous pub/sub mechanism. Key benefits include:

- **Message Persistence**: Events are stored and can be replayed
- **Consumer Groups**: Multiple consumers can process events with load balancing
- **Guaranteed Delivery**: Message acknowledgment ensures reliable delivery
- **Scalability**: Better handling of high-throughput scenarios

See [Redis Streams Documentation](docs/redis-streams.md) for detailed implementation information.