package tparse

// Token represents a token in a program
type Token struct {
	Type int
	Data string
	Line int
	Char int
}
