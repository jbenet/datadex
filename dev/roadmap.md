# datadex roadmap

This document briefly outlines desired features to implement.

## feature list

### package publishing

Publish package onto the index using data, git, or http.

    data publish
    git push datadex.io/<author>/<dataset>.git
    http POST datadex.io/publish Content-Type:application/yaml @Datafile

### package listing

List all the packages, both for users browsing and for data and other clients.

    datadex.io/all

### package viewing

View information about a package, for users + clients.

    datadex.io/<author>/<dataset>
    datadex.io/<author>/<dataset>/Packagefile
    datadex.io/<author>/<dataset>/Datafile       # latest

### package refs/versions

List all the refs/branches/tags/versions in a dataset.

    datadex.io/<author>/<dataset>/refs

Allow looking up/posting a ref:

    datadex.io/<author>/<dataset>/refs/<ref>

Allow browsing each ref, like github:

    datadex.io/<author>/<dataset>/tree/<ref>/
    datadex.io/<author>/<dataset>/blob/<ref>/Datafile

### package archives

Host downloadable archives of each released dataset.

    datadex.io/<author>/<dataset>/archive/<ref>.tar.gz
    datadex.io/<author>/<dataset>/archive/<ref>.zip

(Can 302 redirect elswhere.)

### package search

Allow searching by name, author, keywords, etc. Use info in Datafile.

    datadex.io/search?q=foo+bar

## wish list

github collaboration:

- discussion
- issues
