package main

import (
	"bufio"
	"fmt"
	// recDesParser "mygo_c_compiler/rec_des_parser"
    lRParser "mygo_c_compiler/lr_parser"
	"mygo_c_compiler/lexer"
	"os"
	"strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("请提供源代码文件路径")
        return
    }

    filePath := os.Args[1]
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("无法打开文件:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var sourceCode strings.Builder
    lineNum := 1
    inComment := false

    // 读取整个源文件内容，处理注释
    for scanner.Scan() {
        line := scanner.Text()

        if inComment {
            if strings.Contains(line, "*/") {
                line = line[strings.Index(line, "*/")+2:]
                inComment = false
            } else {
                lineNum++
                continue
            }
        }

        if strings.Contains(line, "//") {
            line = strings.Split(line, "//")[0]
        }
        if strings.Contains(line, "/*") {
            inComment = true
            line = line[:strings.Index(line, "/*")]
        }

        if strings.TrimSpace(line) != "" {
            sourceCode.WriteString(line + "\n")
        }
        lineNum++
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("读取文件错误:", err)
        return
    }

    // 词法分析
    fmt.Println("\n词法分析结果:")

    tokens := []lexer.Token{}
    lines := strings.Split(sourceCode.String(), "\n")
    for i, line := range lines {
        fmt.Printf("%d: %s\n", i+1, line)
        l := lexer.NewLexer(line)
        for {
            tok := l.NextToken()
            if tok.Type == lexer.UNKNOWN && tok.Value == "" {
                break
            }
            tokens = append(tokens, tok)

            if tok.Error != "" {
                fmt.Printf("错误: (%s, %s) - %s\n", tok.Type, tok.Value, tok.Error)
            } else {
                fmt.Printf("(%s, %s)\n", tok.Type, tok.Value)
            }
        }
    }

    // 语法分析
    // fmt.Println("\n递归下降语法分析结果:")
    // grammar := recDesParser.New()
    // defer func() {
    //     if r := recover(); r != nil {
    //         fmt.Printf("语法错误: %v\n", r)
    //     }
    // }()
    // grammar.Parse(sourceCode.String())


    fmt.Println("\nLR(0)语法分析结果:")
    lrParser := lRParser.New()
    lrParser.PrintItemSets("items.txt")
    lrParser.PrintParsingTable()
    lrParser.Parse(tokens)

    fmt.Println("语法分析通过")
}
