<!-- title: F.A.Q. -->
<!-- description: F.A.Q about Data & Datadex -->

# Frequently Asked Questions


### I have another question, not listed here. How should I ask it?

Great! Ask us through any of these channels:

<i class="icon-twitter"></i> [@juanbenet](https://twitter.com/intent/user?screen_name=juanbenet)
<i class="icon-github"></i> [jbenet/datadex](https://github.com/jbenet/datadex/issues/new)
<i class="icon-envelope"></i> [data-discuss@googlegroups.com](https://groups.google.com/forum/#!forum/data-discuss)


### Where is my dataset stored?

It is stored on our S3 bucket. Later on, it will be replicated to other hosts (in a bit-torrent fashion), but that's not yet implemented.

### I found a bug, what do I do?

Ah, please report it here:

- [Click here to report `data` (commandline tool) issues](https://github.com/jbenet/data/issues/new).
- [Click here to report `datadex` (website) issues](https://github.com/jbenet/datadex/issues/new).

We'll try to get to it as soon as possible. Thanks!

### Do I have to leave my server on at all times?

No. For regular usage, you only need the `data` tool, not `datadex`.
If you want to run *your own* index (your own `datadex` website with your own private packages), then yes.

### How can I run my own `datadex`?

Yes. Use the `datadex` server at https://github.com/jbenet/datadex/
We will be writing a guide on how to setup your own index, but that won't come for a while. [Make an issue](https://github.com/jbenet/datadex/issues/new) if you want it sooner.

### Does it have to be public?

No, you can put your `datadex` instance behind any HTTP authentication.
Just make sure clients can connect to it.

### Can I "un-publish" a dataset?

At the moment no. But this is planned. See
https://github.com/jbenet/datadex/issues/28

### I downloaded and edited a dataset. How do I give the original author credit? Am I now the "author"?

You're not the new author. Authorship constitutes the main creators of a dataset. Add yourself to the "contributors" list. You will also be the "maintainer" of datasets published under your account. We know it is a bit unclear what makes up an "Author" vs a "Contributor", but we're hoping to define a good guideline for this ([see the conversation here](https://github.com/jbenet/datadex/issue/29)). Perhaps ask the original authors what they think?

### How do I know what license my dataset has?

If you're creating a dataset from scratch, you can pick whatever you want (see below).

If you're modifying another dataset, we recommend keeping the same license. Please make sure modification and re-publishing is allowed by the original license. This is usually the case, but someone may publish datasets with non-permissive licenses that shouldn't be violated.

It's best to check. Consider:
- asking the original author
- asking [our mailing list](https://groups.google.com/forum/#!forum/data-discuss)
- asking [on github](https://github.com/jbenet/datadex/issues/new)

### I'm publishing a Dataset. What license should I pick?

We recommend using one of the popular open-source licenses. These websites can help you find out more and pick the right license for you:
- http://creativecommons.org/choose/
- http://choosealicense.com/
- http://opendatacommons.org/
- http://www.dcc.ac.uk/sites/default/files/documents/how-to-license-data.pdf


### I made a subset of a dataset. Can I publish that and who is the author?

It depends on the original dataset's license. Some licenses are very permissive. Some are not. If you're not sure, we recommend
- asking the original author
- asking [our mailing list](https://groups.google.com/forum/#!forum/data-discuss)

We hope to iron out the licensing concerns with a straight-foward guide. [See the conversation here](https://github.com/jbenet/datadex/issues/30).
