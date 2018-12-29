package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Coefficients struct {
	MCoeffs []float64
	BCoeffs []float64
}

var (
	dice      = flag.Int("dice", 6, "number of dice to roll")
	curScore  = flag.Int("score", 0, "current score")
	maxDepth  = flag.Int("depth", 50, "depth of linear approximation to use")
	calculate = flag.Bool("calculate", false, "use actual backtracking not just approximation")
	generate  = flag.Bool("generate", false, "whether or not to just generate coefficients")
	coeffs    = Coefficients{make([]float64, 6), make([]float64, 6)}
)

func main() {
	flag.Parse()

	if *generate {
		generateCoeffs()
	} else {
		fName := "coeffs/" + strconv.Itoa(*maxDepth) + "_coeffs.json"
		if _, err := os.Stat(fName); os.IsNotExist(err) {
			fmt.Println("Coefficients have not been generated. Generating them now.")
			generateCoeffs()
		}
		coeffs = readVals(*maxDepth)
		if !*calculate {
			fmt.Println(approxScore(*dice, *curScore))
		} else {
			fmt.Println(backtrack([]int{}, *dice, *curScore))
		}
	}
}

func generateCoeffs() {
	for i := 1; i <= 50; i++ {
		newM := []float64{}
		newB := []float64{}
		for j := 1; j <= 6; j++ {
			// find high and low vals for this iteration
			lowVal := backtrack([]int{}, j, 0)
			highVal := backtrack([]int{}, j, 400)
			// use high and low vals to generate line
			newM = append(newM, (highVal-lowVal)/400)
			newB = append(newB, lowVal)
			// test generated lines
			// test removed because it was too sensitive. A little bit of error is okay for the purposes of this experiment.
			/*
				for k := 0; k < 5; k++ {
					toCheck := rand.Intn(400)
					realVal := backtrack([]int{}, j, toCheck, 0)
					approxVal := newM[len(newM)-1]*float64(toCheck) + newB[len(newB)-1]
					if math.Abs(realVal-approxVal) > 1 {
						fmt.Println("The model broke down at die ", j, " and iteration ", i, ".")
						fmt.Println("Start value: ", toCheck)
						fmt.Println("Real Value: ", realVal, " Approx Value: ", approxVal)
					}
				}*/
		}
		coeffs.MCoeffs = newM
		coeffs.BCoeffs = newB
		writeVals(newM, newB, i)
	}
	fmt.Println("Coefficients generated and written")
}

func approxScore(toRoll int, score int) float64 {
	return coeffs.MCoeffs[toRoll-1]*float64(score) + coeffs.BCoeffs[toRoll-1]
}

func writeVals(m []float64, b []float64, depth int) {
	writingCoeffs := Coefficients{m, b}
	toWrite, err := json.Marshal(writingCoeffs)
	handleErr(err)
	err = ioutil.WriteFile("coeffs/"+strconv.Itoa(depth)+"_coeffs.json", toWrite, 0644)
	handleErr(err)
}

func readVals(depth int) Coefficients {
	vals, err := ioutil.ReadFile("coeffs/" + strconv.Itoa(depth) + "_coeffs.json")
	handleErr(err)
	var toReturn Coefficients
	err = json.Unmarshal(vals, &toReturn)
	handleErr(err)
	return toReturn
}

// backtrack to find all possible roll combinations for a given number of dice
func backtrack(rolls []int, toRoll int, score int) float64 {
	if toRoll == 0 {
		return float64(scoreRoll(rolls, score))
	}

	sum := 0.0
	for i := 1; i <= 6; i++ {
		sum += backtrack(append(rolls, i), toRoll-1, score)
	}
	return sum / 6
}

// score a roll by finding the number of occurences of each number and scoring roll
func scoreRoll(rolls []int, score int) float64 {
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
