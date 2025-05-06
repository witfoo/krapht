// Package mock provides mock implementations of the key pipeline data models
package mock

import (
	"errors"

	"github.com/witfoo/krapht/pkg/pipeline"
)

var _ pipeline.Readable = (*ReadableImpl)(nil)

// ReadableImpl is a mock implementation of the Readable interface
type ReadableImpl struct {
	data []byte
}

// NewReadableImpl creates a new ReadableImpl
func NewReadableImpl(data []byte) ReadableImpl {
	return ReadableImpl{
		data: data,
	}
}

// Read returns the data and a nil error
func (m ReadableImpl) Read() ([]byte, error) {
	return m.data, nil
}

// Ensure DataRawReadableImpl implements the DataRawReadable interface
var _ pipeline.DataRawReadable = (*DataRawReadableImpl)(nil)

// DataRawReadableImpl is a mock implementation of the DataRawReadable interface
// It implements the DataReadable and RawReadable interfaces
type DataRawReadableImpl struct {
	data ReadableImpl
	raw  ReadableImpl
}

// NewDataRawReadableImpl creates a new DataRawReadableImpl
// with the given data and raw data
func NewDataRawReadableImpl(data ReadableImpl, raw ReadableImpl) DataRawReadableImpl {
	return DataRawReadableImpl{
		data: data,
		raw:  raw,
	}
}

// Data returns the data that implements the Readable interface
func (m DataRawReadableImpl) Data() pipeline.Readable {
	return m.data
}

// Raw returns the raw data that implements the Readable interface
func (m DataRawReadableImpl) Raw() pipeline.Readable {
	return m.raw
}

// Ensure ReadableBad implements the Readable interface
var _ pipeline.Readable = (*ReadableBad)(nil)

// ReadableBad is a mock implementation of the RawReadable interface.
// It throws an error on Read to simulate a failure in reading data.
type ReadableBad struct {
}

// Read returns nil and an error
func (r ReadableBad) Read() ([]byte, error) {
	return nil, errors.New("error from mock ReadableBad")
}
