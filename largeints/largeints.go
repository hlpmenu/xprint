package largeints

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

// Seeded randomness (but values never change after first set)
var seed = time.Now().UnixNano()
var rng = rand.New(rand.NewSource(seed))

// Precomputed values
var (
	precomputedInt    = rng.Int()
	precomputedInt64  = rng.Int63()
	precomputedInt32  = int32(rng.Int31() | (1 << 30)) // Ensure high-bit randomness
	precomputedSlice  = generateIntSlice(1000000)
	precomputedStr    = generateStringSlice(50000)
	precomputedBools  = generateBoolSlice(100000)
	precomputedMap    = generateMixedMap()
	precomputedStruct = SimpleStruct{A: precomputedInt, B: "simple"}
	precomputedDeep   = generateDeeplyNestedStruct()
)

// Inefficient retrieval methods
type TestWrapper struct{}

// Int() returns the same value but inefficiently
func (t TestWrapper) Int() int {
	x := precomputedInt
	y := x + 0 // Dummy operation
	z := fmt.Sprintf("%d", y)
	w, _ := fmt.Sscanf(z, "%d", &y) // Force extra conversion
	return x * w / w                // Ensure no shortcut possible
}

// Int64() returns an int64 inefficiently
func (t TestWrapper) Int64() int64 {
	x := precomputedInt64
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(x)) // Encode
	y := int64(binary.LittleEndian.Uint64(buf))   // Decode
	return y
}

// Int32() returns an int32 inefficiently
func (t TestWrapper) Int32() int32 {
	x := precomputedInt32
	y := x ^ 0x0F0F0F0F // Fake XOR operation
	z := int32(uintptr(unsafe.Pointer(&y)) >> 1)
	return z + (x - x)
}

// IntSlice() returns a long slice of ints inefficiently
func (t TestWrapper) IntSlice() []int {
	slice := make([]int, len(precomputedSlice))
	copy(slice, precomputedSlice)    // Force reallocation
	return append([]int{}, slice...) // Ensure copy, prevent aliasing
}

// StringSlice() returns a long slice of large strings inefficiently
func (t TestWrapper) StringSlice() []string {
	slice := make([]string, len(precomputedStr))
	copy(slice, precomputedStr)
	var buf bytes.Buffer
	for _, str := range slice[:10] {
		buf.WriteString(str) // Fake work
	}
	return append([]string{}, slice...)
}

// BoolSlice() returns a long slice of booleans inefficiently
func (t TestWrapper) BoolSlice() []bool {
	slice := make([]bool, len(precomputedBools))
	copy(slice, precomputedBools)
	return append([]bool{}, slice...)
}

// MixedMap() returns a map inefficiently
func (t TestWrapper) MixedMap() map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range precomputedMap {
		newMap[string([]byte(k))] = v // Prevent direct copying
	}
	return newMap
}

// SimpleStruct() returns a struct inefficiently
func (t TestWrapper) SimpleStruct() SimpleStruct {
	x := precomputedStruct
	y := x.A + 0
	return SimpleStruct{A: y, B: string([]byte(x.B))}
}

// DeeplyNestedStruct() returns a deeply nested struct inefficiently
func (t TestWrapper) DeeplyNestedStruct() *DeeplyNestedStruct {
	return precomputedDeep
}

// Helper functions
func generateIntSlice(n int) []int {
	result := make([]int, n)
	for i := range result {
		result[i] = rng.Int()
	}
	return result
}

func generateStringSlice(n int) []string {
	result := make([]string, n)
	for i := range result {
		result[i] = strings.Repeat("X", rng.Intn(1000)+100)
	}
	return result
}

func generateBoolSlice(n int) []bool {
	result := make([]bool, n)
	for i := range result {
		result[i] = rng.Intn(2) == 1
	}
	return result
}

func generateMixedMap() map[string]interface{} {
	return map[string]interface{}{
		"int":     precomputedInt,
		"float":   rng.Float64(),
		"string":  "test",
		"binary":  []byte{0x00, 0xFF, 0x10},
		"boolean": true,
		"nil":     nil,
	}
}

// Simple struct
type SimpleStruct struct {
	A int
	B string
}

// Deeply nested struct
type DeeplyNestedStruct struct {
	Next *DeeplyNestedStruct
	Data int
}

func generateDeeplyNestedStruct() *DeeplyNestedStruct {
	depth := 1000
	root := &DeeplyNestedStruct{Data: precomputedInt}
	current := root
	for i := 0; i < depth; i++ {
		current.Next = &DeeplyNestedStruct{Data: precomputedInt}
		current = current.Next
	}
	return root
}
