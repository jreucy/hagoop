package client

import "fmt"

type Client interface {
	Map(chunk string)
	Reduce(keyValues map[string][]string)
}

func Emit(vals ...interface{}) {
	for i := 0; i < len(vals); i++ {
		switch i {
		case 0:
			fmt.Print(vals[i].(string))
		case 1:
			fmt.Print(", " + vals[i].(string))
		default:
			fmt.Print(" " + vals[i].(string))
		}
	}
	fmt.Print("\n")
}