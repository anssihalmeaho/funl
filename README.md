# FunL
FunL is simple dynamically typed, functional programming language.
It's interpreted language with support for concurrency and immutable data.

## Get started
### Install
#### Install in Linux (or Cygwin, Mac)
    git clone https://github.com/anssihalmeaho/funl.git
    cd funl
    make

Run hello world example:

    ./funla -silent examples/hello.fnl
    Hello World

#### Install in Windows
    git clone https://github.com/anssihalmeaho/funl.git
    cd funl
    go build -trimpath -o funla.exe -v .

Run hello world example:

    funla.exe -silent examples\hello.fnl
    Hello World

There are also other examples in examples folder.

## Getting help and try expressions

### Options: -help, -h
#### In Linux  (or Cygwin, Mac)
    ./funla -help

#### In Windows
    funla.exe -help

### REPL (Read-Eval-Print-Loop)
#### In Linux  (or Cygwin, Mac)
    ./funla -repl

#### In Windows
    funla.exe -repl

In REPL type help for more information.

### Options: -eval
    ./funla -eval "plus(1 2)"
    3

### help operator

help operator can be used to get list of operators:

    ./funla -eval "help('operators')"

help operator provides description for each operator:

    ./funla -eval "help('if')"

in REPL:

    ./funla -repl
    Welcome to FunL REPL (interactive command shell)
    funl> help('if')
