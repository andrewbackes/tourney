package models

// Book represents an opening book in bin format.
type Book struct {
	Id       Id     `json:"id"`
	FilePath string `json:"filePath"`
	MaxDepth int    `json:"maxDepth"`
}
