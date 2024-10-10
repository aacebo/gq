package query

type _TokenKind uint8

const (
	EOF _TokenKind = iota

	// special characters

	COMMA       // ,
	COLON       // :
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // {
	RIGHT_BRACE // }

	// literals

	IDENTIFIER // argument name
	STRING     // "test"
	BYTE       // 'a'
	INT        // 1
	FLOAT      // 1.0
	NULL       // null
	BOOL       // true, false
)

var Keywords = map[string]_TokenKind{
	"true":  BOOL,
	"false": BOOL,
}
