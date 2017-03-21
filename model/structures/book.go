package structures

type Book struct {
	FilePath        string `json:"filePath,omitempty"`
	Depth           int    `json:"depth,omitempty"`
	Randomize       bool   `json:"randomize,omitempty"`
	MirrorPositions bool   `json:"mirrorPositions,omitempty"`
	RepeatPositions bool   `json:"repeatPositions,omitempty"`
}
