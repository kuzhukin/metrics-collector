package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var ErrBadDataHash = errors.New("bad data hash")

var secretKey []byte

func InitSignHandlers(key string) {
	secretKey = []byte(key)
}

// SignCheckHandler - middleware for checking request signature
func SignCheckHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedHash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			zlog.Logger.Warnf("decode HashSHA256 header from hex to string, err=%s", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if len(expectedHash) != 0 {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				zlog.Logger.Warnf("Read all from body path=%v err=%s", r.URL.Path, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := checkDataConsistency(data, secretKey, expectedHash); err != nil {
				if errors.Is(err, ErrBadDataHash) {
					zlog.Logger.Warnf("Bad data signature err=%s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(data))
		}

		h.ServeHTTP(w, r)
	})
}

var _ http.ResponseWriter = &signResponseWriter{}

type signResponseWriter struct {
	http.ResponseWriter
}

func newSignResponseWriter(w http.ResponseWriter) *signResponseWriter {
	return &signResponseWriter{
		ResponseWriter: w,
	}
}

func (c *signResponseWriter) Write(b []byte) (int, error) {
	hash, err := calcHash(b, secretKey)
	if err != nil {
		return 0, err
	}

	c.ResponseWriter.Header().Set("HashSHA256", hex.EncodeToString(hash))

	return c.ResponseWriter.Write(b)
}

func SignCreateHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signResponseWriter := newSignResponseWriter(w)

		h.ServeHTTP(signResponseWriter, r)
	})
}

func checkDataConsistency(data []byte, key []byte, expectedHash []byte) error {
	hash, err := calcHash(data, key)
	if err != nil {
		return fmt.Errorf("calc hash, err=%w", err)
	}

	if !reflect.DeepEqual(hash, expectedHash) {
		return fmt.Errorf("calculated=%v, expected=%v, err=%w", hex.EncodeToString(hash), hex.EncodeToString(expectedHash), ErrBadDataHash)
	}

	return nil
}

func calcHash(data []byte, key []byte) ([]byte, error) {
	hasher := hmac.New(sha256.New, key)
	_, err := hasher.Write(data)
	if err != nil {
		return nil, fmt.Errorf("hasher write, err=%w", err)
	}

	return hasher.Sum(nil), nil

}
