package client

import (
	"strconv"
	"strings"
)

type WordCount struct {}

func New() *WordCount {
	return &WordCount{}
}

func (c *WordCount) Map(line string) {
	if line == "" { return }
	array := strings.Split(line, " ")
	for i := 0; i < len(array); i++ {
		Emit(array[i], "1")
	}
}

func (c *WordCount) Reduce(keyValues map[string][]string) {
	for k, v := range keyValues {
		count := 0
		for i := 0; i < len(v); i++ {
			num, _ := strconv.Atoi(v[i])
			count += num
		}
		Emit(k, strconv.Itoa(count))
	}
}