package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockMessageBroker is a mock implementation of the MessageBroker interface.
type mockMessageBroker struct {
	publishedMessage []byte
	consumeFunc      func([]byte) error
}

func (m *mockMessageBroker) PublishMessage(message []byte) error {
	m.publishedMessage = message
	return nil
}

func (m *mockMessageBroker) ConsumeMessages(processMessage func([]byte) error) {
	m.consumeFunc = processMessage
}

func TestMessageBroker_PublishMessage(t *testing.T) {
	// Create the mock message broker
	mockBroker := &mockMessageBroker{}

	// Publish a test message
	message := []byte("Test Message")
	err := mockBroker.PublishMessage(message)
	assert.NoError(t, err)

	// Assert that the message was published correctly
	assert.Equal(t, message, mockBroker.publishedMessage)
}

//func TestMessageBroker_ConsumeMessages(t *testing.T) {
//	// Create the mock message broker
//	mockBroker := &mockMessageBroker{}
//
//	// Create a test message
//	message := []byte("Test Message")
//
//	// Start consuming messages
//	doneCh := make(chan struct{})
//	go func() {
//		mockBroker.ConsumeMessages(func(msg []byte) error {
//			// Simulate message processing
//			assert.Equal(t, message, msg)
//			doneCh <- struct{}{}
//			return nil
//		})
//	}()
//
//	// Publish the test message to the mock broker
//	err := mockBroker.PublishMessage(message)
//	assert.NoError(t, err)
//
//	// Wait for the message to be processed
//	<-doneCh
//}

//func TestNewMessageBroker(t *testing.T) {
//	// Create a mock RabbitMQ connection
//	mockConn := &amqp.Connection{}
//
//	// Mock the amqp.Dial function to return the mock connection
//	amqpDialFunc := func(amqpURL string) (*amqp.Connection, error) {
//		return mockConn, nil
//	}
//	defer func() { amqpDialFunc = amqp.Dial }()
//
//	// Mock the amqp.Connection.Channel function to return a mock channel
//	mockChannel := &amqp.Channel{}
//	mockChannelFunc = func() (*amqp.Channel, error) {
//		return mockChannel, nil
//	}
//	defer func() { mockChannelFunc = mockChannelFuncDefault }()
//
//	// Create the message broker
//	broker, err := messaging.NewMessageBroker("amqpURL", "queueName")
//	assert.NotNil(t, broker)
//	assert.NoError(t, err)
//
//	// Check that the connection and channel were set correctly
//	assert.Equal(t, mockConn, broker.Conn)
//	assert.Equal(t, mockChannel, broker.Channel)
//}
