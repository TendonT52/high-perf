package main

import (
	"flag"
	"fmt"
	"high-perf/io/usebuffer"
	"high-perf/io/useconcurrent"
	"high-perf/sorting/usemerge"
	"high-perf/sorting/usemergeconcurrent"
	"runtime"
	"sort"

	"github.com/pkg/profile"
)

func main() {
	inputFile := flag.String("input", "", "input file")
	outputFile := flag.String("output", "", "output file")
	profileType := flag.String("profile", "", "profile type (cpu, clock, goroutine, mem, mem-allocs, mem-heap, threadcreation, trace)")
	profilePath := flag.String("profile-path", "", "profile path")
	technique := flag.String("technique", "", "technique (buffer, concurrent)")

	cpu := flag.Int("cpu", runtime.NumCPU(), "number of cpu to use (logical core)")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Input file is required (-input)")
		return
	}
	fmt.Println("Input file :", *inputFile)

	if *outputFile == "" {
		fmt.Println("Output file is required (-output)")
		return
	}
	fmt.Println("Output file :", *outputFile)

	if *profileType != "" {

		if *profilePath == "" {
			fmt.Println("Profile path is required (-profile-path)")
			return
		}

		switch *profileType {
		case "cpu":
			fmt.Println("CPU profile")
			defer profile.Start(
				profile.CPUProfile,
				profile.ProfilePath(*profilePath)).Stop()
		case "clock":
			fmt.Println("Clock profile")
			defer profile.Start(
				profile.ClockProfile,
				profile.ProfilePath(*profilePath)).Stop()
		case "goroutine":
			fmt.Println("Goroutine profile")
			defer profile.Start(
				profile.GoroutineProfile,
				profile.ProfilePath(*profilePath)).Stop()
		case "mem":
			fmt.Println("Mem profile")
			defer profile.Start(
				profile.MemProfile,
				profile.ProfilePath(*profilePath)).Stop()
		case "mem-allocs":
			fmt.Println("Mem allocs profile")
			defer profile.Start(
				profile.MemProfileAllocs,
				profile.ProfilePath(*profilePath)).Stop()
		case "mem-heap":
			fmt.Println("Mem heap profile")
			defer profile.Start(
				profile.MemProfileHeap,
				profile.ProfilePath(*profilePath)).Stop()
		case "threadcreation":
			fmt.Println("Thread creation profile")
			defer profile.Start(
				profile.ThreadcreationProfile,
				profile.ProfilePath(*profilePath)).Stop()
		case "trace":
			fmt.Println("Trace profile")
			defer profile.Start(
				profile.TraceProfile,
				profile.ProfilePath(*profilePath)).Stop()
		default:
			fmt.Println("Unknown profile type (-profile) : ", *profileType)
		}

	} else {
		fmt.Println("No profile (-profile)")
	}

	if *technique == "" {
		fmt.Println("Technique is required (-technique)")
		return
	}

	switch *technique {
	case "buffer":
		fmt.Println("Buffer technique")
		usebuffer.Init(*inputFile, *outputFile)
		data := usebuffer.ReadInput()
		sort.Slice(data, func(i, j int) bool {
			return data[i].Value < data[j].Value
		})
		usebuffer.WriteOutput(data)
	case "concurrent":
		fmt.Printf("Concurrent technique (cpu: %v)\n", *cpu)
		useconcurrent.Init(*inputFile, *outputFile, int64(*cpu))
		data := useconcurrent.ReadInput()
		sort.Slice(data, func(i, j int) bool {
			return data[i].Value < data[j].Value
		})
		useconcurrent.WriteOutput(data)
	case "merge":
		fmt.Printf("Merge technique (cpu: %v)\n", *cpu)
		usemerge.Init(*inputFile, *outputFile, int64(*cpu))
		data := usemerge.ReadInput()
		usemerge.WriteOutput(data)
	case "merge-concurrent":
		fmt.Printf("Merge concurrent technique (cpu: %v)\n", *cpu)
		usemergeconcurrent.Init(*inputFile, *outputFile, int64(*cpu))
		usemergeconcurrent.Run()
	}
}
