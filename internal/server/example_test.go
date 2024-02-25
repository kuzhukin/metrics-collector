package server

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Example() {
	server, err := StartNew()
	noerror(err)
	defer func() {
		err := server.Stop()
		noerror(err)
	}()

	time.Sleep(time.Second * 1)

	// Metric uploading
	respPost, err := http.DefaultClient.Post("http://localhost:8080/update/gauge/metirc/100.1", "text/plain", nil)
	noerror(err)

	defer respPost.Body.Close()
	fmt.Printf("status code = %d\n", respPost.StatusCode)

	// Metric getting
	respGet, err := http.DefaultClient.Get("http://localhost:8080/value/gauge/metirc")
	noerror(err)
	defer respGet.Body.Close()

	data, err := io.ReadAll(respGet.Body)
	noerror(err)

	fmt.Printf("status code = %d\n", respGet.StatusCode)
	fmt.Printf("value = %s\n", string(data))

	time.Sleep(time.Second * 1)
}

func noerror(err error) {
	if err != nil {
		panic(err)
	}
}
