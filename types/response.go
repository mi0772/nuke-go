package types

type NukeResponseCode int

const (
	Ok           NukeResponseCode = 0
	Empty        NukeResponseCode = 100
	NotFound     NukeResponseCode = -900
	DuplicateKey NukeResponseCode = -200
)
