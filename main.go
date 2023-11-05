package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/termenv"
	"github.com/samber/mo"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

type Styles struct {
	Header           lipgloss.Style
	Count            lipgloss.Style
	CurrentDirectory lipgloss.Style
	Prompt           lipgloss.Style
	Foreground       lipgloss.Style
}

var (
	styles = Styles{
		Header: lipgloss.DefaultRenderer().NewStyle().PaddingBottom(1),
		Count:  lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(DisabledColor)),
		CurrentDirectory: lipgloss.DefaultRenderer().NewStyle().
			Bold(true).BorderForeground(lipgloss.Color(BorderColor)).BorderStyle(lipgloss.NormalBorder()).Foreground(lipgloss.Color(CurrentDirectoryColor)),
		Prompt:     lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(PromptColor)),
		Foreground: lipgloss.DefaultRenderer().NewStyle().Foreground(lipgloss.Color(ForegroundColor)),
	}
)

type model struct {
	currentDirectory    Directory
	hasChildDirectory   mo.Option[bool]
	cursor              int
	directories         []Directory
	filteredDirectories []Directory
	textInput           textinput.Model
	height              int
	showAll             bool
	displayIcons        bool
	order               Order
	err                 error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
	order := GetOrderIcon(m.order, m.displayIcons)

	return zone.Scan(
		lipgloss.JoinVertical(
			lipgloss.Top,
			styles.CurrentDirectory.Render(zone.Mark("order", order)+m.selectedDirectoryPath().OrElse(m.currentDirectory.String())+" "),
			m.textInput.View(),
			styles.Count.Render(m.countView()),
			DirPickerView(m.filteredDirectories, m.cursor, int(math.Max(float64(m.height-8), float64(0))), m.displayIcons, m.hasChildDirectory, m.err)))
}

func (m model) countView() string {
	return fmt.Sprintf("  %s/%s", strconv.Itoa(len(m.filteredDirectories)), strconv.Itoa(len(m.directories)))
}

func (m model) selectedDirectoryPath() mo.Option[string] {
	if len(m.filteredDirectories) > m.cursor {
		return mo.Some(m.filteredDirectories[m.cursor].String())
	}

	return mo.None[string]()
}

func filterDirectories(directories []Directory, query string) []Directory {
	var filteredDirectories []Directory

	if strings.Trim(query, "") == "" {
		return directories
	}

	for _, directory := range directories {
		if fuzzy.MatchNormalizedFold(query, directory.String()) {
			filteredDirectories = append(filteredDirectories, directory)
		}
	}

	return filteredDirectories
}

func (m model) changeOrder() (tea.Model, tea.Cmd) {
	switch m.order {
	case ORDER_NAME:
		m.order = ORDER_TIME
	case ORDER_TIME:
		m.order = ORDER_NAME
	}

	m.directories = m.currentDirectory.Dirs(m.showAll, m.order).Map(func(value []Directory) ([]Directory, error) {
		return value, nil
	}).MapErr(func(err error) ([]Directory, error) {
		m.err = err
		return []Directory{}, err
	}).OrElse([]Directory{})
	m.filteredDirectories = m.directories

	return m, nil
}

