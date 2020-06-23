%{
    // Logical rule definition lexer
    // for redsauce.go cellular automata engine
    // Cass Smith, June 2020
    #include "y.tab.h"
    extern int yylval;
%}

variable [a-z]
constant [01]
whitespace [ \t]

%%
{variable}  {yylval = yytext[0] - 'a'; return VARIABLE;}
{constant}  {yylval = atoi(yytext); return CONSTANT;}
{whitespace}    ;
\n  return 0;
.   return yytext[0];

%%

/* main() {
    yylex();
} */
