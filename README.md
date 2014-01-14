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


Or, using docker:

    # first clone data inside (private repo workaround)
    git clone git@github.com:jbenet/data

    # build it
    docker build -t="datadex" .

    # run it
    docker run -p=8080:8080 -d datadex
