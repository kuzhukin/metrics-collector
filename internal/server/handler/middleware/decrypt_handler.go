package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/crypto"
)

func NewDecryptHTTPHandler(keyPath string) (func(http.Handler) http.Handler, error) {
	decryptor, err := crypto.NewDecryptor(keyPath)
	if err != nil {
		return nil, fmt.Errorf("new decryptor, keyPath=%v, err=%w", keyPath, err)
	}

	decryptHandler := func(h http.Handler) http.Handler {
		compressingHandler := func(w http.ResponseWriter, r *http.Request) {
			cryptedBody, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			decryptedBody, err := decryptor.Decrypt(cryptedBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decryptedBody))

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(compressingHandler)
	}

	return decryptHandler, nil
}
