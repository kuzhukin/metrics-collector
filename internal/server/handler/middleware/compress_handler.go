package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

var _ http.ResponseWriter = &compressResponseWriter{}

type compressResponseWriter struct {
	wr http.ResponseWriter
	zw *gzip.Writer
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		wr: w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressResponseWriter) Header() http.Header {
	return c.wr.Header()
}

func (c *compressResponseWriter) Write(b []byte) (int, error) {
	return c.zw.Write(b)
}

func (c *compressResponseWriter) WriteHeader(status int) {
	if status < 300 {
		c.wr.Header().Set("Content-Encoding", "gzip")
	}

	c.wr.WriteHeader(status)
}

func (c *compressResponseWriter) Close() error {
	return c.zw.Close()
}

type decompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newDecompressReader(r io.ReadCloser) (*decompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &decompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c decompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *decompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func CompressingHTTPHandler(h http.Handler) http.Handler {
	compressingHandler := func(w http.ResponseWriter, r *http.Request) {
		resultingWriter := w

		// setting gzip response writer
		if contains(r.Header.Values("Accept-Encoding"), "gzip") {
			cw := newCompressResponseWriter(w)
			resultingWriter = cw
			defer cw.Close()
		}

		// setting ungzip body reader
		if contains(r.Header.Values("Content-Encoding"), "gzip") {
			cr, err := newDecompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(resultingWriter, r)
	}

	return http.HandlerFunc(compressingHandler)
}

func contains[T comparable](slice []T, el T) bool {
	for _, current := range slice {
		if current == el {
			return true
		}
	}

	return false
}
