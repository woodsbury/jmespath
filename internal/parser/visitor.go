package parser

import (
	"fmt"
	"io"
)

type Visitor interface {
	Visit(Node)
}

type Walker interface {
	Walk(Visitor)
}

func WriteTo(w io.Writer, node Node) error {
	v := writeVisitor{
		w: w,
	}

	v.Visit(node)
	return v.err
}

type writeVisitor struct {
	w      io.Writer
	err    error
	indent int
}

var indentBytes = []byte("  ")

func (v *writeVisitor) Visit(node Node) {
	indent := v.indent
	for indent > 0 {
		_, err := v.w.Write(indentBytes)
		if err != nil {
			v.err = err
		}

		indent--
	}

	_, err := fmt.Fprintf(v.w, "%s\n", node.String())
	if err != nil {
		v.err = err
	}

	if v.err != nil {
		return
	}

	if w, ok := node.(Walker); ok {
		v.indent++
		w.Walk(v)
		v.indent--
	}
}
