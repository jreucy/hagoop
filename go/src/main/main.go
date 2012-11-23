package main

import (
	"flag"
	"../client"
)

func main() {
	// var action *int = flag.Int("action", "map", "either 'map' or 'reduce'")
	// var file *string = flag.String("file", "", "file to map/reduce")
	// var range *string = flag.String("range", "", "range of lines to map/reduce")
	flag.Parse()

	var c client.Client
	c = client.New()
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "map":
			c.Map()
		case "reduce":
			c.Reduce()
		}
	} else {
		return
	}
}

