package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewTabs() string {
	labels := []string{"Overview", "Experience", "Projects", "Contact"}
	var rendered []string

	for i, label := range labels {
		if i == m.activeTab {
			rendered = append(rendered, tabActiveStyle.Render(label))
		} else {
			rendered = append(rendered, tabInactiveStyle.Render(label))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
	return tabsRowStyle.Render(row)
}

func (m model) viewTabContent() string {
	var text string

	switch m.activeTab {
	case 0:
		text = contentStyle.Render(m.viewOverview())
	case 1:
		text = contentStyle.Render(m.viewExperience())
	case 2:
		text = contentStyle.Render(m.viewProjects())
	case 3:
		text = "Contact tab — placeholder content."
	}

	return contentStyle.Render(text)
}

func (m model) viewFooter() string {
	helpLine := "h/← & l/→: switch tabs  •  1–4: jump to tab  •  q: quit"
	switch m.activeTab {
	case 1:
		helpLine += "  •  j/k or ↑/↓: switch experience"
	case 2:
		helpLine += "  •  j/k or ↑/↓: switch project"
	}
	return footerStyle.Render(helpLine)
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var content string

	if m.loading {
		// Intro phase: typewriter animation for "Shubhom Srivastava"

		if m.typedChars < 0 {
			m.typedChars = 0
		}
		if m.typedChars > len(m.introText) {
			m.typedChars = len(m.introText)
		}

		visible := m.introText[:m.typedChars]

		// Thicker cursor: full block "█"
		cursorChar := ""
		if m.cursorOn {
			cursorChar = "█"
		}

		// Style the name + cursor separately
		styledName := m.nameStyle.Render(visible)
		styledCursor := ""
		if cursorChar != "" {
			styledCursor = m.cursorStyle.Render(cursorChar)
		}

		line := styledName + styledCursor
		content = line
	} else {
		// 🔹 Main portfolio card view

		tabContent := m.viewTabContent()
		tabsRow := m.viewTabs()
		footer := m.viewFooter()

		body := lipgloss.JoinVertical(
			lipgloss.Left,
			tabContent,
			"",
			tabsRow,
			footer,
		)

		content = cardStyle.Render(body)
	}

	// If we don't know the size yet, just return the raw content.
	if m.width == 0 || m.height == 0 {
		return content + "\n"
	}

	// Center the content (intro or card) in the available space
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // horizontal
		lipgloss.Center, // vertical
		content,
	)
}

const (
	esc = "\x1b"
	bel = "\x07"
)

func termLink(label, url string) string {
	if url == "" {
		return label
	}
	// OSC 8 ; ; url ST   label   OSC 8 ; ; ST
	return esc + "]8;;" + url + bel + label + esc + "]8;;" + bel
}

func centerInContent(s string) string {
	innerWidth := appWidth - 4 // same width you use in contentStyle
	return lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, s)
}

func (m model) viewOverview() string {
	if m.portfolio == nil {
		return "Overview not available."
	}

	p := m.portfolio

	var lines []string

	// 1) Name
	nameLine := centerInContent(m.nameStyle.Render(p.Name))
	lines = append(lines, nameLine)

	// 2) Tagline (slightly dimmer / separate style if you want)
	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DDDDDD"))
	taglineLine := centerInContent(taglineStyle.Render(p.Tagline))
	lines = append(lines, taglineLine)

	// Blank line
	lines = append(lines, "")

	// 3) Intro paragraph
	if intro := strings.TrimSpace(p.Overview.Intro); intro != "" {
		lines = append(lines, intro)
		lines = append(lines, "") // blank after intro
	}

	// 4) Bullets like "Backend: ...", "Infra: ...", etc.
	for _, b := range p.Overview.Bullets {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		lines = append(lines, "• "+b)
	}

	// 5) Social links line at the bottom of the content box
	// inside viewOverview
	var socialParts []string

	if p.Contact.GitHub != "" {
		socialParts = append(socialParts, termLink("GitHub", p.Contact.GitHub))
	}
	if p.Contact.LinkedIn != "" {
		socialParts = append(socialParts, termLink("LinkedIn", p.Contact.LinkedIn))
	}
	if p.Contact.Email != "" {
		socialParts = append(socialParts, termLink("Email", "mailto:"+p.Contact.Email))
	}

	if len(socialParts) > 0 {
		lines = append(lines, "")
		socialLine := strings.Join(socialParts, "  ·  ")
		lines = append(lines, centerInContent(socialLine))
	}

	return strings.Join(lines, "\n")
}

