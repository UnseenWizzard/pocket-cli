package util

import "testing"

func TestIsHttpSuccess(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			"http200_isSuccess",
			200,
			true,
		},
		{
			"http403_isFailue",
			403,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHttpSuccess(tt.statusCode); got != tt.want {
				t.Errorf("IsHttpSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
