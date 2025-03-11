package jmespath_test

import (
	"fmt"

	"github.com/woodsbury/decimal128"
	"github.com/woodsbury/jmespath"
)

func ExampleSearch() {
	value := map[string]any{
		"Field": decimal128.FromUint32(2),
	}

	result, _ := jmespath.Search("Field + `1`", value)
	fmt.Println(result)
	// Output:
	// 3
}

func ExampleCompile() {
	value := map[string]any{
		"Field": decimal128.FromUint32(2),
	}

	expression, _ := jmespath.Compile("Field + `1`")
	result, _ := expression.Search(value)
	fmt.Println(result)
	// Output:
	// 3
}
