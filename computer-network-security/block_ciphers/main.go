package main

import (
	"log"

	"github.com/evgen1067/computer_security/block_ciphers/block"
	"github.com/evgen1067/computer_security/block_ciphers/data"
)

const (
	inpPath = "files/input.txt"
	outPath = "files/output.txt"
)

func main() {
	r := data.NewFilesData(inpPath, outPath)
	input, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	g := block.NewBlockEncoder(input)
	encoded := g.Encrypt()
	g.SetMessage(encoded)
	decoded := g.Decrypt()

	err = r.Write(encoded, decoded)
	if err != nil {
		log.Fatal(err)
	}
}
