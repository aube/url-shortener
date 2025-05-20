package workerpool

// func TestNewWorkDispatcher(t *testing.T) {
// 	t.Run("creates workers and processes work", func(t *testing.T) {
// 		var mu sync.Mutex
// 		var processed []int
// 		processor := func(n int) string {
// 			mu.Lock()
// 			defer mu.Unlock()
// 			processed = append(processed, n)
// 			return ""
// 		}

// 		wd := New(3, processor)

// 		// Add work items
// 		for i := 1; i <= 5; i++ {
// 			wd.AddWork(i)
// 		}

// 		// Give workers time to process
// 		time.Sleep(100 * time.Millisecond)
// 		wd.Close()

// 		mu.Lock()
// 		defer mu.Unlock()
// 		assert.Len(t, processed, 5)
// 		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, processed)
// 	})
// }

// func TestWorkerDispatcher_AddWork(t *testing.T) {
// 	t.Run("sends work to workers", func(t *testing.T) {
// 		var processed []int
// 		var wg sync.WaitGroup
// 		wg.Add(1)

// 		processor := func(n int) string {
// 			defer wg.Done()
// 			processed = append(processed, n)
// 			return ""
// 		}

// 		wd := New(1, processor)
// 		wd.AddWork(42)
// 		wg.Wait()
// 		wd.Close()

// 		assert.Equal(t, []int{42}, processed)
// 	})
// }

// func TestWorkerDispatcher_Close(t *testing.T) {
// 	t.Run("stops workers after close", func(t *testing.T) {
// 		processor := func(n int) string {
// 			return ""
// 		}

// 		wd := New(1, processor)
// 		wd.Close()

// 		// Verify channel is closed
// 		_, ok := <-wd.input
// 		assert.False(t, ok)
// 	})
// }

// func TestWorkerDispatcher_worker(t *testing.T) {
// 	t.Run("processes work items", func(t *testing.T) {
// 		wd := &WorkDispatcher{
// 			input: make(chan int),
// 		}

// 		processor := func(n int) string {
// 			return string(rune(n + 64)) // Convert to letter (1->A, 2->B, etc)
// 		}

// 		output := wd.worker(processor)
// 		go func() {
// 			wd.input <- 1
// 			wd.input <- 2
// 			close(wd.input)
// 		}()

// 		var results []string
// 		for result := range output {
// 			results = append(results, result)
// 		}

// 		assert.Equal(t, []string{"A", "B"}, results)
// 	})
// }

// func TestWorkerDispatcher_fanIn(t *testing.T) {
// 	t.Run("combines multiple channels", func(t *testing.T) {
// 		wd := &WorkDispatcher{}

// 		// Create test channels
// 		ch1 := make(chan string)
// 		ch2 := make(chan string)

// 		// Start fanIn
// 		combined := wd.fanIn(ch1, ch2)

// 		// Send test data
// 		go func() {
// 			ch1 <- "A"
// 			ch2 <- "B"
// 			ch1 <- "C"
// 			close(ch1)
// 			close(ch2)
// 		}()

// 		// Collect results
// 		var results []string
// 		for result := range combined {
// 			results = append(results, result)
// 		}

// 		assert.ElementsMatch(t, []string{"A", "B", "C"}, results)
// 	})
// }
