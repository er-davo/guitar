# 🎸 Guitar Tab Library for Go
[![Tests](https://github.com/er-davo/guitar/actions/workflows/go.yaml/badge.svg)](https://github.com/er-davo/guitar/actions/workflows/go.yaml)

A Go library for generating guitar tablature (ASCII tabs), analyzing fretboard positions, and working with chords, tunings, and playing techniques such as slides and hammer-ons.

## Installation
```bash
go get github.com/er-davo/guitar
```
## Features
- Tab Generation: Build ASCII tabs from notes/chords or sequences.
- Fingering Optimization: Choose the most ergonomic way to play a sequence.
- Tuning Support: Standard, Drop D, and custom tunings.
- Advanced Techniques: Slides (5/7), hammer-ons (2h4), pull-offs(5p3). Harmonics (<12>).
- Note Search: Find multiple fretboard positions for a given note.
- Chord Parsing: Convert simple string-based chords into tab-accurate structures.

# Quick start
## 1. Generate tab
```go
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
				MidiPitch: guitar.NoteToMidi("C", 4),
				Fret:      5,
				String:    3,
				Time:      0.5,
			},
			NoteTo: guitar.Note{
				MidiPitch: guitar.NoteToMidi("D", 4),
				Fret:      7,
				String:    3,
				Time:      0.5,
			},
		},
	})

	fmt.Println(tab.Tab())
}
```
Output
```
e|0----
B|1----
G|2-5/7
D|2----
A|0----
E|-----
```
## 2. Find all positions on the Fretboard
```go
tuning, _ := guitar.ParseTuning(guitar.StandardTuning)
fb, _ := guitar.NewFingerBoard(tuning, 24) // 24-fret board
target := guitar.Note{MidiPitch: guitar.NoteToMidi("C#", 3)}

for _, p := range fb.Find(target) {
	n := p.(guitar.Note)
	name, oct := guitar.MidiToNote(n.MidiPitch)
	fmt.Printf("String %d Fret %d (MIDI %d -> %s%d)\n", n.String, n.Fret, n.MidiPitch, name, oct)
}
```
Output
```
String 4 Fret 4 (MIDI 49 -> C#3)
String 5 Fret 9 (MIDI 49 -> C#3)
```
## 3. Parse Custom Chords
```go
chord := guitar.ParseChord("0 2 2 2 0 -") // A major chord
// Returns tab struture with frets [0, 2, 2, 2, 0] (for each string)
```
## 4. Optimize a sequence using FingeringOptimizer
```go
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

fb, _ := guitar.NewFingerBoard(tuning, 24)
opt := guitar.NewFingeringOptimizer(*fb)

layers, _ := opt.TimeLayers(sequence)
path, _ := opt.OptimizePath(layers)

tab, _ := guitar.NewTabWriter(tuning.NoteNames())
_ = tab.Write(path...)

fmt.Println(tab.Tab())
```
Output
```
e|-------
B|--1/2-5
G|0------
D|2-4----
A|3-5----
E|--5----
```
## Planned Features

This library is actively developed. Future updates may include:

- 📄 **MusicXML Export**: Generate MusicXML representations of tabs for compatibility with notation software.