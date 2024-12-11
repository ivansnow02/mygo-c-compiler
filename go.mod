module mygo_c_compiler

go 1.23.2

require mygo_c_compiler/lexer v0.0.0

replace mygo_c_compiler/lexer => ./lexer

require mygo_c_compiler/rec_des_parser v0.0.0

replace mygo_c_compiler/rec_des_parser => ./rec_des_parser

require mygo_c_compiler/lr_parser v0.0.0

replace mygo_c_compiler/lr_parser => ./lr_parser
