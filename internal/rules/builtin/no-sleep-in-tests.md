---
description: fixed-duration sleeps in tests papering over timing or synchronization issues
severity: error
include: ["**/*_test.go", "**/*_test.py", "**/test_*.py", "**/*.test.js", "**/*.test.ts", "**/*.spec.js", "**/*.spec.ts"]
exclude: []
---
Flag fixed-duration delays used to wait out a concurrent or asynchronous effect before asserting. This includes `time.Sleep`, `sleep()`, `asyncio.sleep`, `setTimeout` used as a wait, and equivalent fixed delays in any language.

Do not flag a sleep that is itself the behavior under test, such as timeout logic, rate limiting, or clock code. Also spare polling loops that have a deadline and check a real condition, and delays in skipped or disabled code.

Flagged:

```go
go refreshCache()
time.Sleep(100 * time.Millisecond)
assertFresh(t, cache)
```

Spared:

```go
deadline := time.Now().Add(time.Second)
for !cache.Fresh() && time.Now().Before(deadline) {
	time.Sleep(10 * time.Millisecond)
}
```
