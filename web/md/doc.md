<!-- title: Documentation -->
<!-- description: Documentation for Data & Datadex -->

# Documentation

## Quick Start

First, install data. See [this page](/doc/install) for instructions.
<br />

### Downloading Datasets

In Short, run:

```
data get <dataset>
```

In Long:

To download a dataset, just run `data get <dataset path>`. You can find the dataset path on the dataset webpage:

![](http://jbenet.static.s3.amazonaws.com/d3a80c0b3a1c8dcc9088e9a4e0097b1f548784f6/example-zipcodes-id.png)

In this case, the path is `jbenet/zipcodes-example@1.0`.
This format (`<owner>/<dataset id>@<version>`) includes:

- the package owner: `jbenet`
- the dataset id: `zipcodes-example`
- the dataset version: `1.0`

(You can omit the version -- `data get jbenet/zipcodes-example` -- which will get the latest published version).

So, simply run:

```
> data get jbenet/zipcodes-example@1.0
Downloading jbenet/zipcodes-example@1.0 from datadex (http://datadex.io/api/v1).
get blob 8001ee9 .data/Manifest
get blob 10dcbf5 Datafile
get blob b6d9e12 README.md
get blob 4b1da58 zipcode_cities.txt
get blob b0f41fc zipcode_states.txt
get blob 8001ee9 .data/Manifest

---------
Installed jbenet/zipcodes-example@1.0 at datasets/jbenet/zipcodes-example@1.0
```
<br />

### Publishing a Dataset

#### Step 1: Set up Username

First you need to tell data your name, so that it can properly attribute the datasets you publish.

```
$ data user add <Your Username Here>
```

#### Step 2: Publish!

```
# in the directory you want to publish
$ data publish
```

## Datafiles

`data` tracks the definition of dataset packages using a file named `Datafile`. This file contains important metadata that `data` uses to find, install, and index the package. The `Datafile` also contains useful publication information (such as a list of authors). It is included when others download the package, and displayed on every dataset's webpage ([see example](http://datadex.io/jbenet/zipcodes-example@1.0)).

A `Datafile` is a [YAML](http://yaml.org) document with several fields, some required and some optional. For example, here is the `Datafile` of [CIFAR-10](http://datadex.io/jbenet/cifar-10@1.0-py):

```
% cat Datafile
dataset: jbenet/cifar-10@1.0-py
tagline: labeled subsets of the 80 million tiny images dataset
website: http://www.cs.toronto.edu/~kriz/cifar.html
authors:
- Alex Krizhevsky
- Vinod Nair
- Geoffrey Hinton
```

Below is a listing of all the fields `Datafile` can use:

```
dataset: <owner>/<dataset id>@<version>
tagline: <a title or one-line description>
description: '<a longer description of the dataset.
  It can span multiple lines.>'
authors:
- Full Name <email@address.org>
- Another Name <another@example.com>
contributors:
- Yeta Nother <person@helping.org>
sources:
- <urls or other references to data sources>
license: <name or url to the license>
website: <a url for the dataset's homepage>
repository: <a url to the dataset's repository, if different>
```

When publishing a dataset, `data publish` will allow you to set the important indexing (owner, dataset id, etc) fields. If you'd like to include more information than the required minimum, make sure to write the `Datafile` before running `data publish`.

### Dataset dependencies

It is possible to specify dataset *dependencies* using a `Datafile`. These tell `data` to download a set of datasets. For example, given:

```
% cat Datafile
dependencies:
- jbenet/mnist@1.0
- jbenet/cifar-10
- jbenet/cifar-100
```

Running `data get` in the directory -- with no arguments -- will download all these datasets:

```
% data get
...
---
Installed jbenet/mnist@1.0 at datasets/jbenet/mnist@1.0
Installed jbenet/cifar-10@1.0 at datasets/jbenet/cifar-10@1.0
Installed jbenet/cifar-100@1.0 at datasets/jbenet/cifar-100@1.0
```

This is useful in two ways:
1. It allows expressing dependency between dataset packages (e.g. A depends on B).
2. It allows easy data distribution with source code: include a `Datafile` in your project, and your users can just run `data get` to install all the data they need. (This works much like `npm's package.json` or `pip's requirements.txt`.)
