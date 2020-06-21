#LEX = lex
#YACC = yacc

CC = gcc

ruledef:	y.tab.o lex.yy.o
	$(CC) -o ruledef y.tab.o lex.yy.o -ly -ll

lex.yy.o:	lex.yy.c y.tab.h

y.tab.c y.tab.h: ruledef.y
	$(YACC) -d ruledef.y

lex.yy.c: ruledef.l
	$(LEX) ruledef.l

.PHONY:	clean

clean:
	rm -f *.o y.tab.h y.tab.c lex.yy.c
