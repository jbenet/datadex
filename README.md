# datadex - dataset index

This is an index of datasets compatible with
[`data`](http://github.com/jbenet/data).


See the [roadmap](dev/roadmap.md).

## Development

Setup:

1. [install go](http://golang.org/doc/install)
2. Run:

    git clone https://github.com/jbenet/datadex
    cd datadex
    make

Run:

    # launch the server
    ./datadex -port 8080

    # make data talk to it
    data config index.datadex.url http://localhost:8080
