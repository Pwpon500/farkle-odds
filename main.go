package main

import "fmt"

var (
	dice     = 3
	curScore = 350
	maxDepth = 2
)

func main() {
	//fmt.Println(backtrack([]int{}, dice, curScore))
	fmt.Println(scoreRoll([]int{1, 2, 3, 4, 5, 6}, 0, 0))
}

// backtrack to find all possible roll combinations for a given number of dice
func backtrack(rolls []int, toRoll int, score int, depth int) float64 {
	if toRoll == 0 {
		return float64(scoreRoll(rolls, score, depth))
	}
	if depth >= maxDepth {
		return 0
	}

	sum := 0.0
	for i := 1; i <= 6; i++ {
		sum += backtrack(append(rolls, i), toRoll-1, score, depth)
	}
	return sum / 6
}

// score a roll by finding the number of occurences of each number and scoring roll
func scoreRoll(rolls []int, score int, depth int) float64 {
	var occurences [6]int
	for _, elem := range rolls {
		occurences[elem-1]++
	}
	var frequencies [7][]int

	for i := 0; i < 6; i++ {
		occurs := occurences[i]
		frequencies[occurs] = append(frequencies[occurs], i)
	}

	toReturn := 0.0
	countIndiv := true
	numUsed := 0
	if len(frequencies[6]) > 0 {
		toReturn += 3000
		numUsed += 6
	} else if len(frequencies[5]) > 0 {
		toReturn += 2000
		numUsed += 5
	} else if len(frequencies[4]) > 0 && len(frequencies[2]) > 0 {
		toReturn += 1500
		countIndiv = false
		numUsed += 6
	} else if len(frequencies[4]) > 0 {
		toReturn += 1000
		numUsed += 4
	} else if len(frequencies[3]) == 2 {
		toReturn += 2500
		numUsed += 6
	} else if len(frequencies[2]) == 3 {
		toReturn += 1500
		countIndiv = false
		numUsed += 6
	} else if len(frequencies[3]) == 1 {
		if frequencies[3][0] == 0 {
			toReturn += 300
			numUsed += 3
		} else {
			toReturn += float64(100 * (frequencies[3][0] + 1))
			numUsed += 3
		}
	} else if len(frequencies[1]) == 6 {
		toReturn += 1500
		countIndiv = false
		numUsed += 6
	}

	if occurences[0] < 3 && countIndiv {
		toReturn += float64(100 * occurences[0])
		numUsed += occurences[0]
	}
	if occurences[4] < 3 && countIndiv {
		toReturn += float64(50 * occurences[4])
		numUsed += occurences[4]
	}

	if toReturn == 0 {
		return 0
	}

	numLeft := len(rolls) - numUsed
	if numLeft == 0 {
		numLeft = 6
	}

	toReturn += float64(score)
	expectedRoll := backtrack([]int{}, numLeft, int(toReturn), depth+1)
	if expectedRoll > float64(toReturn) {
		toReturn = expectedRoll
	}

	return toReturn
}
