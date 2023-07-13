package anneal

import (
	"math"
	"math/rand"
)

type State interface {
	Energy() float64
	Neighbor() State
	PrintMovement()
}

type Config struct {
	Iteration       uint
	AneallingFactor float64
	Temperature     float64
}

func Init(currState State, conf Config) State {
	temperature := conf.Temperature
	currEnergy := currState.Energy()

	// Anonymous function to update current state to new state
	updateState := func(s State, e float64) {
		currState = s
		currEnergy = e
	}

	for i := 0; i < int(conf.Iteration); i++ {
		// Generate neighbor state
		neighbor := currState.Neighbor()
		neighborEnergy := neighbor.Energy()

		// Evaluate neighbor solution
		if neighborEnergy < currEnergy {
			// Update if neighbor is better than current
			updateState(neighbor, neighborEnergy)
		} else if math.Exp((currEnergy-neighborEnergy)/temperature) > rand.Float64() {
			// Update if the acceptance probability is higher than the random number
			updateState(neighbor, neighborEnergy)
		}

		// Anneal the temperature (cooling down)
		temperature = temperature * conf.AneallingFactor
	}

	return currState
}
