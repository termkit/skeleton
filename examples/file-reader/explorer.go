package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/termkit/skeleton"
	"os"
	"path/filepath"
	"time"
)

type explorer struct {
	skeleton *skeleton.Skeleton
	picker   filepicker.Model
}

func (e *explorer) Init() tea.Cmd {
	e.InitializeWidgets()
	return e.picker.Init()
}

func (e *explorer) blinkTwiceBorder(color string) {
	go func() {
		defaultColor := e.skeleton.GetBorderColor()
		for i := 0; i < 2; i++ {
			e.skeleton.SetBorderColor(color)
			time.Sleep(100 * time.Millisecond)
			e.skeleton.SetBorderColor(defaultColor)
			time.Sleep(100 * time.Millisecond)
		}
		e.skeleton.SetBorderColor(defaultColor)
	}()
}

func (e *explorer) InitializeWidgets() {
	e.skeleton.DeleteAllWidgets()
	allowedFiles := fmt.Sprintf("Allowed files: %s", e.picker.AllowedTypes)
	e.skeleton.AddWidget("allowed_types", allowedFiles)
}

func (e *explorer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case skeleton.IAMActivePage:
		e.InitializeWidgets()
	case tea.WindowSizeMsg:
		e.picker.Height = msg.Height - 7
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return e, tea.Quit
		}
	}

	var cmd tea.Cmd
	e.picker, cmd = e.picker.Update(msg)

	if isSelected, path := e.picker.DidSelectFile(msg); isSelected {
		fileName := filepath.Base(path)
		e.skeleton.AddPage(fileName, fileName, newFileReader(e.skeleton, fileName, path))
		e.blinkTwiceBorder("#00ff00") // Visualize the selection
	}

	if isSelected, _ := e.picker.DidSelectDisabledFile(msg); isSelected {
		// Change skeleton's border color to visualize the error
		e.blinkTwiceBorder("#ff0000")
	}

	return e, cmd
}

func (e *explorer) View() string {
	return e.picker.View()
}

func newExplorer(skeleton *skeleton.Skeleton) *explorer {
	fp := filepicker.New()
	fp.AutoHeight = false
	fp.Height = 1
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	return &explorer{
		skeleton: skeleton,
		picker:   fp,
	}
}
