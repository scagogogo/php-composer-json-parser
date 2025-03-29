package archive

import (
	"reflect"
	"testing"
)

func TestNewArchive(t *testing.T) {
	archive := NewArchive()

	// 确保返回的Archive不是nil
	if archive == nil {
		t.Errorf("NewArchive() returned nil")
	}

	// 确保Exclude字段不是nil
	if archive.Exclude == nil {
		t.Errorf("NewArchive().Exclude is nil")
	}

	// 确保Exclude字段包含了默认的排除模式
	expectedExcludes := []string{
		"/.*",
		"/*.md",
		"/composer.json",
		"/composer.lock",
		"/vendor",
		"/tests",
		"/test",
		"/docs",
		"/doc",
	}

	if !reflect.DeepEqual(archive.Exclude, expectedExcludes) {
		t.Errorf("NewArchive().Exclude = %v, want %v", archive.Exclude, expectedExcludes)
	}
}

func TestAddExclusion(t *testing.T) {
	tests := []struct {
		name         string
		archive      *Archive
		pattern      string
		wantExcludes []string
	}{
		{
			name: "Add new pattern to non-empty list",
			archive: &Archive{
				Exclude: []string{"/existing"},
			},
			pattern:      "/new-pattern",
			wantExcludes: []string{"/existing", "/new-pattern"},
		},
		{
			name: "Add existing pattern",
			archive: &Archive{
				Exclude: []string{"/existing"},
			},
			pattern:      "/existing",
			wantExcludes: []string{"/existing"}, // 不应该重复添加
		},
		{
			name: "Add pattern to empty list",
			archive: &Archive{
				Exclude: []string{},
			},
			pattern:      "/new-pattern",
			wantExcludes: []string{"/new-pattern"},
		},
		{
			name: "Add pattern to nil list",
			archive: &Archive{
				Exclude: nil,
			},
			pattern:      "/new-pattern",
			wantExcludes: []string{"/new-pattern"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddExclusion(tt.archive, tt.pattern)

			if !reflect.DeepEqual(tt.archive.Exclude, tt.wantExcludes) {
				t.Errorf("After AddExclusion(), archive.Exclude = %v, want %v", tt.archive.Exclude, tt.wantExcludes)
			}
		})
	}
}

func TestRemoveExclusion(t *testing.T) {
	tests := []struct {
		name         string
		archive      *Archive
		pattern      string
		want         bool
		wantExcludes []string
	}{
		{
			name: "Remove existing pattern from middle",
			archive: &Archive{
				Exclude: []string{"/pattern1", "/pattern2", "/pattern3"},
			},
			pattern:      "/pattern2",
			want:         true,
			wantExcludes: []string{"/pattern1", "/pattern3"},
		},
		{
			name: "Remove existing pattern from beginning",
			archive: &Archive{
				Exclude: []string{"/pattern1", "/pattern2", "/pattern3"},
			},
			pattern:      "/pattern1",
			want:         true,
			wantExcludes: []string{"/pattern2", "/pattern3"},
		},
		{
			name: "Remove existing pattern from end",
			archive: &Archive{
				Exclude: []string{"/pattern1", "/pattern2", "/pattern3"},
			},
			pattern:      "/pattern3",
			want:         true,
			wantExcludes: []string{"/pattern1", "/pattern2"},
		},
		{
			name: "Remove non-existing pattern",
			archive: &Archive{
				Exclude: []string{"/pattern1", "/pattern2"},
			},
			pattern:      "/nonexistent",
			want:         false,
			wantExcludes: []string{"/pattern1", "/pattern2"},
		},
		{
			name: "Remove from single element list",
			archive: &Archive{
				Exclude: []string{"/pattern1"},
			},
			pattern:      "/pattern1",
			want:         true,
			wantExcludes: []string{},
		},
		{
			name: "Remove from empty list",
			archive: &Archive{
				Exclude: []string{},
			},
			pattern:      "/pattern1",
			want:         false,
			wantExcludes: []string{},
		},
		{
			name: "Remove from nil list",
			archive: &Archive{
				Exclude: nil,
			},
			pattern:      "/pattern1",
			want:         false,
			wantExcludes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveExclusion(tt.archive, tt.pattern)

			if got != tt.want {
				t.Errorf("RemoveExclusion() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.archive.Exclude, tt.wantExcludes) {
				t.Errorf("After RemoveExclusion(), archive.Exclude = %v, want %v", tt.archive.Exclude, tt.wantExcludes)
			}
		})
	}
}
