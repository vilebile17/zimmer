package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	database "github.com/vilebile17/zimmer/internal/database"
)

type Classes struct {
	ClassesAsStudent []database.Class `json:"classes_as_student"`
	ClassesAsTeacher []database.Class `json:"classes_as_teacher"`
}

type allClassesViewModel struct {
	classList           []string
	cursor              int
	numClassesAsTeacher int
	err                 error
}

func getClasses(JWTToken string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}

		req, err := http.NewRequest("GET", baseURL+"/api/classes", nil)
		if err != nil {
			return errMsg{err}
		}

		JWTToken = strings.ReplaceAll(JWTToken, "\n", "")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JWTToken))

		res, err := c.Do(req)
		if err != nil {
			return errMsg{err}
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return errMsg{err}
		}

		var classes Classes
		err = json.Unmarshal(data, &classes)
		if err != nil {
			return errMsg{err}
		}

		return classes
	}
}

func (m allClassesViewModel) Init() tea.Cmd {
	token, err := getJWTToken()
	if err != nil {
		return noJWT
	}

	if token == "" {
		return noJWT
	}
	return getClasses(token)
}

func (m allClassesViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Classes:
		for _, class := range msg.ClassesAsTeacher {
			m.classList = append(m.classList, class.Name)
		}
		for _, class := range msg.ClassesAsStudent {
			m.classList = append(m.classList, class.Name)
		}
		m.numClassesAsTeacher = len(msg.ClassesAsTeacher)
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.classList)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m allClassesViewModel) View() tea.View {
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err))
	}

	s := ""
	if m.numClassesAsTeacher != 0 {
		s += "Classes as Teacher: \n"
	}

	for i := 0; i < m.numClassesAsTeacher; i++ {
		if i == m.cursor {
			s += "> "
		} else {
			s += "  "
		}
		s += m.classList[i]
		s += "\n"
	}

	if m.numClassesAsTeacher != len(m.classList) {
		s += "\nClasses as Student: \n"
	}

	for i := m.numClassesAsTeacher; i < len(m.classList); i++ {
		if i == m.cursor {
			s += "> "
		} else {
			s += "  "
		}
		s += m.classList[i]
		s += "\n"
	}

	// Send off whatever we came up with above for rendering.
	return tea.NewView("\n" + s + "\n\n")
}
