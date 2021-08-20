package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTextFromInput(t *testing.T) {
	text:="<abcde84y505> hey"
	text = ParseTextFromInput(text)
	assert.Equal(t, text, " hey")

	text = "hey"
	text = ParseTextFromInput(text)
	assert.Equal(t, text, "hey")
}
