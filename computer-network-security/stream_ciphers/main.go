package main

import (
	"fmt"

	"github.com/evgen1067/computer_security/stream_ciphers/stream"
)

func main() {
	var (
		a       = 3
		b       = 5
		e       = 32
		message = "Hello, World! It is GoLang!"
		gamma   = "z"
		t       = 2
		p       = 11
		q       = 23
		r       = 1024
	)
	fmt.Println("Синхронный поточный шифр")

	enc := stream.NewEncryptInterface(a, b, e, t, p, q, r, gamma)

	code := enc.SyncEncrypt(message, false)
	decoding := enc.SyncEncrypt(string(code), false)
	printResults(decoding, code)

	fmt.Println("\n\nСинхронный поточный шифр, BBS")

	code = enc.SyncEncrypt(message, true)
	decoding = enc.SyncEncrypt(string(code), true)
	printResults(decoding, code)

	fmt.Println("\n\nАсинхронный поточный шифр")

	code = enc.AsyncEncrypt(message, false, false)
	decoding = enc.AsyncEncrypt(string(code), true, false)
	printResults(decoding, code)

	fmt.Println("\n\nАсинхронный поточный шифр, BBS")

	code = enc.AsyncEncrypt(message, false, true)
	decoding = enc.AsyncEncrypt(string(code), true, true)
	printResults(decoding, code)
}

func printResults(decoding, code []byte) {
	fmt.Printf(
		"\nДлина расшифрованного сообщения: %v\nДлина зашифрованного сообщения: %v",
		len(decoding),
		len(code),
	)
	fmt.Printf(
		"\nРасшифрованное сообщение (в байтах): %v\nЗашифрованное сообщение (в байтах): %v",
		decoding,
		code,
	)
	fmt.Printf(
		"\nРасшифрованное сообщение `%v`\nЗашифрованное сообщение: `%v`",
		string(decoding),
		string(code),
	)
}
