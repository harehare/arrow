package main

import (
	"testing"
)

func TestColorFromString(t *testing.T) {
	tests := []struct {
		color string
		want  bool
	}{
		{
			color: "0",
			want:  true,
		},
		{
			color: "255",
			want:  true,
		},
		{
			color: "#000000",
			want:  true,
		},
		{
			color: "#FFFFFF",
			want:  true,
		},
		{
			color: "-1",
			want:  false,
		},
		{
			color: "256",
			want:  false,
		},
		{
			color: "#0000000",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			if got := ColorFromString(tt.color).IsPresent(); got != tt.want {
				t.Errorf("ColorFromString(%v).IsPresent() = %v, want %v", tt.color, got, tt.want)
			}
		})
	}
}
