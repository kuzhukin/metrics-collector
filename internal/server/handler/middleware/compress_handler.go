package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var _ http.ResponseWriter = &compressResponseWriter{}

type compressResponseWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *compressResponseWriter) Write(b []byte) (int, error) {
	return c.ResponseWriter.Write(b)
}

func (c *compressResponseWriter) WriteHeader(status int) {
	if status < 300 {
		c.Header().Set("Content-Encoding", "gzip")
	}

	c.ResponseWriter.WriteHeader(status)
}

func (c *compressResponseWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func CompressingHTTPHandler(h http.Handler) http.Handler {
	compressingHandler := func(w http.ResponseWriter, r *http.Request) {
		resultingWriter := w

		// setting ungzip response writer
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isSupportGzip := strings.Contains(acceptEncoding, "gzip")
		if isSupportGzip {
			cw := newCompressResponseWriter(w)
			resultingWriter = cw
			defer cw.Close()
		}

		// setting gzip body writer
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
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
