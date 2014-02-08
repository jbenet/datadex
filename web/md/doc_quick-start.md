## Quick Start

### Download and Install

Download and install [the latest version of data](/doc/install).
<br />

### Download a Dataset

For this tutorial, we'll be using the [Zipcodes Example](http://datadex.io/jbenet/zipcodes-example) dataset.

![](http://jbenet.static.s3.amazonaws.com/d3a80c0b3a1c8dcc9088e9a4e0097b1f548784f6/example-zipcodes-id.png)

Download the latest published version of the dataset:

```
$ data get jbenet/zipcodes-example
```

Download a specific version of the dataset:

```
$ data get jbenet/zipcodes-example@1.0
```

[Learn More: Downloading](TODO)

### Publish a Dataset

#### Step 1: Set up Username

First you need to tell `data` your name, so that it can properly attribute the datasets you publish.

```
$ data user add <Your Username Here>
```

#### Step 2: Publish!

`data` walks you through the steps required to publish a dataset.

Run the following code:
```
# in the directory you want to publish
$ data publish
```

[Learn More: Publishing](TODO)

### Manage Dependencies with a Datafile

Specify your dependencies in a `Datafile` in your project's root:
```
% cat Datafile
dependencies:
- jbenet/mnist@1.0
- jbenet/cifar-10
- jbenet/cifar-100
```

To download all of the datasets from your specified sources, run the following code:
```
% data get
...
---
Installed jbenet/mnist@1.0 at datasets/jbenet/mnist@1.0
Installed jbenet/cifar-10@1.0 at datasets/jbenet/cifar-10@1.0
Installed jbenet/cifar-100@1.0 at datasets/jbenet/cifar-100@1.0
```
[Learn More: Datafiles](TODO)

