// Given a private key, creates a random generator and public commitment to the
// key for a fresh epoch of DLEQ proofs.
//
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alxdavids/btd/crypto"
)

func main() {
	var keyFile, outFile string
	var defaultFilename = fmt.Sprintf("dleq_proof_%s", time.Now().Format(time.RFC3339))
	flag.StringVar(&keyFile, "key", "", "path to a PEM-encoded EC PRIVATE KEY")
	flag.StringVar(&outFile, "out", defaultFilename, "output path for the commitment")
	flag.Parse()

	if keyFile == "" {
		flag.Usage()
		return
	}

	curve, key, err := crypto.ParseKeyFile(keyFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, G, err := crypto.NewRandomPoint(curve)
	if err != nil {
		fmt.Println(err)
		return
	}
	Hx, Hy := curve.ScalarMult(G.X, G.Y, key)
	H, err := crypto.NewPoint(curve, Hx, Hy)
	if err != nil {
		fmt.Println(err)
		return
	}

	C := &crypto.Commitment{
		G: G,
		H: H,
	}
	cBytes, err := json.Marshal(C)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(outFile, cBytes, os.FileMode(0644))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(outFile)
}
