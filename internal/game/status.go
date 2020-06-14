package game

type Status int

const (
	Running Status = iota
	XWin
	OWin
	Tie
)

func (s Status) String() string {
	return [...]string{"Running", "XWin", "OWin", "Tie"}[s]
}

func (s Status) GetPlayer() Player {
	switch s {
	case XWin:
		return XPlayer
	case OWin:
		return OPlayer
	default:
		return EPlayer
	}
}
