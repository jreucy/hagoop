package client

import "fmt"

type Client interface {
	Map(chunk string)
	Reduce(keyValues map[string][]string)
}

func Emit(vals ...interface{}) {
	fmt.Print(vals[0])
	for i := 1; i < len(vals); i++ {
		fmt.Print(", " + vals[i].(string))
	}
	fmt.Print("\n")
}