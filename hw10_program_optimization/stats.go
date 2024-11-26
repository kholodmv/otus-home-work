package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	suffix := []byte(domain)
	var user User
	for i := 0; scanner.Scan(); i++ {
		if err := jsoniter.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}
		if user.Email == "" {
			continue
		}

		if bytes.HasSuffix([]byte(user.Email), suffix) {
			domainPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainPart]++
		}
	}
	return result, nil
}
