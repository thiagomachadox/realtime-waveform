package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thiagomachadox/realtime-waveform/audio"
	"github.com/thiagomachadox/realtime-waveform/downloader"
)

type Server struct {
	mux      *http.ServeMux
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("GET /ws", s.handleWebSocket)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket connection attempt")
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connection established")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		youtubeURL := string(message)
		log.Printf("Received YouTube URL: %s", youtubeURL)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result of the audio processing
		resultCh := make(chan struct{})

		go func() {
			audioStream, err := downloader.DownloadAudio(youtubeURL)
			if err != nil {
				log.Printf("Error downloading audio: %v", err)
				errorMsg := fmt.Sprintf("Error: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				resultCh <- struct{}{}
				return
			}
			defer audioStream.Close()

			processor := audio.NewProcessor(audioStream, 4096, 100*time.Millisecond)
			ch := make(chan string)
			go processor.Process(ch)

			for ascii := range ch {
				log.Println("Sending waveform data")
				err := conn.WriteMessage(websocket.TextMessage, []byte(ascii))
				if err != nil {
					log.Printf("Error sending message: %v", err)
					break
				}
			}

			resultCh <- struct{}{}
		}()

		// Wait for either the processing to complete or the context to timeout
		select {
		case <-resultCh:
			// Processing completed
		case <-ctx.Done():
			log.Println("Processing timed out")
			conn.WriteMessage(websocket.TextMessage, []byte("Error: Processing timed out"))
		}
	}
}