func (m model) viewExperience() string {
	if m.portfolio == nil || len(m.portfolio.Experiences) == 0 {
		return "No experience data yet."
	}

	exps := m.portfolio.Experiences
	idx := m.expList.Page
	if idx < 0 {
		idx = 0
	}
	if idx >= len(exps) {
		idx = len(exps) - 1
	}
	exp := exps[idx]

	var lines []string

	// Title line: Company — Role
	header := fmt.Sprintf("%s — %s", exp.Company, exp.Role)
	header = lipgloss.NewStyle().Bold(true).Render(header)
	lines = append(lines, header)

	// Period / location
	var metaParts []string
	if exp.Period != "" {
		metaParts = append(metaParts, exp.Period)
	}
	if exp.Location != "" {
		metaParts = append(metaParts, exp.Location)
	}
	if len(metaParts) > 0 {
		meta := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Render(strings.Join(metaParts, " · "))
		lines = append(lines, meta)
	}

	lines = append(lines, "") // blank line

	// Bullets (cap at 3 for readability)
	maxBullets := 3
	for i, b := range exp.Bullets {
		if i >= maxBullets {
			break
		}
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		lines = append(lines, "• "+b)
	}

	// Stack line
	if strings.TrimSpace(exp.Stack) != "" {
		lines = append(lines, "")
		stackLine := "Stack: " + exp.Stack
		lines = append(lines, stackLine)
	}

	// Paginator indicator at bottom
	lines = append(lines, "")
	pagerLine := fmt.Sprintf("(%d/%d)", idx+1, len(exps))
	// use paginator's dots + our numeric info
	pagerView := m.expList.View()
	lines = append(lines, pagerLine+"  "+pagerView)

	return strings.Join(lines, "\n")
}

func (m model) viewProjects() string {
	if len(m.portfolio.Projects) == 0 || m.portfolio.Projects == nil {
		return "no projects data"
	}

	projs := m.portfolio.Projects

	var lines []string

	idx := m.projList.Page

	if idx < 0 {
		idx = 0
	}
	if idx >= len(projs) {
		idx = len(projs) - 1
	}
	proj := projs[idx]

	header := lipgloss.NewStyle().Bold(true).Render(proj.Name)
	lines = append(lines, header)

	lines = append(lines, "")

	// Bullets (cap at 3 for readability)
	maxBullets := 3
	for i, b := range proj.Bullets {
		if i >= maxBullets {
			break
		}
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		lines = append(lines, "• "+b)
	}

	// Stack line
	if strings.TrimSpace(proj.Stack) != "" {
		lines = append(lines, "")
		stackLine := "Stack: " + proj.Stack
		lines = append(lines, stackLine)
	}
	lines = append(lines, " ")

	linksParts := []string{}

	if proj.Links.Code != "" {
		linksParts = append(linksParts, termLink("Code", proj.Links.Code))
	}
	if proj.Links.Demo != "" {
		linksParts = append(linksParts, termLink("Demo", proj.Links.Demo))
	}

	if len(linksParts) > 0 {
		lines = append(lines, "")
		linkLine := strings.Join(linksParts, "  ·  ")
		lines = append(lines, linkLine)
	}

	// Paginator indicator at bottom
	lines = append(lines, "")
	pagerLine := fmt.Sprintf("(%d/%d)", idx+1, len(projs))
	// use paginator's dots + our numeric info
	pagerView := m.projList.View()
	lines = append(lines, pagerLine+"  "+pagerView)

	return strings.Join(lines, "\n")
}
