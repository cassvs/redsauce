# redsauce 🍅

A linear cellular automata exploration tool

```
$ ./redsauce -alive 'o' -dead ' ' -random 80 -wolfram 150 -gen 20 -wrap
o  oo o o           o oo    o o      o       o   o  ooo  o  ooooo  o     o ooo
ooo   o oo         oo   o  oo oo    ooo     ooo oooo o ooooo ooo oooo   oo  o oo
oo o oo   o       o  o oooo     o  o o o   o o   oo  o  ooo   o   oo o o  ooo  o
o  o   o ooo     ooooo  oo o   ooooo o oo oo oo o  ooooo o o ooo o   o ooo o oo
ooooo oo  o o   o ooo oo   oo o ooo  o          ooo ooo  o o  o  oo oo  o  o
 ooo    ooo oo oo  o    o o   o  o oooo        o o   o ooo oooooo     ooooooo  o
  o o  o o       oooo  oo oo ooooo  oo o      oo oo oo  o   oooo o   o ooooo ooo
ooo oooo oo     o oo oo       ooo oo   oo    o        oooo o oo  oo oo  ooo   o
 o   oo    o   oo      o     o o    o o  o  ooo      o oo  o   oo     oo o o oo
ooo o  o  ooo o  o    ooo   oo oo  oo oooooo o o    oo   oooo o  o   o   o o   o
oo  oooooo o  ooooo  o o o o     oo    oooo  o oo  o  o o oo  ooooo ooo oo oo o
  oo oooo  ooo ooo ooo o o oo   o  o  o oo ooo   oooooo o   oo ooo   o        o
 o    oo oo o   o   o  o o   o oooooooo     o o o oooo  oo o    o o ooo      ooo
 oo  o      oo ooo ooooo oo oo  oooooo o   oo o o  oo oo   oo  oo o  o o    o o
o  oooo    o    o   ooo       oo oooo  oo o   o ooo     o o  oo   oooo oo  oo oo
 oo oo o  ooo  ooo o o o     o    oo oo   oo oo  o o   oo ooo  o o oo    oo    o
       ooo o oo o  o o oo   ooo  o     o o     ooo oo o    o ooo o   o  o  o  oo
o     o o  o    oooo o   o o o oooo   oo oo   o o     oo  oo  o  oo oooooooooo
oo   oo ooooo  o oo  oo oo o o  oo o o     o oo oo   o  oo  ooooo    oooooooo oo
o o o    ooo ooo   oo      o ooo   o oo   oo      o oooo  oo ooo o  o oooooo   o
  o oo  o o   o o o  o    oo  o o oo   o o  o    oo  oo oo    o  oooo  oooo o o
```

## Dependencies

To build the `redsuace` binary, all you need is a recent version of the Go environment. I used gccgo, but the standard Go package also works. Check https://golang.org/doc/install for installation instructions for your operating system.

To build the logical ruledef parser, you will need:
- `make`: preferably GNU Make
- A C compiler: GCC or similar
- `lex` or compatible (eg. `flex` with -l or -X option)
- `yacc` or compatible (eg. `bison` with -y option)

Most UNIX/Linux distributions include these tools by default. Check your distribution's documentation or package system for installation instructions.

## Build Instructions

Clone and enter the repository:

    git clone https://github.com/cassvs/redsauce
    cd redsauce

To build `ruledef`, the logical rule definition parser:

    make
    make clean

To build `redsauce`, the main program binary:

    go build

## Options

See https://golang.org/pkg/flag/#hdr-Command_line_syntax for the syntax of command-line options.

### `--help`
Usage information.

### `--state=<state-string>`
The initial state of the 'world'. For example, the state-string `"0001000"` represents a world 7 cells wide, with one live cell in the middle. By default, `1` represents a living cell and `0` represents a dead cell.

### `--gen=<generations>`
The number of generations or iterations to develop the world. Defaults to 10.

### `--alive=<char>`, `--dead=<char>`
Characters used to represent living and dead cells. Defaults are `1` for living and `0` for dead. Be aware that the characters specified here must match the ones used in the `--state` string. Also, only single-width (ASCII) characters will work. For instance, a multi-byte Unicode sequence like `🍅` will cause an error. (You can always pipe the output through `sed` or something similar to replace the default representations with arbitrary strings.)

