package parser

type Types uint8

const (
	TypeArray Types = 1 << iota
	TypeObject
	TypeBoolean
	TypeNumber
	TypeString

	TypeNull Types = 0
	TypeAny  Types = TypeArray | TypeObject | TypeBoolean | TypeNumber | TypeString
)
