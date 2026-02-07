package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

const (
	sampleRate = 44100
	clickFreq  = 1200 // Hz - frequency of the click sound
	clickDur   = 0.01 // seconds - duration of the click
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	beatBoxStyle = lipgloss.NewStyle().
			Width(12).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	activeBeatStyle = lipgloss.NewStyle().
			Width(12).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Background(lipgloss.Color("235")).
			Bold(true)

	downbeatStyle = lipgloss.NewStyle().
			Width(12).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("205")).
			Background(lipgloss.Color("235")).
			Bold(true)

	bpmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Space key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "+", "="),
		key.WithHelp("↑/+", "increase BPM"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "-"),
		key.WithHelp("↓/-", "decrease BPM"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "play/pause"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Space, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Space, k.Quit},
	}
}

type tickMsg time.Time
type animateMsg time.Time

type model struct {
	bpm          int
	currentBeat  int
	playing      bool
	help         help.Model
	keys         keyMap
	animatePhase int // 0-10 for wave animation
}

func initialModel(bpm int) model {
	return model{
		bpm:          bpm,
		currentBeat:  0,
		playing:      true,
		help:         help.New(),
		keys:         keys,
		animatePhase: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tick(m.bpm)
}

func tick(bpm int) tea.Cmd {
	interval := time.Duration(60.0/float64(bpm)*1000) * time.Millisecond
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func animate() tea.Cmd {
	return tea.Tick(30*time.Millisecond, func(t time.Time) tea.Msg {
		return animateMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Space):
			m.playing = !m.playing
			if m.playing {
				return m, tick(m.bpm)
			}
		case key.Matches(msg, m.keys.Up):
			if m.bpm < 300 {
				m.bpm += 5
			}
		case key.Matches(msg, m.keys.Down):
			if m.bpm > 20 {
				m.bpm -= 5
			}
		}

	case tickMsg:
		if m.playing {
			m.currentBeat = (m.currentBeat + 1) % 4
			m.animatePhase = 0
			speaker.Play(generateClick())
			return m, tea.Batch(tick(m.bpm), animate())
		}

	case animateMsg:
		if m.animatePhase < 10 {
			m.animatePhase++
			return m, animate()
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🎵 Terminal Metronome"))
	b.WriteString("\n\n")

	// Beat boxes
	beats := make([]string, 4)
	for i := 0; i < 4; i++ {
		beatNum := fmt.Sprintf("%d", i+1)
		var beatContent string

		// Calculate wave effect
		var symbol string
		if i == 0 {
			symbol = "♩"
		} else {
			symbol = "♪"
		}

		// Add wave animation on active beat
		if i == m.currentBeat && m.playing && m.animatePhase < 10 {
			// Wave effect using different sizes
			waves := []string{
				fmt.Sprintf("\n%s\n\n%s", symbol, beatNum),
				fmt.Sprintf("%s\n\n\n%s", symbol, beatNum),
				fmt.Sprintf("\n%s\n\n%s", symbol, beatNum),
				fmt.Sprintf("\n\n%s\n%s", symbol, beatNum),
				fmt.Sprintf("\n%s\n\n%s", symbol, beatNum),
			}
			waveIndex := m.animatePhase / 2
			if waveIndex >= len(waves) {
				waveIndex = len(waves) - 1
			}

			if i == 0 {
				// Downbeat - special styling with animation
				beatContent = downbeatStyle.Render(waves[waveIndex])
			} else {
				beatContent = activeBeatStyle.Render(waves[waveIndex])
			}
		} else if i == m.currentBeat && m.playing {
			// Active beat (no animation)
			if i == 0 {
				beatContent = downbeatStyle.Render(fmt.Sprintf("\n%s\n\n%s", symbol, beatNum))
			} else {
				beatContent = activeBeatStyle.Render(fmt.Sprintf("\n%s\n\n%s", symbol, beatNum))
			}
		} else {
			// Inactive beat
			beatContent = beatBoxStyle.Render(fmt.Sprintf("\n%s\n\n%s", symbol, beatNum))
		}

		beats[i] = beatContent
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, beats...))
	b.WriteString("\n")

	// BPM info
	status := "Playing"
	if !m.playing {
		status = "Paused"
	}
	b.WriteString(bpmStyle.Render(fmt.Sprintf("BPM: %d  •  Status: %s", m.bpm, status)))
	b.WriteString("\n")

	// Help
	b.WriteString(helpStyle.Render(m.help.ShortHelpView(m.keys.ShortHelp())))

	return b.String()
}

// generateClick creates a short click sound using a sine wave with envelope
func generateClick() beep.Streamer {
	samples := int(sampleRate * clickDur)
	data := make([][2]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / sampleRate
		// Sine wave
		wave := math.Sin(2 * math.Pi * clickFreq * t)
		// Exponential decay envelope for a sharper click sound
		envelope := math.Exp(-t * 200)
		sample := wave * envelope * 0.3 // 0.3 for volume control

		data[i][0] = sample
		data[i][1] = sample
	}

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if len(data) == 0 {
			return 0, false
		}
		n = copy(samples, data)
		data = data[n:]
		return n, true
	})
}

func main() {
	// Default BPM
	bpm := 120
	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil && val > 0 && val <= 300 {
			bpm = val
		} else {
			fmt.Println("Usage: metronome [BPM]")
			fmt.Println("BPM must be between 20 and 300. Defaulting to 120.")
		}
	}

	// Initialize audio
	sr := beep.SampleRate(sampleRate)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		fmt.Printf("Error initializing audio: %v\n", err)
		os.Exit(1)
	}

	// Start Bubble Tea program
	p := tea.NewProgram(initialModel(bpm))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}