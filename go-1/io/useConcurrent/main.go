package useconcurrent

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type KeyValuePair struct {
	Key   int64
	Value float64
}

var (
	cpu        int64
	inputFile  string
	outputFile string
)

func Init(input string, output string, core int64) {
	cpu = core
	inputFile = input
	outputFile = output
}

func read(offset int64, limit int64) []KeyValuePair {
	var data []KeyValuePair

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(offset, 0)
	reader := bufio.NewReader(file)

	var cumulativeSize int64
	for cumulativeSize < limit {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if line[0] == 's' {
			parts := strings.Split(line, ": ")
			keyStr := parts[0][4:]
			valStr := parts[1][:len(parts[1])-1]
			key, _ := strconv.ParseInt(keyStr, 10, 64)
			value, _ := strconv.ParseFloat(valStr, 64)
			data = append(data, KeyValuePair{Key: key, Value: value})
		}
		cumulativeSize += int64(len([]byte(line)))
	}
	return data
}

func ReadInput() []KeyValuePair {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	fileSize, err := file.Stat()
	if err != nil {
		panic(err)
	}
	file.Close()

	var data []KeyValuePair
	wg := sync.WaitGroup{}

	limit := fileSize.Size() / int64(cpu)
	var current int64
	current = 0
	wg.Add(int(cpu))
	for i := 0; i < int(cpu); i++ {

		go func(offset int64, limit int64) {
			data = append(data, read(offset, limit)...)
			wg.Done()
		}(current, limit)

		current += int64(limit)
	}

	wg.Wait()

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
