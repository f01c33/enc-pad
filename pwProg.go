package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type pwModel struct {
	pw       textinput.Model
	pass     string
	quitting bool
	err      error
	file     string
	data     []byte
}

func initialModelGetPW(file string) pwModel {
	ti := textinput.New()
	ti.Placeholder = "password for file"
	ti.Focus()
	return pwModel{
		pw:   ti,
		file: file,
	}
}

func (m pwModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m pwModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.pw.SetValue(m.pass)
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		if key.Matches(msg, enter) {
			_, err := os.Stat(m.file)
			if err != nil {
				return m, tea.Quit
			}
			m.data, m.err = os.ReadFile(m.file)
			if m.err != nil {
				m.err = err
				return m, nil
			}
			data, err := decryptAES(m.data, []byte(m.pass))
			if err != nil {
				m.err = err
				return m, nil
			}
			m.data = data
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	var cmd tea.Cmd
	m.pw, cmd = m.pw.Update(msg)
	m.pass = m.pw.Value()
	m.pw.SetValue(hide(m.pass))
	if m.quitting {
		return m, tea.Quit
	}
	return m, cmd
}

func hide(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		out[i] = '*'
	}
	return string(out)
}

func (m pwModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("%s\n%s\n", m.pw.View(), quitKeys.Help().Desc)
	if m.quitting {
		return str + "\n"
	}
	return str
}
