package types

type NukeResponseCode int

const (
	OK            NukeResponseCode = 0
	EMPTY         NukeResponseCode = 100
	NOT_FOUND     NukeResponseCode = -900
	DUPLICATE_KEY NukeResponseCode = -200
)
