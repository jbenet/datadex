<!-- title: Documentation -->
<!-- description: Documentation for Data & Datadex -->

# Documentation

<div class="alert alert-success">
<i class="icon-flag"></i>
We recommend you first see: <a href="/doc/quick-start">Quick Start</a>.
</div>

### Downloading a Dataset

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

To upload and publish a dataset, follow these steps:

0. Make sure you're logged in. Either:
        # make an account, if it's your first time:
        data user add

        # or, authenticate if you already have one:
        data user auth

1. Create a directory and place all the files you want to publish there.
        mkdir zipcodes-example
        cd zipcodes-example
        cp ~/zipcodes_db.txt zipcodes_states.txt

2. Run `data publish`. It will ask you to fill out `Datafile` information, to define your package and make it searchable in the index. These fields are mostly optional, and include things like:
  - `owner id`: defaults to your user id.
  - `dataset id`: to be used in the dataset path
  - `version`: the version of this data
  - `tagline`: a title or one-line description
  - `description`: a longer (optional) description.

    Note: the current default value is displayed in [brackets]. If you enter nothing, that value will be used. See more about [`Datafiles`](/doc#Datafile) below.

    Once done, data will upload your dataset files (to our S3 bucket).
    This may take several minutes if the datasets are large (100MB+).

Below is the output of an example run of `data publish`:

```
% data publish
==> Guided Data Package Publishing.

Welcome to Data Package Publishing. You should read these short
messages carefully, as they contain important information about
how data works, and how your data package will be published.

First, a 'data package' is a collection of files, containing:
- various files with your data, in any format.
- 'Datafile', a file with descriptive information about the package.
- 'Manifest', a file listing the other files in the package and their checksums.

This tool will automatically:
1. Create the package
  - Generate a 'Datafile', with information you will provide.
  - Generate a 'Manifest', with all the files in the current directory.
2. Upload the package contents
3. Publish the package to the index

(Note: to specify which files are part of the package, and other advanced
 features, use the 'data pack' command directly. See 'data pack help'.)


==> Step 1/3: Creating the package.

First, let's write the package's Datafile, which contains important
information about the package. The 'owner id' is the username of the
package's owner (usually your username). The 'dataset id' is the identifier
which defines this dataset. Good 'dataset ids' are like names: short, unique,
and memorable. For example: "mnist" or "cifar". Choose it carefully.

Writing Datafile fields...
'Field description [current value]'
Enter owner id (required) [jbenet]:
Enter dataset id (required) [zipcodes-example]:
Enter dataset version (required) [1.0]:
Enter tagline description (required) [Example dataset with zipcodes]:
Enter long description (optional) []:
Enter license name (optional) []: MIT
Generating Manifest file...
data manifest: added Datafile
data manifest: added README.md
data manifest: added zipcode_cities.txt
data manifest: added zipcode_states.txt
data manifest: hashed 10dcbf5 Datafile
data manifest: hashed b6d9e12 README.md
data manifest: hashed 4b1da58 zipcode_cities.txt
data manifest: hashed b0f41fc zipcode_states.txt
4 files in Manifest.

==> Step 2/3: Uploading the package contents.

Now, data will upload the contents of the package (this directory) to the index
sotrage service. This may take a while, if the files are large (over 100MB).

put blob 10dcbf5 Datafile - uploading
put blob b6d9e12 README.md - uploading
put blob 4b1da58 zipcode_cities.txt - uploading
put blob b0f41fc zipcode_states.txt - uploading
put blob 8001ee9 .data/Manifest - uploading

==> Step 3/3: Publishing the package to the index.

Finally, data will publish the package to the index, where others can find
and download your package. The index is available through data, and on the web.

data pack: published jbenet/zipcodes-example@1.0 (8001ee9).
Webpage at http://datadex.io/jbenet/zipcodes-example@1.0
```

<br />

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
