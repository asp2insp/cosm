package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)
import (
	"encoding/base64"
	"encoding/pem"
)

var trustedKey *rsa.PublicKey

func initKey(path string) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(keyBytes)
	fmt.Printf("block.Type: %s\n", block.Type)
	pubkeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	pubkey, ok := pubkeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal("Fatal error")
	}
	trustedKey = pubkey
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	b64data := r.FormValue("binary")
	if b64data == "" {
		http.Error(w, "Missing Base-64 Encoded Binary", http.StatusBadRequest)
		return
	}
	signature := r.FormValue("signature")
	if signature == "" {
		http.Error(w, "Missing Signature", http.StatusBadRequest)
		return
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	bindata, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rsa.VerifyPSS(trustedKey, crypto.MD5, bindata, sig, nil)
	if err != nil {
		http.Error(w, "Not Trusted", http.StatusForbidden)
		return
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "cosm-runner"
	app.Usage = "run/supervise a (micro)cosm"
	app.Action = func(c *cli.Context) {
		println("Creating a universe...")
		initKey(c.Args().First())
		http.HandleFunc("/", saveHandler)
		http.ListenAndServe(":1205", nil)
	}

	app.Run(os.Args)
}
