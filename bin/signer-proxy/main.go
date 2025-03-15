package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	serverKey = flag.String("key", "", "Server private key.")
	tsCert    = flag.String("ts-cert", "", "Server certificate.")
	kmsCert   = flag.String("kms-cert", "", "Server certificate.")
	tsAddr    = flag.String("ts-addr", ":8080", "Address to serve timestamp requests.")
	kmsAddr   = flag.String("kms-addr", ":8081", "Address to server Cloud KMS requests.")
)

func timestampHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Timestamp request: %v", req.URL.String())
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func tsServer(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/", timestampHandler)
	log.Printf("Launching timestamp server at %v.", *tsAddr)
	if err := http.ListenAndServeTLS(*tsAddr, *tsCert, *serverKey, nil); err != nil {
		log.Fatalf("Error starting timestamp server: %v", err)
	}
}

func kmsHandler(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("KMS request %v %v: %v", c.Request.Method, c.Request.URL.String(), string(body))
	c.String(http.StatusNotImplemented, "KMS handler is not implemented")
}

func kmsServer(wg *sync.WaitGroup) {
	defer wg.Done()
	r := gin.Default()
	r.GET("/", kmsHandler)
	r.POST("/", kmsHandler)
	log.Printf("Launching Cloud KMS server at %v.", *kmsAddr)
	if err := r.RunTLS(*kmsAddr, *kmsCert, *serverKey); err != nil {
		log.Fatalf("Error starting Cloud KMS server: %v", err)
	}
}

func main() {
	flag.Parse()
	if len(*serverKey) == 0 {
		log.Fatalf("--key is required.")
	}
	if len(*kmsCert) == 0 {
		log.Fatalf("--kms-cert is required.")
	}
	if len(*tsCert) == 0 {
		log.Fatalf("--ts-cert is required.")
	}
	if len(*tsAddr) == 0 {
		log.Fatalf("--ts-addr is required.")
	}
	if len(*kmsAddr) == 0 {
		log.Fatalf("--kms-addr is required.")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go tsServer(&wg)
	go kmsServer(&wg)
	wg.Wait()
}
