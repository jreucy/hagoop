package client

import (
	"strconv"
	"strings"
)

type TestClient struct {}

func New() *TestClient {
	return &TestClient{}
}

func (c *TestClient) Map(line string) {
	if line == "" { return }
	array := strings.Split(line, " ")
	for i := 0; i < len(array); i++ {
		Emit(array[i], "1")
	}
}

func (c *TestClient) Reduce(keyValues map[string][]string) {
	for k, v := range keyValues {
		count := 0
		for i := 0; i < len(v); i++ {
			num, _ := strconv.Atoi(v[i])
			count += num
		}
		Emit(k, strconv.Itoa(count))
	}
}