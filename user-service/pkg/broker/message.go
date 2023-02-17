package broker

// Message is a collection of elements representing a message of the broker.
type Message struct {
	// Topic the broker topic for this message.
	Topic string

	// Value the actual message to use in broker.
	Value []byte
}
