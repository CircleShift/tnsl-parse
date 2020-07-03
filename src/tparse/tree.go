package tparse

// Node represents a group of nodes or a directive
type Node struct {
	SubNodes []Node

	Dir Directive
}

// Directive represents a block or single directive
type Directive struct {
	Type string
	ID   string

	Data []string
}

func handleCode(tokens *[]Token, start int) (Node, int) {
	out := Node{}

	return out, start
}

func handlePre(tokens *[]Token, start int) (Node, int) {
	out := Node{}

	return out, start
}

// CreateTree takes a series of tokens and converts them into an AST
func CreateTree(tokens *[]Token, start int) Node {
	out := Node{}
	out.Dir = Directive{Type: "root", ID: "root"}

	var tmp Node

	for i, t := range *tokens {
		switch t.Type {
		case LINESEP:
			if t.Data == ";" {
				tmp, i = handleCode(tokens, i)
			} else if t.Data == ":" {
				tmp, i = handlePre(tokens, i)
			}
			break
		case DELIMIT:
			if t.Data == "/;" {
				tmp, i = handleCode(tokens, i)
			} else if t.Data == "/:" {
				tmp, i = handlePre(tokens, i)
			}
			break
		}
		out.SubNodes = append(out.SubNodes, tmp)
	}

	return out
}
