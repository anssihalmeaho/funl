
# Possible contributions
There are many ways to contribute, for example:

* Syntax highlighting for IDE's (VSC, Atom etc.)
* Extending standard libraries (for example, networking/sockets, HTTP extending etc.)
* Make projects based on FunL, applications or libraries
* Improve (or make better) testing framework
* Make automated tests for standard libraries
* More tooling (code formatter, code analyzers etc.)

Possible new (standard) libraries:

* support for compression (gzip etc.)
* cryptography support
* data formats: csv, excel, base64, yaml etc.
* sockets (TCP/UDP), websockets, Server-Sent Events, SMTP etc.
* regular expressions
* more advanced command line options parser
* random number generation
* bit operations

# Asking for more information
I'm happy to answer questions if needed.
You can use my email for that: anssi.halmeaho@hotmail.com

# Technical details of developing FunL
Here are some details of developing FunL.

## Developing FunL native (in FunL itself) standard libraries
Some of FunL standard libraries are implemented in FunL itself (like **stdfu**).
As those are embedded into single executable (_funla_) source code of native libraries 
are included as Go raw strings (FunL executable contains all standard libraries).
Source files of those libraries are placed in stdfun directory and following
command will update **stdfunfiles.go** file:

```
go generate
```

This needs to be done when adding or modifying some standard libraries written in FunL itself.

In _main.go_ there's line ```//go:generate go run ./stdfun/stdfun_generator.go``` which causes 
execution of ```/stdfun/stdfun_generator.go``` in ```go generate``` which produces _stdfunfiles.go_.

In addition name of new module needs to be added to __funmodNames__ table in __initFunSourceSTD__ function(__funl.go__)
so that module is included in __funla__ binary.

## Testing
There is testing tool and tests implemented in FunL. Run those tests when doing changes:

### Testing with tester.fnl
Most language functionality is tested with **tester.fnl**.

Short status printed:

```
./funla -args="'tst'" tester.fnl
```

Longer status printed:

```
./funla -args="'tst'" -name=all tester.fnl
```

Note. this tool can be used as unit testing tool also otherwise.

### Unit testing Go code
For some parts of code there are also Go unit tests:

```
go test ./...
```

## Debugging runtime errors
See [guide to read runtime error printouts](https://github.com/anssihalmeaho/funl/wiki/rteprint)

