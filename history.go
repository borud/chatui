package chatui

// History represents a command history
type History struct {
	index   int
	maxSize int
	history []string
}

// NewHistory creates new History instance.
func NewHistory(maxHistorySize int) *History {
	return &History{
		index:   0,
		maxSize: maxHistorySize,
		history: []string{},
	}
}

// Append a new history item
func (h *History) Append(item string) {
	// avoid recording when you just repeat the last command
	if len(h.history) > 0 && h.history[len(h.history)-1] == item {
		return
	}

	// if history has reached max size, we lop off the front of the history
	if len(h.history) == h.maxSize {
		h.history = h.history[1:]
	}

	h.history = append(h.history, item)
	h.index = len(h.history) - 1
}

// Up represents going up the history, returning what's there
func (h *History) Up() string {
	if len(h.history) == 0 {
		return ""
	}

	ret := h.history[h.index]

	h.index--

	if h.index < 0 {
		h.index = 0
	}

	return ret
}

// Down represents going down the history, returning what's there
func (h *History) Down() string {
	if len(h.history) == 0 {
		return ""
	}

	h.index++

	if h.index > (len(h.history) - 1) {
		h.index = len(h.history) - 1
		return ""
	}

	return h.history[h.index]
}

// History returns the last n elements of the history or all of the history if equal to 0.
func (h *History) History(n int) []string {
	if n == 0 {
		return h.history
	}

	if n > (len(h.history) - 1) {
		n = len(h.history) - 1
	}

	return h.history[:n]
}
