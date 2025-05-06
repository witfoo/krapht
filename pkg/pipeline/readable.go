package pipeline

// Readable is an interface that represents a source of data.
type Readable interface {
	// Read returns the data and an error if any.
	Read() ([]byte, error)
}

// DataReadable represents a message that contains the structured data
type DataReadable interface {
	// Data returns a Readable interface for accessing the structured data
	Data() Readable
}

// RawReadable represents a message that contains only the raw message
type RawReadable interface {
	// Raw returns a Readable interface for accessing the raw message
	Raw() Readable
}

// DataRawReadable represents a message that contains both the structured data and the raw message
type DataRawReadable interface {
	DataReadable
	RawReadable
}
