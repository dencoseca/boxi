package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var yellowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")).Render
var greenStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render
var redStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render

func Green(message string, args ...any) {
	fmt.Println(greenStyle(fmt.Sprintf(message, args...)))
}

func Red(message string, args ...any) {
	fmt.Println(redStyle(fmt.Sprintf(message, args...)))
}

func Yellow(message string, args ...any) {
	fmt.Println(yellowStyle(fmt.Sprintf(message, args...)))
}
