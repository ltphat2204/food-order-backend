package service

import (
	"encoding/json"
	"testing"
	"time"

	"food-order-backend/internal/app/order/model"
)

// BenchmarkReadPerformance so sánh hiệu suất đọc giữa Event Replay vs Read Model
func BenchmarkReadPerformance(b *testing.B) {
	// Setup test data - tạo một order với nhiều events
	orderID := "ORD_BENCHMARK_001"
	events := createLargeEventSequence(orderID, 50) // 50 events cho 1 order
	
	b.Run("EventReplay", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := replayEventsForBenchmark(events)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Simulate read model (pre-computed state)
	cachedState := precomputeReadModel(events)
	
	b.Run("ReadModel", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = readFromCachedModel(cachedState)
		}
	})
}

// createLargeEventSequence tạo một chuỗi events lớn để benchmark
func createLargeEventSequence(orderID string, eventCount int) []model.EventStore {
	events := make([]model.EventStore, 0, eventCount)
	baseTime := time.Now()
	
	// Always start with OrderCreated
	events = append(events, model.EventStore{
		AggregateID: orderID,
		EventType:   "OrderCreated",
		EventData:   `{"user_id":123,"restaurant_id":456,"status":"PENDING","note":"Benchmark order"}`,
		CreatedAt:   baseTime,
	})
	
	// Add various events to simulate real usage
	eventTypes := []string{
		"ShipperAssigned",
		"RestaurantAccepted", 
		"CookingStarted",
		"OrderPicked",
		"OrderDelivered",
	}
	
	for i := 1; i < eventCount; i++ {
		eventType := eventTypes[(i-1)%len(eventTypes)]
		var eventData string
		
		switch eventType {
		case "ShipperAssigned":
			eventData = `{"shipper_id":"S001"}`
		case "RestaurantAccepted":
			eventData = `{"merchant_id":"M001","time":"10:30"}`
		case "CookingStarted":
			eventData = `{"estimated_time":"20 minutes"}`
		case "OrderPicked":
			eventData = `{"pickup_time":"11:15"}`
		case "OrderDelivered":
			eventData = `{"delivery_time":"11:45","receiver_info":"Customer"}`
		}
		
		events = append(events, model.EventStore{
			AggregateID: orderID,
			EventType:   eventType,
			EventData:   eventData,
			CreatedAt:   baseTime.Add(time.Duration(i) * time.Minute),
		})
	}
	
	return events
}

// replayEventsForBenchmark simulates event replay (CPU intensive)
func replayEventsForBenchmark(events []model.EventStore) (model.Order, error) {
	return simulateComplexEventReplay(events)
}

// precomputeReadModel simulates creating a read model from events
func precomputeReadModel(events []model.EventStore) model.Order {
	// In real scenario, this would be pre-computed and stored in database
	state, _ := replayEventsForBenchmark(events)
	return state
}

// readFromCachedModel simulates reading from pre-computed read model (realistic database read)
func readFromCachedModel(cachedState model.Order) model.Order {
	// Simulate realistic database operations:
	
	// 1. Database connection overhead (very small)
	time.Sleep(100 * time.Nanosecond)
	
	// 2. Simple SQL query: SELECT * FROM orders WHERE order_id = ?
	// This is much faster than complex event replay but still has some cost
	time.Sleep(50 * time.Nanosecond)
	// 3. Single row lookup with index (very fast but not instant)
	result := model.Order{
		OrderID:      cachedState.OrderID,
		UserID:       cachedState.UserID,
		RestaurantID: cachedState.RestaurantID,
		Status:       cachedState.Status,
		Note:         cachedState.Note,
	}
	
	// 4. Object creation and field copying (minimal overhead)
	return result
}

// simulateComplexEventReplay adds more realistic complexity to event replay
func simulateComplexEventReplay(events []model.EventStore) (model.Order, error) {
	if len(events) == 0 {
		return model.Order{}, nil
	}

	state := model.Order{
		OrderID: events[0].AggregateID,
	}

	// 1. Database query to get all events (with JOINS, WHERE, ORDER BY)
	// Simulate database query overhead (already sorted by ORDER BY created_at)
	time.Sleep(time.Duration(len(events)) * 50 * time.Nanosecond)

	// 2. Process each event (JSON parsing + business logic)
	// Events are already sorted by database ORDER BY clause
	for _, evt := range events {
		// JSON unmarshaling (CPU intensive)
		var data model.OrderEventData
		if err := json.Unmarshal([]byte(evt.EventData), &data); err != nil {
			return model.Order{}, err
		}

		// Business logic application (state transitions)
		switch evt.EventType {
		case "OrderCreated":
			state.UserID = data.UserID
			state.RestaurantID = data.RestaurantID
			state.Status = data.Status
			state.Note = data.Note
		case "OrderCanceled":
			state.Status = "CANCELED"
		case "ShipperAssigned":
			state.Status = "SHIPPER_ASSIGNED"
		case "RestaurantAccepted":
			state.Status = "RESTAURANT_ACCEPTED"
		case "ShipperConfirmedWithRestaurant":
			state.Status = "SHIPPER_CONFIRMED"
		case "CookingStarted":
			state.Status = "COOKING"
		case "OrderPicked":
			state.Status = "PICKED"
		case "OrderDelivered":
			state.Status = "DELIVERED"
		}
		
		// Simulate additional business logic complexity
		time.Sleep(10 * time.Nanosecond)
	}

	return state, nil
}
