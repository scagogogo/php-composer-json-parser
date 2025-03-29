// Package archive provides functionality related to PHP Composer archiving
package archive

// Archive defines how the package should be archived
type Archive struct {
	Exclude []string `json:"exclude,omitempty"`
}

// NewArchive creates a new Archive with commonly excluded paths
func NewArchive() *Archive {
	return &Archive{
		Exclude: []string{
			"/.*",
			"/*.md",
			"/composer.json",
			"/composer.lock",
			"/vendor",
			"/tests",
			"/test",
			"/docs",
			"/doc",
		},
	}
}

// AddExclusion adds a path pattern to the exclude list
func AddExclusion(a *Archive, pattern string) {
	// Check if the pattern already exists
	for _, p := range a.Exclude {
		if p == pattern {
			return
		}
	}
	a.Exclude = append(a.Exclude, pattern)
}

// RemoveExclusion removes a path pattern from the exclude list
func RemoveExclusion(a *Archive, pattern string) bool {
	for i, p := range a.Exclude {
		if p == pattern {
			// Remove the element at index i
			a.Exclude = append(a.Exclude[:i], a.Exclude[i+1:]...)
			return true
		}
	}
	return false
}
