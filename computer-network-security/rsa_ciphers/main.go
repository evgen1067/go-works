package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/evgen1067/computer_security/rsa_ciphers/file"
	"github.com/evgen1067/computer_security/rsa_ciphers/rsa"
	"github.com/evgen1067/computer_security/rsa_ciphers/sha256"
)

const (
	inpPath = "file/input.txt"
	outPath = "file/output.txt"
	bits    = 512
)

func main() {
	var result string
	f := file.NewFilesData(inpPath, outPath)
	input, err := f.Read()
	if err != nil {
		log.Fatal(err)
	}
	sha := sha256.NewSHA256()
	hash := sha.Write(input)
	result += fmt.Sprintf("Хэш: %v\n", hash)
	checkSum := fromUint32(hash)
	result += fmt.Sprintf("Контрольная сумма: %v\n", checkSum)
	p, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		log.Fatal(err)
	}
	q, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		log.Fatal(err)
	}
	keys, err := rsa.GenerateKeyPair(p, q, 65537)
	if err != nil {
		log.Fatal(err)
	}
	result += fmt.Sprintf(
		"Ключи:\nОткрытый ключ:\ne:%v\nn:%v\nЗакрытый ключ:\nd:%v\n",
		keys.PubK.E,
		keys.PubK.N,
		keys.PrivK.D,
	)
	encripted := make([]*big.Int, len(hash))
	decripted := make([]*big.Int, len(hash))
	for i, v := range hash {
		encripted[i] = rsa.Encrypt(big.NewInt(int64(v)), keys.PubK)
		decripted[i] = rsa.Decrypt(encripted[i], keys.PrivK)
	}
	result += fmt.Sprintf("Цифровая подпись до шифра: %v\n", hash)
	result += fmt.Sprintf("Зашифрованная цифровая подпись: %v\n", encripted)
	result += fmt.Sprintf("Расшифрованная цифровая подпись: %v\n", decripted)
	f.Write([]byte(result))
}

func fromUint32(hash []uint32) []byte {
	b := make([]byte, 32)
	var j int
	for _, val := range hash {
		b[j] = byte((val >> 0) & 255)
		b[j+1] = byte((val >> 8) & 255)
		b[j+2] = byte((val >> 16) & 255)
		b[j+3] = byte((val >> 24) & 255)
		j += 4
	}
	return b
}
