package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"../client"
	"log"
)

func unpack(keyVal string) []string {
	return strings.Split(strings.TrimSpace(keyVal), ", ")
}

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

		for startLine != endLine {
			line, _ := fileBuf.ReadString('\n')
			c.Map(strings.TrimSpace(line))
			startLine++
		}
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
			keyVal := unpack(line)
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

		c.Reduce(keyValues)
	}
}