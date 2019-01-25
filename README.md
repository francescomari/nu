# Node Utils

Node Utils is a command line application for processing the data format
generated by [Export Nodes](https://github.com/francescomari/export-nodes),
which is a serialization of a tree whose nodes can have properties, and each
property can have multiple values.

## Build

To build this project you need a recent version of Go. From the root of the
project, run

    go build

The output will be an executable called `nu`. Otherwise, you can run

    go get github.com/francescomari/nu

This command will download the project, build it, and create an executable in
your `GOPATH`.

## Usage

### Help

`nu` is composed of a bunch of commands. You can see the list of supported
commands by running

    nu help

You can also get more information about a specific command by running

    nu help [command]

where `command` is the name of a command.

### Extract node paths

    nu nodes <export.txt

You can extract the fully qualifed paths of every node in an export with the
`nodes` command. The command reads the export from stdin and prints the list of
fully qualifed paths on stdout.

### Extract property paths

    nu properties <export.txt

You can extract the fully qualified paths of every property in the export with
the `properties` command. The command reads the export from stdin and prints
both the fully qualified paths and the type of every properties to stdout.

### Compute statistics

    nu stats <export.txt

You can extract some statistics about the content tree within the export with
the `stats` command. The command reads the export from stdin and prints the
statistics on stdout.

### Shrink to a subtree

    nu subtree [path] <export.txt

Sometimes you are interested only in a part of the export. If you need to reduce
the focus of the export to a particular subtree, you can use the `subtree`
command. The command receives a mandatory argument `path` that identifies the
subtree you are interested in. The `path` needs to be an absolute path, e.g.
`/path/to/tree`. The command prints to stdout a new export, whose root is the
subtree at `path` from the original export.

Since `subtree` outputs a valid export, you can easily pipe its output into
other commands, like

    cat export.txt | nu subtree /path/to/tree | nu stats

## License

This software is released under the MIT license.