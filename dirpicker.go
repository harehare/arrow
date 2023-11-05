package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/mo"
)

type DirPickerStyles struct {
	Cursor         lipgloss.Style
	Symlink        lipgloss.Style
	Text           lipgloss.Style
	Selected       lipgloss.Style
	Error          lipgloss.Style
	EmptyDirectory lipgloss.Style
}

var dirPickerStyles = DirPickerStyles{
	Cursor:         lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(CursorColor)),
	Symlink:        lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(SymLinkColor)),
	Text:           lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(ForegroundColor)),
	Selected:       lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(HighlightColor)).Bold(true),
	Error:          lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color("9")).PaddingLeft(2),
	EmptyDirectory: lipgloss.DefaultRenderer().NewStyle().Background(lipgloss.Color(DisabledColor)).MarginLeft(2).SetString(" No directory found."),
}

func DirPickerView(directories []Directory, selectedIndex, height int, displayIcons bool, hasChildDirectory mo.Option[bool], err error) string {
	if err != nil {
		return dirPickerStyles.Error.Render(err.Error())
	}

	if len(directories) == 0 {
		return dirPickerStyles.EmptyDirectory.String()
	}
	var displayDirectories []Directory
	var displayStart int

	if height == 0 {
		return ""
	}

	if selectedIndex+1 > height {
		if len(directories) > selectedIndex {
			displayStart, displayDirectories = selectedIndex-height+1, directories[selectedIndex-height+1:selectedIndex+1]
		} else {
			displayStart, displayDirectories = selectedIndex-height+1, directories[selectedIndex-height+1:]
		}
	} else if len(directories) > height {
		displayStart, displayDirectories = 0, directories[0:height]
	} else {
		displayStart, displayDirectories = 0, directories
	}

	var lines []string

	for i, directory := range displayDirectories {
		selected := selectedIndex == i+displayStart
		name := directory.Name()
		icon := GetIcon(directory, selected, displayIcons)
		msg := ""

		if !hasChildDirectory.OrElse(true) {
			msg = dirPickerStyles.EmptyDirectory.String()
		}

		line := directory.SymLink().Match(
			func(symLink string) (string, bool) {
				s := fmt.Sprintf("%s%s → %s", icon, name, symLink)
				if selected {
					return fmt.Sprintf("%s %s%s", dirPickerStyles.Cursor.Render("❯"), dirPickerStyles.Selected.Render(s), msg), true
				} else {
					return fmt.Sprintf("  %s", dirPickerStyles.Text.Render(s)), true
				}
			},
			func() (string, bool) {
				if selected {
					return fmt.Sprintf("%s %s%s%s", dirPickerStyles.Cursor.Render("❯"), dirPickerStyles.Selected.Render(icon), dirPickerStyles.Selected.Render(name), msg), true
				} else {
					return fmt.Sprintf("  %s%s", icon, dirPickerStyles.Text.Render(name)), true
				}
			},
		).OrElse("")

		lines = append(lines, zone.Mark(directory.String(), line))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lines...,
	)
}
