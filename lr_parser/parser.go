package lr_parser

import (
    "fmt"
    "regexp"
    "strings"
    "mygo_c_compiler/lexer"
    "os"
    "bufio"
)

// 产生式结构
type Production struct {
    Left  string   // 左部
    Right []string // 右部
}

// LR(1)项目
type Item struct {
    Prod  Production // 产生式
    Dot   int        // 点的位置
    Lookahead string // 展望符
}

// 项目集
type ItemSet struct {
    Items []Item
}

// LR分析表
type ActionTable map[int]map[string]string
type GotoTable map[int]map[string]int

type Parser struct {
    Productions []Production
    ItemSets    []ItemSet
    Action      ActionTable
    Goto        GotoTable
}

// 创建新的解析器
func New() *Parser {
    parser := &Parser{}
    // load grammar.md
    file, err := os.Open("./lr_parser/grammar.md")
    if (err != nil) {
        fmt.Println("无法打开文件:", err)
        return nil
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var grammar strings.Builder

    for scanner.Scan() {
        line := scanner.Text()
        grammar.WriteString(line + "\n")
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("读取文件错误:", err)
        return nil
    }

    if err := parser.ParseGrammar(grammar.String()); err != nil {
        fmt.Println("解析文法错误:", err)
        return nil
    }
    parser.GenerateCanonicalCollection()
    parser.BuildParsingTable()

    return parser
}

// 解析产生式
func (p *Parser) ParseGrammar(grammar string) error {
    lines := strings.Split(grammar, "\n")
    regex := regexp.MustCompile(`(\w+)\s*->\s*(.+)`)

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        matches := regex.FindStringSubmatch(line)
        if len(matches) != 3 {
            continue
        }

        left := matches[1]
        right := strings.Fields(matches[2])

        p.Productions = append(p.Productions, Production{
            Left:  left,
            Right: right,
        })
    }

    return nil
}

// 生成项目集规范簇
func (p *Parser) GenerateCanonicalCollection() {
    // 初始项目集
    initialItem := Item{
        Prod: p.Productions[0],
        Dot:  0,
        Lookahead: "$",
    }

    initialSet := ItemSet{
        Items: []Item{initialItem},
    }

    p.ItemSets = append(p.ItemSets, p.closure(initialSet))

    // 继续生成其他项目集
    for i := 0; i < len(p.ItemSets); i++ {
        set := p.ItemSets[i]
        // 获取所有可能的下一个符号
        symbols := p.getNextSymbols(set)

        for _, symbol := range symbols {
            newSet := p.goto_(set, symbol)
            if len(newSet.Items) == 0 {
                continue
            }

            // 检查是否已经存在相同的项目集
            existingIndex := p.findItemSetIndex(newSet)
            if (existingIndex == -1) {
                p.ItemSets = append(p.ItemSets, newSet)
            }
        }
    }
}

// 构建LR(1)分析表
func (p *Parser) BuildParsingTable() {
    p.Action = make(ActionTable)
    p.Goto = make(GotoTable)

    for i := range p.ItemSets {
        p.Action[i] = make(map[string]string)
        p.Goto[i] = make(map[string]int)

        set := p.ItemSets[i]
        for _, item := range set.Items {
            if item.Dot < len(item.Prod.Right) {
                // 移进动作
                symbol := item.Prod.Right[item.Dot]
                nextSet := p.goto_(set, symbol)
                nextIndex := p.findItemSetIndex(nextSet)
                if nextIndex >= 0 {
                    if p.isTerminal(symbol) {
                        p.Action[i][symbol] = fmt.Sprintf("s%d", nextIndex)
                    } else {
                        p.Goto[i][symbol] = nextIndex
                    }
                }
            } else {
                // 规约动作
				prodIndex := p.findProductionIndex(item.Prod)
				if prodIndex == 0 && item.Lookahead == "$" {
					p.Action[i]["$"] = "accept"
				} else {
					if _, exists := p.Action[i][item.Lookahead]; !exists {
						p.Action[i][item.Lookahead] = fmt.Sprintf("r%d", prodIndex)
					}
				}
            }
        }
    }
}

// 获取所有终结符
func (p *Parser) getTerminals() []string {
    terminals := make(map[string]bool)

    // 遍历所有产生式的右部
    for _, prod := range p.Productions {
        for _, symbol := range prod.Right {
            if p.isTerminal(symbol) {
                terminals[symbol] = true
            }
        }
    }

    // 转换为切片
    result := make([]string, 0)
    for terminal := range terminals {
        result = append(result, terminal)
    }
    return result
}

