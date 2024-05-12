// Package chatui provides a very simple chat like interface for applications.
package chatui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChatUI is a simple text UI that looks kind of like a chat window.
type ChatUI struct {
	outputCh  <-chan string
	commandCh chan<- string
	output    *tview.TextView
	status    *tview.TextView
	input     *tview.InputField
	flex      *tview.Flex
	app       *tview.Application
	history   *History
}

// Config is the configuration of the chat UI
type Config struct {
	OutputCh     <-chan string
	CommandCh    chan<- string
	Theme        *tview.Theme
	DynamicColor bool
	BlockCtrlC   bool
	HistorySize  int
}

const (
	defaultHistorySize = 1000
)

// New creates a new Text UI
func New(config Config) *ChatUI {
	if config.OutputCh == nil {
		panic("you must define an InputCh")
	}

	if config.CommandCh == nil {
		panic("you must define a CommandCh")
	}

	if config.HistorySize == 0 {
		config.HistorySize = defaultHistorySize
	}

	// Override with supplied theme if present
	if config.Theme != nil {
		tview.Styles = *config.Theme
	}

	tui := &ChatUI{
		history:   NewHistory(config.HistorySize),
		outputCh:  config.OutputCh,
		commandCh: config.CommandCh,
	}

	// Create the output textview
	tui.output = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true).
		SetDynamicColors(config.DynamicColor)

	tui.output. // returns tview.Box
			SetBorder(false)

	// Create the status bar
	tui.status = tview.NewTextView().
		SetText("").
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tview.Styles.SecondaryTextColor).
		SetScrollable(false)

	tui.status. // returns *tview.Box
			SetBackgroundColor(tview.Styles.MoreContrastBackgroundColor).
			SetBorder(false)

	// Create the input field
	inputStyle := tcell.Style{}.
		Background(tview.Styles.PrimitiveBackgroundColor).
		Foreground(tview.Styles.PrimaryTextColor)

	tui.input = tview.NewInputField().
		SetLabel("> ").
		SetPlaceholderStyle(inputStyle).
		SetFieldStyle(inputStyle).
		SetLabelStyle(inputStyle)

	// Handle input by handing it off to the commandCh
	tui.input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			t := tui.input.GetText()
			if t != "" {
				tui.history.Append(t)
				tui.commandCh <- t
				tui.input.SetText("")
			}

		default:
			return
		}
	})

	tui.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.output, 0, 1, false).
		AddItem(tui.status, 1, 0, false).
		AddItem(tui.input, 1, 0, true).
		SetFullScreen(true)

	tui.app = tview.NewApplication().SetRoot(tui.flex, true)

	// History function
	tui.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			tui.input.SetText(tui.history.Up())
			return nil

		case tcell.KeyDown:
			tui.input.SetText(tui.history.Down())
			return nil

		default:
			return event
		}
	})

	// Are we going to block Ctrl-C?
	if config.BlockCtrlC {
		tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlC:
				return nil

			default:
				return event
			}
		})
	}

	return tui
}

// Run the ChatUI
func (t *ChatUI) Run() error {
	go func() {
		for s := range t.outputCh {
			t.app.QueueUpdateDraw(func() {
				fmt.Fprintf(t.output, "\n%s", s)
			})
			t.output.ScrollToEnd()
		}
	}()

	return t.app.Run()
}

// Stop stops the ChatUI
func (t *ChatUI) Stop() {
	t.app.Stop()
}

// History retrieves command history
func (t *ChatUI) History(n int) []string {
	return t.history.History(n)
}

// SetStatus changes the text of the status line
func (t *ChatUI) SetStatus(status string) {
	t.app.QueueUpdateDraw(func() {
		t.status.SetText(status)
	})
}
