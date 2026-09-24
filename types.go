package jmespath

type Types uint8

const (
	TypeArray Types = 1 << iota
	TypeObject
	TypeBoolean
	TypeNumber
	TypeString

	TypeNull Types = 0
	TypeAny  Types = TypeBoolean | TypeNumber | TypeString | TypeArray | TypeObject
)

var typeNames = [...]string{
	"array",
	"object",
	"boolean",
	"number",
	"string",
}

func (t Types) String() string {
	if t&TypeAny == TypeAny {
		return "any"
	}

	s := ""
	for i, name := range typeNames {
		if t&(1<<i) != 0 {
			if s != "" {
				s += "|"
			}

			s += name
		}
	}

	if s == "" {
		return "null"
	}

	return s
}
