package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
)

func run(c *cli.Context) {
	keyBytes, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		log.Fatal(err)
		return
	}
	block, _ := pem.Decode(keyBytes)
	fmt.Printf("block.Type: %s\n", block.Type)
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	binBytes, err := ioutil.ReadFile(c.Args().Get(1))

	sig, err := privKey.Sign(rand.Reader, binBytes, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	sig64 := base64.StdEncoding.EncodeToString(sig)
	bin64 := base64.StdEncoding.EncodeToString(binBytes)
	data := url.Values{}
	data.Set("signature", sig64)
	data.Add("binary", bin64)
	http.PostForm("localhost:1205", data)
}

func main() {
	app := cli.NewApp()
	app.Name = "cosm"
	app.Usage = "manage your own little (micro)cosm"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "send a binary to the server: run <sig path> <binary path>",
			Action: run,
		},
	}

	app.Run(os.Args)
}
