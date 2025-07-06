package guitar

import (
	"fmt"
	"math"
	"slices"
)

const (
	maxFretSpan = 8
)

type Fingering []Playable

func (f *Fingering) TransitionCostTo(to Fingering) float64 {
	// maybe add hungarian algorithm ?

	usedTo := make([]bool, len(to))
	totalCost := 0.0
	missingPenalty := 3.0

	for _, fromNote := range *f {
		bestScore := math.MaxFloat64
		bestIdx := -1

		for j, toNote := range to {
			if usedTo[j] {
				continue
			}
			score := fromNote.ScoreTo(toNote)
			if score < bestScore {
				bestScore = score
				bestIdx = j
			}
		}

		if bestIdx >= 0 {
			totalCost += bestScore
			usedTo[bestIdx] = true
		} else {
			totalCost += missingPenalty
		}
	}

	// штраф за toNote, которые вообще не использовались
	for _, used := range usedTo {
		if !used {
			totalCost += missingPenalty
		}
	}

	return totalCost
}

func (f *Fingering) NewTransitionCostTo(to Fingering) float64 {
	cmp := func(left, right Playable) int {
		if left.StringNumber() < right.StringNumber() {
			return -1
		} else if left.StringNumber() > right.StringNumber() {
			return 1
		}
		return 0
	}
	slices.SortStableFunc(*f, cmp)
	slices.SortStableFunc(to, cmp)
	// toOpen := []Playable{}
	// toPinch := []Playable{}

	totalCost := 0.0

	return totalCost
}

func (f *Fingering) hasString(stringNum int) bool {
	for _, note := range *f {
		if note.StringNumber() == stringNum {
			return true
		}
	}

	return false
}

func (f Fingering) FretSpan() int {
	if len(f) == 0 {
		return 0
	}
	minFret := f[0].FretPosition()
	maxFret := f[0].FretPosition()

	for _, n := range f {
		if n.FretPosition() < minFret {
			minFret = n.FretPosition()
		}
		if n.FretPosition() > maxFret {
			maxFret = n.FretPosition()
		}
	}
	return maxFret - minFret
}

type TimeLayer []Fingering

type nodeInfo struct {
	cost      float64
	prevIndex int
}

type fingeringOptimizer struct {
	fb FingerBoard
}

func NewFingeringOptimizer(fb FingerBoard) fingeringOptimizer {
	return fingeringOptimizer{fb: fb}
}

func (opt *fingeringOptimizer) TimeLayer(notes []Playable) (TimeLayer, error) {
	if len(notes) > len(opt.fb.tuning) {
		return TimeLayer{}, fmt.Errorf("to many notes (%d) for fingerboard with %d strings", len(notes), len(opt.fb.tuning))
	}

	if err := validate(&notes); err != nil {
		return TimeLayer{}, err
	}

	tl := TimeLayer{}

	possibleNotes := make([][]Playable, len(notes))

	for i, note := range notes {
		possibleNotes[i] = opt.fb.Find(note)
	}

	var generate func(depth int, current Fingering)
	unplayables := []Fingering{}

	generate = func(depth int, current Fingering) {
		if depth == len(possibleNotes) {
			if current.FretSpan() > maxFretSpan {
				unplayables = append(unplayables, current)
				return
			}

			combo := make(Fingering, len(current))
			copy(combo, current)
			tl = append(tl, combo)

			return
		}

		for _, cand := range possibleNotes[depth] {
			if current.hasString(cand.StringNumber()) {
				continue
			}

			generate(depth+1, append(current, cand))
		}
	}

	generate(0, Fingering{})

	if len(tl) == 0 {
		bestFrmWorst := Fingering{Note{Time: notes[0].StartTime()}}
		minSpan := 1000000

		for _, worst := range unplayables {
			if worst.FretSpan() < minSpan {
				minSpan = worst.FretSpan()
				bestFrmWorst = worst
			}
		}

		tl = append(tl, bestFrmWorst)
	}

	return tl, nil
}

func (opt *fingeringOptimizer) TimeLayers(sequence [][]Playable) ([]TimeLayer, error) {
	layers := make([]TimeLayer, 0, len(sequence))
	for _, events := range sequence {
		layer, err := opt.TimeLayer(events)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	return layers, nil
}

func (opt *fingeringOptimizer) OptimizePath(layers []TimeLayer) ([]TabFrame, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("no layers provided")
	}

	dp := make([][]nodeInfo, len(layers))

	dp[0] = make([]nodeInfo, len(layers[0]))
	for j := range dp[0] {
		dp[0][j].cost = 0
		dp[0][j].prevIndex = -1
	}

	const inf = 1e9

	for i := 1; i < len(layers); i++ {
		dp[i] = make([]nodeInfo, len(layers[i]))
		for j := range dp[i] {
			dp[i][j].cost = inf
			dp[i][j].prevIndex = -1
		}
	}

	for i := 1; i < len(layers); i++ {
		for toIndex, to := range layers[i] {
			minCost := inf
			prev := -1

			for fromIdx, from := range layers[i-1] {
				cost := dp[i-1][fromIdx].cost + from.TransitionCostTo(to)
				if cost < minCost {
					minCost = cost
					prev = fromIdx
				}
			}

			dp[i][toIndex].cost = minCost
			dp[i][toIndex].prevIndex = prev
		}
	}

	minPathIndex := 0
	for i := 0; i < len(dp[len(dp)-1]); i++ {
		if dp[len(dp)-1][minPathIndex].cost > dp[len(dp)-1][i].cost {
			minPathIndex = i
		}
	}

	optimizedPath := make([]TabFrame, len(dp))

	forwardIndex := minPathIndex
	for i := len(optimizedPath) - 1; i >= 0; i-- {
		if forwardIndex < 0 {
			return nil, fmt.Errorf("invalid path: disconnected layer at %d", i)
		}
		optimizedPath[i] = TabFrame(layers[i][forwardIndex])
		forwardIndex = dp[i][forwardIndex].prevIndex
	}

	return optimizedPath, nil
}

func validate(events *[]Playable) error {
	for i := range *events {
		switch v := (*events)[i].(type) {
		case Note:
			return v.Validate()
		case Harmonic:
			return v.Validate()
		case Slide:
			if err := v.NoteFrom.Validate(); err != nil {
				return err
			}
			if err := v.NoteTo.Validate(); err != nil {
				return err
			}
		case HammerOn:
			if err := v.NoteFrom.Validate(); err != nil {
				return err
			}
			if err := v.NoteTo.Validate(); err != nil {
				return err
			}
		case PullOff:
			if err := v.NoteFrom.Validate(); err != nil {
				return err
			}
			if err := v.NoteTo.Validate(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type in Playable interface %t", (*events)[i])
		}
	}

	return nil
}
