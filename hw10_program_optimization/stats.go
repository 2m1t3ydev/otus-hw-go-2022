package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	s := bufio.NewScanner(r)

	for s.Scan() {
		email := fastjson.GetString(s.Bytes(), "Email")

		i := strings.Index(email, "@")
		if i == -1 {
			continue
		}

		if strings.HasSuffix(email, domain) {
			result[strings.ToLower(email[i+1:])]++
		}
	}

	return result, nil
}
