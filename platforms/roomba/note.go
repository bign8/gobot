package roomba

// TODO: Generate constats for the MusicScale supported by robot

// Note is an internal structure to store music.
type Note struct {
	Number   uint8 `json:"number"`
	Duration uint8 `json:"duration"`
}