// 执行语法分析
func (p *Parser) Parse(tokens []lexer.Token) bool {
    stack := []int{0}      // 状态栈
    symbols := []string{}   // 符号栈
    actions := []string{}   // 动作序列

    // 将词法单元转换为语法符号的辅助函数
    tokenToSymbol := func(tok lexer.Token) string {
        switch tok.Type {
        case lexer.MAIN:
            return "main"
        case lexer.IDENT:
            return "id"
        case lexer.NUMBER:
            return "num"
        case lexer.LBRACE:
            return "{"
        case lexer.RBRACE:
            return "}"
        case lexer.LPAREN:
            return "("
        case lexer.RPAREN:
            return ")"
        case lexer.ASSIGN:
            return "="
        case lexer.SEMICOLON:
            return ";"
        case lexer.WHILE:
            return "while"
        case lexer.LTE:
            return "<="
        case lexer.GTE:
            return ">="
        case lexer.PLUS:
            return "+"
        case lexer.ASTERISK:
            return "*"
        default:
            return string(tok.Type)
        }
    }

    i := 0
    for {
        state := stack[len(stack)-1]
        var symbol string
        if i < len(tokens) {
            symbol = tokenToSymbol(tokens[i])
        } else {
            symbol = "$"
        }

        action, exists := p.Action[state][symbol]
        if !exists {
            fmt.Printf("\n语法错误: 状态%d下无法处理符号%s\n", state, symbol)
            fmt.Printf("当前令牌: %v\n", tokens[i])
            fmt.Printf("当前分析栈: %v\n符号栈: %v\n输入: %s\n", stack, symbols, symbol)
            return false
        }

        if action[0] == 's' { // 移进
            nextState := 0
            if _, err := fmt.Sscanf(action, "s%d", &nextState); err != nil {
                fmt.Printf("\n语法错误: 无效的移进动作 %s\n", action)
                return false
            }
            stack = append(stack, nextState)
            symbols = append(symbols, symbol)

            // 记录移进动作
            actionStr := fmt.Sprintf("移进: 输入 %s, 移进到状态 %d", symbol, nextState)
            actions = append(actions, actionStr)
            fmt.Printf("%s\n当前状态栈: %v\n符号栈: %v\n\n", actionStr, stack, symbols)

            i++
        } else if action[0] == 'r' { // 规约
            prodIndex := 0
            if _, err := fmt.Sscanf(action, "r%d", &prodIndex); err != nil {
                fmt.Printf("\n语法错误: 无效的规约动作 %s\n", action)
                return false
            }
            prod := p.Productions[prodIndex]

            // 初始化 nextState 变量
            nextState := 0

            // 检查是否是空产生式的归约
            if len(prod.Right) == 1 && prod.Right[0] == "ε" {
                // 直接在符号栈中压入左部符号
                symbols = append(symbols, prod.Left)
            } else {
                // 弹出右部长度个状态和符号
                stack = stack[:len(stack)-len(prod.Right)]
                symbols = symbols[:len(symbols)-len(prod.Right)]

                // 压入左部符号
                symbols = append(symbols, prod.Left)

                // 查找GOTO表确定下一个状态
                state = stack[len(stack)-1]
                nextState = p.Goto[state][prod.Left]
                stack = append(stack, nextState)
            }

            // 记录规约动作
            actionStr := fmt.Sprintf("规约: 使用产生式 %s -> %s", prod.Left, strings.Join(prod.Right, " "))
            actions = append(actions, actionStr)
            fmt.Printf("%s\nGOTO(I%d, %s) = %d\n当前状态栈: %v\n符号栈: %v\n\n",
                actionStr, state, prod.Left, nextState, stack, symbols)

        } else if action == "accept" {
            fmt.Printf("分析成功: 输入串已被接受!\n\n分析过程:\n")
            for i, act := range actions {
                fmt.Printf("%d. %s\n", i+1, act)
            }
            return true
        } else {
            fmt.Printf("\n语法错误: 无效的动作 %s\n", action)
            return false
        }
    }
}

// 判断是否为终结符
func (p *Parser) isTerminal(symbol string) bool {
    // 检查symbol是否出现在任何产生式的左部
    for _, prod := range p.Productions {
        if prod.Left == symbol {
            return false
        }
	}
    return true
}

