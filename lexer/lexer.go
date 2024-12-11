package lexer

import (
    "fmt"
    "unicode"
)

func NewLexer(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition++
}

func (l *Lexer) NextToken() Token {
    if l.prevToken != nil {
        tok := *l.prevToken
        l.prevToken = nil
        return tok
    }

    var tok Token
    l.skipWhitespace()

    if l.ch == 0 {
        return Token{Type: UNKNOWN, Value: ""}
    }

    switch l.ch {
    case '+':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: PLUS_ASSIGN, Value: "+="}
        } else {
            tok = Token{Type: PLUS, Value: string(l.ch)}
        }
    case '-':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: MINUS_ASSIGN, Value: "-="}
        } else {
            tok = Token{Type: MINUS, Value: string(l.ch)}
        }
    case '*':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: ASTERISK_ASSIGN, Value: "*="}
        } else {
            tok = Token{Type: ASTERISK, Value: string(l.ch)}
        }
    case '/':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: SLASH_ASSIGN, Value: "/="}
        } else {
            tok = Token{Type: SLASH, Value: string(l.ch)}
        }
    case '=':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: EQ, Value: "=="}
        } else {
            tok = Token{Type: ASSIGN, Value: string(l.ch)}
        }
    case '!':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: NEQ, Value: "!="}
        } else {
            tok = Token{Type: NOT, Value: string(l.ch)}
        }
    case '<':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: LTE, Value: "<="}
        } else {
            tok = Token{Type: LT, Value: string(l.ch)}
        }
    case '>':
        if l.peekChar() == '=' {
            l.readChar()
            tok = Token{Type: GTE, Value: ">="}
        } else {
            tok = Token{Type: GT, Value: string(l.ch)}
        }
    case '&':
        if l.peekChar() == '&' {
            l.readChar()
            tok = Token{Type: AND, Value: "&&"}
        } else {
            tok = Token{Type: AND, Value: string(l.ch)}
        }
    case '|':
        if l.peekChar() == '|' {
            l.readChar()
            tok = Token{Type: OR, Value: "||"}
        } else {
            tok = Token{Type: OR, Value: string(l.ch)}
        }
    case ';':
        tok = Token{Type: SEMICOLON, Value: string(l.ch)}
    case '{':
        tok = Token{Type: LBRACE, Value: string(l.ch)}
    case '}':
        tok = Token{Type: RBRACE, Value: string(l.ch)}
    case '(':
        tok = Token{Type: LPAREN, Value: string(l.ch)}
    case ')':
        tok = Token{Type: RPAREN, Value: string(l.ch)}
    case '\'':
        tok = l.readCharConstant()
    case '"':
        tok = l.readString()
    case '0':
        peek := l.peekChar()
        if peek == 'x' || peek == 'X' {
            l.readChar()
            l.readChar()
            tok = l.readHex()
            return tok
        } else if peek == 'b' || peek == 'B' {
            l.readChar()
            l.readChar()
            tok = l.readBinary()
            return tok
        } else if peek == 'o' || peek == 'O' {
            l.readChar()
            l.readChar()
            tok = l.readOctal()
            return tok
        } else {
            tok = l.readNumber()
            return tok
        }
    case '.':
        tok = l.readFloat()
        return tok
    default:
        if isLetter(l.ch) {
            ident := l.readIdentifier()
            if tokType, ok := reservedWords[ident]; ok {
                tok = Token{Type: tokType, Value: ident}
            } else {
                tok = Token{Type: IDENT, Value: ident}
            }
            return tok
        } else if isDigit(l.ch) {
            tok = l.readNumber()
            return tok
        } else {
            tok = Token{Type: UNKNOWN, Value: string(l.ch), Error: fmt.Sprintf("未知字符: '%c'", l.ch)}
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) UnreadToken(tok Token) {
    l.prevToken = &tok
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readNumber() Token {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    // 检查是否为浮点数
    if l.ch == '.' {
        return l.readFloatFrom(position)
    }
    // 检查是否有紧跟的字母或下划线，表示无效的标识符
    if isLetter(l.ch) || l.ch == '_' {
        invalidPart := ""
        for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
            invalidPart += string(l.ch)
            l.readChar()
        }
        return Token{
            Type:    UNKNOWN,
            Value:   l.input[position:l.position],
            Error:   fmt.Sprintf("无效的标识符: %s", l.input[position:l.position]),
        }
    }
    value := l.input[position:l.position]
    return Token{Type: NUMBER, Value: value}
}

