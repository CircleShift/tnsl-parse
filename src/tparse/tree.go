package tparse

// Node represents a group of nodes or a directive
type Node struct {
	SubNodes []Node

	Dir Directive
}

// Directive represents a block or single directive
type Directive struct {
	Type string
	Data string

	Param Paramaters
}

// Paramaters represents a set of paramaters for a directive
type Paramaters struct {
	In  []string
	Out []string
}

func handleCode(tokens *[]Token, start int) (Node, int) {
	out := Node{}

	return out, start
}

func handleBlock(tokens *[]Token, start int) (Node, int) {
	var out Node
	var tmp Node

	l := len(*tokens)

	if start >= l {
		panic((*tokens)[l-1])
	}

	for ; start < l; start++ {
		t := (*tokens)[start]
		switch t.Type {
		case LINESEP:
			if t.Data == ";" {
				tmp, start = handleCode(tokens, start+1)
			}
			break
		case DELIMIT:
			if t.Data == "/;" {
				tmp, start = handleCode(tokens, start+1)
			}
			break
		default:
			panic(t)
		}
		out.SubNodes = append(out.SubNodes, tmp)
	}

	return out, start
}

func handlePre(tokens *[]Token, start int) (Node, int) {
	out := Node{}

	return out, start
}

// CreateTree takes a series of tokens and converts them into an AST
func CreateTree(tokens *[]Token, start int) Node {
	out := Node{}
	out.Dir = Directive{Type: "root"}

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
