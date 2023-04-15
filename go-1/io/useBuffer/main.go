package usebuffer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type KeyValuePair struct {
	Key   int64
	Value float64
}
var (
	inputFile  string
	outputFile string
)

func Init(input string, output string) {
	inputFile = input
	outputFile = output
}


func ReadInput() []KeyValuePair {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data []KeyValuePair
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading file:", err)
			return nil
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		keyStr := parts[0][4:]
		key, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			fmt.Println("Error parsing key:", err)
			continue
		}

		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			fmt.Println("Error parsing value:", err)
			continue
		}

		data = append(data, KeyValuePair{Key: key, Value: value})
	}
	return data
}

func WriteOutput(data []KeyValuePair) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)

	for _, kv := range data {
		bw.Write([]byte(fmt.Sprintf("std-%d: %v\n", kv.Key, kv.Value)))
	}
	err = bw.Flush()
	if err != nil {
		return err
	}
	return nil
}
