# Description

Command line tool which accepts uris and prints request address along with MD5 hash of the response.

# How to run

1. Clone the repository
2. Run `go build myhttp.go`
3. Run `./myhttp -parallel <no. of workers> <uri1> <uri2> ...`
4. Run `go test` to run the unit tests.

Examples:
`./myhttp -parallel 2 https://facebook.com https://google.com`
`./myhttp https://facebook.com`