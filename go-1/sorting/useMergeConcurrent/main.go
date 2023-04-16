package usemergeconcurrent

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
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
	sort.Slice(data, func(i, j int) bool {
		return data[i].Value < data[j].Value
	})
	return data
}

func merge(data ...[]KeyValuePair) []KeyValuePair {
	var result []KeyValuePair
	cursor := make([]int, len(data))

	for {
		var min KeyValuePair
		var minIndex = -1
		for i, d := range data {
			if cursor[i] < len(d) {
				if minIndex == -1 || d[cursor[i]].Value < min.Value {
					min = d[cursor[i]]
					minIndex = i
				}
			}
		}
		if minIndex == -1 {
			break
		}
		result = append(result, min)
		cursor[minIndex]++
	}
	return result
}

func LowerBound(array []KeyValuePair, target float64) int {
	low, high, mid := 0, len(array)-1, 0
	for low <= high {
		mid = (low + high) / 2
		if array[mid].Value >= target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low
}

func UpperBound(array []KeyValuePair, target float64) int {
	low, high, mid := 0, len(array)-1, 0

	for low <= high {
		mid = (low + high) / 2
		if array[mid].Value > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low
}

func MergeInBlock(min, max float64, data [][]KeyValuePair) []KeyValuePair {
	var result []KeyValuePair
	start := make([]int, len(data))
	stop := make([]int, len(data))

	for i, d := range data {
		start[i] = LowerBound(d, min)
		stop[i] = UpperBound(d, max)
	}

	for {
		var min KeyValuePair
		var minIndex = -1
		for i, d := range data {
			if start[i] < stop[i] {
				if minIndex == -1 || d[start[i]].Value < min.Value {
					min = d[start[i]]
					minIndex = i
				}
			}
		}
		if minIndex == -1 {
			break
		}
		result = append(result, min)
		start[minIndex]++
	}
	return result
}

func ReadInput() [][]KeyValuePair {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	fileSize, err := file.Stat()
	if err != nil {
		panic(err)
	}
	file.Close()

	data := make([][]KeyValuePair, cpu)
	wg := sync.WaitGroup{}

	limit := fileSize.Size() / int64(cpu)
	var current int64
	current = 0
	wg.Add(int(cpu))
	for i := 0; i < int(cpu); i++ {

		go func(offset int64, limit int64, i int) {
			data[i] = read(offset, limit)
			wg.Done()
		}(current, limit, i)

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

func FindMaxMin(data [][]KeyValuePair) (float64, float64) {
	max := data[0][len(data[0])-1].Value
	min := data[0][0].Value
	for _, d := range data {
		if d[0].Value < min {
			min = d[0].Value
		}
		if d[len(d)-1].Value > max {
			max = d[len(d)-1].Value
		}
	}
	return max, min
}

func Run() {
	data := ReadInput()
	max, min := FindMaxMin(data)
	size := max - min
	block := size / float64(cpu)
	var start, stop float64
	start = min
	stop = min + block
	result := make([][]KeyValuePair, cpu)

	wg1 := sync.WaitGroup{}
	for i := 0; i < int(cpu); i++ {
		if i == int(cpu)-1 {
			stop = max
		}
		wg1.Add(1)
		go func(i int, start float64, stop float64, data [][]KeyValuePair) {
			result[i] = MergeInBlock(start, stop, data)
			wg1.Done()
		}(i, start, stop, data)
		start = stop + math.SmallestNonzeroFloat64
		stop = start + block
	}

	wg1.Wait()

	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	file.Close()

	wg3 := sync.WaitGroup{}
	blockSize := make([]int64, cpu)
	for i := 0; i < int(cpu); i++ {
		wg3.Add(1)
		go func(i int) {
			for _, kv := range result[i] {
				blockSize[i] += int64(len([]byte(fmt.Sprintf("std-%d: %v\n", kv.Key, kv.Value))))
			}
			wg3.Done()
		}(i)
	}

	wg3.Wait()

	wg2 := sync.WaitGroup{}

	var offset int64
	for i := 0; i < int(cpu); i++ {
		wg2.Add(1)
		go func(offset int64, data []KeyValuePair) {
			write(offset, data)
			wg2.Done()
		}(offset, result[i])
		offset += blockSize[i]
	}

	wg2.Wait()
}

func write(offset int64, data []KeyValuePair) {
	file, err := os.OpenFile(outputFile, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(offset, 0)
	writer := bufio.NewWriter(file)

	for _, kv := range data {
		writer.Write([]byte(fmt.Sprintf("std-%d: %v\n", kv.Key, kv.Value)))
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