// 判断是否为非终结符
func (p *Parser) isNonTerminal(symbol string) bool {
    return !p.isTerminal(symbol)
}

func (p *Parser) closure(set ItemSet) ItemSet {
    result := ItemSet{
        Items: make([]Item, len(set.Items)),
    }
    copy(result.Items, set.Items)

    changed := true
    for changed {
        changed = false
        size := len(result.Items)

        // 遍历当前项目集中的所有项目
        for i := 0; i < size; i++ {
            item := result.Items[i]

            // 如果点后面没有符号，继续下一个项目
            if item.Dot >= len(item.Prod.Right) {
                continue
            }

            // 获取点后面的符号
            nextSymbol := item.Prod.Right[item.Dot]

            // 如果点后面的符号是非终结符
            if p.isNonTerminal(nextSymbol) {
                // 计算FIRST(βa)
                lookaheads := p.computeFirst(item.Prod.Right[item.Dot+1:], item.Lookahead)


                // 查找以该非终结符为左部的所有产生式
                for _, prod := range p.Productions {
                    if prod.Left == nextSymbol {
                        for _, lookahead := range lookaheads {
                            // 创建新项目
                            newItem := Item{
                                Prod: prod,
                                Dot:  0,
                                Lookahead: lookahead,
                            }

                            // 检查新项目是否已存在
                            exists := false
                            for _, existingItem := range result.Items {
                                if p.itemsEqual(existingItem, newItem) {
                                    exists = true
                                    break
                                }
                            }

                            // 如果不存在则添加
                            if !exists {
                                result.Items = append(result.Items, newItem)
                                changed = true
                            }
                        }
                    }
                }

                // // 如果E可推导出空字符串，生成E->#
                // if p.canDeriveEmpty(nextSymbol) {
                //     newItem := Item{
                //         Prod: Production{Left: nextSymbol, Right: []string{"#"}},
                //         Dot:  1,
                //         Lookahead: item.Lookahead,
                //     }
                //     result.Items = append(result.Items, newItem)
                //     changed = true
                // }
            }
        }
    }

    return result
}

// // 判断非终结符是否可推导出空字符串
// func (p *Parser) canDeriveEmpty(symbol string) bool {
//     for _, prod := range p.Productions {
//         if prod.Left == symbol && len(prod.Right) == 1 && prod.Right[0] == "ε" {
//             return true
//         }
//     }
//     return false
// }
// 计算FIRST集合
func (p *Parser) computeFirst(symbols []string, lookahead string) []string {
    if len(symbols) == 0 {
        return []string{lookahead}
    }

    firstSet := make(map[string]bool)
    for _, symbol := range symbols {
        if p.isTerminal(symbol) {
            firstSet[symbol] = true
            break
        }

        nullable := false
        for _, prod := range p.Productions {
            if prod.Left == symbol {
                if len(prod.Right) == 1 && prod.Right[0] == "ε" {
                    nullable = true
                } else {
                    for _, sym := range p.computeFirst(prod.Right, lookahead) {
                        firstSet[sym] = true
                    }
                }
            }
        }

        if !nullable {
            break
        }
    }

    // 将 lookahead 添加到 firstSet 中
    firstSet[lookahead] = true

    result := make([]string, 0, len(firstSet))
    for sym := range firstSet {
        result = append(result, sym)
    }
    return result
}
// 计算GOTO函数
func (p *Parser) goto_(set ItemSet, symbol string) ItemSet {
    resultSet := ItemSet{
        Items: make([]Item, 0),
    }

    // 遍历项目集中的所有项目
    for _, item := range set.Items {
        // 如果点已经到达末尾，继续下一个项目
        if item.Dot >= len(item.Prod.Right) {
            continue
        }

        // 如果点后面的符号匹配
        if item.Prod.Right[item.Dot] == symbol {
            // 创建新项目，将点向后移动一位
            newItem := Item{
                Prod: item.Prod,
                Dot:  item.Dot + 1,
                Lookahead: item.Lookahead,
            }
            resultSet.Items = append(resultSet.Items, newItem)
        }
    }

    // 如果结果集非空，计算其闭包
    if len(resultSet.Items) > 0 {
        return p.closure(resultSet)
    }

    return resultSet
}

// 获取项目集中所有可能的下一个符号
func (p *Parser) getNextSymbols(set ItemSet) []string {
    symbols := make(map[string]bool)
    for _, item := range set.Items {
        if item.Dot < len(item.Prod.Right) {
            symbols[item.Prod.Right[item.Dot]] = true
        }
    }

    result := make([]string, 0)
    for symbol := range symbols {
        result = append(result, symbol)
    }
    return result
}

