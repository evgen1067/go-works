package sha256

type Digest struct {
	h   [8]uint32
	x   [BlockSize]byte
	nx  int
	len uint64
}

type SHA256 interface {
	Reset()
	Write(p []byte) []uint32
}

// Сброс к стандартным установкам
func (d *Digest) Reset() {
	d.h[0] = init0
	d.h[1] = init1
	d.h[2] = init2
	d.h[3] = init3
	d.h[4] = init4
	d.h[5] = init5
	d.h[6] = init6
	d.h[7] = init7
}

func (d *Digest) Write(p []byte) []uint32 {
	d.len += uint64(len(p))
	for len(p)%BlockSize != 0 {
		p = append(p, byte(0))
	}
	block(d, p)
	return d.h[:]
}

func NewSHA256() SHA256 {
	sha := &Digest{}
	sha.Reset()
	return sha
}
