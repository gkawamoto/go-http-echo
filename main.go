package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	statusCodeResponse, err := getStatusCodeResponse()
	if err != nil {
		log.Fatal(err)
	}

	useStructuredLogs := getUseStructuredLogs()

	flag.IntVar(&port, "port", port, "port to listen on")
	flag.IntVar(&statusCodeResponse, "status-code-response", statusCodeResponse, "status code to respond with")
	flag.BoolVar(&useStructuredLogs, "structured-logs", useStructuredLogs, "enable structured logs")
	flag.Parse()

	if useStructuredLogs {
		log.SetFlags(log.Lmsgprefix)
		log.SetOutput(&structuredLogger{
			json.NewEncoder(os.Stdout),
		})
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Print(string(data))
		_, _ = w.Write(data)
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		log.Print("shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("listening on port %d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getUseStructuredLogs() bool {
	useStructuredLogs := os.Getenv("USE_STRUCTURED_LOGS")
	if useStructuredLogs == "" {
		return false
	}

	return useStructuredLogs == "true"
}

func getPort() (int, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return 8080, nil
	}
	return strconv.Atoi(port)
}

func getStatusCodeResponse() (int, error) {
	statusCodeResponse := os.Getenv("STATUS_CODE_RESPONSE")
	if statusCodeResponse == "" {
		return 200, nil
	}
	return strconv.Atoi(statusCodeResponse)
}

type structuredLogger struct {
	encoder *json.Encoder
}

func (l *structuredLogger) Write(p []byte) (n int, err error) {
	l.encoder.Encode(map[string]string{
		"message":   string(p),
		"timestamp": time.Now().Format(time.RFC3339Nano),
	})
	return len(p), nil
}
