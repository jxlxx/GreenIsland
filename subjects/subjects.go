package subjects

type Subject string

const (
	WorldTick Subject = "world.time.tick"
)

func (s Subject) String() string {
	return string(s)
}
