package Utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	ReturnData map[string]interface{}
	Process    float32
	Error      error
}

// Stream helps transmit each streams withing a channel.
type Stream struct {
	stream chan Entry
}

// NewJSONStream returns a new `Stream` type.
func NewJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// User object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

// Start starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
func (s Stream) Start(path string) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	// Open file to read.
	file, err := os.Open(path)
	fi, _ := file.Stat()
	totalsize := fi.Size()
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		returnData := make(map[string]interface{})
		offset := decoder.InputOffset()
		if err := decoder.Decode(&returnData); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode line %d: %w", i, err)}
			return
		}
		processInt := float32(offset) / float32(totalsize*2)
		process, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (processInt*100)), 64)
		s.stream <- Entry{ReturnData: returnData, Process: float32(process)}
		i++
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}