func (l *Lexer) readHex() Token {
    position := l.position
    for isHexDigit(l.ch) {
        l.readChar()
    }
    value := l.input[position:l.position]
    // 检查后续字符是否为字母或数字，表示无效的十六进制数
    if isLetter(l.ch) || isDigit(l.ch) {
        invalidPart := ""
        for isLetter(l.ch) || isDigit(l.ch) {
            invalidPart += string(l.ch)
            l.readChar()
        }
        return Token{
            Type:    UNKNOWN,
            Value:   "0x" + value + invalidPart,
            Error:   fmt.Sprintf("无效的十六进制数: 0x%s%s", value, invalidPart),
        }
    }
    return Token{Type: HEX, Value: "0x" + value}
}

func (l *Lexer) readOctal() Token {
    position := l.position
    for isOctalDigit(l.ch) {
        l.readChar()
    }
    value := l.input[position:l.position]
    // 检查后续字符是否为字母或数字，表示无效的八进制数
    if isLetter(l.ch) || isDigit(l.ch) {
        invalidPart := ""
        for isLetter(l.ch) || isDigit(l.ch) {
            invalidPart += string(l.ch)
            l.readChar()
        }
        return Token{
            Type:    UNKNOWN,
            Value:   "0o" + value + invalidPart,
            Error:   fmt.Sprintf("无效的八进制数: 0o%s%s", value, invalidPart),
        }
    }
    return Token{Type: OCTAL, Value: "0o" + value}
}

func (l *Lexer) readBinary() Token {
    position := l.position
    for isBinaryDigit(l.ch) {
        l.readChar()
    }
    value := l.input[position:l.position]
    // 检查后续字符是否为字母或数字，表示无效的二进制数
    if isLetter(l.ch) || isDigit(l.ch) {
        invalidPart := ""
        for isLetter(l.ch) || isDigit(l.ch) {
            invalidPart += string(l.ch)
            l.readChar()
        }
        return Token{
            Type:    UNKNOWN,
            Value:   "0b" + value + invalidPart,
            Error:   fmt.Sprintf("无效的二进制数: 0b%s%s", value, invalidPart),
        }
    }
    return Token{Type: BINARY, Value: "0b" + value}
}

func (l *Lexer) readFloat() Token {
    return l.readFloatFrom(l.position)
}

func (l *Lexer) readFloatFrom(start int) Token {
    hasDot := false
    hasExp := false

    for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-' {
        if l.ch == '.' {
            if hasDot {
                break
            }
            hasDot = true
        }
        if l.ch == 'e' || l.ch == 'E' {
            if hasExp {
                break
            }
            hasExp = true
        }
        l.readChar()
    }
    value := l.input[start:l.position]
    return Token{Type: FLOAT, Value: value}
}

func (l *Lexer) readCharConstant() Token {
    l.readChar()
    start := l.position
    for l.ch != '\'' && l.ch != 0 {
        l.readChar()
    }
    if l.ch != '\'' {
        return Token{Type: CHAR, Value: l.input[start:l.position], Error: "字符常量未闭合"}
    }
    value := l.input[start:l.position]
    l.readChar()
    return Token{Type: CHAR, Value: value}
}

func (l *Lexer) readString() Token {
    l.readChar()
    start := l.position
    for l.ch != '"' && l.ch != 0 {
        l.readChar()
    }
    if l.ch != '"' {
        return Token{Type: STRING, Value: l.input[start:l.position], Error: "字符串常量未闭合"}
    }
    value := l.input[start:l.position]
    l.readChar()
    return Token{Type: STRING, Value: value}
}

func isLetter(ch byte) bool {
    return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
    return unicode.IsDigit(rune(ch))
}

func isHexDigit(ch byte) bool {
    return isDigit(ch) ||
        (ch >= 'a' && ch <= 'f') ||
        (ch >= 'A' && ch <= 'F')
}

func isOctalDigit(ch byte) bool {
    return ch >= '0' && ch <= '7'
}

func isBinaryDigit(ch byte) bool {
    return ch == '0' || ch == '1'
}

func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    }
    return l.input[l.readPosition]
}
