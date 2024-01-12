package main

import (
	"fmt"
	"time"

	"sysu.com/task3/algo"
)

func main() {
	startTime := time.Now()
	ga := &algo.Ga{}
	ga.FindSeed()

	elapsedTime := time.Since(startTime)
	fmt.Printf("运行时间: %s\n", elapsedTime)
}
