# FunL
FunL is simple dynamically typed, functional programming language.
It's interpreted language with support for concurrency and immutable data.

## Get started
### Install
#### Install in Linux (or Cygwin, Mac)
    git clone https://github.com/anssihalmeaho/funl.git
    make

#### Install in Windows
    git clone https://github.com/anssihalmeaho/funl.git
    go build -trimpath -o funla.exe -v .

### Run Hello World
#### In Linux  (or Cygwin, Mac)
    ./funla -silent examples/hello.fnl
    Hello World

#### In Windows
    funla.exe -silent examples\hello.fnl
    Hello World

There are also other examples in examples folder.

### Options (-help, -h)
#### In Linux  (or Cygwin, Mac)
    funla.exe -help

#### In Windows
    ./funla -help

### REPL (Read-Eval-Print-Loop)
#### In Linux  (or Cygwin, Mac)
    funla.exe -repl

#### In Windows
    ./funla -repl

In REPL type help for more information.
