package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingResponseWriterRecordsStatusAndBytes(t *testing.T) {
	underlying := httptest.NewRecorder()
	writer := &loggingResponseWriter{ResponseWriter: underlying}

	writer.WriteHeader(http.StatusCreated)
	writer.WriteHeader(http.StatusInternalServerError)
	n, err := writer.Write([]byte("ok"))
	if err != nil {
		t.Fatalf("Write error = %v", err)
	}

	if n != 2 {
		t.Errorf("Write bytes = %d, want 2", n)
	}
	if writer.status() != http.StatusCreated {
		t.Errorf("status = %d, want %d", writer.status(), http.StatusCreated)
	}
	if writer.bytesWritten != 2 {
		t.Errorf("bytesWritten = %d, want 2", writer.bytesWritten)
	}
	if underlying.Code != http.StatusCreated {
		t.Errorf("underlying status = %d, want %d", underlying.Code, http.StatusCreated)
	}
}

func TestLoggingResponseWriterPreservesFlushAndUnwrap(t *testing.T) {
	underlying := &flushTrackingWriter{ResponseRecorder: httptest.NewRecorder()}
	writer := &loggingResponseWriter{ResponseWriter: underlying}

	writer.Flush()

	if !underlying.flushed {
		t.Error("Flush did not reach the underlying response writer")
	}
	if writer.status() != http.StatusOK {
		t.Errorf("status = %d, want %d", writer.status(), http.StatusOK)
	}
	if writer.Unwrap() != underlying {
		t.Error("Unwrap did not return the underlying response writer")
	}
}

func TestRequestLoggingRecordsStructuredFields(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	server := New(":0", logger, nil, nil)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health?api_key=secret", nil)

	server.httpServer.Handler.ServeHTTP(recorder, request)

	requestID := recorder.Header().Get("x-request-id")
	if !strings.HasPrefix(requestID, "req_") {
		t.Fatalf("x-request-id = %q, want req_ prefix", requestID)
	}

	var record struct {
		Message       string `json:"msg"`
		RequestID     string `json:"request_id"`
		Method        string `json:"method"`
		Path          string `json:"path"`
		Status        int    `json:"status"`
		ResponseBytes int    `json:"response_bytes"`
		DurationMS    int64  `json:"duration_ms"`
	}
	if err := json.NewDecoder(&output).Decode(&record); err != nil {
		t.Fatalf("decode request log: %v", err)
	}

	if record.Message != "HTTP request completed" {
		t.Errorf("message = %q, want request completion message", record.Message)
	}
	if record.RequestID != requestID {
		t.Errorf("request_id = %q, want %q", record.RequestID, requestID)
	}
	if record.Method != http.MethodGet || record.Path != "/health" {
		t.Errorf("request = %s %s, want GET /health", record.Method, record.Path)
	}
	if record.Status != http.StatusOK {
		t.Errorf("status = %d, want %d", record.Status, http.StatusOK)
	}
	if record.ResponseBytes != recorder.Body.Len() {
		t.Errorf("response_bytes = %d, want %d", record.ResponseBytes, recorder.Body.Len())
	}
	if record.DurationMS < 0 {
		t.Errorf("duration_ms = %d, want non-negative duration", record.DurationMS)
	}
	if strings.Contains(output.String(), "secret") {
		t.Error("request log contains query parameter value")
	}
}

type flushTrackingWriter struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (w *flushTrackingWriter) Flush() {
	w.flushed = true
}
