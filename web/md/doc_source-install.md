<!-- title: Install data from Source -->
<!-- description: How to install data from Source -->


# Install data from Source

<span style="color: red;">
**WARNING: Installing from source is not easy. It is much easier (and  recommended) to install with some of the precompiled binaries and installers  available [here](/doc/install).**
</span>

Installing from source is divided in parts:

1. Installing Go (1.2)
2. Installing version control systems
3. Installing data

## Installing Go (1.2)

This is the hardest part. Follow the official instructions for your platform available at [the Go install page](http://golang.org/doc/install).

Make sure you setup your **Go workspace**, your `$GOPATH` and `$PATH`, as  described.

You'll know you're done if you can run Go. For example, for me:

```
> go version
go version go1.2 darwin/amd64
```

And if you have your `$GOPATH` set. For example, for me:

```
> echo $GOPATH
/Users/jbenet/go
```

## Installing version control systems

`data` has many library dependencies. Thankfully, these can all be installed using the Go tool. Go tracks dependencies with repository URLs, fetches their code, compiles it, and installs it all automatically.

You'll need to install the various version control systems that Go uses to fetch the dependencies. As of this writing, you'll need:

1. `git`
1. `hg` (mercurial)
1. `bzr` (bazaar)

You may already have them. Try running these commands to see which one you're missing (the version number shouldn't matter):

```
> git version
git version 1.8.3.4
> hg version
Mercurial Distributed SCM (version 2.7.2)
> bzr version
Bazaar (bzr) 2.6.0
```

If any of those commands failed, then you need to install that one (probably `bzr`).

### From a package manager

These are available in your favorite package manager. Here are two popular ones:

MacOSX (if you don't have `brew` installed, do yourself a favor and get it from [the Homebrew webpage](http://brew.sh/).

```
brew install git hg bzr
```

Ubuntu

```
sudo apt-get -y install git mercurial bzr
```


### Other (Windows)

If your platform does not use package managers, then install from their websites:

1. http://git-scm.com/
1. http://mercurial.selenic.com/
1. http://bazaar.canonical.com/en/


## Installing data

Ok, you should have Go and the version control systems installed (if not, see above). The rest is easy. Simply:

```
# get the code
git clone https://github.com/jbenet/data

# move into the data directory
cd data

# fetch the dependencies
make deps

# compile + install data
make install
```

You should now be able to run `data`:

    > data version
    data version 0.1.0
