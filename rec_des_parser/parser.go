package rec_des_parser

import (
	"fmt"
	"mygo_c_compiler/lexer"
)

func New() *Parser {
	return &Parser{}
}

func (g *Parser) Parse(input string) {
	g.lexer = lexer.NewLexer(input)
	g.program()
}

func (g *Parser) match(tokenType lexer.TokenType) lexer.Token {
	token := g.lexer.NextToken()
	if token.Type != tokenType {
		panic(fmt.Sprintf("Syntax Error: expected %s, got %s", tokenType, token.Type))
	}
	return token
}

func (g *Parser) program() {
	g.block()
}

func (g *Parser) block() {
	fmt.Println("Entering block")
	fmt.Println("block -> { stmts }")
	g.match(lexer.LBRACE)
	g.stmts()
	g.match(lexer.RBRACE)
}

func (g *Parser) stmts() {
	fmt.Println("Entering stmts")
	token := g.lexer.NextToken()
	switch token.Type {
	case lexer.IF, lexer.WHILE, lexer.DO, lexer.BREAK, lexer.IDENT, lexer.LBRACE:
		g.lexer.UnreadToken(token)
		fmt.Println("stmts -> stmt stmts")
		g.stmt()
		g.stmts()
	default:
		g.lexer.UnreadToken(token)
		fmt.Println("stmts -> ε")
		return
	}
}

func (g *Parser) stmt() {
	fmt.Println("Entering stmt")
	token := g.lexer.NextToken()
	g.lexer.UnreadToken(token)
	switch token.Type {
	case lexer.IF:
		fmt.Println("stmt -> if ( bool ) stmt stmt'")
		g.match(lexer.IF)
		g.match(lexer.LPAREN)
		g.boolExpr()
		g.match(lexer.RPAREN)
		g.stmt()
		g.stmtPrime()
	case lexer.IDENT:
		fmt.Println("stmt -> id = expr ;")
		g.match(lexer.IDENT)
		g.match(lexer.ASSIGN)
		g.expr()
		g.match(lexer.SEMICOLON)
	case lexer.WHILE:
		fmt.Println("stmt -> while ( bool ) stmt")
		g.match(lexer.WHILE)
		g.match(lexer.LPAREN)
		g.boolExpr()
		g.match(lexer.RPAREN)
		g.stmt()
	case lexer.DO:
		fmt.Println("stmt -> do stmt while ( bool ) ;")
		g.match(lexer.DO)
		g.stmt()
		g.match(lexer.WHILE)
		g.match(lexer.LPAREN)
		g.boolExpr()
		g.match(lexer.RPAREN)
		g.match(lexer.SEMICOLON)
	case lexer.BREAK:
		fmt.Println("stmt -> break ;")
		g.match(lexer.BREAK)
		g.match(lexer.SEMICOLON)
	case lexer.LBRACE:
		fmt.Println("stmt -> block")
		g.block()
	default:
		panic("Syntax Error")
	}
}

func (g *Parser) stmtPrime() {
	fmt.Println("Entering stmt'")
	token := g.lexer.NextToken()
	if token.Type == lexer.ELSE {
		fmt.Println("stmt' -> else stmt")
		g.match(lexer.ELSE)
		g.stmt()
	} else {
		g.lexer.UnreadToken(token)
		fmt.Println("stmt' -> ε")
		return
	}
}

func (g *Parser) boolExpr() {
	fmt.Println("Entering bool")
	fmt.Println("bool -> expr bool'")
	g.expr()
	g.boolPrime()
}

func (g *Parser) boolPrime() {
	fmt.Println("Entering bool'")
	token := g.lexer.NextToken()
	switch token.Type {
	case lexer.LT:
		g.lexer.UnreadToken(token)
		g.match(lexer.LT)
		g.expr()
	case lexer.LTE:
		g.lexer.UnreadToken(token)
		g.match(lexer.LTE)
		g.expr()
	case lexer.GT:
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.GT)
		g.expr()
	case lexer.GTE:
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.GTE)
		g.expr()
	case lexer.EQ:
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.EQ)
		g.expr()
	case lexer.NEQ:
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.NEQ)
		g.expr()
	default:
		g.lexer.UnreadToken(token)

		fmt.Println("bool' -> ε")
		return
	}
}

func (g *Parser) expr() {
	fmt.Println("Entering expr")
	fmt.Println("expr -> term expr'")
	g.term()
	g.exprPrime()
}

func (g *Parser) exprPrime() {
	fmt.Println("Entering expr'")
	token := g.lexer.NextToken()
	switch token.Type {
	case lexer.PLUS:
		g.lexer.UnreadToken(token)
		g.match(lexer.PLUS)
		g.term()
		g.exprPrime()
	case lexer.MINUS:
		g.lexer.UnreadToken(token)
		g.match(lexer.MINUS)
		g.term()
		g.exprPrime()
	default:
		g.lexer.UnreadToken(token)
		fmt.Println("expr' -> ε")
		return
	}
}

func (g *Parser) term() {
	fmt.Println("Entering term")
	fmt.Println("term -> factor term'")
	g.factor()
	g.termPrime()
}

func (g *Parser) termPrime() {
	fmt.Println("Entering term'")
	token := g.lexer.NextToken()
	switch token.Type {
	case lexer.ASTERISK:
		g.lexer.UnreadToken(token)
		g.match(lexer.ASTERISK)
		g.factor()
		g.termPrime()
	case lexer.SLASH:
		g.lexer.UnreadToken(token)
		g.match(lexer.SLASH)
		g.factor()
		g.termPrime()
	default:
		g.lexer.UnreadToken(token)
		fmt.Println("term' -> ε")
		return
	}
}

func (g *Parser) factor() {
	fmt.Println("Entering factor")
	token := g.lexer.NextToken()
	switch token.Type {
	case lexer.LPAREN:
		fmt.Println("factor -> ( expr )")
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.LPAREN)
		g.expr()
		g.match(lexer.RPAREN)
	case lexer.IDENT:
		fmt.Println("factor -> id")
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.IDENT)
	case lexer.NUMBER:
		fmt.Println("factor -> num")
		g.lexer.UnreadToken(token) // 先放回token
		g.match(lexer.NUMBER)
	default:
		panic(fmt.Sprintf("Syntax Error: unexpected token %s", token.Type))
	}
}
