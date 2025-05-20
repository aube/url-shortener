package workerpool

import (
	"context"
	"sync"
)

// Work represents a unit of work to be processed by the worker pool.
// It contains the context for cancellation/timeout, the hash to process,
// and the user ID associated with the work.
type Work struct {
	ctx    context.Context // Context for cancellation and deadlines
	hash   string          // The hash to be processed
	userID string          // User ID associated with the work
}

// WorkDispatcher manages a pool of workers that process Work items.
// It provides methods to add work, close the dispatcher, and manages
// the distribution of work across workers and collection of results.
type WorkDispatcher struct {
	input chan Work // Channel for receiving work items
}

// New creates and initializes a new WorkDispatcher with a worker pool.
// Parameters:
//   - numWorkers: Number of worker goroutines to create
//   - processor: Function that will process each work item. It should accept
//     a context, hash string, and userID string, and return an error.
//
// The WorkDispatcher:
//   - Creates worker goroutines that listen for work
//   - Implements fan-in pattern to collect results from all workers
//   - Starts a goroutine to handle processing results
//
// Example:
//
//	processor := func(ctx context.Context, hash, userID string) error {
//	    // Process the work
//	    return nil
//	}
//	dispatcher := workerpool.New(4, processor)
func New(numWorkers int, processor func(context.Context, string, string) error) *WorkDispatcher {
	wd := &WorkDispatcher{
		input: make(chan Work),
	}

	// Start workers
	workerChannels := make([]<-chan error, numWorkers)
	for i := range numWorkers {
		workerChannels[i] = wd.worker(processor)
	}

	// Start fan-in goroutine to combine results from all workers
	output := wd.fanIn(workerChannels...)

	// Start processing results
	go func() {
		for result := range output {
			println(result)
		}
	}()

	return wd
}

// AddWork submits a new work item to the worker pool.
// Parameters:
//   - ctx: Context for cancellation/timeout
//   - hash: The hash to be processed
//   - userID: User ID associated with the work
//
// This method is thread-safe and can be called from multiple goroutines.
func (wd *WorkDispatcher) AddWork(ctx context.Context, hash string, userID string) {
	wd.input <- Work{
		ctx:    ctx,
		hash:   hash,
		userID: userID,
	}
}

// Close shuts down the WorkDispatcher by closing the input channel.
// This should be called when no more work will be added to ensure
// proper cleanup of worker goroutines.
func (wd *WorkDispatcher) Close() {
	close(wd.input)
}

// worker creates a worker goroutine that processes work items using the provided processor function.
// Parameters:
//   - processor: Function that handles the actual work processing
//
// Returns:
//   - <-chan error: Channel that emits processing results/errors
//
// Each worker:
//   - Listens for work on the shared input channel
//   - Processes work using the provided function
//   - Sends results/errors to its output channel
func (wd *WorkDispatcher) worker(processor func(context.Context, string, string) error) <-chan error {
	output := make(chan error)
	go func() {
		for n := range wd.input {
			output <- processor(n.ctx, n.hash, n.userID)
		}
		close(output)
	}()
	return output
}

// fanIn combines multiple error channels into a single output channel.
// Parameters:
//   - channels: Variadic list of error channels to combine
//
// Returns:
//   - <-chan error: Single channel that receives from all input channels
//
// This implements the fan-in pattern to collect results from multiple workers.
// It uses a WaitGroup to ensure all channels are properly drained before
// closing the output channel.
func (wd *WorkDispatcher) fanIn(channels ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	output := make(chan error)

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan error) {
			for item := range c {
				output <- item
			}
			wg.Done()
		}(ch)
	}

	// Start goroutine to close output when all inputs are done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
