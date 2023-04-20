package main

import (
	"os"
	"runtime"
)

func main() {
	args := os.Args[1:]
	inputFile := args[0]
	outputFile := args[1]

	cpu := runtime.NumCPU()
	Init(inputFile, outputFile, int64(cpu))
	Run()
}
