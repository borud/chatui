package chatui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistor(t *testing.T) {
	h := NewHistory(4)

	assert.Equal(t, "", h.Up())
	assert.Equal(t, 0, h.index)
	assert.Equal(t, "", h.Down())
	assert.Equal(t, 0, h.index)

	h.Append("foo")
	h.Append("bar")
	h.Append("baz")

	assert.Equal(t, 2, h.index)

	up := h.Up()
	assert.Equal(t, "baz", up)
	assert.Equal(t, 1, h.index)

	up = h.Up()
	assert.Equal(t, "bar", up)
	assert.Equal(t, 0, h.index)

	up = h.Up()
	assert.Equal(t, "foo", up)
	assert.Equal(t, 0, h.index)

	up = h.Up()
	assert.Equal(t, "foo", up)
	assert.Equal(t, 0, h.index)

	down := h.Down()
	assert.Equal(t, "bar", down)
	assert.Equal(t, 1, h.index)

	down = h.Down()
	assert.Equal(t, "baz", down)
	assert.Equal(t, 2, h.index)

	down = h.Down()
	assert.Equal(t, "", down)
	assert.Equal(t, 2, h.index)

	h.Append("one")
	assert.Equal(t, 3, h.index)

	down = h.Down()
	assert.Equal(t, "", down)
	assert.Equal(t, 3, h.index)

	down = h.Down()
	assert.Equal(t, "", down)
	assert.Equal(t, 3, h.index)

	up = h.Up()
	assert.Equal(t, "one", up)
	assert.Equal(t, 2, h.index)

	up = h.Up()
	assert.Equal(t, "baz", up)
	assert.Equal(t, 1, h.index)

	h.Append("two")
	assert.Equal(t, 3, h.index)

	h.Append("three")
	assert.Equal(t, 3, h.index)
}

func TestHistoryOverflow(t *testing.T) {
	h := NewHistory(3)
	h.Append("foo")
	h.Append("bar")
	h.Append("baz")
	h.Append("gazonk")
	h.Append("bazonk")

	require.Equal(t, 2, h.index)
	h.Up()
	require.Equal(t, 1, h.index)
	h.Up()
	require.Equal(t, 0, h.index)
	h.Up()
	require.Equal(t, 0, h.index)
	h.Up()
	require.Equal(t, 0, h.index)
}
