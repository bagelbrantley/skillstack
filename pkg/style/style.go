package style

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Cyan      = lipgloss.Color("#06b6d4")
	Emerald   = lipgloss.Color("#10b981")
	Rose      = lipgloss.Color("#f43f5e")
	Amber     = lipgloss.Color("#f59e0b")
	SlateDark = lipgloss.Color("#0f172a")
	SlateGray = lipgloss.Color("#64748b")
	SlateLite = lipgloss.Color("#cbd5e1")
	White     = lipgloss.Color("#ffffff")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan).
			MarginBottom(1)

	BannerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(White).
				Background(Cyan).
				Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(Emerald)
	ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(Rose)
	WarnStyle    = lipgloss.NewStyle().Bold(true).Foreground(Amber)
	MutedStyle   = lipgloss.NewStyle().Foreground(SlateGray)
	AccentStyle  = lipgloss.NewStyle().Bold(true).Foreground(Cyan)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Cyan).
			Padding(1, 2)

	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Emerald).
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Cyan).
			MarginBottom(1)

	TagStyle = lipgloss.NewStyle().
			Foreground(SlateDark).
			Background(Cyan).
			Padding(0, 1).
			Bold(true)

	BadgeSuccess = lipgloss.NewStyle().
			Foreground(White).
			Background(Emerald).
			Padding(0, 1).
			Bold(true)
)

// PrintBanner prints a prominent welcome banner
func PrintBanner(out io.Writer, title string, subtitle string) {
	bannerContent := BannerTitleStyle.Render(" "+title+" ") + "\n" +
		MutedStyle.Render(subtitle)
	fmt.Fprintln(out, BoxStyle.Render(bannerContent))
}

// PrintHeader prints a styled header
func PrintHeader(out io.Writer, title string) {
	fmt.Fprintln(out, HeaderStyle.Render("=== "+title+" ==="))
}

// PrintSuccess prints success message
func PrintSuccess(out io.Writer, msg string) {
	fmt.Fprintln(out, SuccessStyle.Render("✓ [SUCCESS] ")+msg)
}

// PrintError prints error message
func PrintError(out io.Writer, msg string) {
	fmt.Fprintln(out, ErrorStyle.Render("✗ [ERROR] ")+msg)
}

// PrintWarn prints warning message
func PrintWarn(out io.Writer, msg string) {
	fmt.Fprintln(out, WarnStyle.Render("! [WARNING] ")+msg)
}

// PrintKeyValue prints a key-value pair aligned neatly
func PrintKeyValue(out io.Writer, key string, val string) {
	fmt.Fprintf(out, "%s %s\n", AccentStyle.Render(fmt.Sprintf("%-18s:", key)), val)
}

// PrintJSON prints structured data as formatted JSON
func PrintJSON(out io.Writer, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("json format error: %w", err)
	}
	fmt.Fprintln(out, string(bytes))
	return nil
}

// FormatTable renders a simple formatted table for CLI output
func FormatTable(headers []string, rows [][]string) string {
	var sb strings.Builder

	// Header row
	headerStr := ""
	for i, h := range headers {
		if i == 0 {
			headerStr += fmt.Sprintf("%-25s", h)
		} else if i == 1 {
			headerStr += fmt.Sprintf("%-12s", h)
		} else if i == 2 {
			headerStr += fmt.Sprintf("%-12s", h)
		} else {
			headerStr += fmt.Sprintf("%-30s", h)
		}
	}
	sb.WriteString(AccentStyle.Render(headerStr) + "\n")
	sb.WriteString(MutedStyle.Render(strings.Repeat("-", 75)) + "\n")

	// Data rows
	for _, row := range rows {
		rowStr := ""
		for i, cell := range row {
			if i == 0 {
				rowStr += fmt.Sprintf("%-25s", cell)
			} else if i == 1 {
				rowStr += fmt.Sprintf("%-12s", cell)
			} else if i == 2 {
				rowStr += fmt.Sprintf("%-12s", cell)
			} else {
				rowStr += fmt.Sprintf("%-30s", cell)
			}
		}
		sb.WriteString(rowStr + "\n")
	}

	return sb.String()
}
