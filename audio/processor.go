package audio

import (
	"io"
	"time"

	"github.com/thiagomachadox/realtime-waveform/waveform"
)

type Processor struct {
	reader    io.Reader
	chunkSize int
	interval  time.Duration
}

func NewProcessor(reader io.Reader, chunkSize int, interval time.Duration) *Processor {
	return &Processor{
		reader:    reader,
		chunkSize: chunkSize,
		interval:  interval,
	}
}

func (p *Processor) Process(ch chan<- string) {
	buffer := make([]byte, p.chunkSize)
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for range ticker.C {
		n, err := p.reader.Read(buffer)
		if err == io.EOF {
			close(ch)
			return
		}
		if err != nil {
			log.Printf("Error reading audio: %v", err)
			continue
		}

		samples := bytesToFloat32(buffer[:n])
		ascii := waveform.GenerateASCII(samples, 80, 20)
		ch <- ascii
	}
}

func bytesToFloat32(bytes []byte) []float32 {
	floats := make([]float32, len(bytes)/4)
	for i := 0; i < len(floats); i++ {
		floats[i] = float32(int32(bytes[i*4]) | int32(bytes[i*4+1])<<8 | int32(bytes[i*4+2])<<16 | int32(bytes[i*4+3])<<24)
	}
	return floats
}