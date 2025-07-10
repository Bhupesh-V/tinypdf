package entities

type Image struct {
	Size   int64  // Size in bytes
	Color  string // Color space (e.g., RGB, CMYK)
	Width  int    // Width in pixels
	Height int    // Height in pixels
}
