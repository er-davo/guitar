package guitar

import (
	"errors"
)

type FingerBoard struct {
	tuning Tuning
	frets  int
	cache  map[int][]Playable
}

func NewFingerBoard(tun Tuning, frets int) (*FingerBoard, error) {
	if frets < 0 {
		return nil, errors.New("frets value can not be negative")
	}

	return &FingerBoard{
		tuning: tun,
		frets:  frets,
		cache:  make(map[int][]Playable, frets*len(tun)),
	}, nil
}

func (fb *FingerBoard) GetTuningNotes() []string {
	return fb.tuning.NoteNames()
}

func (fb *FingerBoard) Find(target Playable) []Playable {
	switch v := target.(type) {
	case Note:
		return fb.FindNotes(v)
	case Harmonic:
		return fb.FindHarmonics(v)
	case Slide:
		return fb.FindSlides(v)
	case HammerOn:
		return fb.FindHammerOns(v)
	case PullOff:
		return fb.FindPullOffs(v)
	default:
		return nil
	}
}

func (fb *FingerBoard) FindNotes(target Note) []Playable {
	if notes, ok := fb.cache[target.MidiPitch]; ok {
		return notes
	}

	notes := []Playable{}
	currentNote := Note{}

	for i := range fb.tuning {
		currentNote = fb.tuning[i]

		for fret := 0; fret < fb.frets; fret++ {
			if currentNote.MidiPitch == target.MidiPitch {
				currentNote.Time = target.Time
				notes = append(notes, currentNote)
			}
			currentNote.AddFret()
		}
	}

	fb.cache[target.MidiPitch] = notes

	return notes
}

func (fb *FingerBoard) FindHarmonics(target Harmonic) []Playable {
	notes := fb.FindNotes(target.Note)
	harmonics := make([]Playable, len(notes))

	for i, note := range notes {
		harmonics[i] = Harmonic{Note: note.(Note)}
	}

	return harmonics
}

func (fb *FingerBoard) FindSlides(target Slide) []Playable {
	slides := []Playable{}

	start := fb.FindNotes(target.NoteFrom)
	end := fb.FindNotes(target.NoteTo)

	for _, from := range start {
		for _, to := range end {
			if from.StringNumber() == to.StringNumber() &&
				from.FretPosition() != to.FretPosition() {
				slides = append(slides, Slide{
					NoteFrom: from.(Note),
					NoteTo:   to.(Note),
				})
			}
		}
	}

	return slides
}

func (fb *FingerBoard) FindHammerOns(target HammerOn) []Playable {
	hammerons := []Playable{}

	start := fb.FindNotes(target.NoteFrom)
	end := fb.FindNotes(target.NoteTo)

	for _, from := range start {
		for _, to := range end {
			if from.StringNumber() == to.StringNumber() &&
				from.FretPosition() != to.FretPosition() {
				hammerons = append(hammerons, HammerOn{
					NoteFrom: from.(Note),
					NoteTo:   to.(Note),
				})
			}
		}
	}

	return hammerons
}

func (fb *FingerBoard) FindPullOffs(target PullOff) []Playable {
	pulloffs := []Playable{}

	start := fb.FindNotes(target.NoteFrom)
	end := fb.FindNotes(target.NoteTo)

	for _, from := range start {
		for _, to := range end {
			if from.StringNumber() == to.StringNumber() &&
				from.FretPosition() != to.FretPosition() {
				pulloffs = append(pulloffs, PullOff{
					NoteFrom: from.(Note),
					NoteTo:   to.(Note),
				})
			}
		}
	}

	return pulloffs
}
