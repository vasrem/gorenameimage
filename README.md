# gorenameimage
A cli tool to rename a bunch of images according to the date of their creation.

## Installation
Simply run the following commands:

* `go get github.com/vasrem/gorenameimages`
* `cd $GOPATH/src/github.com/vasrem/gorenameimages`
* `go install`

## Using it

You just have to put the absolute paths of the folders as flags like the command below:

`$ gorenameimages --input "/path/to/input/folder" --output "/path/to/output/folder"`

There are 2 modes. `--mode=copy` and `--mode=move`. Default is `copy`.
You can also define the prefix by adding the flag `--prefix=my_prefix`.

It only supports `.jpg` pictures.

## Why to use this?

If you are bored of watching pictures which you have taken in your trips with the wrong order, just use this tool to rename them all according to the date created.
