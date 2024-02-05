package algo

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"

	"sysu.com/task3/im"
)

type Ga struct {
	wg     sync.WaitGroup
	Pop    []im.Seed
	NewPop []im.Seed
}

func (g *Ga) Init() {
	g.Pop = make([]im.Seed, im.PopSize)
	for i := 0; i < im.PopSize; i++ {
		g.Pop[i] = *im.NewSeed()
	}
}

// 交叉
func (g *Ga) Crossover() {
	g.NewPop = im.DeepCopyPop(g.Pop)
	rand.Shuffle(len(g.NewPop), func(i, j int) {
		g.NewPop[i], g.NewPop[j] = g.NewPop[j], g.NewPop[i]
	})
	for i := 0; i < im.PopSize; i += 2 {
		if rand.Float32() < im.PC {
			g.doCrossover(&g.NewPop[i], &g.NewPop[i+1])
			g.wg.Add(2)
			go im.EvaluteSeedSync(&g.NewPop[i], &g.wg)
			go im.EvaluteSeedSync(&g.NewPop[i+1], &g.wg)
		}
	}
	g.wg.Wait()
}

// 均匀交叉
func (g *Ga) doCrossover(seed1, seed2 *im.Seed) {
	for i := 0; i < im.SeedSize; i++ {
		if rand.Float32() < 0.5 {
			seed1.Nodes[i], seed2.Nodes[i] = seed2.Nodes[i], seed1.Nodes[i]
		}
	}
}

// 变异
func (g *Ga) Mutate() {
	for i := range g.NewPop {
		g.wg.Add(1)
		g.doMutate(&g.NewPop[i])
		go im.EvaluteSeedSync(&g.NewPop[i], &g.wg)
	}
	g.wg.Wait()
}

// 单点变异
func (g *Ga) doMutate(seed *im.Seed) {
	for i := range seed.Nodes {
		if rand.Float32() < im.PM {
			seed.Nodes[i] = rand.Intn(im.NetworkSize)
		}
	}
}

func (g *Ga) Select() {
	mergerdPop := im.DeepCopyPop(g.Pop)
	mergerdPop = append(mergerdPop, im.DeepCopyPop(g.NewPop)...)

	g.Pop = im.RouletteSelection(mergerdPop)
	sort.Sort(im.BySeed(g.Pop))
}

func (g *Ga) FindSeed() {
	g.Init()
	for i := 0; i < im.MaxGen; i++ {
		g.Crossover()
		g.Mutate()
		g.Select()
		g.ExportEvolutionInfo(i)
	}
}

func (g *Ga) ExportPop() {
	for idx, elem := range g.Pop {
		fmt.Printf("Pop--%d:%+v\n", idx, elem)
	}
}

func (g *Ga) ExportNewPop() {
	for idx, elem := range g.NewPop {
		fmt.Printf("NewPop--%d:%+v\n", idx, elem)
	}
}

func (g *Ga) ExportEvolutionInfo(gen int) {
	fmt.Printf("gen--%d: [ ", gen)
	for idx, elem := range g.Pop {
		fmt.Printf("%d:%.3f ", idx, elem.Fit)
	}
	fmt.Printf("]\n")
}
