package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go" //nolint:depguard
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var user User
	for scanner.Scan() {
		err := jsoniter.ConfigFastest.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			result[strings.ToLower(user.Email[strings.LastIndex(user.Email, "@")+1:])]++
		}
	}

	return result, nil
}
