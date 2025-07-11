package service

import (
	"encoding/json"
	"testing"
	"time"

	"food-order-backend/internal/app/order/model"
)

// TestAggregateConsistency kiểm tra tính nhất quán của aggregate khi áp dụng events
func TestAggregateConsistency(t *testing.T) {
	tests := []struct {
		name           string
		events         []model.EventStore
		expectedState  model.Order
		shouldBeValid  bool
		description    string
	}{
		{
			name: "Valid order creation maintains consistent state",
			events: []model.EventStore{
				{
					AggregateID: "ORD_001",
					EventType:   "OrderCreated",
					EventData:   `{"user_id":123,"restaurant_id":456,"status":"PENDING","note":"Test order"}`,
					CreatedAt:   time.Now(),
				},
			},
			expectedState: model.Order{
				OrderID:      "ORD_001",
				UserID:       123,
				RestaurantID: 456,
				Status:       "PENDING",
				Note:         "Test order",
			},
			shouldBeValid: true,
			description:   "Aggregate should maintain valid state after OrderCreated event",
		},
		{
			name: "Valid state progression through order lifecycle",
			events: []model.EventStore{
				{
					AggregateID: "ORD_002",
					EventType:   "OrderCreated",
					EventData:   `{"user_id":123,"restaurant_id":456,"status":"PENDING"}`,
					CreatedAt:   time.Now(),
				},
				{
					AggregateID: "ORD_002",
					EventType:   "ShipperAssigned",
					EventData:   `{"shipper_id":"S001"}`,
					CreatedAt:   time.Now().Add(5 * time.Minute),
				},
				{
					AggregateID: "ORD_002",
					EventType:   "RestaurantAccepted",
					EventData:   `{"merchant_id":"M001"}`,
					CreatedAt:   time.Now().Add(10 * time.Minute),
				},
				{
					AggregateID: "ORD_002",
					EventType:   "CookingStarted",
					EventData:   `{"estimated_time":"20 minutes"}`,
					CreatedAt:   time.Now().Add(15 * time.Minute),
				},
				{
					AggregateID: "ORD_002",
					EventType:   "OrderPicked",
					EventData:   `{"pickup_time":"11:15"}`,
					CreatedAt:   time.Now().Add(45 * time.Minute),
				},
				{
					AggregateID: "ORD_002",
					EventType:   "OrderDelivered",
					EventData:   `{"delivery_time":"11:45"}`,
					CreatedAt:   time.Now().Add(75 * time.Minute),
				},
			},
			expectedState: model.Order{
				OrderID:      "ORD_002",
				UserID:       123,
				RestaurantID: 456,
				Status:       "DELIVERED",
			},
			shouldBeValid: true,
			description:   "Aggregate should maintain consistent state through complete lifecycle",
		},
		{
			name: "Cancellation maintains state consistency",
			events: []model.EventStore{
				{
					AggregateID: "ORD_003",
					EventType:   "OrderCreated",
					EventData:   `{"user_id":123,"restaurant_id":456,"status":"PENDING"}`,
					CreatedAt:   time.Now(),
				},
				{
					AggregateID: "ORD_003",
					EventType:   "ShipperAssigned",
					EventData:   `{"shipper_id":"S001"}`,
					CreatedAt:   time.Now().Add(5 * time.Minute),
				},
				{
					AggregateID: "ORD_003",
					EventType:   "RestaurantAccepted",
					EventData:   `{"merchant_id":"M001"}`,
					CreatedAt:   time.Now().Add(10 * time.Minute),
				},
				{
					AggregateID: "ORD_003",
					EventType:   "OrderCanceled",
					EventData:   `{"reason":"Customer request","canceled_by":"customer"}`,
					CreatedAt:   time.Now().Add(12 * time.Minute),
				},
			},
			expectedState: model.Order{
				OrderID:      "ORD_003",
				UserID:       123,
				RestaurantID: 456,
				Status:       "CANCELED",
			},
			shouldBeValid: true,
			description:   "Aggregate should maintain consistent state after cancellation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply events to aggregate and check consistency
			finalState, err := applyEventsToAggregate(tt.events)
			
			if !tt.shouldBeValid {
				if err == nil {
					t.Errorf("Expected error but got none for test: %s", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for test %s: %v", tt.name, err)
				return
			}

			// Verify aggregate consistency
			if !isAggregateStateValid(finalState) {
				t.Errorf("Aggregate state is invalid for test: %s", tt.name)
			}

			// Verify expected state
			if finalState.OrderID != tt.expectedState.OrderID {
				t.Errorf("OrderID mismatch for %s: expected %s, got %s", 
					tt.name, tt.expectedState.OrderID, finalState.OrderID)
			}

			if finalState.UserID != tt.expectedState.UserID {
				t.Errorf("UserID mismatch for %s: expected %d, got %d", 
					tt.name, tt.expectedState.UserID, finalState.UserID)
			}

			if finalState.RestaurantID != tt.expectedState.RestaurantID {
				t.Errorf("RestaurantID mismatch for %s: expected %d, got %d", 
					tt.name, tt.expectedState.RestaurantID, finalState.RestaurantID)
			}

			if finalState.Status != tt.expectedState.Status {
				t.Errorf("Status mismatch for %s: expected %s, got %s", 
					tt.name, tt.expectedState.Status, finalState.Status)
			}

			if tt.expectedState.Note != "" && finalState.Note != tt.expectedState.Note {
				t.Errorf("Note mismatch for %s: expected %s, got %s", 
					tt.name, tt.expectedState.Note, finalState.Note)
			}

			t.Logf("✅ %s: %s", tt.name, tt.description)
		})
	}
}

// TestAggregateInvariants kiểm tra các invariant của aggregate được bảo toàn
func TestAggregateInvariants(t *testing.T) {
	tests := []struct {
		name        string
		events      []model.EventStore
		expectError bool
		description string
	}{
		{
			name: "Aggregate must have OrderID after OrderCreated",
			events: []model.EventStore{
				{
					AggregateID: "ORD_INV_001",
					EventType:   "OrderCreated",
					EventData:   `{"user_id":123,"restaurant_id":456,"status":"PENDING"}`,
					CreatedAt:   time.Now(),
				},
			},
			expectError: false,
			description: "OrderID invariant should be maintained",
		},
		{
			name: "Aggregate must have UserID after OrderCreated",
			events: []model.EventStore{
				{
					AggregateID: "ORD_INV_002",
					EventType:   "OrderCreated",
					EventData:   `{"restaurant_id":456,"status":"PENDING"}`, // Missing user_id
					CreatedAt:   time.Now(),
				},
			},
			expectError: true,
			description: "Missing UserID should violate invariant",
		},
		{
			name: "Aggregate must have RestaurantID after OrderCreated",
			events: []model.EventStore{
				{
					AggregateID: "ORD_INV_003",
					EventType:   "OrderCreated",
					EventData:   `{"user_id":123,"status":"PENDING"}`, // Missing restaurant_id
					CreatedAt:   time.Now(),
				},
			},
			expectError: true,
			description: "Missing RestaurantID should violate invariant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finalState, err := applyEventsToAggregate(tt.events)
			
			// Check if error occurred as expected
			if tt.expectError {
				if err == nil && !hasInvariantViolation(finalState) {
					t.Errorf("Expected invariant violation for test: %s", tt.name)
				}
				t.Logf("✅ %s: %s", tt.name, tt.description)
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for test %s: %v", tt.name, err)
				return
			}

			// Verify all invariants are satisfied
			if hasInvariantViolation(finalState) {
				t.Errorf("Invariant violation detected for test: %s", tt.name)
				return
			}

			t.Logf("✅ %s: %s", tt.name, tt.description)
		})
	}
}

