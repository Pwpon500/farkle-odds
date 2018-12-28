package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
)

type Coefficients struct {
	MCoeffs []float64
	BCoeffs []float64
}

var (
	dice     = flag.Int("dice", 6, "number of dice to roll")
	curScore = flag.Int("score", 0, "current score")
	maxDepth = flag.Int("depth", 1, "maximum depth for more rolls")
	//mVals    = []float64{0.3333333333333333, 0.5555555555555556, 0.7222222222222222, 0.8425925925925926, 0.9228395061728395, 0.9768518518518519}
	//bVals    = []float64{25, 50, 83.56481481481481, 132.71604938271605, 203.2857510288066, 388.5352366255144}
	mVals = []float64{0.26891084419308525, 0.3863731313184077, 0.5395296126895001, 0.6460222781185294, 0.7322907852593873, 0.8067325325792554}
	bVals = []float64{187.7691217568772, 161.72911616661224, 177.01463179559335, 224.53924356873253, 303.21985133373, 502.80242532718756}
)

func main() {
	flag.Parse()

	//fmt.Println(backtrack([]int{}, *dice, *curScore, 0))
	fmt.Println(approxScore(*dice, *curScore))
	//writeVals(mVals, bVals, 1)

	//generateCoeffs()

	fmt.Println("done")
}

func generateCoeffs() {
	for i := 2; i <= 50; i++ {
		newM := []float64{}
		newB := []float64{}
		for j := 1; j <= 6; j++ {
			// find high and low vals for this iteration
			lowVal := backtrack([]int{}, j, 0, 0)
			highVal := backtrack([]int{}, j, 400, 0)
			// use high and low vals to generate line
			newM = append(newM, (highVal-lowVal)/400)
			newB = append(newB, lowVal)
			// test generated line
			for k := 0; k < 5; k++ {
				toCheck := rand.Intn(400)
				realVal := backtrack([]int{}, j, toCheck, 0)
				approxVal := newM[len(newM)-1]*float64(toCheck) + newB[len(newB)-1]
				if math.Abs(realVal-approxVal) > .00001 {
					fmt.Println("The model broke down at die ", j, " and iteration ", i, ".")
					fmt.Println("Real Value: ", realVal, " Approx Value: ", approxVal)
				}
			}
		}
		mVals = newM
		bVals = newB
		writeVals(mVals, bVals, i)
		fmt.Println("Finished iteration ", i, ".")
	}
}

func approxScore(toRoll int, score int) float64 {
	return mVals[toRoll-1]*float64(score) + bVals[toRoll-1]
}

func writeVals(m []float64, b []float64, depth int) {
	coeffs := Coefficients{m, b}
	toWrite, err := json.Marshal(coeffs)
	handleErr(err)
	err = ioutil.WriteFile(strconv.Itoa(depth)+"_coeffs.json", toWrite, 0644)
	handleErr(err)
}

// backtrack to find all possible roll combinations for a given number of dice
func backtrack(rolls []int, toRoll int, score int, depth int) float64 {
	if toRoll == 0 {
		return float64(scoreRoll(rolls, score, depth))
	}
	if depth >= *maxDepth {
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
	expectedRoll := approxScore(numLeft, int(toReturn))
	//expectedRoll := backtrack([]int{}, numLeft, int(toReturn), depth+1)
	if expectedRoll > float64(toReturn) {
		toReturn = expectedRoll
	}

	return toReturn
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
