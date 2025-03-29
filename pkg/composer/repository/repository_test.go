package repository

import (
	"testing"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name     string
		repoType string
		url      string
		want     *Repository
	}{
		{
			name:     "Create git repository",
			repoType: "git",
			url:      "https://github.com/example/repo.git",
			want: &Repository{
				Type:    "git",
				URL:     "https://github.com/example/repo.git",
				Package: map[string]interface{}{},
				Options: map[string]interface{}{},
			},
		},
		{
			name:     "Create composer repository",
			repoType: "composer",
			url:      "https://packagist.org",
			want: &Repository{
				Type:    "composer",
				URL:     "https://packagist.org",
				Package: map[string]interface{}{},
				Options: map[string]interface{}{},
			},
		},
		{
			name:     "Create with empty values",
			repoType: "",
			url:      "",
			want: &Repository{
				Type:    "",
				URL:     "",
				Package: map[string]interface{}{},
				Options: map[string]interface{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRepository(tt.repoType, tt.url)

			if got == nil {
				t.Errorf("NewRepository() returned nil")
				return
			}

			if got.Type != tt.want.Type {
				t.Errorf("NewRepository().Type = %v, want %v", got.Type, tt.want.Type)
			}

			if got.URL != tt.want.URL {
				t.Errorf("NewRepository().URL = %v, want %v", got.URL, tt.want.URL)
			}

			// 确保Package和Options是初始化过的map
			if got.Package == nil {
				t.Errorf("NewRepository().Package is nil")
			} else if len(got.Package) != 0 {
				t.Errorf("NewRepository().Package should be empty, got %v", got.Package)
			}

			if got.Options == nil {
				t.Errorf("NewRepository().Options is nil")
			} else if len(got.Options) != 0 {
				t.Errorf("NewRepository().Options should be empty, got %v", got.Options)
			}
		})
	}
}

func TestIsVCS(t *testing.T) {
	tests := []struct {
		name string
		repo *Repository
		want bool
	}{
		{
			name: "Git repository",
			repo: &Repository{Type: "git", URL: "https://github.com/example/repo.git"},
			want: true,
		},
		{
			name: "SVN repository",
			repo: &Repository{Type: "svn", URL: "https://svn.example.com/repo"},
			want: true,
		},
		{
			name: "Mercurial repository",
			repo: &Repository{Type: "hg", URL: "https://hg.example.com/repo"},
			want: true,
		},
		{
			name: "Composer repository",
			repo: &Repository{Type: "composer", URL: "https://packagist.org"},
			want: false,
		},
		{
			name: "Path repository",
			repo: &Repository{Type: "path", URL: "../local/path"},
			want: false,
		},
		{
			name: "Empty type",
			repo: &Repository{Type: "", URL: "https://example.com"},
			want: false,
		},
		{
			name: "Nil repository",
			repo: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repo == nil && tt.want == true {
				t.Fatalf("Test case error: Cannot expect true for nil repository")
			}

			got := false
			if tt.repo != nil {
				got = IsVCS(tt.repo)
			}

			if got != tt.want {
				t.Errorf("IsVCS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPackagist(t *testing.T) {
	tests := []struct {
		name string
		repo *Repository
		want bool
	}{
		{
			name: "Official Packagist with repo.packagist.org",
			repo: &Repository{Type: "composer", URL: "https://repo.packagist.org"},
			want: true,
		},
		{
			name: "Official Packagist with packagist.org",
			repo: &Repository{Type: "composer", URL: "https://packagist.org"},
			want: true,
		},
		{
			name: "Custom Packagist",
			repo: &Repository{Type: "composer", URL: "https://custom.packagist.org"},
			want: false,
		},
		{
			name: "Packagist with wrong type",
			repo: &Repository{Type: "git", URL: "https://packagist.org"},
			want: false,
		},
		{
			name: "Git repository",
			repo: &Repository{Type: "git", URL: "https://github.com/example/repo.git"},
			want: false,
		},
		{
			name: "Empty repository",
			repo: &Repository{Type: "", URL: ""},
			want: false,
		},
		{
			name: "Nil repository",
			repo: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.repo == nil && tt.want == true {
				t.Fatalf("Test case error: Cannot expect true for nil repository")
			}

			got := false
			if tt.repo != nil {
				got = IsPackagist(tt.repo)
			}

			if got != tt.want {
				t.Errorf("IsPackagist() = %v, want %v", got, tt.want)
			}
		})
	}
}
