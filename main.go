package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Preset int

func (p *Preset) String() string {
	presetLimit := Preset(len(PresetStringArr))
	if *p < PresetRandom || *p >= presetLimit {
		return "NA"
	}
	return PresetStringArr[*p]
}

func (p *Preset) Set(presetStr string) error {
	presetInt := 0
	presetLimit := len(PresetStringArr)
	if i64, err := strconv.ParseInt(presetStr, 0, 0); err == nil {
		presetInt = int(i64)
		if presetInt < 0 || presetInt >= presetLimit {
			return fmt.Errorf("invalid preset number: %d", presetInt)
		}
	} else {
		pStr := ""
		for presetInt = range PresetStringArr {
			if PresetStringArr[presetInt] == presetStr {
				pStr = presetStr
				break
			}
		}
		if pStr == "" {
			return fmt.Errorf("invalid preset name: %q", presetStr)
		}
	}
	*p = Preset(presetInt)
	return nil
}

const (
	PresetRandom Preset = iota
	PresetBeehive
	PresetGlider
	PresetPulsar
)

var (
	_ flag.Value = new(Preset)
)

var (
	presetVar       Preset
	PresetStringArr = []string{
		"random", "beehive", "glider", "pulsar"}
)

type Life [64][256]int

func printTheGrid(grid Life) {
	fmt.Print("\033[H\033[2J")

	liveChar := "\u2588"
	deadChar := " "

	liveColor := "\033[48;5;227m"
	deadColor := "\033[48;5;235m"
	resetColor := "\033[0m"

	for _, row := range grid {
		for _, cell := range row {
			if cell == 0 {
				fmt.Printf("%s%s%s", deadColor, deadChar, resetColor)
			} else if cell == 1 {
				fmt.Printf("%s%s%s", liveColor, liveChar, resetColor)
			}
		}
		fmt.Println("")
	}
}

func randomise(grid *Life) {
	for i, v := range grid {
		for j := range v {
			rand.Seed(time.Now().UnixNano())
			grid[i][j] = rand.Intn(2)
		}
	}
}

func presetPulsar(grid *Life) {
	// Initialize the grid with the Pulsar pattern
	// Period: 3
	for i := 26; i <= 38; i += 11 {
		for j := 18; j <= 30; j++ {
			grid[i][j] = 1
			grid[j][i] = 1
		}
	}
}

func presetGlider(grid *Life) {
	// Initialize the grid with the Glider pattern
	grid[5][4] = 1
	grid[6][5] = 1
	grid[6][6] = 1
	grid[5][6] = 1
	grid[4][6] = 1
}

func presetBeehive(grid *Life) {
	// Initialize the grid with the Beehive pattern
	grid[20][20] = 1
	grid[20][21] = 1
	grid[19][22] = 1
	grid[21][22] = 1
	grid[20][23] = 1
	grid[20][24] = 1
}

func init() {
	flag.Var(&presetVar, "preset",
		"choose life pattern: 0|random (default), 1|beehive, 2|glider, 3|pulsar")
}

func main() {
	flag.Parse()
	grid := Life{}

	// presetGlider(&grid)
	switch presetVar {
	case PresetRandom:
		randomise(&grid)
	case PresetBeehive:
		presetBeehive(&grid)
	case PresetGlider:
		presetGlider(&grid)
	case PresetPulsar:
		presetPulsar(&grid)
	default:
		randomise(&grid)
	}

	// run the grid
	for {
		printTheGrid(grid)
		var newGrid Life

		for i, v := range grid {
			for j, k := range v {
				bi := (i - 1 + len(grid)) % len(grid)
				fi := (i + 1) % len(grid)

				bj := (j - 1 + len(v)) % len(v)
				fj := (j + 1) % len(v)

				neighbours := grid[i][bj] + grid[bi][j] + grid[bi][bj] + grid[fi][j] + grid[i][fj] + grid[fi][fj] + grid[bi][fj] + grid[fi][bj]
				if k == 1 {
					if neighbours < 2 || neighbours > 3 {
						newGrid[i][j] = 0
					} else {
						newGrid[i][j] = 1
					}
				}
				if k == 0 && neighbours == 3 {
					newGrid[i][j] = 1
				}
			}
		}
		copy(grid[:], newGrid[:])
		time.Sleep(20 * time.Millisecond)
	}
}
