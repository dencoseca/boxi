package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var yellowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")).Render
var greenStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render
var redStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render

// Green formats the given message with arguments, applies green styling, and
// prints it to the standard output.
func Green(message string, args ...any) {
	fmt.Println(greenStyle(fmt.Sprintf(message, args...)))
}

// Red applies red styling to the formatted message and prints it to the standard
// output.
func Red(message string, args ...any) {
	fmt.Println(redStyle(fmt.Sprintf(message, args...)))
}

// Yellow prints a formatted message with yellow styling using the specified
// format and arguments.
func Yellow(message string, args ...any) {
	fmt.Println(yellowStyle(fmt.Sprintf(message, args...)))
}
