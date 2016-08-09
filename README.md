# wu

A minimal **w**atch **u**tility to watch changes on file system and run specified
command on events.

This utility is intended to provide a simple tool to automate the Edit-Build-Run
loop of development. Although it is much similar to watch tasks of Grunt or Gulp,
`wu` is designed to be just a single command with clean interfaces to work with.

# Install

To install `wu` from source code, you have to install Golang's tool chain first.
Then run:

```
go get github.com/shanzi/wu
go install github.com/shanzi/wu
```


# Examples

You just run you command with `wu`, `wu` will run your command at the start,
try to terminate previous when new changes take place and then start running a new one.
You can stop the process by sending a `SIGINT` signal (typically by `CTRL-C`):

```
wu sleep 10
```

Output:
```
Start watching...
- Running command: sleep 10
- Terminated.
File changed: /path/to/changed/file.txt
- Running command: sleep 10
- Done.
File changed: /path/to/changed/file.txt
- Running command: sleep 10
- Done.
^C
Shutting down...
```

Usually you can only run one command with `wu`, but it doesn't prevent you from
running complex command with `sh`, `bash` or other shell command:

```
wu sh -c 'echo "START" && sleep 5 && echo "END"'
```

Output:
```
Start watching...
- Running command: sh -c echo "START" && sleep 5 && echo "END"
START
END
- Done.
```

You can specified a pattern to filter the files to watch. Multiple patterns
should be seperated by spaces or commas:

```
wu -pattern="*.js, *.html"
```

Output (If no command specified, `wu` just log changed files):

```
Start watching...
File changed: /path/to/changed/file.js
File changed: /path/to/changed/file.html
```

One practical use case is to use `wu` with some light weight web frameworks,
For example, you can start a server written in go by:

```
wu -pattern="*.go" go run main.go
```

Use `wu -h` to view help for all options.

# LICENSE

See [LICENSE](./LICENSE)