// applyEventsToAggregate applies events to aggregate in memory (without database)
func applyEventsToAggregate(events []model.EventStore) (model.Order, error) {
	if len(events) == 0 {
		return model.Order{}, nil
	}

	// Initialize aggregate
	state := model.Order{
		OrderID: events[0].AggregateID,
	}

	// Sort events by timestamp (ensuring correct order)
	for i := 0; i < len(events)-1; i++ {
		for j := i + 1; j < len(events); j++ {
			if events[i].CreatedAt.After(events[j].CreatedAt) {
				events[i], events[j] = events[j], events[i]
			}
		}
	}

	// Apply events sequentially
	for _, evt := range events {
		var data model.OrderEventData
		if err := json.Unmarshal([]byte(evt.EventData), &data); err != nil {
			return model.Order{}, err
		}

		// Apply event to state (same logic as ReplayOrderState)
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
	}

	return state, nil
}

// isAggregateStateValid checks if aggregate state is valid
func isAggregateStateValid(state model.Order) bool {
	// Basic invariants that must always be true
	if state.OrderID == "" {
		return false
	}
	
	// If status is set, it should be one of valid statuses
	if state.Status != "" {
		validStatuses := map[string]bool{
			"PENDING":             true,
			"RESTAURANT_ACCEPTED": true,
			"SHIPPER_ASSIGNED":    true,
			"SHIPPER_CONFIRMED":   true,
			"COOKING":             true,
			"PICKED":              true,
			"DELIVERED":           true,
			"CANCELED":            true,
		}
		if !validStatuses[state.Status] {
			return false
		}
	}

	return true
}

// hasInvariantViolation checks for business invariant violations
func hasInvariantViolation(state model.Order) bool {
	// After OrderCreated event, these fields must be set
	if state.Status != "" { // If any event has been applied
		if state.UserID == 0 {
			return true // UserID is required
		}
		if state.RestaurantID == 0 {
			return true // RestaurantID is required
		}
	}
	
	return false
}
