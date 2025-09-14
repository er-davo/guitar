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

	// Add a slide from fret 5 to 7 on the G string (G3 open -> fret5 = C4, fret7 = D4)
	_ = tab.Write(guitar.TabFrame{
		guitar.Slide{
			NoteFrom: guitar.Note{
				MidiPitch: guitar.NoteToMidi("C", 4), // G3 + 5 => C4
				Fret:      5,
				String:    3,
				Time:      0.5,
			},
			NoteTo: guitar.Note{
				MidiPitch: guitar.NoteToMidi("D", 4), // G3 + 7 => D4
				Fret:      7,
				String:    3,
				Time:      0.5,
			},
		},
	})

	fmt.Println("=== A minor + slide ===")
	fmt.Println(tab.Tab())

	// --------------------------------------
	// Fingerboard lookup example
	fb, _ := guitar.NewFingerBoard(tuning, 24)
	target := guitar.Note{MidiPitch: guitar.NoteToMidi("C#", 3)}

	fmt.Println("=== Positions for C#3 ===")
	for _, p := range fb.Find(target) {
		n := p.(guitar.Note)
		name, oct := guitar.MidiToNote(n.MidiPitch)
		fmt.Printf("String %d Fret %d (MIDI %d -> %s%d)\n", n.String, n.Fret, n.MidiPitch, name, oct)
	}

	// --------------------------------------
	// Optimized tab generation from chords and techniques
	fmt.Println("=== OptimizePath: Chords + Slide ===")

	sequence := [][]guitar.Playable{
		{
			guitar.Note{MidiPitch: guitar.NoteToMidi("C", 3), Time: 0},
			guitar.Note{MidiPitch: guitar.NoteToMidi("E", 3), Time: 0},
			guitar.Note{MidiPitch: guitar.NoteToMidi("G", 3), Time: 0},
		},
		{
			guitar.Note{MidiPitch: guitar.NoteToMidi("D", 3), Time: 0.5},
			guitar.Note{MidiPitch: guitar.NoteToMidi("F#", 3), Time: 0.5},
			guitar.Note{MidiPitch: guitar.NoteToMidi("A", 2), Time: 0.5},
			guitar.Slide{
				NoteFrom: guitar.Note{MidiPitch: guitar.NoteToMidi("C", 4), Time: 0.5},
				NoteTo:   guitar.Note{MidiPitch: guitar.NoteToMidi("C#", 4), Time: 0.5},
			},
		},
		{
			guitar.Note{MidiPitch: guitar.NoteToMidi("E", 4), Time: 1.0},
		},
	}

	opt := guitar.NewFingeringOptimizer(fb)
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
		{guitar.Note{MidiPitch: guitar.NoteToMidi("F", 2), Time: 0.0}},
		{guitar.Note{MidiPitch: guitar.NoteToMidi("C", 3), Time: 0.5}},
		{guitar.Note{MidiPitch: guitar.NoteToMidi("F", 3), Time: 1.0}},
		{guitar.Note{MidiPitch: guitar.NoteToMidi("A", 3), Time: 1.5}},
		{guitar.Note{MidiPitch: guitar.NoteToMidi("C", 4), Time: 2.0}},
		{guitar.Note{MidiPitch: guitar.NoteToMidi("F", 4), Time: 2.5}},
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
