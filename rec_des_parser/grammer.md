program -> block

block -> { stmts }

stmts -> stmt stmts | ε

<!-- stmt -> id = expr ;
      | if ( bool ) stmt
      | if ( bool ) stmt else stmt
      | while ( bool ) stmt
      | do stmt while ( bool )
      | break
      | block -->

stmt → if ( bool ) stmt stmt'
      | id = expr ;
      | while ( bool ) stmt
      | do stmt while ( bool )
      | break
      | block

stmt' → else stmt | ε

<!-- bool -> expr < expr
      | expr <= expr
      | expr > expr
      | expr >= expr
      | expr -->

bool -> expr bool'

bool' -> < expr
       | <= expr
       | > expr
       | >= expr
       | ε

<!-- expr -> expr + term
      | expr - term
      | term -->

expr -> term expr'

expr' -> + term expr'
       | - term expr'
       | ε

<!-- term -> term * factor
      | term / factor
      | factor -->

term -> factor term'

term' -> * factor term'
       | / factor term'
       | ε

factor -> ( expr )
        | id
        | num
