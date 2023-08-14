package rsa

import (
	"math/big"
)

// type RSA struct {
// 	n         uint32
// 	phi       uint32
// 	publicKey uint32
// 	secretKey uint32
// }

// func (r *RSA) SetPublicKey(pk uint32) error {
// 	if !(1 < pk && pk <= r.phi) {
// 		return ErrPublicKeyEuler
// 	}
// 	if !NOD(pk, r.phi) {
// 		return ErrPublicKeyNod
// 	}
// 	r.publicKey = pk
// 	r.SetSecretKey()
// 	return nil
// }

// func (r *RSA) SetSecretKey() {
// 	var (
// 		u1 uint32 = 0
// 		u2 uint32 = 1
// 		u3 uint32 = r.phi
// 		v1 uint32 = 1
// 		v2 uint32 = 0
// 		v3 uint32 = r.publicKey
// 		q  uint32
// 	)

// 	for {
// 		if u3 == 1 {
// 			break
// 		}
// 		q = u3 / v3
// 		t1, t2, t3 := u1-v1*q, u2-v2*q, u3-v3*q
// 		u1, u2, u3 = v1, v2, v3
// 		v1, v2, v3 = t1, t2, t3
// 	}
// 	r.secretKey = u1
// }

// func (r *RSA) Encrypt(message []byte) []uint32 {
// 	chars := []byte(message)
// 	encryptMsg := make([]uint32, len(chars))
// 	for i, val := range chars {
// 		encryptMsg[i] = pow(uint32(val), r.publicKey) % r.n
// 	}
// 	return encryptMsg
// }

// func (r *RSA) Decrypt(message []uint32) []byte {
// 	decryptMsg := make([]byte, len(message))
// 	for i, val := range message {
// 		decryptMsg[i] = byte(pow(uint32(val), r.secretKey) % r.n)
// 	}
// 	return decryptMsg
// }

// func pow(x, y uint32) uint32 {
// 	if y == 0 {
// 		return 1
// 	}
// 	tmp := pow(x, y/2)
// 	if y%2 == 0 { // обработка чётного y
// 		return tmp * tmp
// 	} else { // обработка нечётного y
// 		return tmp * tmp * x
// 	}
// }

// func NewRSA(p, q uint32) *RSA {
// 	return &RSA{n: p * q, phi: (p - 1) * (q - 1)}
// }

// var (
// 	ErrPublicKeyEuler = errors.New("Значение открытого ключа должно быть от 1 до значения функции Эйлера")
// 	ErrPublicKeyNod   = errors.New("Значение НОД (Открытый ключ, значение ф-и Эйлера) должно быть равно единице")
// )

// func NOD(a, b uint32) bool {
// 	for a != 0 && b != 0 {
// 		if a > b {
// 			a = a % b
// 		} else {
// 			b = b % a
// 		}
// 	}
// 	return (a + b) == 1
// }

const (
	bits = 512 // 2048
)

var bigOne = big.NewInt(int64(1))

type PublicKey struct {
	E *big.Int
	N *big.Int
}

type PrivateKey struct {
	D *big.Int
	N *big.Int
}

type Key struct {
	PubK  PublicKey
	PrivK PrivateKey
}

func Encrypt(m *big.Int, pubK PublicKey) *big.Int {
	c := new(big.Int).Exp(m, pubK.E, pubK.N)
	return c
}

// Decrypt deencrypts a ciphertext c with given PrivateKey
func Decrypt(c *big.Int, privK PrivateKey) *big.Int {
	m := new(big.Int).Exp(c, privK.D, privK.N)
	return m
}

func GenerateKeyPair(p, q *big.Int, e int) (key Key, err error) {
	n := new(big.Int).Mul(p, q)       // n := p * q
	p1 := new(big.Int).Sub(p, bigOne) // p1 := p - 1
	q1 := new(big.Int).Sub(q, bigOne) // q1 := q - 1
	phi := new(big.Int).Mul(p1, q1)   // phi := p1 * q1
	var pubK PublicKey
	pubK.E = big.NewInt(int64(e))
	pubK.N = n

	d := new(big.Int).ModInverse(big.NewInt(int64(e)), phi)

	var privK PrivateKey
	privK.D = d
	privK.N = n

	key.PubK = pubK
	key.PrivK = privK
	return key, nil
}
