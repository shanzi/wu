# wu (呜~)

A minimal **W**atch **U**tility who can run and restart specified command in
response to file changes automatically.

This utility is intended to provide a tiny tool to automate the Edit-Build-Run
loop of development. Although it is quite similar to watch tasks of Grunt or Gulp,
`wu` is designed to be just a single command with simplest interfaces to work with.

# Install

To install `wu` from source code, you have to install Golang's tool chain first.

Install [godep](https://github.com/tools/godep):
```
go get github.com/tools/godep
```
Install packages dependencies:
```
godep get
```
Build:
```
make
```

Precompiled version can be found [here](https://github.com/shanzi/wu/releases).

# Usage

Run `wu -h` for help message:

```
Usage: wu [options] [command]
  -config string
        Config file (default ".wu.json")
  -dir string
        Directory to watch
  -pattern string
        Patterns to filter filenames
  -save
        Save options to conf
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

`wu` will try to read config file under current directory at the start (default: `.wu.json`),
you can user `-config` flag to specify the config file by hand. Use `-save` flag to save
current options.

```
wu -pattern="*.go" -save go build
# A `.wu.json` has been created

wu # Running `wu` without any options, it will read from `.wu.json`
```

# LICENSE

See [LICENSE](./LICENSE)
