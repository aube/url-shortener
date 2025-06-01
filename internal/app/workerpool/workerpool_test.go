package workerpool

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestProcessor is a mock processor function for testing
type TestProcessor struct {
	mock.Mock
}

func (p *TestProcessor) Process(ctx context.Context, hash, userID string) error {
	args := p.Called(ctx, hash, userID)
	return args.Error(0)
}

func TestNewWorkDispatcher(t *testing.T) {
	t.Run("creates dispatcher with specified number of workers", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		numWorkers := 3
		dispatcher := New(numWorkers, processor.Process)

		// Verify the dispatcher is created
		assert.NotNil(t, dispatcher)
		assert.NotNil(t, dispatcher.input)

		// Add some work to ensure workers are running
		for i := 0; i < 5; i++ {
			dispatcher.AddWork(context.Background(), "hash1", "user1")
		}

		// Give workers time to process
		time.Sleep(100 * time.Millisecond)

		// Clean up
		dispatcher.Close()

		// Verify processor was called
		processor.AssertExpectations(t)
	})
}

func TestAddWork(t *testing.T) {
	t.Run("adds work to the pool", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, "testHash", "testUser").Return(nil)

		dispatcher := New(1, processor.Process)
		defer dispatcher.Close()

		dispatcher.AddWork(context.Background(), "testHash", "testUser")

		// Give worker time to process
		time.Sleep(100 * time.Millisecond)

		processor.AssertExpectations(t)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, "testHash", "testUser").Return(context.Canceled)

		dispatcher := New(1, processor.Process)
		defer dispatcher.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dispatcher.AddWork(ctx, "testHash", "testUser")

		// Give worker time to process
		time.Sleep(100 * time.Millisecond)

		processor.AssertExpectations(t)
	})
}

func TestWorker(t *testing.T) {
	t.Run("processes work and returns result", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, "hash1", "user1").Return(nil)
		processor.On("Process", mock.Anything, "hash2", "user2").Return(errors.New("processing error"))

		dispatcher := &WorkDispatcher{
			input: make(chan Work),
		}

		output := dispatcher.worker(processor.Process)

		// Send test work
		go func() {
			dispatcher.input <- Work{context.Background(), "hash1", "user1"}
			dispatcher.input <- Work{context.Background(), "hash2", "user2"}
			close(dispatcher.input)
		}()

		// Collect results
		var results []error
		for err := range output {
			results = append(results, err)
		}

		// Verify results
		assert.Len(t, results, 2)
		assert.Nil(t, results[0])
		assert.EqualError(t, results[1], "processing error")

		processor.AssertExpectations(t)
	})
}

func TestFanIn(t *testing.T) {
	t.Run("combines multiple channels into one", func(t *testing.T) {
		dispatcher := &WorkDispatcher{}

		// Create test channels
		ch1 := make(chan error)
		ch2 := make(chan error)

		// Start fan-in
		output := dispatcher.fanIn(ch1, ch2)

		// Send test data
		go func() {
			ch1 <- nil
			ch2 <- errors.New("error from ch2")
			ch1 <- errors.New("error from ch1")
			close(ch1)
			close(ch2)
		}()

		// Collect results and cast into string
		var results []string
		for err := range output {
			results = append(results, fmt.Sprint(err))
		}

		// Verify results
		assert.Len(t, results, 3)
		assert.Equal(t, results[0], "<nil>")
		assert.True(t, slices.Contains(results, "error from ch2"))
		assert.True(t, slices.Contains(results, "error from ch1"))
	})
}

func TestClose(t *testing.T) {
	t.Run("closes input channel and stops workers", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		dispatcher := New(2, processor.Process)

		// Add some work
		dispatcher.AddWork(context.Background(), "hash1", "user1")
		dispatcher.AddWork(context.Background(), "hash2", "user2")

		// Close the dispatcher
		dispatcher.Close()

		// Verify we can't add more work
		assert.Panics(t, func() {
			dispatcher.AddWork(context.Background(), "hash3", "user3")
		})

		// Give workers time to finish
		time.Sleep(100 * time.Millisecond)

		processor.AssertExpectations(t)
	})
}

func TestConcurrentAddWork(t *testing.T) {
	t.Run("handles concurrent work additions", func(t *testing.T) {
		processor := &TestProcessor{}
		processor.On("Process", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(100)

		dispatcher := New(5, processor.Process)
		defer dispatcher.Close()

		// Add work concurrently
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				dispatcher.AddWork(context.Background(), "hash", "user")
			}(i)
		}
		wg.Wait()

		// Give workers time to process
		time.Sleep(200 * time.Millisecond)

		processor.AssertExpectations(t)
	})
}
