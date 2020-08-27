package funl

// OperatorInfo contains information about one operator
type OperatorInfo struct{}

// NewOperatorDocs returns documentation for operators
func NewOperatorDocs() map[string]string {
	return map[string]string{
		"and": `
Operator: and
  Performs logical and -operation for arguments.
  All arguments are assumed to be boolean expressions.
  Number of arguments need to be at least 1.
  Return value is boolean type:
    - true: if all arguments are true
    - false: otherwise

Usage: and(<boolean-expr> <boolean-expr> <boolean-expr> ...)  

Note. and -operator stops evaluating any more input arguments
      after any input expression evaluates to false.
`,
		"or": `
Operator: or
  Performs logical or -operation for arguments.
  All arguments are assumed to be boolean expressions.
  Number of arguments need to be at least 1.
  Return value is boolean type:
    - true: if any of arguments is true
    - false: otherwise

Usage: or(<boolean-expr> <boolean-expr> <boolean-expr> ...)  

Note. or -operator stops evaluating any more input arguments
      after any input expression evaluates to true.
`,
		"call": `
Operator: call
  Calls function or procedure given as 1st argument.
  Arguments of function/procedure call are following arguments.
  Number of arguments needs to be at least 1 (func/proc to be called).
  Arguments are evaluated before calling the function/procedure.
  Function/procedure can be also external procedure or function.
  Return value is return value of function/procedure. 

Usage: call(<func/proc> <arg-1> <arg-2> ...)  
`,
		"not": `
Operator: not
  Performs logical not -operation for argument.
  Argument is assumed to be boolean expression.
  Number of arguments need to be 1.
  Return value is boolean type:
    - true: if argument is false
    - false: if argument is true

Usage: not(<boolean-expr>)
`,
		"eq": `
Operator: eq
  Evaluates two arguments and compares resulting values.
  Number of arguments need to be 2.
  Return value is boolean type:
    - true: if values are equal
    - false: if values differ

Usage: eq(<expr-1> <expr-2>)

Note. not all types are comparable (function/procedure values)
`,
		"if": `
Operator: if
  Evaluates 1st argument and based on value evaluates eiher 2nd
  or 3rd argument. It's assumed that 1st argument evalutes to
  boolean value:
    - true: 2nd argument is evaluated and returned as value
    - false: 3rd argument is evaluated and returned as value

  Number of arguments need to be 3.
  Return value is evaluated value of either 2nd or 3rd argument
  expression.

Usage: if(<condition-expression> <expr-1> <expr-2>)
`,
		"plus": `
Operator: plus
  Performs either arithmetic sum operation or string concatenation
  depending on input argument types.
  All arguments need to be of same type.
  Input arguments can be of type:
    - int: evaluates to arithmetic sum
    - float: evaluates to arithmetic sum
    - string: concatenation of argument strings

  Number of arguments need to be at least 2.
  Return value is of type int/float/string depending on type
  of input arguments.

Usage: plus(<expr-1> <expr-2> <expr-3> ...)
`,
		"minus": `
Operator: minus
  Performs arithmetic subtraction operation for two input
  arguments so that 2nd argument is subtracted from 1st one.
  Both arguments need to be of same type.
  Input arguments can be of type:
    - int
    - float

  Number of arguments need to be 2.
  Return value is result of subtraction of type int/float
  depending on type of input arguments.

Usage: minus(<expr-1> <expr-2>)
`,
		"mul": `
Operator: mul
  Performs arithmetic multiplication operation.
  Input arguments can be of type:
    - int
    - float

  Number of arguments need to be at least 2.
  Return value is multiplication result of input arguments.
  If any of arguments is of type float then result value type is float.

Usage: mul(<expr-1> <expr-2> <expr-3> ...)
`,
		"div": `
Operator: div
  Performs arithmetic division operation for two input arguments so
  that 1st argument is dividend and 2nd argument is divisor.
  Both arguments need to be of same type.
  Input arguments can be of type:
    - int
    - float
  If both arguments are of type int then result is of type int,
  otherwise result is of type float.
  In case of division of two int's result is quotient of division
  operation.

  Number of arguments need to be 2.
  Return value is result of division of type int/float
  depending on type of input arguments.

Note. If divisor is zero (int or float) runtime error
      is generated.

Usage: div(<expr-1> <expr-2>)
`,
		"mod": `
Operator: mod
  Performs modulo operation for two input arguments, result is
  remainder of division operation where 1st argument is dividend
  and 2nd argument is divisor.
  Input arguments need to be of type int.

  Number of arguments need to be 2.
  Return value is remainder value of division of input arguments.

Note. If divisor is zero (int or float) runtime error
      is generated.

Usage: mod(<expr-1> <expr-2>)
`,
		"list": `
Operator: list
  Creates list value. Arguments can be of any type.
  Arguments are evaluated to values which are put to list.
  Order of items in list is order of arguments.

  Number of arguments is not restricted (if none, empty list is created).
  Return value is list value.

Usage: list(<expr-1> <expr-2> <expr-2> ...)
`,
		"empty": `
Operator: empty
  Returns true if list/map is empty, false otherwise.

  Number of arguments need to be 1. Argument is assumed
  to be either list or map.
  Return value is boolean value.

Usage: empty(<expr>)
`,
		"head": `
Operator: head
  Returns 1st item from list which is given as argument.

  Number of arguments need to be 1. Argument is assumed
  to be list.
  Return value may be any type of value.

Note. if list is empty runtime error is generated.

Usage: head(<list-expr>)
`,
		"last": `
Operator: last
  Returns last item from list which is given as argument.

  Number of arguments need to be 1. Argument is assumed
  to be list.
  Return value may be any type of value.

Note. if list is empty runtime error is generated.

Usage: last(<list-expr>)
`,
		"rest": `
Operator: rest
  Returns rest of the list given as argument, excluding 1st item
  (head) of list.

  Number of arguments need to be 1. Argument is assumed
  to be list.
  Return value is list.

Note. if list is empty runtime error is generated.

Usage: rest(<list-expr>)
`,
		"append": `
Operator: append
  Appends value(s) to the end of the list.

  Number of arguments need to be at least 1. First argument is assumed
  to be list. Following arguments are added to the end of list in that
  same order (2nd argument, 3rd argument etc.)
  Return value is list (with appended values).

Note. if there's just one argument (list) then equal list is returned.

Usage: append(<list-expr> <expr> <expr> ...)
`,
		"add": `
Operator: add
  Adds value(s) in the front of the list.

  Number of arguments need to be at least 1. First argument is assumed
  to be list. Following arguments are added to the end of list in that
  same order (2nd argument, 3rd argument etc.)
  Return value is list (with added values).

Note. if there's just one argument (list) then equal list is returned.

Usage: add(<list-expr> <expr> <expr> ...)
`,
		"len": `
Operator: len
  Returns length of list/map/string.

  Number of arguments need to be 1. Argument is assumed
  to be list, map or string. Length is evaluated as follows:
    - list: number of items in list
    - map: number of key-value pairs in map
    - string: number of characters in string
  Return value is type of int (length of argument value).

Usage: len(<expr>)
`,
		"type": `
Operator: type
  Returns type of value evaluated from argument expression.

  Number of arguments need to be 1.
  Return value is type of string, returning following values
  depending on argument value type:
    - int: 'int'
    - float: 'float'
    - bool: 'bool'
    - string: 'string'
    - function: 'function' (also for procedure)
    - list: 'list'
    - channel: 'channel'
    - map: 'map'
    - opaque value: 'opaque:' + opaque type specific name (string)
    - external procedure/function: 'ext-proc'

Note. Runtime error is generated in case type -operator is called
      for opaque value in function call.
      (could otherwise cause impure side-effects in theory)

Usage: type(<expr>)
`,
		"in": `
Operator: in
  Returns true if 2nd argument is included in list, map or string given
  as 1st argument. Inclusion means:
    - for list: value equals to some value in list
    - for map: value equals to some key value of map
    - for string: (string) value is substring of string

  Number of arguments need to be 2.
  Return value is boolean value:
    - true: if 2nd argument is included in value of 1st argument
    - false: if 2nd argument is not included in value of 1st argument

Note. In case 1st argument is type of string then it's assumed
      that 2nd argument is also type of string
      (otherwise runtime error is generated)

Usage: in(<list/map/string-expr> <expr>)
`,
		"ind": `
Operator: ind
  Returns value in list/string (given as 1st argument) located
  in location defined by index value given as 2nd argument (int).
  This means that if 1st argument is:
    - list: returns item in list in location defined by index
    - string: returns character (string) located in location defined by index

  Number of arguments need to be 2.
  Return value is:
    - in case of list: value of item in given location -> any type
    - in case of string: string (one character)

Note. If index defined by 2nd argument points outside the limits
      of list/string then runtime error is generated

Usage: ind(<list/string-expr> <index-expr>)
`,
		"find": `
Operator: find
  Returns list of index values (int) which point the locations
  of values (defined by 2nd argument) in list/string (1st argument).
  Finds equal values in list and returns indexes to those in list.
  Finds substring locations in string and and returns indexes to those in list.

  Number of arguments need to be 2.
  Return value is list, contains index values (of type int).

Note. if value is not found in list/string then empty list is returned

Usage: find(<list/string-expr> <index-expr>)
`,
		"slice": `
Operator: slice
  Returns "slice" (sub-list/sub-string) of list/string (1st argument)
  defined by 2nd argument, and optionally by 3rd argument.
  2nd argument defines starting index from where rest of items
  in list or characters in string are included.
  3rd argument is optional and defines end location until which
  items/characers are included.

  Number of arguments need to be 2 or 3. 1st argument is assumed
  to be list or string. 2nd and 3rd argument are assumed to be type
  of int.

Examples:
  slice('12345' 3) -> '45'
  slice('12345' 2) -> '345'
  slice('12345' 5) -> Runtime error:  slice: Index out of range (2nd: 5)
  slice('12345678' 3 4) -> '45'
  slice('12345' 3 3) -> '4'
  slice('12345' 3 2) -> Runtime error:  slice: 3rd index (2) should not be less than 2nd one (3)
  slice(list(1 2 3 4 5) 2 3) -> list(3, 4)

Note. It's assumed that 2nd argument is equal or less than 3rd argument
      (otherwise runtime error is generated)
Note. If 2nd argument is out of range of list/string then runtime
      error is generated

Usage: slice(<list/string-expr> <begin-index-expr>)
       slice(<list/string-expr> <begin-index-expr> <end-index-expr>)
`,
		"rrest": `
Operator: rrest
  Returns list which is equal to list given as argument but last item
  not included.

  Number of arguments need to be 1. Argument is assumed to be list.

Examples:
  rrest(list(1 2 3 4)) -> list(1, 2, 3)
  rrest(list(1)) -> list()
  rrest(list()) -> Runtime error:  Attempt to access empty list in rrest operator

Note. If list given as argument is empty runtime error is generated

Usage: rrest(<list-expr>)
`,
		"reverse": `
Operator: reverse
  Returns list given as argument in reverse order.

  Number of arguments need to be 1. Argument is assumed to be list.

Examples:
  reverse(list(1 2 3 4)) -> list(4, 3, 2, 1)
  reverse(list(1)) -> list(1)
  reverse(list()) -> list()

Usage: reverse(<list-expr>)
`,
		"extend": `
Operator: extend
  Returns list containing items which are items of lists which are
  given as arguments.

  Number of arguments can be anything from zero to upwards.
  Arguments are assumed to be list types.

Examples:
  extend(list(1 2) list(3 4) list(5 6)) -> list(1, 2, 3, 4, 5, 6)
  extend(list(1 2) list() list(5 6)) -> list(1, 2, 5, 6)
  extend(list(1 2)) -> list(1, 2)
  extend(list()) -> list()
  extend() -> list()

Usage: extend(<list-expr> <list-expr> ...)
`,
		"split": `
Operator: split
  Returns list of substrings that result from splitting string given
  as 1st argument by string given as 2nd argument.
  If there's just one argument then splitting is made by one or more
  consecutive whitespace characters.

  Number of arguments can be 1 or 2.
  Arguments are assumed to be string types.

Examples:
  split('abcd,abcd,abcd' ',') -> list('abcd', 'abcd', 'abcd')
  split('first and second and third' 'and') -> list('first ', ' second ', ' third')
  split('first and second and third' ' and ') -> list('first', 'second', 'third')
  split('some text') -> list('some', 'text')
  split('sometext') -> list('sometext')
  split('some text' '') -> list('s', 'o', 'm', 'e', ' ', 't', 'e', 'x', 't')

Note. splitting with empty string ('') results list containing all characters of
      1st argument as strings

Usage: split(<string-expr> <string-expr>)
       split(<string-expr>)
`,
		"gt": `
Operator: gt
  Returns true if 1st argument is greater than 2nd argument (in arithmetic sense),
  false otherwise.

  Number of arguments is assumed to be 2.
  Arguments are assumed to be int or float type.

Usage: gt(<expr> <expr>)
`,
		"lt": `
Operator: lt
  Returns true if 1st argument is less than 2nd argument (in arithmetic sense),
  false otherwise.

  Number of arguments is assumed to be 2.
  Arguments are assumed to be int or float type.

Usage: lt(<expr> <expr>)
`,
		"le": `
Operator: le
  Returns true if 1st argument is less than or equal to 2nd argument (in arithmetic sense),
  false otherwise.

  Number of arguments is assumed to be 2.
  Arguments are assumed to be int or float type.

Usage: le(<expr> <expr>)
`,
		"ge": `
Operator: ge
  Returns true if 1st argument is greater than or equal to 2nd argument (in arithmetic sense),
  false otherwise.

  Number of arguments is assumed to be 2.
  Arguments are assumed to be int or float type.

Usage: ge(<expr> <expr>)
`,
		"str": `
Operator: str
  Return string representation of argument.

  Number of arguments is assumed to be 1.
  Argument can be of any type, return value is string.

Usage: str(<expr>)
`,
		"conv": `
Operator: conv
  Converts 1st argument to type defined by 2nd argument (string), if possible.

  Number of arguments is assumed to be 2.
  First argument can be of any type, 2nd argument is assumed to string
  having one of follwong values:
    - 'string' : converts 1st argument to string
    - 'list' : converts string to list containing all string characters as items
    - 'float' : converts int value to float value
    - 'int' : converts float value to int or string value to int
  Return value is converted value.

Examples:
  conv(100 'string') -> '100'
  conv(list(1 2 3) 'string') -> 'list(1, 2, 3)'
  conv('abcd' 'list') -> list('a', 'b', 'c', 'd')
  conv(100 'float') -> 100 (float)
  conv(10.5 'int') -> 10
  conv('100' 'int') -> 100
  conv('abc' 'int') -> 'Not able to convert to int'

Note. If conversion cannot be done Runtime error is generated, exception
      is failing conversion from string to int:
      string 'Not able to convert to int' is returned

Usage: conv(<expr> <expr>)
`,
		"case": `
Operator: case
  Evaluates 1st argument and compares it to first matching value and returns
  corresponding value in case of match.
  Last argument can be used for returning default value in case no
  match was found.

  Number of arguments is assumed to be at least 2.

  Arguments that are compared can be of any comparable type.
  (int/bool/string/float/list/map, opaque type if compared in function call)
  Other arguments can be of any type.
  Last argument (default value) is optional, if it's not given and there is
  no match then Runtime error is generated.

Usage: case( <expr>
         <compare-expr-1> <value-1>
         <compare-expr-2> <value-2>
         ...
         <default-value>
       )

       case( <expr>
         <compare-expr-1> <value-1>
         <compare-expr-2> <value-2>
         ...
       )
`,
		"name": `
Operator: name
  Return string representation of symbol given as argument.

  Number of arguments is assumed to be 1.
  Argument is assumed to be symbol (not value nor operator call).
  Return value is string.

Example:
  name(some-symbol) -> 'some-symbol'

Note. argument is not evaluated

Usage: name(<symbol>)
`,
		"error": `
Operator: error
  Generates runtime error. Prints string representation
  of possible arguments.

  Number of arguments can be anything from 0 upwards.
  Arguments can be of any type, string representation is printed.
  There is no return value as execution is not continued.

Example:
  error() -> Runtime error:
  error('...some error...') -> Runtime error:  ...some error...
  error('...some error...' list(1 2 3)) -> Runtime error:  ...some error...list(1, 2, 3)
  error('...some error...' list(1 2 3):) -> Runtime error:  ...some error...123

Usage: error()
       error(<expr>)
       error(<expr> <expr> <expr> ...)
`,
		"print": `
Operator: print
  Concatenates string representations of arguments and prints that
  to screen/stdout.

  Number of arguments can be anything from 0 upwards.
  Arguments can be of any type, string representation is concatenated and printed.
  Return value is always true (boolean).

Note. As print is meant to be used for debugging purposes (stdio having better
      procedures for printing to console) it's allowed in functions only
      if printing in functions -mode is enabled.
      That's because pure functions should not have I/O side-effects.
      This mode is controlled by -noprint option.
      (default mode is that printing is allowed in functions)
      If printing is not allowed in functions and that is done then
      runtime error is generated.

Usage: print(<expr> <expr> <expr> ...)
`,
		"spawn": `
Operator: spawn
  Starts fiber (lightweight unit of execution thread) for each argument
  to evaluate that argument.

  Number of arguments can be anything from 0 upwards.
  Arguments can be of any type.
  Return value is always true (boolean).

Note. spawn is not allowed to be called from function (only from procedure).

Usage: spawn(<expr> <expr> <expr> ...)
`,
		"chan": `
Operator: chan
  Creates channel value. If no arguments are given then
  channel is unbuffered (meaning send waits until there's someone
  receiving on channel). If argument is given it's assumed to be
  int -value which defines buffer size for channel (in which case channel
  is buffered one, meaning that send may not block until there's reader).
  Return value is channel value.

Usage: chan()
       chan(<int: buffer-size>)
`,
		"send": `
Operator: send
  Evaluates and sends value (given as 2nd argument) to channel
  (given as 1st argument). By default blocks execution until receiving
  fiber reads channel, however it can be defined with optional 3rd argument
  that execution does not block to writing to channel (so that if channel is
  full current fiber does not block to wait other fibers to read from it)

  Requires 2 or 3 arguments.

  Optional 3rd argument is map (options map) which can have following name-values:
  - 'wait' : bool value:
      - true -> block to wait if channel is full
      - false -> returns if channel is full (no waiting)
      Default is true (blocks to writing to channel)

  Return value is true if value was written to channel, false if value was not
  written to channel.

Note. send is not allowed to be called from function (only from procedure).

Example:
  ch = chan(100)
  was-added = send(ch 'some value')
  was-added = send(ch 'some value' map('wait' false))

Usage: send(<channel-expr> <expr>)
       send(<channel-expr> <expr> <options-map>)
`,
		"recwith": `
Operator: recwith
  Receives value from channel (given as 1st argument).

  Requires 2 argument, 1st argument is assumed to be channel.
  Second argument is map (options map) which can have following name-values:
  - 'wait' : bool value: 
      - true -> block to wait if channel is empty
      - false -> returns if channel is empty (no waiting)
      Default is true (blocks to receiving from channel)
  - 'limit-sec' : int-value, number of seconds to wait in channel
                  (returns after time limit if no items received from channel)
  - 'limit-nanosec' : int-value, number of nanoseconds to wait in channel
                  (returns after time limit if no items received from channel)
  Default is that there is no time limit (waits forever).

  Return value is list of two items.

  List returned has following items:
  1) First item (bool) is true if value was received from channel.
     If no value was received then it's false.
  2) Second item is value received from channel ('' if value not received)

Note. recwith is not allowed to be called from function (only from procedure).

Example:
  ch = chan()
  value-received, value = recwith(ch map('wait' false)):

Usage: recwith(<channel-expr> <options-map>)
`,
		"recv": `
Operator: recv
  Receives value from channel (given as argument). Blocks until
  there's value available in channel.

  Requires 1 argument, argument is assumed to be channel.
  Return value is value received from channel.

Note. recv is not allowed to be called from function (only from procedure).

Usage: recv(<channel-expr>)
`,
		"symval": `
Operator: symval
  Symval returns value represented by symbol name given as string argument.

  Requires 1 argument, argument is assumed to be string type.
  Return value is value represented by symbol in scope.

Note. if symbol is not found from current scope then runtime error is generated.
Note. symval is not allowed to be called from function (only procedure allowed),
      otherwise runtime error is generated.

Example:
  call(proc() some-sym = 10 symval('some-sym') end) -> 10
  call(proc() some-sym = 10 symval('rubbish') end) -> Runtime error:  symval: symbol not found (rubbish)
  call(func() some-sym = 10 symval('some-sym') end) -> Runtime error:  symval not allowed in function

Usage: symval(<string-expr>)
`,
		"try": `
Operator: try
  Catches runtime error if such happens during evaluation of 1st argument. If no
  runtime error happens then try returns evaluated value of 1st argument.
  If runtime error happens then runtime error text is returned (as string value)
  unless 2nd argument is given in which case 2nd argument is evaluated and value
  of that is returned.

  Requires 1 or 2 arguments, arguments can be of any type.
  Return value is value of 1st argument evaluated, or in runtime error case
  2nd argument if such exists, otherwise runtime error text as string.

Note. try is not allowed to be called from function (only procedure allowed),
      otherwise runtime error is generated.

Usage: try(<expr>)
       try(<expr> <expr>)
`,
		"select": `
Operator: select
  Receives input value from any of several channels and calls channel specific
  handler function/procedure. Channels and handler functions/procedures can be given as separate
  arguments or as two lists containing channels and handlers. Return value of
  handler is returned as value from select.

  Can have two kind of arguments:
    1) channel and handler (proc) pairs:
       as channel N:th argument and corresponding handler N+1:th argument
    2) two lists:
       as 1st list containing all channels and 2nd list containing
       all handlers so that channel and related handler are in same index
       in lists
  Number of arguments must not be zero and must be even number.
  Return value is value returned from handler.

Note. select is not allowed to be called from function (only procedure allowed),
      otherwise runtime error is generated.

Usage: select(
         <channel-expr> <handler-expr>
         <channel-expr> <handler-expr>
         <channel-expr> <handler-expr>
         ...
       )
       select(<list-of-channels-expr> <list-of-handlers-expr>)
`,
		"eval": `
Operator: eval
  Evaluates expression given as string argument and returns resulted
  value of evaluation.

  Requires 1 argument, argument is assumed to be string.
  Return value result of evaluated expression represented as argument.

Example:
  eval('plus(1 2 3)') -> 6
  eval(sprintf('plus(%d %d %d)' list(1 2 3):)) -> 6

Usage: eval(<string-expression>)
`,
		"while": `
Operator: while
  Similar to call -operator but can be used for "tail call optimization".
  if condition (1st argument) is true while reconstructs current call frame
  by re-evaluating all current frame call arguments with ones following
  1st argument (2nd, 3rd etc.), as many as current frame has call arguments.
  Also all let -definitions are re-evaluated in innermost frame.
  Last argument is returned when condition becomes false.

  Use cases:
    - recursive call with tail call optimization (not consuming call stack)
    - in procedure to implement some I/O event loop kind of handling (or
      channel reading)

  Requires at least 2 arguments.

Note. while -operator can only be used in function/procedure body, not
      in let -definitions, nor as function/procedure calls as arguments.

Usage: while(<condition-expr> <arg-1> <arg-2> ... <result-expr>)
`,
		"float": `
Operator: float
  Constructs float value. Converts int value given as argument to corresponding
  float value (if float value is given then same is returned).

  Requires 1 argument. Return value type is float.

Usage: float(<expr>)
`,
		"map": `
Operator: map
  Creates (persistent) map. Arguments are interpreted so that n:th
  (0, 2, 4, ...) argument is key and following argument is corresponding value
  (n+1:th: 1, 3, 5, ...). If no arguments are given then empty map
  is created. Map item values can be of any type but map keys
  can be only:
    - int
    - string
    - float
    - list
    - boolean
  If same key is given twice then runtime error is generated.

  There must be even number of arguments otherwise runtime error is
  generated (or no arguments).
  Return value is map value that was created.

Example:
  map() -> map()
  map(1 2 3 4) -> map(1 : 2, 3 : 4)
  map(1 2 3) -> Runtime error:  map: uneven amount of arguments (3)

Usage: map(<key> <value> <key> <value> ...)
`,
		"put": `
Operator: put
  Puts key-value pair to map. First argument is map, 2nd argument
  is key and 3rd argument is value.
  Map item value can be of any type but map key can be only:
    - int
    - string
    - float
    - list
    - boolean
  If key is already in map then runtime error is generated.

  There must be 3 arguments.
  Return value is map value with added key-value pair.

Usage: put(<map> <key> <value>)
`,
		"get": `
Operator: get
  Gets value for given key from map. If key is not found
  runtime error is generated.

  First argument is map, 2nd argument is key.

  There must be 2 arguments.
  Return value is value corresponding to key.

Example:
  get(map(1 2 3 4) 1) -> 2
  get(map(1 2 3 4) 10) -> Runtime error:  get: key not found (10)

Usage: get(<map> <key>)
`,

		"getl": `
Operator: getl
  Gets value for given key from map. Returns result as
  list of two items where 1st item (boolean) is true if key was
  found in map, otherwise false. Second item in list is corresponding
  value if key was found, otherwise value is false.

  First argument is map, 2nd argument is key.

  There must be 2 arguments.
  Return value is list of 2 items:
    - First item true if key found, false otherwise
    - Second item value corresponding to key if key was found,
      otherwise false value

Example:
  getl(map(1 2 3 4) 1) -> list(true, 2)
  getl(map(1 2 3 4) 10) -> list(false, false)
  is-key-found value = getl(map(1 2 3 4) 1):

Usage: getl(<map> <key>)
`,
		"keys": `
Operator: keys
  Returns all keys of map (given as argument) as list.
  There's not any particular order guaranteed in list.

  There must be one argument which is map.
  Return value is list (conatins all keys of map).

Example:
  keys(map(1 2 3 4)) -> list(1, 3)
  keys(map()) -> list()

Usage: keys(<map>)
`,
		"vals": `
Operator: vals
  Returns all values of map (given as argument) as list.
  There's not any particular order guaranteed in list.

  There must be one argument which is map.
  Return value is list (conatins all values of map).

Example:
  vals(map(1 2 3 4)) -> list(2, 4)
  vals(map()) -> list()

Usage: vals(<map>)
`,
		"keyvals": `
Operator: keyvals
  Returns all key-value pairs of map (given as argument) as lists.
  There's not any particular order guaranteed in list.
  Each key-value pair is represented as its own list (of 12 items):
    - 1st item is key
    - 2nd item is value

  There must be one argument which is map.
  Return value is list containing lists representing
  key-value pairs.

Example:
  keyvals(map(1 2 3 4)) -> list(list(1, 2), list(3, 4))
  keyvals(map()) -> list()

Usage: keyvals(<map>)
`,
		"let": `
Operator: let
  Let-definition by using operator let.
  First argument is symbol for which value is assigned
  from evaluating expression given as 2nd argument.
  Symbol value is assgined to in current scope.

  x = 100
  is identical to:
  _ = let(x 100)

  Can be used in REPL (option -repl) to set some let-definitions.

  There must be 2 arguments: 1st argument needs to be symbol and
  2nd argument is any expression.
  Return value is evaluation result (value) from 2nd argument
  (same value which is assigned to symbol).

Example:
  keyvals(map(1 2 3 4)) -> list(list(1, 2), list(3, 4))
  keyvals(map()) -> list()

Usage: let(<symbol> <expr>)
`,
		"imp": `
Operator: imp
  Imports module by using operator imp.
  First argument is symbol which is name of module.
  Module is returned as map in which:
    - keys are names of functions/procedures/other values (as string)
    - corresponding values are function/procedure/other values

  import stdfiles
  is equivalent to:
  my-file-mod = imp(stdfiles)
  usage is via map:
  call(get(my-file-mod 'cwd'))
  same as -> call(stdfiles.cwd)

  There must be 1 argument which is symbol representing
  module name to be imported.
  Return value is map containing symbol names (strings)
  as keys and corresponding values.

Note. if module of given symbol is not found then
      runtime error is generated.

Example:
  imp(stdlog) -> map('get-logger' : ext-proc, 'get-default-logger' : ext-proc)
  imp(not-to-found) -> Runtime error:  Module not found: not-to-found

Usage: imp(<symbol>)
`,
		"del": `
Operator: del
  Returns map value based on map given as 1st argument from
  which key-value pair defined by key given as 2nd argument
  is removed.
  If key does not exist is map given as 1st argument then
  runtime error is generated.

  There must be 2 arguments: 1st argument is source map and
  2nd argument defines key for which key-value pair is to be removed.

Example:
  del(map(1 2 3 4) 3) -> map(1 : 2)
  del(map(1 2 3 4) 9) -> Runtime error:  del: key not found

Usage: del(<map> <key>)
`,
		"dell": `
Operator: dell
  Similar to del-operator that produces map (based on map
  given as 1st argument) from which key-value pair defined
  by key given as 2nd argument
  Returns list (of two items):
    - first item (boolean): true if key was found, false otherwise
    - second item (map):
        1. if 1st item is true, then map without given key-value pair
        2. if 1st item is false, original map (1st argument)
  (no runtime error is generated when key is not found)

  There must be 2 arguments: 1st argument is source map and
  2nd argument defines key for which key-value pair is to be removed.

Example:
  dell(map(1 2 3 4) 3) -> list(true, map(1 : 2))
  dell(map(1 2 3 4) 9) -> list(false, map(1 : 2, 3 : 4))
  found newmap = dell(map(1 2 3 4) 3):

Usage: dell(<map> <key>)
`,
		"sprintf": `
Operator: sprintf
  Formats according to a format specifier and returns the resulting string.
  First argument is format string and following arguments are operands
  for formation.

  There must be at least 1 argument.
  First argument is format string (type of string) and following
  ones operands.
  Return value is string (formatted).

Example:
  sprintf('%d : %v : %s : %f : %v' 10 true 'some text' 0.5 list(1 2 3))
    -> '10 : true : some text : 0.500000 : list(1, 2, 3)'

Usage: sprintf(<format-string> <expr> <expr> ...)
`,
		"argslist": `
Operator: argslist
  Returns list of arguments given to current (innermost) function/procedure
  call context.

  Requires no arguments. Return value is list type. List contains
  all argument values in same order as given in function/procedure call.

Example:
  call(func() argslist() end 1 2 3) -> list(1, 2, 3)
  call(func(p1 p2 p3) argslist() end 1 2 3) -> list(1, 2, 3)
  call(func(p1 _ p3) argslist() end 1 2 3) -> list(1, 2, 3)
  call(func(_ _ _) argslist() end 1 2 3) -> list(1, 2, 3)

Usage: argslist()
`,
		"cond": `
Operator: cond
  Multiway conditional expression (multiway if).
  Arguments are interpreted so that there's multiple
  conditional pairs and last argument defines value in
  case no other pair matches.
  First condition expression that returns true causes
  corresponding expression to be evaluated and returned
  from cond -operator.

  cond(
    <1st condition (boolean expression)>
    <1st expression: evaluated if 1st condition is true>
    <2nd condition (boolean expression)>
    <2nd expression: evaluated if 2nd condition is true>
    <3rd condition (boolean expression)>
    <3rd expression: evaluated if 3rd condition is true>
    ...
    <default expression: evaluated if no other condition is true>
  )

  Requires at least 3 arguments. Assumes that there's always
  (2 * n) + 1 arguments as there needs to be 1...n condition-expression
  pairs and one default expression.
  Return value is value corresponding to matching condition or default.

Example:
  cond(eq(1 1) 10 eq(3 4) 20 'default') -> 10
  cond(eq(1 2) 10 eq(3 3) 20 'default') -> 20
  cond(eq(1 2) 10 eq(3 4) 20 'default') -> 'default'

Usage: cond(
         <condition (bool)> <expr>
         <condition (bool)> <expr>
         <condition (bool)> <expr>
         ...
         <default/else-expr>
       )
`,
		"help": `
Operator: help
  Returns documentation about certain language topic
  as string or in some cases list of strings.
  See more information by help().

  Assumes no arguments or one argument.
  Argument type is assumed to be string (defining topic).
  Return type can be:
    - string
    - list (of strings)

Usage: help()
       help(<topic-as-string>)
`,
	}
}

