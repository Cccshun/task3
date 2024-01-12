package im

import (
	"math/rand"
	"sync"
	"time"
)

type Seed struct {
	Nodes []int
	Fit   float32
}

func NewSeed() *Seed {
	nodes := make([]int, PopSize)
	for i := 0; i < PopSize; i++ {
		nodes[i] = rand.Intn(NetworkSize)
	}
	return &Seed{nodes, 0}
}

func EvaluteSeed(seed *Seed, wg *sync.WaitGroup) {
	// TODO
	defer wg.Done()
	seed.Fit = rand.Float32() * 100
	time.Sleep(10 * time.Millisecond)
}

type ByFit []Seed

func (s ByFit) Len() int           { return len(s) }
func (s ByFit) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFit) Less(i, j int) bool { return s[i].Fit > s[j].Fit } // 逆序
