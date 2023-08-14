package stream

type BBS struct {
	gamma byte
	n     int
}

type BBSInterface interface {
	Gamma() byte
	NextGamma()
	SetGamma(gamma byte)
}

func NewBBSInterface(r, p, q int) BBSInterface {
	return &BBS{
		gamma: byte(r * r % (p * q)),
		n:     p * q,
	}
}

func (B *BBS) Gamma() byte {
	return B.gamma
}

func (B *BBS) NextGamma() {
	B.gamma = byte(int(B.gamma) * int(B.gamma) % B.n)
}

func (B *BBS) SetGamma(gamma byte) {
	B.gamma = gamma
}
