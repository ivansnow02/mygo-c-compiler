program_prime -> program
program -> main block
block -> { stmts }
stmts -> stmt stmts
stmts -> ε
stmt -> id = E ;
stmt -> while ( bool )  stmt
stmt -> block
E -> E + F
E -> F
F -> F * G
F -> G
G -> ( E )
G -> T
bool -> T <= T
bool -> T >= T
bool -> T
T -> id
T -> num