func (m model) moveTo(d Directory) (tea.Model, tea.Cmd) {
	m.hasChildDirectory = mo.None[bool]()
	if len(m.filteredDirectories)-1 < m.cursor {
		return m, nil
	}

	m.textInput.SetValue("")
	m.currentDirectory = d
	m.directories = m.currentDirectory.Dirs(m.showAll, m.order).Map(func(value []Directory) ([]Directory, error) {
		if len(value) == 0 {
			m.hasChildDirectory = mo.Some(false)
			return m.directories, nil
		}

		m.cursor = 0
		return value, nil
	}).MapErr(func(err error) ([]Directory, error) {
		m.cursor = 0
		return []Directory{}, err
	}).OrElse([]Directory{})
	m.filteredDirectories = m.directories
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.err = nil

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		return m, nil

	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		if zone.Get("order").InBounds(msg) {
			return m.changeOrder()
		}

		for _, d := range m.filteredDirectories {
			z := zone.Get(d.String())

			if z != nil && z.InBounds(msg) {
				return m.moveTo(d)
			}
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			wd, err := os.Getwd()

			if err != nil {
				panic(err)
			}

			fmt.Println(wd)
			return m, tea.Quit

		case tea.KeyUp:
			m.hasChildDirectory = mo.None[bool]()
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case tea.KeyDown:
			m.hasChildDirectory = mo.None[bool]()
			if m.cursor < len(m.filteredDirectories)-1 {
				m.cursor++
			}
			return m, nil

		case tea.KeyLeft:
			m.hasChildDirectory = mo.None[bool]()
			m.textInput.SetValue("")
			m.currentDirectory.Parent().ForEach(func(value Directory) {
				currentPath := m.currentDirectory.String()
				m.currentDirectory = value
				m.directories = m.currentDirectory.Dirs(m.showAll, m.order).Map(func(value []Directory) ([]Directory, error) {
					return value, nil
				}).MapErr(func(err error) ([]Directory, error) {
					m.err = err
					return []Directory{}, err
				}).OrElse([]Directory{})
				m.filteredDirectories = m.directories

				for i, d := range m.filteredDirectories {
					if d.String() == currentPath {
						m.cursor = i
						break
					}
					m.cursor = 0
				}
			})
			return m, nil

		case tea.KeyRight:
			return m.moveTo(m.filteredDirectories[m.cursor])

		case tea.KeyShiftDown:
			return m.changeOrder()

		case tea.KeyEnter:
			if len(m.filteredDirectories) == 0 {
				return m, nil
			} else {
				fmt.Println(m.filteredDirectories[m.cursor].String())
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	ct := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)
	m.filteredDirectories = filterDirectories(m.directories, m.textInput.Value())

	if ct != m.textInput.Value() {
		m.cursor = 0
	}

	return m, cmd
}

func initialModel(query string, showAll bool, displayIcons bool) model {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	currentDirectory := NewDirectory(wd)
	directories := currentDirectory.Dirs(showAll, ORDER_NAME).MapErr(func(err error) ([]Directory, error) {
		slog.Error(err.Error())
		return []Directory{}, nil
	}).OrElse([]Directory{})

	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.Prompt = "❯ "
	ti.PromptStyle = styles.Prompt
	ti.TextStyle = styles.Foreground

	if query != "" {
		ti.SetValue(query)
	}

	return model{
		currentDirectory:    currentDirectory,
		hasChildDirectory:   mo.None[bool](),
		cursor:              0,
		directories:         directories,
		filteredDirectories: directories,
		textInput:           ti,
		showAll:             showAll,
		displayIcons:        displayIcons,
		order:               ORDER_NAME,
		err:                 nil,
	}
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	cli.AppHelpTemplate = `USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
{{end}}`

	app := &cli.App{
		Name:    "arrow",
		Version: "v0.1.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Show hidden files.",
			},
			&cli.BoolFlag{
				Name:    "icons",
				Aliases: []string{"i"},
				Usage:   "Display icons.",
			},
			&cli.StringFlag{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "Specifies a query to search the directory.",
			},
		},
		Action: func(ctx *cli.Context) error {
			zone.NewGlobal()
			output := termenv.NewOutput(os.Stderr)
			lipgloss.SetColorProfile(output.ColorProfile())
			p := tea.NewProgram(initialModel(ctx.String("query"), ctx.Bool("all"), ctx.Bool("icons")), tea.WithOutput(os.Stderr), tea.WithAltScreen(), tea.WithMouseCellMotion())
			if _, err := p.Run(); err != nil {
				fmt.Printf("error: %v", err)
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
