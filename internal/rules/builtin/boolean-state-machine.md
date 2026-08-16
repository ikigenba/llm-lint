---
description: multiple booleans encoding one multi-state concept that should be a single enum
severity: error
exclude: []
---
Flag groups of two or more boolean fields, variables, or columns that together encode one concept with more than two states, where some combinations are invalid or meaningless. Look for mutual dependence: setting one flag forces another false, or changing state requires writing several flags together. Recommend one enum or status value as the source of truth.

For example, a ticket represented by `held` and `sold` is really an open/held/sold state: `held && sold` is meaningless, and transitions must coordinate both booleans.

Flagged:

```go
type Ticket struct {
	Held bool
	Sold bool
}
```

Do not flag a single boolean, genuinely independent booleans for which every combination is valid, booleans derived on the fly from one source of truth, or bit flags in an established flags or bitmask type.

Spared:

```go
type User struct {
	IsAdmin       bool
	EmailVerified bool
}
```
