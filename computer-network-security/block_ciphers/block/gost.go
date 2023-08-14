package block

type Sequence []uint8 // последовательность

type nv uint32 // блок в 32

var (
	encryptSequence = Sequence([]uint8{
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		7, 6, 5, 4, 3, 2, 1, 0,
	})
	decryptSequence = Sequence([]uint8{
		0, 1, 2, 3, 4, 5, 6, 7,
		7, 6, 5, 4, 3, 2, 1, 0,
		7, 6, 5, 4, 3, 2, 1, 0,
		7, 6, 5, 4, 3, 2, 1, 0,
	})
) // последовательности применения ключей

func (n nv) shift() nv {
	return n << 11
} // сдвиг

func toUint32(b []byte) (n1, n2 nv) {
	n1 = nv(b[0]) | nv(b[1])<<8 | nv(b[2])<<16 | nv(b[3])<<24
	n2 = nv(b[4]) | nv(b[5])<<8 | nv(b[6])<<16 | nv(b[7])<<24
	return
} // перевод в блоки 32

func fromUint32(n1, n2 nv) []byte {
	b := make([]byte, 8)
	b[0] = byte((n2 >> 0) & 255)
	b[1] = byte((n2 >> 8) & 255)
	b[2] = byte((n2 >> 16) & 255)
	b[3] = byte((n2 >> 24) & 255)
	b[4] = byte((n1 >> 0) & 255)
	b[5] = byte((n1 >> 8) & 255)
	b[6] = byte((n1 >> 16) & 255)
	b[7] = byte((n1 >> 24) & 255)
	return b
} // перевод из блоков 32

type CipherGostV2 struct {
	sbox *Sbox
	key [8]nv
}

type CipherGost interface {
	Encrypt(src []byte) []byte
	Decrypt(src []byte) []byte
}

func (c *CipherGostV2) Encrypt(src []byte) (dst []byte) {
	n1, n2 := toUint32(src)
	n1, n2 = c.xcrypt(encryptSequence, n1, n2)
	return fromUint32(n1, n2)
}

func (c *CipherGostV2) Decrypt(src []byte) (dst []byte) {
	n1, n2 := toUint32(src)
	n1, n2 = c.xcrypt(decryptSequence, n1, n2)
	return fromUint32(n1, n2)
}

func (c *CipherGostV2) xcrypt(seq Sequence, n1, n2 nv) (nv, nv) {
	for _, i := range seq {
		n1, n2 = c.sbox.k(n1+c.key[i]).shift()^n2, n1
	}
	return n1, n2
}

func NewCipherGost() CipherGost {
	c := &CipherGostV2{sbox: &TestSBox}
	key := TestKey
	c.key = [8]nv{
		nv(key[0]) | nv(key[1])<<8 | nv(key[2])<<16 | nv(key[3])<<24,
		nv(key[4]) | nv(key[5])<<8 | nv(key[6])<<16 | nv(key[7])<<24,
		nv(key[8]) | nv(key[9])<<8 | nv(key[10])<<16 | nv(key[11])<<24,
		nv(key[12]) | nv(key[13])<<8 | nv(key[14])<<16 | nv(key[15])<<24,
		nv(key[16]) | nv(key[17])<<8 | nv(key[18])<<16 | nv(key[19])<<24,
		nv(key[20]) | nv(key[21])<<8 | nv(key[22])<<16 | nv(key[23])<<24,
		nv(key[24]) | nv(key[25])<<8 | nv(key[26])<<16 | nv(key[27])<<24,
		nv(key[28]) | nv(key[29])<<8 | nv(key[30])<<16 | nv(key[31])<<24,
	}
	return c
}

