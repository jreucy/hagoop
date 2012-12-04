package main

import (
	"../client"
	"strings"
	"strconv"
	"bufio"
	"log"
	"os"
)

func unpack(keyVal string) []string {
	return strings.Split(strings.TrimSpace(keyVal), ", ")
}

func main() {
	if len(os.Args) != 6 { return }

	var c client.Client
	c = client.New()

	file, err := os.Open(os.Args[2])
	if err != nil { log.Fatal("main : ", err) }
	startLine, err := strconv.Atoi(os.Args[3])
	if err != nil { log.Fatal("main : ", err) }
	endLine, err := strconv.Atoi(os.Args[4])
	if err != nil { log.Fatal("main : ", err) }
	offset, err := strconv.ParseInt(os.Args[5], 10, 64)
	if err != nil { log.Fatal("main : ", err) }
	_, err = file.Seek(offset, 0)
	if err != nil { log.Fatal("main : ", err) }

	switch os.Args[1] {
	case "map":
		fileBuf := bufio.NewReader(file)
		for startLine != endLine {
			line, _ := fileBuf.ReadString('\n')
			c.Map(strings.TrimSpace(line))
			startLine++
		}
	case "reduce":
		preMap := make(map[string]string)
		keyValues := make(map[string][]string)

		fileBuf := bufio.NewReader(file)
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