## Quick Start

### Download and Install Data

Download and install [the latest version of Data](/doc/install).
<br />

### Download a Dataset

For this tutorial, we'll be using the [Zipcodes Example](http://datadex.io/jbenet/zipcodes-example) dataset.

![](http://jbenet.static.s3.amazonaws.com/d3a80c0b3a1c8dcc9088e9a4e0097b1f548784f6/example-zipcodes-id.png)

Download the latest published version:

```
data get jbenet/zipcodes-example
```

Download a specific version:

```
data get jbenet/zipcodes-example@1.0
```

### Publish a Dataset

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
