package im

import (
	"math/rand"
	"sort"
)

// 深拷贝种子
func DeepCopySeed(src Seed) Seed {
	nodes := make([]int, SeedSize)
	for i := range nodes {
		nodes[i] = src.Nodes[i]
	}
	return Seed{Nodes: nodes, Fit: src.Fit}
}

// 深拷贝种群
func DeepCopyPop(src []Seed) []Seed {
	dis := make([]Seed, PopSize)
	for i := range dis {
		dis[i] = DeepCopySeed(src[i])
	}
	return dis
}

func RouletteSelection(src []Seed, cnt int) []Seed {
	sort.Sort(ByFit(src))
	totalFit := float32(0)
	for _, seed := range src {
		totalFit += seed.Fit
	}

	dist := make([]Seed, cnt)
	selectedSeed := make(map[*Seed]bool)
	for i := range dist {
		randomNumber := rand.Float32() * totalFit
		accumulatedFit := float32(0)
		for _, seed := range src {
			accumulatedFit += seed.Fit
			if accumulatedFit >= randomNumber && !selectedSeed[&seed] {
				dist[i] = DeepCopySeed(seed)
				selectedSeed[&seed] = true
				break
			}
		}
	}
	return dist
}
