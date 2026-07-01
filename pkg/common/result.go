package common

// Result is anything that can render itself as a human-readable string.
type Result interface {
	String() string
}
