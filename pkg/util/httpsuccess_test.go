package util

import (
	"testing"
)

func TestIsHttpSuccess(t *testing.T) {
	type statusCodeRange struct {
		firstCode       int
		lastCode        int
		additionalCodes []int
	}
	tests := []struct {
		name  string
		codes statusCodeRange
		want  bool
	}{
		{
			"http1xx_isNotSuccess",
			statusCodeRange{
				100,
				103,
				[]int{},
			},
			false,
		},
		{
			"http2xx_isSuccess",
			statusCodeRange{
				200,
				208,
				[]int{226},
			},
			true,
		},
		{
			"http3xx_isNotSuccess",
			statusCodeRange{
				300,
				308,
				[]int{},
			},
			false,
		},
		{
			"http4xx_isNotSuccess",
			statusCodeRange{
				400,
				418,
				[]int{421, 422, 423, 424, 425, 426, 428, 429, 431, 451},
			},
			false,
		},
		{
			"http5xx_isNotSuccess",
			statusCodeRange{
				500,
				508,
				[]int{510, 511},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for statusCode := tt.codes.firstCode; statusCode <= tt.codes.lastCode; statusCode++ {
				testIsSuccess(t, statusCode, tt.want)
			}
			for _, statusCode := range tt.codes.additionalCodes {
				testIsSuccess(t, statusCode, tt.want)
			}
		})
	}
}

func testIsSuccess(t *testing.T, statusCode int, want bool) {
	if got := IsHttpSuccess(statusCode); got != want {
		t.Errorf("IsHttpSuccess(%v) = %v, want %v", statusCode, got, want)
	}
}
