package lexer
// Token类型
type TokenType string

const (
	IDENT       TokenType = "IDENTIFIER"
	NUMBER      TokenType = "NUMBER"
	HEX         TokenType = "HEX_NUMBER"
	OCTAL       TokenType = "OCTAL_NUMBER"
	BINARY      TokenType = "BINARY_NUMBER"
	FLOAT       TokenType = "FLOAT_NUMBER"
	CHAR        TokenType = "CHAR_CONSTANT"
	STRING      TokenType = "STRING_CONSTANT"
	LPAREN      TokenType = "LPAREN"      // (
	RPAREN      TokenType = "RPAREN"      // )
	LBRACE      TokenType = "LBRACE"      // {
	RBRACE      TokenType = "RBRACE"      // }
	SEMICOLON   TokenType = "SEMICOLON"   // ;
	ASSIGN      TokenType = "ASSIGN"      // =
	PLUS        TokenType = "PLUS"        // +
	MINUS       TokenType = "MINUS"       // -
	ASTERISK    TokenType = "ASTERISK"    // *
	SLASH       TokenType = "SLASH"       // /
	LT          TokenType = "LT"          // <
	GT          TokenType = "GT"          // >
	LTE         TokenType = "LTE"         // <=
	GTE         TokenType = "GTE"         // >=
	EQ          TokenType = "EQ"          // ==
	NEQ         TokenType = "NEQ"         // !=
	AND         TokenType = "AND"         // &&
	OR          TokenType = "OR"          // ||
	NOT         TokenType = "NOT"         // !
	PLUS_ASSIGN  TokenType = "PLUS_ASSIGN"  // +=
	MINUS_ASSIGN TokenType = "MINUS_ASSIGN" // -=
	ASTERISK_ASSIGN TokenType = "ASTERISK_ASSIGN" // *=
	SLASH_ASSIGN TokenType = "SLASH_ASSIGN" // /=
	UNKNOWN     TokenType = "UNKNOWN"
	IF          TokenType = "IF"
	ELSE        TokenType = "ELSE"
	WHILE       TokenType = "WHILE"
	DO          TokenType = "DO"
	MAIN        TokenType = "MAIN"
	INT         TokenType = "INT"
	FLOAT_TYPE  TokenType = "FLOAT"
	DOUBLE      TokenType = "DOUBLE"
	RETURN      TokenType = "RETURN"
	CONST       TokenType = "CONST"
	VOID        TokenType = "VOID"
	CONTINUE    TokenType = "CONTINUE"
	BREAK       TokenType = "BREAK"
	CHAR_TYPE   TokenType = "CHAR"
	UNSIGNED    TokenType = "UNSIGNED"
	ENUM        TokenType = "ENUM"
	LONG        TokenType = "LONG"
	SWITCH      TokenType = "SWITCH"
	CASE        TokenType = "CASE"
	AUTO        TokenType = "AUTO"
	STATIC      TokenType = "STATIC"
)


// 保留字表
var reservedWords = map[string]TokenType{
    "if":       IF,
    "else":     ELSE,
    "while":    WHILE,
    "do":       DO,
    "main":     MAIN,
    "int":      INT,
    "float":    FLOAT_TYPE,
    "double":   DOUBLE,
    "return":   RETURN,
    "const":    CONST,
    "void":     VOID,
    "continue": CONTINUE,
    "break":    BREAK,
    "char":     CHAR_TYPE,
    "unsigned": UNSIGNED,
    "enum":     ENUM,
    "long":     LONG,
    "switch":   SWITCH,
    "case":     CASE,
    "auto":     AUTO,
    "static":   STATIC,
    "+=": PLUS_ASSIGN,
    "-=": MINUS_ASSIGN,
    "*=": ASTERISK_ASSIGN,
    "/=": SLASH_ASSIGN,
}

// Token结构
type Token struct {
    Type       TokenType
    Value      string
    Error      string
}

// Lexer结构
type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           byte
    prevToken    *Token
}
