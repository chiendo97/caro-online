package game

type Player int

const (
	EPlayer Player = iota
	XPlayer
	OPlayer
)

func (p Player) String() string {
	return [...]string{"_", "X", "O"}[p]
}

func (p Player) swi() Player {
	switch p {
	case XPlayer:
		return OPlayer
	case OPlayer:
		return XPlayer
	}
	return EPlayer
}
