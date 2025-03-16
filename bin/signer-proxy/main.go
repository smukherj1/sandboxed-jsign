package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/smukherj1/sandboxed-jsign/pkg/signrsa"
)

var (
	signKey   = flag.String("sign-key", "", "Signing private key.")
	serverKey = flag.String("key", "", "Server private key.")
	tsCert    = flag.String("ts-cert", "", "Server certificate.")
	kmsCert   = flag.String("kms-cert", "", "Server certificate.")
	tsAddr    = flag.String("ts-addr", ":8080", "Address to serve timestamp requests.")
	kmsAddr   = flag.String("kms-addr", ":8081", "Address to server Cloud KMS requests.")
)

const (
	signRequestPattern = "/v1/projects/{project}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}/cryptoKeyVersions/{keyversion}:{method}"
)

func timestampHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Timestamp request: %v", req.URL.String())
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func runTSServer(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/", timestampHandler)
	log.Printf("Launching timestamp server at %v.", *tsAddr)
	if err := http.ListenAndServeTLS(*tsAddr, *tsCert, *serverKey, nil); err != nil {
		log.Fatalf("Error starting timestamp server: %v", err)
	}
}

type kmsServer struct {
	s *signrsa.Signer
}

type digestRequest struct {
	SHA256  string `json:"sha256,omitempty"`
	payload []byte
}

func (d *digestRequest) validate() error {
	if len(d.SHA256) == 0 {
		return errors.New("missing required field 'sha256'")
	}
	blob, err := base64.StdEncoding.DecodeString(d.SHA256)
	if err != nil {
		return fmt.Errorf("error decoding 'sha256' digest as base64: %w", err)
	}
	if len(blob) != 32 {
		return fmt.Errorf("unexpected length for decoded 'sha256' digest, got %v, want 32", len(blob))
	}
	d.payload = blob
	return nil
}

type signRequestBody struct {
	Digest *digestRequest `json:"digest,omitempty"`
}

func (s *signRequestBody) validate() error {
	if s.Digest == nil {
		return errors.New("missing required field 'digest'")
	}
	if err := s.Digest.validate(); err != nil {
		return fmt.Errorf("error validating field 'digest': %w", err)
	}
	return nil
}

type signRequest struct {
	project    string
	location   string
	keyring    string
	key        string
	keyVersion string
	body       *signRequestBody
}

func (s *signRequest) name() string {
	return fmt.Sprintf("projects/%v/locations/%v/keyRings/%v/cryptoKeys/%v/cryptoKeyVersions/%v",
		s.project, s.location, s.keyring, s.key, s.keyVersion,
	)
}

func signRequestFrom(c *gin.Context) *signRequest {
	p := c.Param("project")
	l := c.Param("location")
	kr := c.Param("keyring")
	k := c.Param("key")
	kv := c.Param("keyversion")
	if len(p) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing {project} in URL, got URL %v, want pattern %v", c.Request.URL.Path, signRequestPattern))
		return nil
	}
	if len(l) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing {location} in URL, got URL %v, want pattern %v", c.Request.URL.Path, signRequestPattern))
		return nil
	}
	if len(kr) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing {keyring} in URL, got URL %v, want pattern %v", c.Request.URL.Path, signRequestPattern))
		return nil
	}
	if len(k) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing {key} in URL, got URL %v, want pattern %v", c.Request.URL.Path, signRequestPattern))
		return nil
	}
	if len(kv) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing {keyversion}:{method} in URL, got URL %v, want pattern %v", c.Request.URL.Path, signRequestPattern))
		return nil
	}
	splitKV := strings.Split(kv, ":")
	if len(splitKV) != 2 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bad {keyversion}:{method} in URL, got element %v in URL %v, want pattern {keyversion}:{method}", kv, c.Request.URL.Path))
		return nil
	}
	kv = splitKV[0]
	method := splitKV[1]
	if method != "asymmetricSign" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unsupported {method} in URL, got method %v in URL %v, want asymmetricSign", method, c.Request.URL.Path))
		return nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("error reading request body: %w", err))
		return nil
	}
	srb := &signRequestBody{}
	if err := json.Unmarshal(body, srb); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("error parsing request body: %w", err))
		return nil
	}
	if err := srb.validate(); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("error validating request body: %w", err))
		return nil
	}
	return &signRequest{
		project:    p,
		location:   l,
		keyring:    kr,
		key:        k,
		keyVersion: kv,
		body:       srb,
	}
}

type signResponse struct {
	Signature string `json:"signature,omitempty"`
}

func (s *kmsServer) sign(c *gin.Context) {
	req := signRequestFrom(c)
	if req == nil {
		return
	}
	sblob, err := s.s.Sign(req.body.Digest.payload)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating signature: %w", err))
		return
	}
	c.JSON(http.StatusOK, signResponse{
		Signature: base64.StdEncoding.EncodeToString(sblob),
	})
}

func runKMSServer(wg *sync.WaitGroup) {
	defer wg.Done()
	sg, err := signrsa.NewSigner(*signKey)
	if err != nil {
		log.Fatalf("Error loading server: %v", err)
	}
	s := kmsServer{s: sg}
	r := gin.Default()
	r.POST("/v1/projects/:project/locations/:location/keyRings/:keyring/cryptoKeys/:key/cryptoKeyVersions/:keyversion", s.sign)
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
	if len(*signKey) == 0 {
		log.Fatalf("--sign-key is required.")
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
	go runTSServer(&wg)
	go runKMSServer(&wg)
	wg.Wait()
}
