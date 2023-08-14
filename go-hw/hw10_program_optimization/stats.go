package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/evgen1067/hw10_program_optimization/models"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var user models.User
	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return result, err
		}

		err = user.UnmarshalJSON(line)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.Split(user.Email, "@")[1])]++
		}
	}
	return result, nil
}
