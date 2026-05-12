package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type recordingTransport struct {
	dir   string
	inner http.RoundTripper
	seq   atomic.Uint32
}

func newRecordingTransport(dir string, inner http.RoundTripper) *recordingTransport {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &recordingTransport{dir: dir, inner: inner}
}

type exchange struct {
	Seq          uint32              `json:"seq"`
	StartedAt    time.Time           `json:"startedAt"`
	DurationMS   int64               `json:"durationMs"`
	ReqMethod    string              `json:"reqMethod"`
	ReqURL       string              `json:"reqUrl"`
	ReqHeader    map[string][]string `json:"reqHeader"`
	ReqBodyB64   string              `json:"reqBodyB64,omitempty"`
	ReqBodyText  string              `json:"reqBodyText,omitempty"`
	RespStatus   int                 `json:"respStatus,omitempty"`
	RespHeader   map[string][]string `json:"respHeader,omitempty"`
	RespBodyB64  string              `json:"respBodyB64,omitempty"`
	RespBodyText string              `json:"respBodyText,omitempty"`
	Err          string              `json:"err,omitempty"`
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	seq := r.seq.Add(1)
	started := time.Now().UTC()

	reqBody, err := drainBody(req.Body)
	if err != nil {
		return nil, fmt.Errorf("wirecapture: read req body: %w", err)
	}
	if reqBody != nil {
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	ex := exchange{
		Seq:       seq,
		StartedAt: started,
		ReqMethod: req.Method,
		ReqURL:    req.URL.String(),
		ReqHeader: cloneHeader(req.Header),
	}
	setBody(&ex.ReqBodyB64, &ex.ReqBodyText, reqBody, req.Header.Get("Content-Type"))

	resp, rtErr := r.inner.RoundTrip(req)
	ex.DurationMS = time.Since(started).Milliseconds()

	if rtErr != nil {
		ex.Err = rtErr.Error()
		r.write(seq, req.Method, req.URL.Path, &ex)
		return resp, rtErr
	}

	respBody, readErr := drainBody(resp.Body)
	if readErr != nil {
		ex.Err = "wirecapture: read resp body: " + readErr.Error()
	}
	if respBody != nil {
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}
	ex.RespStatus = resp.StatusCode
	ex.RespHeader = cloneHeader(resp.Header)
	setBody(&ex.RespBodyB64, &ex.RespBodyText, respBody, resp.Header.Get("Content-Type"))
	r.write(seq, req.Method, req.URL.Path, &ex)
	return resp, nil
}

func (r *recordingTransport) write(seq uint32, method, path string, ex *exchange) {
	name := fmt.Sprintf("%03d-%s-%s.json", seq, strings.ToLower(method), sanitizePath(path))
	f, err := os.Create(filepath.Join(r.dir, name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "wirecapture: open %s: %v\n", name, err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ex); err != nil {
		fmt.Fprintf(os.Stderr, "wirecapture: encode %s: %v\n", name, err)
	}
}

func drainBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	return io.ReadAll(body)
}

func cloneHeader(h http.Header) map[string][]string {
	out := make(map[string][]string, len(h))
	for k, v := range h {
		c := make([]string, len(v))
		copy(c, v)
		out[k] = c
	}
	return out
}

func setBody(b64Field, textField *string, body []byte, contentType string) {
	if len(body) == 0 {
		return
	}
	if isTextual(contentType, body) {
		*textField = string(body)
		return
	}
	*b64Field = base64Encode(body)
}

func isTextual(contentType string, body []byte) bool {
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "text/"),
		strings.Contains(ct, "json"),
		strings.Contains(ct, "xml"),
		strings.Contains(ct, "x-www-form-urlencoded"),
		strings.Contains(ct, "csv"):
		return true
	}
	for _, b := range body {
		if b == 0 {
			return false
		}
	}
	return true
}

func sanitizePath(p string) string {
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		p = "root"
	}
	out := make([]rune, 0, len(p))
	for _, r := range p {
		switch {
		case r == '/' || r == '\\':
			out = append(out, '_')
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '.':
			out = append(out, r)
		default:
			out = append(out, '_')
		}
	}
	s := string(out)
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}
