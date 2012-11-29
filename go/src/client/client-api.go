package client

import "fmt"

type Client interface {
	Map(chunk string)
	Reduce(keyValues map[string][]string)
}

func Emit(key string, vals ...string) {
	fmt.Print(key)
	for i := 0; i < len(vals); i++ {
		switch i {
		case 0:
			fmt.Print(", " + vals[i])
		default:
			fmt.Print(" " + vals[i])
		}
	}
	fmt.Print("\n")
}