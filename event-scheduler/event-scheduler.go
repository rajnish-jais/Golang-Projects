package main

import (
	"fmt"
	"sync"
)

// Event represents a generic event type
type Event struct {
	Type string
	Data interface{}
}

// EventBroker represents the event broker component
type EventBroker struct {
	subscribers map[string][]chan Event
	mu          sync.Mutex
}

// NewEventBroker creates a new EventBroker instance
func NewEventBroker() *EventBroker {
	return &EventBroker{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe adds a new subscriber to the event broker
func (eb *EventBroker) Subscribe(eventType string) chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 10) // Buffer channel to avoid blocking publishers
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)

	return ch
}

// Publish sends an event to the event broker
func (eb *EventBroker) Publish(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := eb.subscribers[event.Type]
	for _, ch := range subscribers {
		go func(ch chan Event) {
			ch <- event
		}(ch)
	}
}

// EventDispatcher represents the event dispatcher component
type EventDispatcher struct {
	eventBroker *EventBroker
}

// NewEventDispatcher creates a new EventDispatcher instance
func NewEventDispatcher(eventBroker *EventBroker) *EventDispatcher {
	return &EventDispatcher{
		eventBroker: eventBroker,
	}
}

// Start starts the event dispatcher to handle incoming events
func (ed *EventDispatcher) Start() {
	for eventType, subscribers := range ed.eventBroker.subscribers {
		go func(eventType string, subscribers []chan Event) {
			for {
				select {
				case event := <-subscribers[0]:
					// Process the event based on the subscriber's requirements or business logic
					fmt.Printf("Received event of type '%s': %+v\n", eventType, event)
				}
			}
		}(eventType, subscribers)
	}
}

func main() {
	eventBroker := NewEventBroker()
	eventDispatcher := NewEventDispatcher(eventBroker)

	// Start the event dispatcher in a separate goroutine
	go eventDispatcher.Start()

	// Create subscribers for different event types
	subscriber1 := eventBroker.Subscribe("type1")
	subscriber2 := eventBroker.Subscribe("type2")

	// Publish events
	eventBroker.Publish(Event{Type: "type1", Data: "Event 1"})
	eventBroker.Publish(Event{Type: "type2", Data: "Event 2"})
	eventBroker.Publish(Event{Type: "type1", Data: "Event 3"})

	// Wait for events to be processed
	fmt.Println(<-subscriber1)
	fmt.Println(<-subscriber2)
	fmt.Println(<-subscriber1)
}
