package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// The number of experiments to run
	experiments = 1000
	// The number of prisoners and boxes. Note that scalability is pretty shit, so anything over 2K will take a lot of time.
	prisoners = 100
)

func main() {
	var e experiment
	rand.Seed(time.Now().UTC().UnixNano())
	nFreedoms := 0
	for i := 0; i < experiments; i++ {
		e = initialize()
		wins := naiveSolve(e)
		if wins == prisoners {
			nFreedoms++
		}
	}
	fmt.Printf("Prisoners got their freedom, using the naive approach: %d / %d (%.2f%%)\n",
		nFreedoms, experiments, float64(100*nFreedoms)/float64(experiments))
	sFreedoms := 0
	for i := 0; i < experiments; i++ {
		e = initialize()
		wins := smartSolve(e)
		// fmt.Printf("wins: %d\n", wins)
		if wins == prisoners {
			sFreedoms++
		}
	}
	fmt.Printf("Prisoners got their freedom, using the linked approach: %d / %d (%.2f%%)\n",
		sFreedoms, experiments, float64(100*sFreedoms)/float64(experiments))

}

// newPool creates a new pool of random numbers. It's a slice of numbers and each number is unique
// and is not repeated.
func newPool() []int {
	numbers := make([]int, prisoners)
	pool := make([]int, prisoners)
	for i := 0; i < prisoners; i++ {
		numbers[i] = i
	}
	for i := 0; i < prisoners; i++ {
		idx := rand.Intn(len(numbers))
		num := numbers[idx]
		pool[i] = num
		// delete the number from the list
		numbers = append(numbers[:idx], numbers[idx+1:]...)
	}
	return pool
}

type experiment struct {
	boxes []box
}

// naiveSolve tries to see if the prisoner will succeed by grabbing 50 random
// boxes. Returns true if the prisoner succeeds.
func naiveSolve(exp experiment) int {
	const boxes = prisoners / 2
	successes := 0
	// iterate over all the prisoners
	for pid := 0; pid < prisoners; pid++ {
		e := experiment{
			boxes: make([]box, prisoners),
		}
		copy(e.boxes, exp.boxes)
	boxLoop:
		// for each prisoner iterate over all the boxes
		for i := 0; i < boxes; i++ {
			if len(e.boxes) == 0 {
				break boxLoop
			}
			boxIdx := rand.Intn(len(e.boxes)) // pick a random box
			box := e.boxes[boxIdx]
			// remove the box from the experiment
			e.boxes = append(e.boxes[:boxIdx], e.boxes[boxIdx+1:]...)
			if box.prisoner == pid {
				successes++
				break boxLoop
			}
		}

	}
	return successes
}

// smartSolve tries to be smart. It implements the following strategy:
// Each prisoner starts with the box matching their id.
// the next number is the number of the box that the prisoner is holding.
// Repeat this until the prisoner has found their ID or opened half the boxes.
func smartSolve(e experiment) int {
	const boxes = prisoners / 2
	successes := 0
	// iterate over all the prisoners
	for pid := 0; pid < prisoners; pid++ {
		i := 0
		curBoxIdx := pid
		curBox := e.boxes[curBoxIdx]
	innerLoop:
		for {
			if i > boxes {
				break innerLoop
			}
			if curBox.prisoner == pid {
				successes++
				break innerLoop
			}
			curBoxIdx = curBox.prisoner
			curBox = e.boxes[curBoxIdx]
			i++
		}
	}
	return successes
}

type box struct {
	prisoner int
}

func initialize() experiment {
	pool := newPool()
	e := experiment{
		boxes: make([]box, prisoners),
	}
	for i := 0; i < prisoners; i++ {
		box := box{
			prisoner: pool[i],
		}
		e.boxes[i] = box
	}
	return e
}
