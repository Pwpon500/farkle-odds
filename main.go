package main

import "fmt"

var (
	dice     = 3
	curScore = 350
)

func main() {
	fmt.Println(backtrack([]int{}, dice))
}

// backtrack to find all possible roll combinations for a given number of dice
func backtrack(rolls []int, toRoll int) float64 {
	if toRoll == 0 {
		return float64(maxScore(rolls))
	}

	sum := 0.0
	for i := 1; i <= 6; i++ {
		sum += backtrack(append(rolls, i), toRoll-1)
	}
	return sum / 6
}

// find the maximum score for a given roll set
func maxScore(rolls []int) int {
	max := 0
	partitions := generatePartitions(rolls)
	for _, partition := range partitions {
		score := scoreSuperSet(partition)
		if score > max {
			max = score
		}
	}

	if max > 0 {
		return curScore + max
	}
	return 0
}

// find all the ways to partition a given roll in order to score each partition
func generatePartitions(rolls []int) [][][]int {
	if len(rolls) == 1 {
		return [][][]int{{{rolls[0]}}}
	}

	firstElem := rolls[0]
	rest := generatePartitions(rolls[1:])
	toReturn := [][][]int{}

	for _, elem := range rest {
		toReturn = append(toReturn, append(elem, []int{firstElem}))

		for i, set := range elem {
			removed := make([][]int, len(elem))
			copy(removed, elem[:])
			removed = append(removed[:i], removed[i+1:]...)
			toReturn = append(toReturn, append(removed, append(set, firstElem)))
		}
	}

	return toReturn
}

// score all the partitions
func scoreSuperSet(rolls [][]int) int {
	if len(rolls) == 2 {
		if len(rolls[0]) == 3 && len(rolls[1]) == 3 {
			if rolls[0][0] == rolls[0][1] && rolls[0][1] == rolls[0][2] && rolls[1][0] == rolls[1][1] && rolls[1][1] == rolls[1][2] {
				return curScore + 2500
			}
		}
		if len(rolls[0]) == 4 && len(rolls[1]) == 2 {
			if rolls[0][0] == rolls[0][1] && rolls[0][1] == rolls[0][2] && rolls[0][2] == rolls[0][3] && rolls[1][0] == rolls[1][1] {
				return curScore + 1500
			}
		}
	}

	if len(rolls) == 3 {
		if len(rolls[0]) == 2 && len(rolls[1]) == 2 && len(rolls[2]) == 2 {
			if rolls[0][0] == rolls[0][1] && rolls[1][0] == rolls[1][1] && rolls[2][0] == rolls[2][1] {
				return curScore + 1500
			}
		}
	}

	sum := 0
	for _, elem := range rolls {
		sum += scoreSet(elem)
	}

	return sum
}

// score a single partition
func scoreSet(rolls []int) int {
	if allEqual(rolls) {
		if len(rolls) == 6 {
			return 3000
		}
		if len(rolls) == 5 {
			return 2000
		}
		if len(rolls) == 4 {
			return 1000
		}
		if len(rolls) == 3 {
			if rolls[0] == 1 {
				return 300
			}

			return rolls[0] * 100
		}
	}

	sum := 0
	for _, roll := range rolls {
		if roll == 5 {
			sum += 50
		}
		if roll == 1 {
			sum += 100
		}
	}

	return sum
}

// check if all rolls in a partition are the same
func allEqual(rolls []int) bool {
	for i := 0; i < len(rolls)-1; i++ {
		if rolls[i] != rolls[i+1] {
			return false
		}
	}
	return true
}
