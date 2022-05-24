package otf

import "time"

type StateVersionOutput struct {
	id        string
	createdAt time.Time
	Name      string
	Sensitive bool
	Type      string
	Value     string
	// StateVersionOutput belongs to StateVersion
	StateVersionID string
}

func (svo *StateVersionOutput) ID() string     { return svo.id }
func (svo *StateVersionOutput) String() string { return svo.id }

type StateVersionOutputList []*StateVersionOutput
