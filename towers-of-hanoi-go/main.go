package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorDefault = "\033[0m"
)

const welcomeInfo = `Welcome to the game Tower of Hanoi!
The rules of the Tower of Hanoi are:
1. You should move all the disks to the last rod.
2. Only one disk can be moved at a time.
3. A disk can only be placed on a larger disk or an empty rod.
Game can be solved in 2^N - 1 moves for N disks.`

// todo: let user define N
func main() {
	game := NewGame()
	scan := bufio.NewScanner(os.Stdin)
	var lastError error

	for !game.IsFinished() {
		render(game, lastError)

		fmt.Print("Move disk (format '{from} {to}'): ")
		scan.Scan()
		input := scan.Text()
		split := strings.Split(input, " ")

		from, err := strconv.Atoi(split[0])
		if err != nil {
			lastError = fmt.Errorf("incorrect input for parameter 'from'")
			continue
		}
		to, err := strconv.Atoi(split[1])
		if err != nil {
			lastError = fmt.Errorf("incorrect input for parameter 'to'")
			continue
		}

		lastError = game.Move(from, to)
	}
	render(game, lastError)
	fmt.Println("Game finished. You won!")
}

func render(game *Game, err error) {
	// clear screen
	fmt.Print("\033[H\033[2J")
	// print welcome info
	fmt.Println(welcomeInfo)
	// print rods and piles
	printRodsAndPiles(game)
	// print moves
	fmt.Printf("Moves: %d\n", game.Moves)
	// print error if it is not nil
	if err != nil {
		fmt.Printf("%sInvalid input:%s %s\n", ColorRed, ColorDefault, err)
	}
}

func printRodsAndPiles(game *Game) {
	height := 4

	for i := 0; i < height; i++ {
		// first row always empty
		var rows []string
		for _, p := range game.Piles {
			mappedIndex := i - (height - len(p))
			if mappedIndex >= 0 {
				rows = append(rows, drawDisk(p[mappedIndex]))
			} else {
				rows = append(rows, drawRod())
			}
		}
		fmt.Printf(" %s \n", strings.Join(rows, " "))
	}

	fmt.Printf("xxxxxxxxxxxxxxxxxxxxxxxxx\n")
}

func drawDisk(disk Disk) string {

	switch disk {
	case Small:
		return "  " + fmt.Sprintf(ColorRed+"ooo"+ColorDefault) + "  "
	case Medium:
		return " " + fmt.Sprintf(ColorGreen+"ooooo"+ColorDefault) + " "
	case Large:
		return fmt.Sprintf(ColorYellow + "ooooooo" + ColorDefault)
	}

	return ""
}

func drawRod() string {
	return "   |   "
}