// 查找项目集的索引
func (p *Parser) findItemSetIndex(set ItemSet) int {
    for i, existingSet := range p.ItemSets {
        if p.itemSetsEqual(existingSet, set) {
            return i
        }
    }
    return -1
}

// 比较两个项目集是否相等
func (p *Parser) itemSetsEqual(set1, set2 ItemSet) bool {
    if len(set1.Items) != len(set2.Items) {
        return false
    }

    for _, item1 := range set1.Items {
        found := false
        for _, item2 := range set2.Items {
            if p.itemsEqual(item1, item2) {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}

// 比较两个项目是否相等
func (p *Parser) itemsEqual(item1, item2 Item) bool {
    return item1.Dot == item2.Dot &&
        item1.Prod.Left == item2.Prod.Left &&
        strings.Join(item1.Prod.Right, " ") == strings.Join(item2.Prod.Right, " ") &&
        item1.Lookahead == item2.Lookahead
}

// 查找产生式的索引
func (p *Parser) findProductionIndex(prod Production) int {
    for i, existingProd := range p.Productions {
        if existingProd.Left == prod.Left &&
            strings.Join(existingProd.Right, " ") == strings.Join(prod.Right, " ") {
            return i
        }
    }
    return -1
}

// 获取FOLLOW集
func (p *Parser) getFollowSet(symbol string) []string {
    follows := make(map[string]bool)

    // 初始化FOLLOW集
    if symbol == p.Productions[0].Left {
        follows["$"] = true
    }

    changed := true
    for changed {
        changed = false

        // 遍历所有产生式
        for _, prod := range p.Productions {
            for i, sym := range prod.Right {
                if sym == symbol {
                    // 如果是最后一个符号
                    if i == len(prod.Right)-1 {
                        // 添加产生式左部的FOLLOW集
                        if prod.Left != symbol { // 避免左递归
                            for _, follow := range p.getFollowSet(prod.Left) {
                                if !follows[follow] {
                                    follows[follow] = true
                                    changed = true
                                }
                            }
                        }
                    } else {
                        // 添加后面符号的FIRST集
                        nextSymbols := prod.Right[i+1:]
                        first := p.computeFirst(nextSymbols, "$")
                        for _, f := range first {
                            if f != "ε" && !follows[f] {
                                follows[f] = true
                                changed = true
                            }
                        }

                        // 如果FIRST集包含ε，添加产生式左部的FOLLOW集
                        if containsEpsilon(first) {
                            for _, follow := range p.getFollowSet(prod.Left) {
                                if !follows[follow] {
                                    follows[follow] = true
                                    changed = true
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    result := make([]string, 0)
    for follow := range follows {
        result = append(result, follow)
    }
    return result
}

// 检查FIRST集是否包含ε
func containsEpsilon(first []string) bool {
    for _, f := range first {
        if f == "ε" {
            return true
        }
    }
    return false
}

// 打印项目集规范族（用于调试）
func (p *Parser) PrintItemSets(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    for i, set := range p.ItemSets {
        fmt.Fprintf(file, "I%d:\n", i)
        for _, item := range set.Items {
            // 构建带点的右部字符串
            rightPart := make([]string, len(item.Prod.Right)+1)
            copy(rightPart, item.Prod.Right[:item.Dot])
            rightPart[item.Dot] = "·"
            copy(rightPart[item.Dot+1:], item.Prod.Right[item.Dot:])

            fmt.Fprintf(file, "    %s -> %s, %s\n", item.Prod.Left, strings.Join(rightPart, " "), item.Lookahead)
        }
        fmt.Fprintln(file)
    }

    return nil
}
// 打印分析表（用于调试）
func (p *Parser) PrintParsingTable() {
    fmt.Println("ACTION TABLE:")
    for state := range p.Action {
        fmt.Printf("State %d:\n", state)
        for symbol, action := range p.Action[state] {
            fmt.Printf("    %s -> %s\n", symbol, action)
        }
    }

    fmt.Println("\nGOTO TABLE:")
    for state := range p.Goto {
        fmt.Printf("State %d:\n", state)
        for symbol, nextState := range p.Goto[state] {
            fmt.Printf("    %s -> %d\n", symbol, nextState)
        }
    }
}
