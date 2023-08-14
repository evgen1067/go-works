package stream

type LinearCongruent struct {
	a     int
	b     int
	e     int
	gamma byte
}

type LinearCongruentGenerator interface {
	Gamma() byte
	NextGamma()
	SetGamma(gamma byte)
}

func (l *LinearCongruent) Gamma() byte {
	return l.gamma
}

func (l *LinearCongruent) NextGamma() {
	l.gamma = byte((l.a*int(l.gamma) + l.b) % l.e)
}

func (l *LinearCongruent) SetGamma(gamma byte) {
	l.gamma = gamma
}

func NewLinearCongruentGenerator(a, b, e int, gamma string) LinearCongruentGenerator {
	return &LinearCongruent{
		a:     a,
		b:     b,
		e:     e,
		gamma: []byte(gamma)[0],
	}
}
