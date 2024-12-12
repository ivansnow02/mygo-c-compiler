package rec_des_parser

import (
	"mygo_c_compiler/lexer"
)

type Parser struct {
	Result string
	lexer  *lexer.Lexer
}
