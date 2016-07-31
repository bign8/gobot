package roomba

// Note is an internal structure to store music.
type Note struct {
	Number   uint8 `json:"n"`
	Duration uint8 `json:"d"`
}

// Valid does a quick check to make sure the note is valid
func (n Note) Valid() bool {
	return n.Number > 30 && n.Number < 128 // Valid: [31, 127]
}
