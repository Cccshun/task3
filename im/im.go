package im

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Seed struct {
	Nodes []int
	Fit   float64
}

var (
	AdjMatrix   [][]int //邻接矩阵
	AdjList     [][]int //邻接表
	NetworkSize int
)

// 加载顺序不可改变!
func Init(path string) {
	LoadNetwork(path)
	NetworkSize = len(AdjMatrix)
	LoadGraph(AdjMatrix)
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
	AdjMatrix = [][]int{}
	row := 0
	for scanner.Scan() {
		if str := scanner.Text(); str != "\n" {
			if len(AdjMatrix) == row {
				AdjMatrix = append(AdjMatrix, []int{})
			}
			b, _ := strconv.ParseInt(str, 10, 64) // b={0, 1}
			AdjMatrix[row] = append(AdjMatrix[row], int(b))
		} else {
			row++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read file error:[%v]\n", err)
	}
}

func LoadGraph(network [][]int) {
	AdjList = make([][]int, NetworkSize)
	for row := 0; row < NetworkSize; row++ {
		for col := 0; col < len(network[row]); col++ {
			if network[row][col] == 1 {
				AdjList[row] = append(AdjList[row], col)
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
	defer wg.Done()
	seed.Fit = CalRobustInfluence(seed.Nodes)
}

// 评估种子适应度，并保存在map[seed]fit中。异步计算
func EvaluteSeedAsync(seed *Seed, seedMap map[*Seed]float64, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	fit := CalRobustInfluence(seed.Nodes)
	mu.Lock()
	seedMap[seed] = fit
	mu.Unlock()
}

type BySeed []Seed

func (s BySeed) Len() int           { return len(s) }
func (s BySeed) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySeed) Less(i, j int) bool { return s[i].Fit > s[j].Fit } // 逆序
