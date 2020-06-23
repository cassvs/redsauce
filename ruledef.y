%{
    /*
    Logical rule definition parser
    for redsauce.go cellular automata engine
    Cass Smith, June 2020
    */
#include <stdbool.h>
#include <stdio.h>
#define VARLEN 26

int yylex(void);
void yyerror(char *);

typedef enum expr_type_t {
    EXPR,
    CONST,
    VAR,
    INV
} expr_type_t;

typedef struct expr_t {
    expr_type_t type;
    int var;
    bool val;
    struct expr_t *ls;
    struct expr_t *rs;
    char op;
} expr_t;

bool eval(expr_t *expr);
expr_t *new_expr_t();
expr_t *new_expression(char op, expr_t *ls, expr_t *rs);
expr_t *new_constant(bool val);
expr_t *new_variable(int var);
expr_t *new_inversion(expr_t *rs);

bool variables[VARLEN];
int highestVariable = 0;

expr_t *tree;

%}

%token VARIABLE CONSTANT
%left '|'
%left '^'
%left '&'
%nonassoc '!'

%%
rule:   expression  {tree = (expr_t*)$1;
/*printf(" = %i. HV = %i\n", eval(tree), highestVariable);*/};

expression: expression '&' expression   {$$ = (int)new_expression('&', (expr_t*)$1, (expr_t*)$3);}
    |   expression '|' expression   {$$ = (int)new_expression('|', (expr_t*)$1, (expr_t*)$3);}
    |   expression '^' expression   {$$ = (int)new_expression('^', (expr_t*)$1, (expr_t*)$3);}
    |   '!' expression  {$$ = (int)new_inversion((expr_t*)$2);}
    |   '(' expression ')'  {$$ = $2;}
    |   CONSTANT    {$$ = (int)new_constant((bool)$1);}
    |   VARIABLE    {highestVariable = ($1 + 1 > highestVariable) ? $1 + 1 : highestVariable;
                        $$ = (int)new_variable($1);}
    ;

%%
bool eval(expr_t *expr) {
    bool l;
    bool r;
    switch(expr->type) {
        case EXPR:
            switch(expr->op) {
                case '&':
                    return eval(expr->ls) && eval(expr->rs);
                    break;
                case '|':
                    return eval(expr->ls) || eval(expr->rs);
                    break;
                case '^':
                    l = eval(expr->ls);
                    r = eval(expr->rs);
                    return (l && !r) || (!l && r);
                    break;
                default:
                    break;
            }
            break;
        case CONST:
            return expr->val;
            break;
        case VAR:
            return variables[expr->var];
            break;
        case INV:
            return !eval(expr->rs);
        default:
            break;
    }
}

expr_t *new_expr_t(){
    expr_t *new = malloc(sizeof(expr_t));
    return new;
}

expr_t *new_expression(char op, expr_t *ls, expr_t *rs) {
    expr_t *new = new_expr_t();
    new->type = EXPR;
    new->op = op;
    new->ls = ls;
    new->rs = rs;
    //printf("New expression %p, op %c, ls %p, rs %p\n", new, new->op, new->ls, new->rs);
    return new;
}

expr_t *new_constant(bool val) {
    expr_t *new = new_expr_t();
    new->type = CONST;
    new->val = val;
    //printf("New constant %p, val %u\n", new, new->val);
    return new;
}

expr_t *new_variable(int var) {
    expr_t *new = new_expr_t();
    new->type = VAR;
    new->var = var;
    //printf("New variable %p, slot %c\n", new, (char)new->var + 'a');
    return new;
}

expr_t *new_inversion(expr_t *rs) {
    expr_t *new = new_expr_t();
    new->type = INV;
    new->rs = rs;
    //printf("New inversion %p, rs %p\n", new, new->rs);
    return new;
}

bool increment(bool *array, int limit) {
    for (int i = 0; i < limit; i++) {
        if (array[i]) {
            array[i] = false;
            continue;
        } else {
            array[i] = true;
            break;
        }
    }
    for (int j = 0; j < limit; j++) {
        if (array[j]) return true;
    }
    return false;
}

int main() {
    yyparse();
    for (int i = 0; i < VARLEN; i++) {
        variables[i] = false;
    }
    do {
        printf("%u\n", eval(tree));
    } while (increment(variables, highestVariable));
}
