package main

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	"github.com/samber/mo"
)

func TestParent(t *testing.T) {
	fs := fstest.MapFS{
		".":       {Mode: fs.ModeDir},
		"foo/bar": {Mode: fs.ModeDir},
	}
	tests := []struct {
		name string
		path string
		want mo.Option[Directory]
	}{
		{
			name: "When root directory",
			path: ".",
			want: mo.None[Directory](),
		},
		{
			name: "When subdirectory",
			path: "foo/bar",
			want: mo.Some(Directory{path: "foo", fsys: fs}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Directory{path: tt.path, fsys: fs}).Parent().IsPresent(); got != tt.want.IsPresent() {
				t.Errorf("directory.Parent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHidden(t *testing.T) {
	fs := fstest.MapFS{
		".foo": {Mode: fs.ModeDir},
		"bar":  {Mode: fs.ModeDir},
	}
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "When hidden directory",
			path: ".foo",
			want: true,
		},
		{
			name: "When not a hidden directory",
			path: "bar",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Directory{path: tt.path, fsys: fs}).IsHidden(); got != tt.want {
				t.Errorf("directory.IsHidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirs(t *testing.T) {
	fs := fstest.MapFS{
		"bar1":  {Mode: fs.ModeDir, ModTime: time.Date(2023, 10, 30, 11, 0, 0, 0, time.Local)},
		"bar2":  {Mode: fs.ModeDir, ModTime: time.Date(2023, 10, 30, 10, 0, 0, 0, time.Local)},
		"bar3":  {Mode: fs.ModeDir, ModTime: time.Date(2023, 10, 30, 9, 0, 0, 0, time.Local)},
		".bar4": {Mode: fs.ModeDir, ModTime: time.Date(2023, 10, 30, 8, 0, 0, 0, time.Local)},
	}
	tests := []struct {
		name    string
		path    string
		showAll bool
		order   Order
		want    []string
	}{
		{
			name:    "When sorting by directory name",
			path:    ".",
			showAll: false,
			order:   ORDER_NAME,
			want:    []string{"bar1", "bar2", "bar3"},
		},
		{
			name:    "When sorting by mod time",
			path:    ".",
			showAll: false,
			order:   ORDER_TIME,
			want:    []string{"bar3", "bar2", "bar1"},
		},
		{
			name:    "When sort by modification time and display all directories",
			path:    ".",
			showAll: true,
			order:   ORDER_TIME,
			want:    []string{".bar4", "bar3", "bar2", "bar1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, d := range (Directory{path: tt.path, fsys: fs}).Dirs(tt.showAll, tt.order).OrElse([]Directory{}) {
				got = append(got, d.String())
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("directory.Dirs() = %v, want %v", got, tt.want)
			}
		})
	}
}
