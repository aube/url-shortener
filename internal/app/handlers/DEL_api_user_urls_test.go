package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFanOutFanIn(t *testing.T) {
	t.Run("processes all inputs through workers", func(t *testing.T) {
		doneCh := make(chan struct{})
		defer close(doneCh)

		inputCh := make(chan string)
		go func() {
			defer close(inputCh)
			inputCh <- "a"
			inputCh <- "b"
			inputCh <- "c"
		}()

		channels := fanOut(doneCh, inputCh)
		resultCh := fanIn(doneCh, channels...)

		var results []string
		for val := range resultCh {
			results = append(results, val)
		}

		assert.ElementsMatch(t, []string{"a", "b", "c"}, results)
	})

	t.Run("handles early termination", func(t *testing.T) {
		doneCh := make(chan struct{})
		inputCh := make(chan string)

		go func() {
			defer close(inputCh)
			inputCh <- "a"
			close(doneCh)
			inputCh <- "b" // Should be ignored
		}()

		channels := fanOut(doneCh, inputCh)
		resultCh := fanIn(doneCh, channels...)

		var results []string
		for val := range resultCh {
			results = append(results, val)
		}

		assert.LessOrEqual(t, len(results), 1) // Only "a" might get through
	})
}
