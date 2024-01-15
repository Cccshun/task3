package im

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Seed struct {
	Nodes []int
	Fit   float32
}

var (
	Network     [][]int
	G           [][]int
	NetworkSize int
)

func Init(path string) {
	LoadNetwork(path)
	LoadGraph(Network)
	NetworkSize = len(Network)
}

// 读取数据文件
func LoadNetwork(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open file error:[%v]\n", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)
	Network = [][]int{}
	row := 0
	for scanner.Scan() {
		if str := scanner.Text(); str != "\n" {
			if len(Network) == row {
				Network = append(Network, []int{})
			}
			b, _ := strconv.ParseInt(str, 10, 64) // b={0, 1}
			Network[row] = append(Network[row], int(b))
		} else {
			row++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read file error:[%v]\n", err)
	}
}

func LoadGraph(network [][]int) {
	G = make([][]int, len(network))
	for row := 0; row < len(network); row++ {
		for col := 0; col < len(network[row]); col++ {
			if network[row][col] == 1 {
				G[row] = append(G[row], col)
			}
		}
	}
}

func NewSeed() *Seed {
	nodes := make([]int, PopSize)
	for i := 0; i < PopSize; i++ {
		nodes[i] = rand.Intn(NetworkSize)
	}
	return &Seed{nodes, 0}
}

func EvaluteSeedSync(seed *Seed, wg *sync.WaitGroup) {
	// TODO
	defer wg.Done()
	fit := 0
	for _, node := range seed.Nodes {
		fit += node
	}
	seed.Fit = float32(fit)
	time.Sleep(30 * time.Millisecond)
}

// 评估种子适应度，并保存在map[seed]fit中。异步计算
func EvaluteSeedAsync(seed *Seed, seedMap map[*Seed]float32, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	// TODO
	fit := 0
	for _, node := range seed.Nodes {
		fit += node
	}
	mu.Lock()
	seedMap[seed] = float32(fit)
	mu.Unlock()
	time.Sleep(30 * time.Millisecond)
}

type ByFit []Seed

func (s ByFit) Len() int           { return len(s) }
func (s ByFit) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFit) Less(i, j int) bool { return s[i].Fit > s[j].Fit } // 逆序
