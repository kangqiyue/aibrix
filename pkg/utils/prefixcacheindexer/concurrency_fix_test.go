package prefixcacheindexer

import (
	"sync"
	"testing"
	"time"
)

// Test_ConcurrentMapAccess_Fix verifies that the concurrent map access bug is resolved.
// It simulates the scenario described by the user but uses the thread-safe public APIs
// that are now used in the codebase.
func Test_ConcurrentMapAccess_Fix(t *testing.T) {
	cache := NewLPRadixCache(2)
	unMatchedTokens := []int{9906, 4435, 0}
	model := "m1"
	podName := "p1"

	// Setup: Add a node
	node, _, _ := cache.AddPrefix(unMatchedTokens, model, podName)
	if node == nil {
		t.Fatal("Failed to add prefix")
	}

	var wg sync.WaitGroup
	done := make(chan bool)

	// Simulate the read goroutine (e.g., in prefix_cache_preble.go)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Simulate continuous access similar to the bug report loop
		// But using the now thread-safe API
		for {
			select {
			case <-done:
				return
			default:
				// Simulate: if modelPods := currentNode.GetPodsForModel(ctx.Model); modelPods != nil { ... }
				// Accessing parent to match the bug report scenario
				if node.parent != nil {
					pods := node.parent.GetPodsForModel(model)
					for p := range pods {
						_ = p // Simulate read access
					}
				}
			}
		}
	}()

	time.Sleep(time.Millisecond * 100)

	// Simulate the write operation (Evict) which triggered the bug
	// This will call evictNode, which modifies the parent's maps
	cache.Evict(time.Now().Add(60 * time.Minute))

	close(done)
	wg.Wait()
}
