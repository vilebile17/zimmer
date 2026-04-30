package main

import (
	"errors"
	"log"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const baseURL = "http://localhost:8080"

func noJWT() tea.Msg {
	return errMsg{errors.New("couldn't find the .zimmer_token file in the home directory")}
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func getJWTToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return "", err
	}

	var JWTToken []byte
	for _, file := range files {
		if file.Name() == ".zimmer_token" {
			JWTToken, err = os.ReadFile(homeDir + "/" + ".zimmer_token")
			if err != nil {
				return "", err
			}
			return string(JWTToken), nil
		}
	}

	token := strings.TrimSpace(string(JWTToken))
	return token, nil
}

func main() {
	if _, err := tea.NewProgram(allClassesViewModel{}).Run(); err != nil {
		log.Fatalf("Uh oh, there was an error: %v\n", err)
	}
}
