package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	// распаковываемая строка
	str := []rune(s)
	// флаг на то, встречали ли двойной бэкслеш
	backslash := false
	// результат
	resultString := ""
	var n int
	for i, value := range str {
		// если первый символ цифра
		if unicode.IsDigit(value) && i == 0 {
			return "", ErrInvalidString
		}
		// при встрече числа (но не цифры) без двух бэкслешей выбросится исключение
		if unicode.IsDigit(value) && unicode.IsDigit(str[i-1]) && str[i-2] != '\\' {
			return "", ErrInvalidString
		}
		if value == '\\' && !backslash {
			backslash = true
			continue
		}
		// после двойного бэкслеша нечего экранировать (встречена буква)
		if backslash && unicode.IsLetter(value) {
			return "", ErrInvalidString
		}
		// экранируем символ
		if backslash {
			resultString += string(value)
			backslash = !backslash
			continue
		}
		if unicode.IsDigit(value) {
			n = int(value - '0')
			if n == 0 {
				resultString = resultString[:len(resultString)-1]
				continue
			}
			resultString += strings.Repeat(string(str[i-1]), n-1)
			continue
		}
		resultString += string(value)
	}
	return resultString, nil
}
