# bigmail
A Golang project demonstrating concurrent programming by writing a simple mass email sender.

To use:

* Compile the Cli by running `go get github.com/dndx/bigmail` and `go build -o bigmail github.com/dndx/bigmail/main`. This will
generate `bigmail` at your current working directory.
* `./bigmail` to see usages

bigmail contains a sender library and could be integrated into existing project easily.

# Benchmark
`bigmail` sends email very efficiently thanks to the concurrent model offered by Golang.

# To Do
* Better documentation

# License
GPL v2, see LICENSE for more details.
