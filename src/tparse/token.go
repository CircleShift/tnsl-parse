package tparse

// Token represents a token in a program
type Token struct {
	Type int
	Data string
}

// Container represents a container of data
type Container struct {
	Data  []interface{}
	Holds bool
}