// Operators contains operator information (operator name as key)
type Operators map[string]OperatorInfo

// NewDefaultOperators returns default set of operators
func NewDefaultOperators() Operators {
	return Operators{
		"and":      OperatorInfo{},
		"or":       OperatorInfo{},
		"call":     OperatorInfo{},
		"not":      OperatorInfo{},
		"eq":       OperatorInfo{},
		"if":       OperatorInfo{},
		"plus":     OperatorInfo{},
		"minus":    OperatorInfo{},
		"mul":      OperatorInfo{},
		"div":      OperatorInfo{},
		"mod":      OperatorInfo{},
		"list":     OperatorInfo{},
		"empty":    OperatorInfo{},
		"head":     OperatorInfo{},
		"last":     OperatorInfo{},
		"rest":     OperatorInfo{},
		"append":   OperatorInfo{},
		"add":      OperatorInfo{},
		"len":      OperatorInfo{},
		"type":     OperatorInfo{},
		"in":       OperatorInfo{},
		"ind":      OperatorInfo{},
		"find":     OperatorInfo{},
		"slice":    OperatorInfo{},
		"rrest":    OperatorInfo{},
		"reverse":  OperatorInfo{},
		"extend":   OperatorInfo{},
		"split":    OperatorInfo{},
		"gt":       OperatorInfo{},
		"lt":       OperatorInfo{},
		"le":       OperatorInfo{},
		"ge":       OperatorInfo{},
		"str":      OperatorInfo{},
		"conv":     OperatorInfo{},
		"case":     OperatorInfo{},
		"name":     OperatorInfo{},
		"error":    OperatorInfo{},
		"print":    OperatorInfo{},
		"spawn":    OperatorInfo{},
		"chan":     OperatorInfo{},
		"send":     OperatorInfo{},
		"recv":     OperatorInfo{},
		"symval":   OperatorInfo{},
		"try":      OperatorInfo{},
		"select":   OperatorInfo{},
		"eval":     OperatorInfo{},
		"while":    OperatorInfo{},
		"float":    OperatorInfo{},
		"map":      OperatorInfo{},
		"put":      OperatorInfo{},
		"get":      OperatorInfo{},
		"getl":     OperatorInfo{},
		"keys":     OperatorInfo{},
		"vals":     OperatorInfo{},
		"keyvals":  OperatorInfo{},
		"let":      OperatorInfo{},
		"imp":      OperatorInfo{},
		"del":      OperatorInfo{},
		"dell":     OperatorInfo{},
		"sprintf":  OperatorInfo{},
		"argslist": OperatorInfo{},
		"cond":     OperatorInfo{},
		"help":     OperatorInfo{},
		"recwith":  OperatorInfo{},
	}
}

func (ops Operators) isOperator(operatorName string) bool {
	_, found := ops[operatorName]
	return found
}
