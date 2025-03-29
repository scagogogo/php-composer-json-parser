package autoload

import (
	"reflect"
	"testing"
)

func TestGetPSR4Map(t *testing.T) {
	tests := []struct {
		name     string
		autoload *Autoload
		want     map[string]string
		wantOk   bool
	}{
		{
			name: "Valid PSR-4 map",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Vendor\\Package\\": "src/",
					"Vendor\\Tests\\":   "tests/",
				},
			},
			want: map[string]string{
				"Vendor\\Package\\": "src/",
				"Vendor\\Tests\\":   "tests/",
			},
			wantOk: true,
		},
		{
			name: "Empty PSR-4 map",
			autoload: &Autoload{
				PSR4: map[string]interface{}{},
			},
			want:   map[string]string{},
			wantOk: true,
		},
		{
			name: "Nil PSR-4 map",
			autoload: &Autoload{
				PSR4: nil,
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "PSR-4 map with non-string values",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Vendor\\Package\\": "src/",
					"Vendor\\Tests\\":   []string{"tests/"}, // 不是字符串
				},
			},
			want: map[string]string{
				"Vendor\\Package\\": "src/",
			},
			wantOk: true,
		},
		{
			name: "PSR-4 is not a map",
			autoload: &Autoload{
				PSR4: "not a map", // 这个会在运行时才会被解析
			},
			want:   nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetPSR4Map(tt.autoload)

			if ok != tt.wantOk {
				t.Errorf("GetPSR4Map() ok = %v, wantOk %v", ok, tt.wantOk)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPSR4Map() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetPSR4(t *testing.T) {
	tests := []struct {
		name      string
		autoload  *Autoload
		namespace string
		path      string
		wantPSR4  map[string]interface{}
	}{
		{
			name: "Add to existing map",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Existing\\": "existing/",
				},
			},
			namespace: "Vendor\\Package\\",
			path:      "src/",
			wantPSR4: map[string]interface{}{
				"Existing\\":        "existing/",
				"Vendor\\Package\\": "src/",
			},
		},
		{
			name: "Add to empty map",
			autoload: &Autoload{
				PSR4: map[string]interface{}{},
			},
			namespace: "Vendor\\Package\\",
			path:      "src/",
			wantPSR4: map[string]interface{}{
				"Vendor\\Package\\": "src/",
			},
		},
		{
			name: "Add to nil map",
			autoload: &Autoload{
				PSR4: nil,
			},
			namespace: "Vendor\\Package\\",
			path:      "src/",
			wantPSR4: map[string]interface{}{
				"Vendor\\Package\\": "src/",
			},
		},
		{
			name: "Update existing namespace",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Vendor\\Package\\": "old/path/",
				},
			},
			namespace: "Vendor\\Package\\",
			path:      "new/path/",
			wantPSR4: map[string]interface{}{
				"Vendor\\Package\\": "new/path/",
			},
		},
		{
			name: "PSR4 is not a map",
			autoload: &Autoload{
				PSR4: []string{"not a map"}, // 不是map，但也不是字符串
			},
			namespace: "Vendor\\Package\\",
			path:      "src/",
			wantPSR4: map[string]interface{}{
				"Vendor\\Package\\": "src/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPSR4(tt.autoload, tt.namespace, tt.path)

			// 检查PSR4字段是否为map
			psr4Map, ok := tt.autoload.PSR4.(map[string]interface{})
			if !ok {
				t.Errorf("SetPSR4() PSR4 is not a map[string]interface{}, got %T", tt.autoload.PSR4)
				return
			}

			// 检查map内容是否符合预期
			if !reflect.DeepEqual(psr4Map, tt.wantPSR4) {
				t.Errorf("SetPSR4() PSR4 = %v, want %v", psr4Map, tt.wantPSR4)
			}
		})
	}
}

func TestRemovePSR4(t *testing.T) {
	tests := []struct {
		name      string
		autoload  *Autoload
		namespace string
		want      bool
		checkPSR4 bool // 是否检查PSR4字段的值
		wantPSR4  interface{}
	}{
		{
			name: "Remove existing namespace",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Vendor\\Package\\": "src/",
					"Vendor\\Tests\\":   "tests/",
				},
			},
			namespace: "Vendor\\Package\\",
			want:      true,
			checkPSR4: true,
			wantPSR4: map[string]interface{}{
				"Vendor\\Tests\\": "tests/",
			},
		},
		{
			name: "Remove non-existing namespace",
			autoload: &Autoload{
				PSR4: map[string]interface{}{
					"Vendor\\Tests\\": "tests/",
				},
			},
			namespace: "Vendor\\Package\\",
			want:      false,
			checkPSR4: true,
			wantPSR4: map[string]interface{}{
				"Vendor\\Tests\\": "tests/",
			},
		},
		{
			name: "Remove from empty map",
			autoload: &Autoload{
				PSR4: map[string]interface{}{},
			},
			namespace: "Vendor\\Package\\",
			want:      false,
			checkPSR4: true,
			wantPSR4:  map[string]interface{}{},
		},
		{
			name: "Remove from nil map",
			autoload: &Autoload{
				PSR4: nil,
			},
			namespace: "Vendor\\Package\\",
			want:      false,
			checkPSR4: true,
			wantPSR4:  nil,
		},
		{
			name: "PSR4 is not a map",
			autoload: &Autoload{
				PSR4: []string{"not a map"}, // 使用slice代替字符串
			},
			namespace: "Vendor\\Package\\",
			want:      false,
			checkPSR4: true,
			wantPSR4:  []string{"not a map"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemovePSR4(tt.autoload, tt.namespace)

			if got != tt.want {
				t.Errorf("RemovePSR4() = %v, want %v", got, tt.want)
			}

			// 检查PSR4字段是否符合预期
			if tt.checkPSR4 && !reflect.DeepEqual(tt.autoload.PSR4, tt.wantPSR4) {
				t.Errorf("RemovePSR4() PSR4 = %v, want %v", tt.autoload.PSR4, tt.wantPSR4)
			}
		})
	}
}
