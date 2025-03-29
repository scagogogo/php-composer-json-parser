// Package repository provides functionality related to PHP Composer repositories
package repository

// Repository defines a package repository
type Repository struct {
	Type    string                 `json:"type,omitempty"`
	URL     string                 `json:"url,omitempty"`
	Package map[string]interface{} `json:"package,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// NewRepository creates a new repository with the given type and URL
func NewRepository(repoType, url string) *Repository {
	return &Repository{
		Type:    repoType,
		URL:     url,
		Package: make(map[string]interface{}),
		Options: make(map[string]interface{}),
	}
}

// IsVCS returns true if the repository is a VCS type
func IsVCS(r *Repository) bool {
	return r.Type == "git" || r.Type == "svn" || r.Type == "hg"
}

// IsPackagist returns true if the repository is packagist.org
func IsPackagist(r *Repository) bool {
	return r.Type == "composer" &&
		(r.URL == "https://repo.packagist.org" || r.URL == "https://packagist.org")
}
