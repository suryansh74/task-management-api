# Days 13–14 — Interview prep

## Part A: Find the concurrency bug

Study this broken code. Find the bug(s), explain the failure mode, and fix it.

```go
package main

import (
	"fmt"
	"sync"
)

// Broken counter — classic interview trap
type Counter struct {
	n int
}

func (c *Counter) Inc() { c.n++ }
func (c *Counter) Value() int { return c.n }

func main() {
	var c Counter
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(c.Value()) // often NOT 1000
}
```

### What is wrong?
1. **Data race** on `c.n` — concurrent read/write without synchronization.
2. Expected 1000, actual is random lower number (lost updates).

### How to detect
```bash
go run -race .
```

### Fixes (pick one and justify)
- `sync.Mutex` around Inc/Value
- `atomic.AddInt64` / `atomic.LoadInt64`
- Channel-based single owner (overkill here)

### Harder variant (also practice)

```go
// Cache with double-checked locking — often still wrong
type Cache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *Cache) GetOrLoad(key string, load func() string) string {
	if v, ok := c.data[key]; ok { // BUG: read without lock
		return v
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.data[key]; ok {
		return v
	}
	if c.data == nil {
		c.data = make(map[string]string)
	}
	v := load()
	c.data[key] = v
	return v
}
```

Bug: first map read is unlocked → data race. Fix: lock before any map access, or use `sync.Map` with care.

---

## Part B: Production-judgment stories (STAR)

Prepare 2–3 stories from **your** projects. Template:

| Section | Content |
|---------|---------|
| **Situation** | Context (project, scale, constraint) |
| **Task** | What you owned |
| **Action** | What you did (and alternatives you rejected) |
| **Result** | Measurable or qualitative outcome |
| **What I'd change now** | Senior judgment — hindsight |

### Story ideas from your repos

**1. task-management-api — caching / consistency**
- Situation: Cache-aside on tasks; risk of stale reads after update.
- Task: Make reads fast without wrong data.
- Action: Invalidate on write; TTL as safety net.
- Change now: metric on hit/miss ratio; maybe request coalescing; versioned keys.

**2. task-management-api — session auth vs gRPC**
- Situation: Cookie sessions work for REST; gRPC needed user_id in request.
- Task: Dual protocol without rewriting core.
- Action: Shared `ports.TaskService`; thin gRPC adapter.
- Change now: metadata-based auth interceptor; one auth model for both.

**3. rabbitmq_implementation_with_dlx — retries / DLX**
- Situation: Failed messages need retry without poison-pill loops.
- Task: Reliable processing path.
- Action: DLX + retry queue pattern (document your actual design).
- Change now: max retry count metric; dead-letter dashboard; idempotent consumers.

**4. concurrent-rate-limited-fetcher — RPS vs workers**
- Situation: Bulk fetch without bans or hangs.
- Task: Control concurrency and rate.
- Action: Worker pool + token bucket + per-request timeout + context cancel.
- Change now: per-host limiters; adaptive RPS; structured result metrics.

### Practice aloud
For each story, 90–120 seconds. End with: “Knowing what I know now, I would …”
