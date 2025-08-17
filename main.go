package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/eiannone/keyboard"
)

//2048 game

func randIndexes() (int, int) {
	return rand.IntN(4), rand.IntN(4)
}

func initiateGrid(grid *[4][4]int) {
	for i := range grid {
		for j := range grid[i] {
			grid[i][j] = 0
		}
	}
}

func printGrid(grid [4][4]int) {
	fmt.Println("+------+------+------+------+")
	for _, row := range grid {
		for _, val := range row {
			if val == 0 {
				fmt.Printf("|%4s ", ".") // show . for empty cells
			} else {
				fmt.Printf("|%4d ", val) // align numbers to 4 spaces
			}
		}
		fmt.Println("|")
		fmt.Println("+------+------+------+------+")
	}
}

func gameWin(grid [4][4]int) bool {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 2048 {
				return true
			}
		}
	}
	return false
}

func gameLose(grid [4][4]int) bool {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				return false
			}
			if i < 3 && grid[i][j] == grid[i+1][j] {
				return false
			}
			if j < 3 && grid[i][j] == grid[i][j+1] {
				return false
			}
		}
	}
	return true
}

func compress(grid *[4]int) {
	temp := [4]int{}
	var sliceIndex int
	var doubled bool
	for i := 0; i < len(grid); i++ {
		if grid[i] == 0 {
			continue
		}

		if temp[sliceIndex] == 0 {
			temp[sliceIndex] = grid[i]
			doubled = false
		} else {
			if !doubled && grid[i] == temp[sliceIndex] {
				temp[sliceIndex] = temp[sliceIndex] * 2
				doubled = true
				sliceIndex++
			} else if grid[i] != temp[sliceIndex] {
				sliceIndex++
				temp[sliceIndex] = grid[i]
				doubled = false
			}
		}
		// fmt.Println(temp)
	}
	*grid = temp
}

func reverse(grid *[4][4]int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			grid[i][j], grid[i][3-j] = grid[i][3-j], grid[i][j]
		}
	}
}

func transpose(grid *[4][4]int) {
	for i := 0; i < 4; i++ {
		for j := i; j < 4; j++ {
			grid[i][j], grid[j][i] = grid[j][i], grid[i][j]
		}
	}
}

func moveToDown(grid *[4][4]int) {
	transpose(grid)
	moveToRight(grid)
	transpose(grid)
}

func moveToUp(grid *[4][4]int) {
	transpose(grid)
	moveToLeft(grid)
	transpose(grid)
}

func moveToRight(grid *[4][4]int) {
	reverse(grid)
	moveToLeft(grid)
	reverse(grid)
}

func moveToLeft(grid *[4][4]int) {
	for i := range grid {
		compress(&grid[i])
	}
}

func toPlaceOnRandIndex(grid *[4][4]int) {
	empty := [][2]int{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if grid[i][j] == 0 {
				empty = append(empty, [2]int{i, j})
			}
		}
	}
	if len(empty) == 0 {
		return
	}
	idx := rand.IntN(len(empty))
	grid[empty[idx][0]][empty[idx][1]] = 2
}

func main() {
	fmt.Println("Welcome to 2048 game")
	firstRandRow, firstRandCol := randIndexes()
	fmt.Println("First indexes", firstRandRow, firstRandCol)
	grid := [4][4]int{}
	initiateGrid(&grid)
	grid[firstRandRow][firstRandCol] = 2
	printGrid(grid)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		fmt.Println("Use arrows to move, q to quit")
		char, key, _ := keyboard.GetKey()
		switch key {
		case keyboard.KeyArrowUp:
			moveToUp(&grid)
		case keyboard.KeyArrowDown:
			moveToDown(&grid)
		case keyboard.KeyArrowLeft:
			moveToLeft(&grid)
		case keyboard.KeyArrowRight:
			moveToRight(&grid)
		case keyboard.KeyEsc:
			return
		}
		if char == 'q' {
			return
		}
		if gameWin(grid) {
			fmt.Println("You won")
			break
		}
		if gameLose(grid) {
			fmt.Println("You lost")
			break
		}
		toPlaceOnRandIndex(&grid)
		printGrid(grid)
	}
}
