package utils

//RingBuffer generic ring bufffer implemention with a fixed capacity
type RingBuffer[T any] struct {
	data        []T
	dataPointer int
	capacity    int
	nElements   int
}

//NewRingBuffer creates a new RingBuffer
func NewRingBuffer[T any](capacity int) RingBuffer[T] {
	return RingBuffer[T]{capacity: capacity, nElements: 0, dataPointer: 0, data: make([]T, capacity)}
}

//Push adds a value into the RingBuffer
func (buffer *RingBuffer[T]) Push(value T) {
	buffer.data[buffer.dataPointer] = value
	buffer.dataPointer = (buffer.dataPointer + 1) % buffer.capacity
	if buffer.nElements < buffer.capacity {
		buffer.nElements++
	}
}

// func (buffer *RingBuffer[T]) Pop() T {
// 	index := (buffer.dataPointer + buffer.capacity - 1) % buffer.capacity
// 	buffer.dataPointer = index
//     buffer.nElements--
// 	return buffer.data[index]
// }

// Peek return the last value added
func (buffer *RingBuffer[T]) Peek() T {
	index := (buffer.dataPointer + buffer.capacity - 1) % buffer.capacity
	return buffer.data[index]
}

// GetAll returns all data inside the buffer, newest data first
func (buffer *RingBuffer[T]) GetAll() []T {
	output := make([]T, buffer.nElements)

	index := (buffer.dataPointer + buffer.capacity - 1) % buffer.capacity
	for i := 0; i < buffer.nElements; i++ {
		output[i] = buffer.data[index]
		index = (index + buffer.capacity - 1) % buffer.capacity
	}

	return output
}

//GetData return the underlying data array and the number of elements
func (buffer *RingBuffer[T]) GetData() (*[]T, int) {
	return &buffer.data, buffer.nElements
}
