/*
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
exp9  ::= exp8 {‘|’ exp8}
exp8  ::= exp7 {‘~’ exp7}
exp7  ::= exp6 {‘&’ exp6}
exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
exp5  ::= exp4 {‘..’ exp5}
exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp2
exp1  ::= exp0 {‘^’ exp2}
exp0  ::= nil | false | true | Numeral | LiteralString
		| ‘...’ | functiondef | prefixexp | tableconstructor
*/

grammar Lua;


lua: exp;

exp: exp12;

exp12: exp11 ('or' exp11)*;

exp11: exp10 ('and' exp10)*;

exp10: exp9 (('<' | '>' | '<=' | '>=' | '~=' | '==') exp9)*;

exp9: exp8 ('|' exp8);

exp8: exp7 ('~' exp7)*;

exp7: exp6 ('&' exp6)*;

exp6: exp5 (('<<' | '>>') exp5)*;

exp5: exp4 | exp4 '..' exp5;

exp4: exp3 (('+' | '-') exp3)*;

exp3: exp2 (( '*' | '/' | '//' | '%') exp2)*;

exp2: ('not' | '#' | '-' | '~' ) exp2 | exp1;

exp1: exp0 ('^' exp1)*;

exp0: prefix2;

prefix2 : prefix1 ( '('   ')' )* |  prefix1 ( '(' (ID ',')*ID ')' )*;
prefix1: prefix0 ( '.' ID)* | prefix0 ('[' exp ']')*;
prefix0: 'nil' | 'false' | 'true' | STRING | NUMBER | '...' | ID | '(' exp ')';


STRING: '"'.*?'"';

ID: ([a-zA-Z] | '_') ([a-zA-Z0-9] | '_' )*;

NUMBER: [0-9]+('.' [0-9]+)?('e'[0-9]+)?;

WS: [\t\n\r ]+ -> skip;