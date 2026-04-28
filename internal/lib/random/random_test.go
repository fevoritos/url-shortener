package random

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	t.Run("correct_length", func(t *testing.T) {
		size := 10
		res := NewRandomString(size)
		assert.Len(t, res, size)
	})

	t.Run("only_allowed_chars", func(t *testing.T) {
		size := 100
		res := NewRandomString(size)

		for _, char := range res {
			assert.True(t, strings.ContainsRune(charset, char),
				"string contains char not from charset: %c", char)
		}
	})

	t.Run("is_random", func(t *testing.T) {
		res1 := NewRandomString(10)
		res2 := NewRandomString(10)
		assert.NotEqual(t, res1, res2)
	})

	t.Run("zero_length", func(t *testing.T) {
		assert.Equal(t, "", NewRandomString(0))
	})
}

func BenchmarkNewRandomString(b *testing.B) {
	for b.Loop() {
		_ = NewRandomString(10)
	}
}
