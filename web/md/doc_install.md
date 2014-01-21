<!-- title: Install data -->
<!-- description: How to install data -->


# Install data

There are many ways to install data. Scroll down to select your preferred method.

If you run into trouble installing, get help via the mailing list, or github issues (see the bottom-right corner of this page).

## Downloadable Installers

data will be available to install with platform-specific downloadable installers. coming soon.

## Package Managers

data will be available to install with your favorite software package manager. coming soon.

## Platform Binaries

Download precompiled binaries for various platforms.

- darwin 32-bit (Mac OS X) - coming soon
- [darwin 64-bit (Mac OS X)](https://github.com/jbenet/data/releases/download/v0.1.0/data-v0.1.0-darwin_amd64.tar.gz)
- [linux 32-bit](https://github.com/jbenet/data/releases/download/v0.1.0/data-v0.1.0-linux_386.tar.gz)
- [linux 64-bit](https://github.com/jbenet/data/releases/download/v0.1.0/data-v0.1.0-linux_amd64.tar.gz)
- windows 32-bit - coming soon
- windows 64-bit - coming soon

Each archive has instructions on how to install. For linux/osx, just put the binary somewhere in your path, e.g.:

    sudo mv data /usr/bin/data




## From Source

Installing from source is actually very easy. The only hard part is installing Go.

1. First, install Go (1.2+). Check out [the Go install page](http://golang.org/doc/install) for instructions.
   Make sure you set your `$GOPATH` and `$PATH`, as described.

2. Get the data source code. Either with git, or wget:
        # clone the repository with git
        git clone https://github.com/jbenet/data

        # download an archive
        wget https://github.com/jbenet/data/archive/latest-release.zip -O data.tar.gz
        tar xzf data.tar.gz

3. Build data (and get dependencies)
        make deps
        make install


You should now be able to run `data`:

    > data version
    data version 0.1.0
