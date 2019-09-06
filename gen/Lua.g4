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

exp5: exp4 ('..' exp5)*;

exp4: exp3 (('+' | '-') exp3)*;

exp3: exp2 (( '*' | '/' | '//' | '%') exp2)*;

exp2: ('not' | '#' | '-' | '~' )+ exp2;

exp1: exp0 ('^' exp1)*;

exp0: 'nil' | 'false' | 'true' | STRING | NUMBER | '(' exp12 ')' | ID '(' (exp12',')* ')';

prefix1 : prefix0 ('.' STRING)* | prefix0 ('[' prefix1 ']')*;
prefix0 : STRING | '(' exp ')';

IDlist: '(' (ID ',')*ID ')';

STRING: '"' .*? '"';

ID: STRING;

NUMBER: '-'?[0-9]+('.' [0-9]+)?;