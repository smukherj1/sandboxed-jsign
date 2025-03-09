package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	serverKey  = flag.String("key", "", "Server private key.")
	serverCert = flag.String("cert", "", "Server certificate.")
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {
	flag.Parse()
	if len(*serverKey) == 0 {
		log.Fatalf("--key is required.")
	}
	if len(*serverCert) == 0 {
		log.Fatalf("--cert is required.")
	}
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8080", *serverCert, *serverKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
