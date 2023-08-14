package stream

type Encrypt struct {
	a     int
	b     int
	e     int
	t     int
	p     int
	q     int
	r     int
	gamma string
}

type EncryptInterface interface {
	SyncEncrypt(message string, bbs bool) []byte
	AsyncEncrypt(message string, decode, bbs bool) []byte
}

func NewEncryptInterface(a, b, e, t, p, q, r int, gamma string) EncryptInterface {
	return &Encrypt{
		a:     a,
		b:     b,
		e:     e,
		t:     t,
		p:     p,
		q:     q,
		r:     r,
		gamma: gamma,
	}
}

func (e *Encrypt) SyncEncrypt(message string, bbs bool) []byte {
	var (
		inc          = make([]byte, 0)
		tmp          = make([]byte, 0)
		str          = []byte(message)
		lCongruent   LinearCongruentGenerator
		bbsInterface BBSInterface
	)

	if !bbs {
		lCongruent = NewLinearCongruentGenerator(e.a, e.b, e.e, e.gamma)
	} else {
		bbsInterface = NewBBSInterface(e.r, e.p, e.q)
	}

	for i, val := range str {
		inc = append(inc, val)
		if !bbs {
			tmp = append(tmp, inc[i]^lCongruent.Gamma())
			lCongruent.NextGamma()
			continue
		}
		tmp = append(tmp, inc[i]^bbsInterface.Gamma())
		bbsInterface.NextGamma()
	}
	return tmp
}

func (e *Encrypt) AsyncEncrypt(message string, decode, bbs bool) []byte {
	var (
		inc          = make([]byte, 0)
		tmp          = make([]byte, 0)
		str          = []byte(message)
		lCongruent   LinearCongruentGenerator
		bbsInterface BBSInterface
		_gamma       byte
	)

	if !bbs {
		lCongruent = NewLinearCongruentGenerator(e.a, e.b, e.e, e.gamma)
	} else {
		bbsInterface = NewBBSInterface(e.r, e.p, e.q)
	}

	for i, val := range str {
		inc = append(inc, val)
		if i < e.t {
			if !bbs {
				tmp = append(tmp, inc[i]^lCongruent.Gamma())
				lCongruent.NextGamma()
				continue
			}
			tmp = append(tmp, inc[i]^bbsInterface.Gamma())
			bbsInterface.NextGamma()
			continue
		}

		if !decode {
			_gamma = tmp[0]
		} else {
			_gamma = inc[0]
		}

		for j := 1; j < i; j++ {
			if !decode {
				_gamma = tmp[j] ^ _gamma
				continue
			}
			_gamma = inc[j] ^ _gamma
		}
		if !bbs {
			_gamma = lCongruent.Gamma() ^ _gamma
			lCongruent.SetGamma(_gamma)
			tmp = append(tmp, inc[i]^lCongruent.Gamma())
		} else {
			_gamma = bbsInterface.Gamma() ^ _gamma
			bbsInterface.SetGamma(_gamma)
			tmp = append(tmp, inc[i]^bbsInterface.Gamma())
		}

	}
	return tmp
}
