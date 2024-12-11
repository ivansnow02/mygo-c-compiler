package rec_des_parser

import (
    "go_compiler/lexer"
)

type Parser struct {
    Result string
    lexer  *lexer.Lexer
}
