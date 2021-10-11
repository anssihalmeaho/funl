![](https://github.com/anssihalmeaho/funl/blob/master/hellow.png)

# FunL examples

## hello.fnl
Hello world example for FunL.

## dohelp.fnl
Demonstrates how to create some files consisting FunL operator helps.

## inspace.fnl
HTTP client example for inquiring people currently in space.

## tictac.fnl
Tic-tac-toe game (program against user).

## fizzbuzz.fnl
FizzBuzz implementation in FunL.

## ToDo Application (todo_client.fnl / todo_server.fnl / todo_common.fnl)
Implementation of HTTP client and server for simple todo-application.
There is also common module (todo_common) for client and server to use.
It demonstrates how to create HTTP (micro)services, using JSON, logging,
environment variables and shutting down HTTP server.

Client and server use **localhost** address with default port number _8003_.
Port number can be redefined with **TODO_SRV_PORT** environment variable.

For example:

    export TODO_SRV_PORT=8009

Server (todo_server.fnl) needs to started first and after that
client(s) (todo_client.fnl) can be started to use server.

Todo-items are JSON objects/maps with one mandatory field (**id**).
Value for **id** is allocated by server when item is added.

### todo_server.fnl
Server implementation for todo-application. Server maintains
(todo-)items in its memory. Current implementation doesn't store those
any permanent storage so those are lost when server is restarted.

#### Service paths provided by server
Inquire all items:

```
GET /items

response body -> JSON: array of JSON objects (object corresponding to todo-item)
```

Put new item:

```
POST /items
request body: JSON object (item)
```

Remove item:

```
DELETE /items/id/<id>
```

#### Starting server
Starting todo server:

```
./funla ./examples/todo_server.fnl
todo-server: 2020/05/10 21:11:40 :'...listening...'
```


### todo_client.fnl
Client implements command-line interface which receives
commands from user and uses todo-server via HTTP API
to implement those commands.

#### Starting client
Starting client:

```
./funla ./examples/todo_client.fnl
Welcome to Todo application client
todo> help
Input can be:
  help         -> prints this help
  ?            -> prints this help
  quit         -> exits repl
  exit         -> exits repl

  put <JSON value for item> -> adds item
  get                       -> prints all items
  del <id of item>          -> removes item

todo>
```

#### put -command
Item is added with **put** -command.
Item data is given directly as JSON object.

```
todo> put {"name": "buy food", "tag": "home"}
ok
todo> put {"name": "clean the house", "tag": "home"}
ok
```

#### get -command
With **get** -command user can list all items.

```
todo> get

 item:
  - id : 100
  - tag : home
  - name : buy food

 item:
  - id : 101
  - tag : home
  - name : clean the house
```

#### del -command
Item is removed with **del** -command.
Item **id** needs to be given as argument for the command.

```
todo> del 100
ok
todo> get

 item:
  - id : 101
  - tag : home
  - name : clean the house

todo>
```

### todo_common.fnl
Common module **todo_common** provides procedure (_get-port_) which provides
port number for clöient and server. It's by default 8003 but can redefined
with **TODO_SRV_PORT** environment value.

