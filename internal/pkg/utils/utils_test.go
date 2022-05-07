package utils_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/pkg/utils"
)

func TestSliceContainsString(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		value string
		want  bool
	}{
		{
			name:  "slice contains value",
			slice: []string{"a", "b", "c"},
			value: "b",
			want:  true,
		},
		{
			name:  "slice does not contain value",
			slice: []string{"a", "b", "c"},
			value: "d",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.SliceContainsString(tt.slice, tt.value); got != tt.want {
				t.Errorf("SliceContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