### `--wrap`
When `--wrap` is set, the ends of the world are connected. This means that cell patterns can cross over the edges as though the world were a tube.

### `--end=<end-state>`
The default state of cells 'outside' the world, used for computing the states of cells near the world boundaries. Unlike `--state`, this option is not affected by `--alive` or `--dead`. Use `true` or `1` to set outside cells to alive, or `false` or `0` to make them to dead (the default). This option is ignored if `--wrap` is set.

### `--random=<width>`
Generate a random initial world of the given width. If `--random` is set, `--state` is ignored.

### `--quiet`
Suppress all output except the final state of the world.

### `--wolfram=<wolfram-code>`
Set the rule used for computing the new cell states of each generation, using the Wolfram rule code. See https://en.wikipedia.org/wiki/Wolfram_code for information on how this system of rule specification works.

### `--rule=<rule-string>`
Use a custom logical rule described by an expression string. Each cell's state is a logical expression in terms of the cells in the previous generation. Cells are represented in the expression by single lowercase letters, and logical AND, OR, XOR, and NOT operations are notated with `&`, `|`, `^`, and `!` respectively. The variables are mapped to cells as follows:

    ... g e c b a d f ...
              _ <- Cell being calculated

For example, the rule `'!b'` produces a pattern in which each generation is the inverse of the previous one, and `'a^c'` is equivalent to the Wolfram code 90.

There are a few drawbacks to using this feature. First, interesting rules tend to have complex logical descriptions which are difficult to simplify. For example, rule 110 can be represented by the expression `'(!a & b) | (!b & c) | (a & b & !c)'`. For rules with neighborhoods wider than three cells, these expressions can be very tedious to derive and type in. Second, the rule parser program is separate from the main program, and uses a truth table exchange format to transfer the rule to the main program. This truth table grows *exponentially* with the furthest-away cell used in the expression. For example, in the rule `'a^c'`, the highest variable used is `c`, the third cell, so the truth table contains 2^3 (or 8) lines. However, the rule `'z'`, which seems fairly simple (and boring), uses `z`, the 26th cell, and generates a truth table 2^26 (or 67 million) lines long, which takes a very long time to compute and a lot of memory to store. So try to keep the neighborhoods of your rules small.

## Examples

### State strings, alive and dead cells, generations

```
$ ./redsauce --alive='^' --dead=' ' --state='       ^       ' -rule='a^c' --gen=7
       ^
      ^ ^
     ^   ^
    ^ ^ ^ ^
   ^       ^
  ^ ^     ^ ^
 ^   ^   ^   ^
^ ^ ^ ^ ^ ^ ^ ^
```

### Random state, Wolfram rule

```
$ ./redsauce --alive='O' --dead='.' --random=30 --wolfram=184
.O.O.OOOOO..OOOOOO.O.O...OOOOO
..O.OOOOO.O.OOOOO.O.O.O..OOOO.
...OOOOO.O.OOOOO.O.O.O.O.OOO.O
...OOOO.O.OOOOO.O.O.O.O.OOO.O.
...OOO.O.OOOOO.O.O.O.O.OOO.O.O
...OO.O.OOOOO.O.O.O.O.OOO.O.O.
...O.O.OOOOO.O.O.O.O.OOO.O.O.O
....O.OOOOO.O.O.O.O.OOO.O.O.O.
.....OOOOO.O.O.O.O.OOO.O.O.O.O
.....OOOO.O.O.O.O.OOO.O.O.O.O.
.....OOO.O.O.O.O.OOO.O.O.O.O.O
```

### Wrap

```
$ ./redsauce --wrap --state='00000100000' --rule='a'
00000100000
00001000000
00010000000
00100000000
01000000000
10000000000
00000000001
00000000010
00000000100
00000001000
00000010000
```

### End-state

