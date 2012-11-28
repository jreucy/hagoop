package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"../client"
	"log"
)

func main() {
	// ./main [map/reduce] [file] [start] [end]
	if len(os.Args) != 5 { return }

	var c client.Client
	c = client.New()

	switch os.Args[1] {
	case "map":
		file, err := os.Open(os.Args[2])
		if err != nil { log.Fatal(err) }
		fileBuf := bufio.NewReader(file)
	
		startLine, _ := strconv.Atoi(os.Args[3])
		endLine, _ := strconv.Atoi(os.Args[4])

		for i := 0; i < startLine; i++ {
			fileBuf.ReadString('\n')
		}

		data := ""
		for startLine != endLine {
			line, _ := fileBuf.ReadString('\n')
			data += line
			startLine++
		}
		data = strings.Replace(data, "\n", " ", -1)
		data = strings.TrimSpace(data)
		fmt.Println(c.Map(data))
	case "reduce":
		preMap := make(map[string]string)
		keyValues := make(map[string][]string)

		file, err := os.Open(os.Args[2])
		if err != nil { log.Fatal(err) }
		fileBuf := bufio.NewReader(file)
	
		startLine, _ := strconv.Atoi(os.Args[3])
		endLine, _ := strconv.Atoi(os.Args[4])

		for i := 0; i < startLine; i++ {
			fileBuf.ReadString('\n')
		}

		// determine the length of the file
		for startLine != endLine {
			line, _ := fileBuf.ReadString('\n')
			keyVal := client.Unpack(line)
			_, ok := preMap[keyVal[0]]
			if ok {
				preMap[keyVal[0]] += " " + keyVal[1]
			} else {
				preMap[keyVal[0]] = keyVal[1]
			}
			startLine++
		}

		for i, v := range preMap {
			keyValues[i] = strings.Split(v, " ")
		}

		fmt.Println(c.Reduce(keyValues))
	}
}