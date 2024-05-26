package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	textarea      textarea.Model
	err           error
	width, height int
}

var quitKeys = key.NewBinding(
	key.WithKeys("esc", "ctrl+c"),
	key.WithHelp("", "press esc to quit"),
)

var enter = key.NewBinding(
	key.WithKeys("enter"),
)

func initialModel(data string) model {
	ta := textarea.New()
	ta.SetValue(data)
	ta.Focus()
	return model{textarea: ta}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			return m, tea.Quit

		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.sizeInputs()
	case errMsg:
		m.err = msg
		return m, nil
	default:
		break
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) sizeInputs() {
	m.textarea.SetHeight(m.height)
	m.textarea.SetWidth(m.width)
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return m.textarea.View() + quitKeys.Help().Desc
}

var pwkey []byte

func main() {
	if len(os.Args) < 2 {
		fmt.Println("select a file to create or open like:\nenc-pad file.txt")
		os.Exit(0)
	}
	p := tea.NewProgram(initialModelGetPW(os.Args[1]), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if m.(pwModel).err != nil {
		fmt.Println(m.(pwModel).err)
		os.Exit(1)
	}
	pwkey = []byte(m.(pwModel).pass)
	p = tea.NewProgram(initialModel(string(m.(pwModel).data)), tea.WithAltScreen())
	m, err = p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data, err := encryptAES([]byte(m.(model).textarea.Value()), pwkey)
	if err != nil {
		fmt.Println(err, "writing file without encryption")
		err = os.WriteFile(os.Args[1], []byte(m.(model).textarea.Value()), os.ModePerm)
		if err != nil {
			fmt.Println(err, m.(model).textarea.Value())
		}
		os.Exit(1)
	}
	err = os.WriteFile(os.Args[1], data, os.ModePerm)
	if err != nil {
		fmt.Println(data)
	}
}
