package hw10programoptimization

import (
	"bufio"
	"io"
	"log/slog"
	"strings"
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
	return countDomain(r, domain)
}

func countDomain(r io.Reader, domain string) (DomainStat, error) {
	reader := bufio.NewReader(r)
	result := make(DomainStat)

	for {
		line, err := reader.ReadSlice('\n')
		if len(line) > 0 {
			var user User
			if err := user.UnmarshalJSON(line); err != nil {
				return nil, err
			}
			addressParts := strings.SplitN(user.Email, "@", 2)
			if len(addressParts) != 2 {
				slog.Error("email not contain @:", "email", user.Email)
				continue
			}
			fullDomain := strings.ToLower(addressParts[1])
			if strings.Contains(fullDomain, "."+domain) {
				result[fullDomain]++
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result, nil
}
