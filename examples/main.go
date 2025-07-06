package main

import (
	"fmt"

	"github.com/er-davo/guitar"
)

func main() {
	// Initialize a tab writer with standard tuning
	tuning, _ := guitar.ParseTuning(guitar.StandardTuning)
	tab, _ := guitar.NewTabWriter(tuning.NoteNames())

	// Add an A minor chord (open position)
	_ = tab.Write(guitar.ParseChord("0 1 2 2 0 -", 0))

	// Add a slide from fret 5 to 7 on the G string
	_ = tab.Write(guitar.TabFrame{
		guitar.Slide{
			NoteFrom: guitar.Note{Fret: 5, String: 2, Time: 0.5},
			NoteTo:   guitar.Note{Fret: 7, String: 2, Time: 0.5},
		},
	})

	fmt.Println("=== A minor + slide ===")
	fmt.Println(tab.Tab())

	// --------------------------------------
	// Fingerboard lookup example
	fb, _ := guitar.NewFingerBoard(tuning, 24)
	fmt.Println("=== Positions for C#3 ===")
	for _, p := range fb.Find(guitar.Note{Name: "C#", Octave: 3}) {
		fmt.Println(p)
	}

	// --------------------------------------
	// Optimized tab generation from chords and techniques
	fmt.Println("=== OptimizePath: Chords + Slide ===")

	sequence := [][]guitar.Playable{
		{
			guitar.Note{Name: "C", Octave: 3, Time: 0},
			guitar.Note{Name: "E", Octave: 3, Time: 0},
			guitar.Note{Name: "G", Octave: 3, Time: 0},
		},
		{
			guitar.Note{Name: "D", Octave: 3, Time: 0.5},
			guitar.Note{Name: "F#", Octave: 3, Time: 0.5},
			guitar.Note{Name: "A", Octave: 2, Time: 0.5},
			guitar.Slide{
				NoteFrom: guitar.Note{Name: "C", Octave: 4, Time: 0.5},
				NoteTo:   guitar.Note{Name: "C#", Octave: 4, Time: 0.5},
			},
		},
		{
			guitar.Note{Name: "E", Octave: 4, Time: 1.0},
		},
	}

	opt := guitar.NewFingeringOptimizer(*fb)
	layers1, err := opt.TimeLayers(sequence)
	if err != nil {
		panic(err)
	}

	path1, err := opt.OptimizePath(layers1)
	if err != nil {
		panic(err)
	}

	tab1, _ := guitar.NewTabWriter(tuning.NoteNames())
	_ = tab1.Write(path1...)
	fmt.Println(tab1.Tab())

	// --------------------------------------
	// Second example: Simple melody
	fmt.Println("=== OptimizePath: Melody ===")

	sequence2 := [][]guitar.Playable{
		// {guitar.Note{Name: "E", Octave: 2, Time: 0.0}},
		// {guitar.Note{Name: "G", Octave: 3, Time: 0.5}},
		// {guitar.Note{Name: "B", Octave: 3, Time: 1.0}},
		// {guitar.Note{Name: "E", Octave: 4, Time: 1.5}},
		// {guitar.Note{Name: "B", Octave: 3, Time: 2.0}},
		// {guitar.Note{Name: "G", Octave: 3, Time: 2.5}},
		// {guitar.Note{Name: "E", Octave: 2, Time: 3.0}},
		// {
		// 	guitar.Note{Name: "F", Octave: 2, Time: 0.0},
		// 	guitar.Note{Name: "C", Octave: 3, Time: 0.0},
		// 	guitar.Note{Name: "F", Octave: 3, Time: 0.0},
		// 	guitar.Note{Name: "A", Octave: 3, Time: 0.0},
		// 	guitar.Note{Name: "C", Octave: 4, Time: 0.0},
		// 	guitar.Note{Name: "F", Octave: 4, Time: 0.0},
		// },
		// {
		// 	guitar.Note{Name: "G#", Octave: 2, Time: 1.0},
		// 	guitar.Note{Name: "D#", Octave: 3, Time: 1.0},
		// 	guitar.Note{Name: "G#", Octave: 3, Time: 1.0},
		// 	guitar.Note{Name: "C", Octave: 4, Time: 1.0},
		// 	guitar.Note{Name: "D#", Octave: 4, Time: 1.0},
		// 	guitar.Note{Name: "G#", Octave: 4, Time: 1.0},
		// },
		{guitar.Note{Name: "F", Octave: 2, Time: 0.0}},
		{guitar.Note{Name: "C", Octave: 3, Time: 0.5}},
		{guitar.Note{Name: "F", Octave: 3, Time: 1.0}},
		{guitar.Note{Name: "A", Octave: 3, Time: 1.5}},
		{guitar.Note{Name: "C", Octave: 4, Time: 2.0}},
		{guitar.Note{Name: "F", Octave: 4, Time: 2.5}},
	}

	layers2, err := opt.TimeLayers(sequence2)
	if err != nil {
		panic(err)
	}

	path2, err := opt.OptimizePath(layers2)
	if err != nil {
		panic(err)
	}

	tab2, _ := guitar.NewTabWriter(tuning.NoteNames())
	_ = tab2.Write(path2...)
	fmt.Println(tab2.Tab())
}
