package client

import (
	"strconv"
	"strings"
)

type TestClient struct {}

func New() *TestClient {
	return &TestClient{}
}

func (c *TestClient) Map(chunk string) string {
	res := ""
	array := strings.Split(chunk, " ")
	for i := 0; i < len(array); i++ {
		res += Pack(array[i], "1")
	}
	return res
}

func (c *TestClient) Reduce(keyValues map[string][]string) string {
	answer := ""
	for k, v := range keyValues {
		count := 0
		for i := 0; i < len(v); i++ {
			num, _ := strconv.Atoi(v[i])
			count += num
		}
		answer += Pack(k, strconv.Itoa(count))
	}
	return answer
}