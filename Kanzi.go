/*
Copyright 2011-2017 Frederic Langlet
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
you may obtain a copy of the License at

                http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kanzi

const (
	ERR_MISSING_PARAM       = 1
	ERR_BLOCK_SIZE          = 2
	ERR_INVALID_CODEC       = 3
	ERR_CREATE_COMPRESSOR   = 4
	ERR_CREATE_DECOMPRESSOR = 5
	ERR_OUTPUT_IS_DIR       = 6
	ERR_OVERWRITE_FILE      = 7
	ERR_CREATE_FILE         = 8
	ERR_CREATE_BITSTREAM    = 9
	ERR_OPEN_FILE           = 10
	ERR_READ_FILE           = 11
	ERR_WRITE_FILE          = 12
	ERR_PROCESS_BLOCK       = 13
	ERR_CREATE_CODEC        = 14
	ERR_INVALID_FILE        = 15
	ERR_STREAM_VERSION      = 16
	ERR_CREATE_STREAM       = 17
	ERR_INVALID_PARAM       = 18
	ERR_CRC_CHECK           = 19
	ERR_UNKNOWN             = 127
)

// An integer function is an operation that takes an array of integers as input and
// and turns it into another array of integers. The size of the returned array
// is not known in advance (by the caller).
// Return index in src, index in dst and error
type IntTransform interface {
	Forward(src, dst []int) (uint, uint, error)

	Inverse(src, dst []int) (uint, uint, error)
}

// A byte function is an operation that takes an array of bytes as input and
// turns it into another array of bytes. The size of the returned array is not
// known in advance (by the caller).
// Return index in src, index in dst and error
type ByteTransform interface {
	Forward(src, dst []byte) (uint, uint, error)

	Inverse(src, dst []byte) (uint, uint, error)
}

// An integer function is an operation that transforms the input int array and writes
// the result in the output int array. The result may have a different size.
// The function may fail if input and output array are the same array.
// The index of input and output arrays are updated appropriately.
// Return index in src, index in dst and error
type IntFunction interface {
	Forward(src, dst []int) (uint, uint, error)

	Inverse(src, dst []int) (uint, uint, error)

	// Return the max size required for the encoding output buffer
	// If the max size of the output buffer is not known, return -1
	MaxEncodedLen(srcLen int) int
}

// A byte function is an operation that transforms the input byte array and writes
// the result in the output byte array. The result may have a different size.
// The function may fail if input and output array are the same array.
// Return index in src, index in dst and error
type ByteFunction interface {
	Forward(src, dst []byte) (uint, uint, error)

	Inverse(src, dst []byte) (uint, uint, error)

	// Return the max size required for the encoding output buffer
	MaxEncodedLen(srcLen int) int
}

type InputBitStream interface {
	// Return 1 or 0
	// Panic if error
	ReadBit() int

	// Length is the number of bits in [1..64]. Return the bits read as an uint64
	// Panic if error
	ReadBits(length uint) uint64

	// Read bits and put them in the byte array. Length is the number of bits
	// Return the number of bits read.
	// Panic if error
	ReadArray(bits []byte, length uint) uint

	Close() (bool, error)

	// Number of bits read
	Read() uint64

	// Return false when the bitstream is closed or the End-Of-Stream has been reached
	HasMoreToRead() (bool, error)
}

type OutputBitStream interface {
	// Write the least significant bit of the input integer
	// Panic if error
	WriteBit(bit int)

	// Length is the number of bits in [1..64]. Return the number of bits written.
	WriteBits(bits uint64, length uint) uint

	// Write bits ouf of the byte array. Length is the number of bits.
	// Return the number of bits written.
	// Panic if error
	WriteArray(bits []byte, length uint) uint

	Close() (bool, error)

	// Number of bits written
	Written() uint64
}

type Predictor interface {
	// Update the probability model
	Update(bit byte)

	// Return the split value representing the probability of 1 in the [0..4095] range.
	// E.G. 410 represents roughly a probability of 10% for 1
	Get() int
}

type EntropyEncoder interface {
	// Encode the array provided into the bitstream. Return the number of byte
	// written to the bitstream
	Write(block []byte) (int, error)

	// Return the underlying bitstream
	BitStream() OutputBitStream

	// Must be called before getting rid of the entropy encoder
	Dispose()
}

type EntropyDecoder interface {
	// Decode the next chunk of data from the bitstream and return in the
	// provided buffer.
	Read(block []byte) (int, error)

	// Return the underlying bitstream
	BitStream() InputBitStream

	// Must be called before getting rid of the entropy decoder
	// Trying to encode after a call to dispose gives undefined behavior
	Dispose()
}