```
$ ./redsauce --end=true --state='0000000000000000000000000' --wolfram=30
0000000000000000000000000
1000000000000000000000001
0100000000000000000000011
0110000000000000000000110
0101000000000000000001100
0101100000000000000011011
0101010000000000000110010
0101011000000000001101110
0101010100000000011001000
0101010110000000110111101
0101010101000001100100001
```

### Quiet

```
$ ./redsauce --wolfram=90 --random=80 --quiet --dead=' '
1    1   1 11111111   1 11 111 1111   11 11  1 1 1111     1 1  1111 1     1111 1
```

### Pretty full-terminal patterns

```
$ ./redsauce --wolfram=15 --random=80 --dead=' ' --alive='O' --wrap --gen 20
 O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO     OOOO   OO
OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOOO    OOO
O  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO     OOOO   O
  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOOO    OOO
 OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO     OOOO
OO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOOO    O
   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO     OOOO
 OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOOO
OO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO     OO
   OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOOO
OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO
O    OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OOOO
  OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO
OOO    OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO  OO
    OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  OO
OOOOO    OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO
O     OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO  O
  OOOOO    OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O   OO
 OO     OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OOO
OO  OOOOO    OOO  OO O OO   OOOOO OO   OO O O O O O O OO OO OOOOO  OO O O O O
O  OO     OOOO   OO  O O  OOO     O  OOO  O O O O O O O  O  O     OO  O O O O OO
```

### Custom rules

```
$ ./redsauce --random=80 --dead=' ' --alive='O' --wrap --gen 20 --rule='(b&!c)|(!a&c&e)|(a&!b&d)|(!a&b&c)|(a&!b&c&!d&e)'
OO     O    O OOOOOO O OO     O   OOO   O   O O  O  O  OOOOO O OO OO  OO OO   OO
 OO    O    OOO    OOOOOOO    O  OO OO  O   O O  O  O OO   OOOOOOOOOOOOOOOOO OO
OOOO   O   OO OO  OO     OO   O OOOOOOO O   O O  O  OOOOO OO               OOOOO
   OO  O  OOOOOOOOOOO   OOOO  OOO     OOO   O O  O OO   OOOOO             OO
  OOOO O OO         OO OO  OOOO OO   OO OO  O O  OOOOO OO   OO           OOOO
 OO  OOOOOOO       OOOOOOOOO  OOOOO OOOOOOO O O OO   OOOOO OOOO         OO  OO
OOOOOO     OO     OO       OOOO   OOO     OOO OOOOO OO   OOO  OO       OOOOOOOO
O    OO   OOOO   OOOO     OO  OO OO OO   OO OOO   OOOOO OO OOOOOO     OO      OO
OO  OOOO OO  OO OO  OO   OOOOOOOOOOOOOO OOOOO OO OO   OOOOOO    OO   OOOO    OO
OOOOO  OOOOOOOOOOOOOOOO OO            OOO   OOOOOOOO OO    OO  OOOO OO  OO  OOOO
    OOOO              OOOOO          OO OO OO      OOOOO  OOOOOO  OOOOOOOOOOO
   OO  OO            OO   OO        OOOOOOOOOO    OO   OOOO    OOOO         OO
  OOOOOOOO          OOOO OOOO      OO        OO  OOOO OO  OO  OO  OO       OOOO
 OO      OO        OO  OOO  OO    OOOO      OOOOOO  OOOOOOOOOOOOOOOOO     OO  OO
OOOO    OOOO      OOOOOO OOOOOO  OO  OO    OO    OOOO               OO   OOOOOOO
   OO  OO  OO    OO    OOO    OOOOOOOOOO  OOOO  OO  OO             OOOO OO
  OOOOOOOOOOOO  OOOO  OO OO  OO        OOOO  OOOOOOOOOO           OO  OOOOO
 OO          OOOO  OOOOOOOOOOOOO      OO  OOOO        OO         OOOOOO   OO
OOOO        OO  OOOO           OO    OOOOOO  OO      OOOO       OO    OO OOOO
O  OO      OOOOOO  OO         OOOO  OO    OOOOOO    OO  OO     OOOO  OOOOO  OO O
OOOOOO    OO    OOOOOO       OO  OOOOOO  OO    OO  OOOOOOOO   OO  OOOO   OOOOOOO
```
