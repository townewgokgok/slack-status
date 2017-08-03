package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapEmoji(t *testing.T) {
	var txt string

	txt = WrapEmoji("")
	assert.Equal(t, "", txt, "")

	txt = WrapEmoji("+1")
	assert.Equal(t, ":+1:", txt, "")
}

func TestLimitStringByLength(t *testing.T) {
	var txt string

	txt = LimitStringByLength(`ASCII`, 4)
	assert.Equal(t, `ASC…`, txt, "")

	txt = LimitStringByLength(`日本語です`, 4)
	assert.Equal(t, `日本語…`, txt, "")

	txt = LimitStringByLength(`?????`, 4)
	assert.Equal(t, `???…`, txt, "")
}

func TestSplitEmoji(t *testing.T) {
	var emj, txt string

	emj, txt = SplitEmoji(` foo bar baz `)
	assert.Equal(t, ``, emj, "")
	assert.Equal(t, `foo bar baz`, txt, "")

	emj, txt = SplitEmoji(` :+1: foo bar baz `)
	assert.Equal(t, `:+1:`, emj, "")
	assert.Equal(t, `foo bar baz`, txt, "")

	emj, txt = SplitEmoji(` :+1:foo bar baz `)
	assert.Equal(t, `:+1:`, emj, "")
	assert.Equal(t, `foo bar baz`, txt, "")

	emj, txt = SplitEmoji(` :: foo bar baz`)
	assert.Equal(t, ``, emj, "")
	assert.Equal(t, `:: foo bar baz`, txt, "")

	emj, txt = SplitEmoji(` foo :+1: bar baz `)
	assert.Equal(t, ``, emj, "")
	assert.Equal(t, `foo :+1: bar baz`, txt, "")
}
