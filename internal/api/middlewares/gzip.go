package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/aube/url-shortener/internal/logger"
)

// compressWriter implements http.ResponseWriter and transparently compresses data.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter creates a new compressWriter instance.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map from the underlying ResponseWriter.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write compresses and writes the data to the underlying ResponseWriter.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets the status code and Content-Encoding header if applicable.
func (c *compressWriter) WriteHeader(statusCode int) {
	ct := c.w.Header().Get("Content-Type")
	if strings.Contains(ct, "application/") || strings.Contains(ct, "text/") {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip writer and flushes any buffered data.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader implements io.ReadCloser and transparently decompresses data.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader creates a new compressReader instance.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read decompresses and reads data from the underlying reader.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes both the underlying reader and gzip reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware handles request/response compression using gzip.
// It checks Accept-Encoding and Content-Encoding headers to determine
// if compression/decompression should be applied.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		ow := w

		// Check if client accepts gzip encoding
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// Check if request body is gzipped
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		log.Debug("GzipMiddleware", "sendsGzip", sendsGzip)
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
