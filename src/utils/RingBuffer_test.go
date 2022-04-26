package utils

import (
	"mjpeg_multiplexer/src/tests"
	"testing"
)

// TestRingBuffer_Peek creates an overflow in data structure and tests correct behavior
func TestRingBuffer_Peek(t *testing.T) {
	size := 20
	buffer := NewRingBuffer[int](size)

	for i := 0; i < size*5; i++ {
		buffer.Push(i)
		tests.Assert(t, i, buffer.Peek(), true)
	}

}

//TestRingBuffer_GetAll tests if the correct amount of array entries in the correct order is returned
func TestRingBuffer_GetAll(t *testing.T) {
	size := 20

	buffer := NewRingBuffer[int](size)
	tests.Assert(t, 0, len(buffer.GetAll()), true)

	// when
	for i := 1; i < size*5; i++ {
		buffer.Push(i)
		// then expect frames stored in framestorage
		if i < size {
			tests.Assert(t, i, len(buffer.GetAll()), true)
		} else {
			tests.Assert(t, size, len(buffer.GetAll()), true)
		}
	}

	// expect GetAll() to return the newest frames first ('newest' = added last)
	data_ := make([]int, size)
	for i := 0; i < size; i++ {
		buffer.Push(i)
		data_[size-i-1] = i
	}

	tests.Assert(t, data_, buffer.GetAll(), true)
}

//TestRingBuffer_GetAll tests if the correct data and the correct size information is returned, even on 'overflow'
func TestRingBuffer_GetData(t *testing.T) {
	size := 20

	buffer := NewRingBuffer[int](size)

	// when
	data := make([]int, size)
	for i := 0; i < size*5; i++ {
		buffer.Push(i)
		data[i%size] = i
		bufferData, bufferSize := buffer.GetData()
		tests.Assert(t, data, *bufferData, true)
		if i < size {
			tests.Assert(t, i+1, bufferSize, true)
		} else {
			tests.Assert(t, size, bufferSize, true)
		}
	}
}
