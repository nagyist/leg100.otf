package main

import (
	"context"
	"fmt"
	"os"

	cmdutil "github.com/leg100/otf/cmd"
	"github.com/leg100/otf/http"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 0)

var colorBGFG = os.Getenv("COLORFGBG")

var items = []list.Item{
	item{title: colorBGFG, desc: "I have ’em all over my house"},
	item{title: "Nutella", desc: "It's good on toast"},
	item{title: "Bitter melon", desc: "It cools you down"},
	item{title: "Nice socks", desc: "And by that I mean socks without holes"},
	item{title: "Eight hours of sleep", desc: "I had this once"},
	item{title: "Cats", desc: "Usually"},
	item{title: "Plantasia, the album", desc: "My plants love it too"},
	item{title: "Pour over coffee", desc: "It takes forever to make though"},
	item{title: "VR", desc: "Virtual reality...what is there to say?"},
	item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
	item{title: "Linux", desc: "Pretty much the best OS"},
	item{title: "Business school", desc: "Just kidding"},
	item{title: "Pottery", desc: "Wet clay is a great feeling"},
	item{title: "Shampoo", desc: "Nothing like clean hair"},
	item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
	item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
	item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
	item{title: "Stickers", desc: "The thicker the vinyl the better"},
	item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
	item{title: "Warm light", desc: "Like around 2700 Kelvin"},
	item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
	item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
	item{title: "Terrycloth", desc: "In other words, towel fabric"},
}

type model struct {
	list list.Model
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, nil
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	// Configure ^C to terminate program
	ctx, cancel := context.WithCancel(context.Background())
	cmdutil.CatchCtrlC(cancel)

	if err := Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, args []string) error {
	cfg, err := http.NewConfig(LoadCredentials)
	if err != nil {
		return err
	}

	cmd := &cobra.Command{
		Use:           "otf",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m := model{list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0)}
			m.list.Title = "My Fave Things"
			m.list.Paginator.Type = paginator.Arabic

			p := tea.NewProgram(m)
			p.EnterAltScreen()

			if err := p.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.Address, "address", http.DefaultAddress, "Address of OTF server")

	cmd.SetArgs(args)

	store, err := NewCredentialsStore()
	if err != nil {
		return err
	}

	cmd.AddCommand(LoginCommand(store, cfg.Address))
	cmd.AddCommand(OrganizationCommand(cfg))
	cmd.AddCommand(WorkspaceCommand(cfg))

	cmdutil.SetFlagsFromEnvVariables(cmd.Flags())

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}
