package main

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/termenv"
	"github.com/samber/mo"
)

func TestDirPickerView(t *testing.T) {
	lipgloss.SetColorProfile(termenv.Ascii)
	zone.NewGlobal()
	fs := fstest.MapFS{
		".":       {Mode: fs.ModeDir},
		"foo/bar": {Mode: fs.ModeDir},
	}
	tests := []struct {
		name              string
		directories       []Directory
		selectedIndex     int
		height            int
		displayIcons      bool
		hasChildDirectory mo.Option[bool]
		err               error
		want              string
	}{
		{
			name:              "When empty directory",
			directories:       []Directory{},
			selectedIndex:     0,
			height:            100,
			displayIcons:      false,
			hasChildDirectory: mo.None[bool](),
			err:               nil,
			want:              "   No directory found.",
		},
		{
			name:              "When only one directory",
			directories:       []Directory{{path: "foo", fsys: fs}},
			selectedIndex:     0,
			height:            100,
			displayIcons:      false,
			hasChildDirectory: mo.None[bool](),
			err:               nil,
			want: lipgloss.JoinVertical(
				lipgloss.Top, zone.Mark("foo", "❯ foo")),
		},
		{
			name:              "When multiple directories",
			directories:       []Directory{{path: "foo", fsys: fs}, {path: "foo/bar", fsys: fs}},
			selectedIndex:     0,
			height:            100,
			displayIcons:      false,
			hasChildDirectory: mo.None[bool](),
			err:               nil,
			want: lipgloss.JoinVertical(
				lipgloss.Top, zone.Mark("foo", "❯ foo"), zone.Mark("foo/bar", "  bar")),
		},
		{
			name:              "When multiple directories and change the cursor position",
			directories:       []Directory{{path: "foo", fsys: fs}, {path: "foo/bar", fsys: fs}},
			selectedIndex:     1,
			height:            100,
			displayIcons:      false,
			hasChildDirectory: mo.None[bool](),
			err:               nil,
			want: lipgloss.JoinVertical(
				lipgloss.Top, zone.Mark("foo", "  foo"), zone.Mark("foo/bar", "❯ bar")),
		},
		{
			name:              "When has error",
			directories:       []Directory{{path: "foo", fsys: fs}, {path: "foo/bar", fsys: fs}},
			selectedIndex:     0,
			height:            100,
			displayIcons:      false,
			hasChildDirectory: mo.None[bool](),
			err:               errors.New("Error"),
			want:              "  Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirPickerView(tt.directories, tt.selectedIndex, tt.height, tt.displayIcons, tt.hasChildDirectory, tt.err); got != tt.want {
				fmt.Println(got)
				t.Errorf("DirPickerView = %v, want = %v", got, tt.want)
			}
		})
	}
}
