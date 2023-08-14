package hw09structvalidator

import "regexp"

func minRule(_min, _v int) bool {
	return _v > _min
}

func maxRule(_max, _v int) bool {
	return _v < _max
}

func inRule(_elems []string, _in string) bool {
	for _, s := range _elems {
		if _in == s {
			return true
		}
	}
	return false
}

func lenRule(_len int, _v string) bool {
	return len(_v) == _len
}

func regexpRule(_pattern, _v string) (matched bool, err error) {
	return regexp.MatchString(_pattern, _v)
}
