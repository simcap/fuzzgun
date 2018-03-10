[![Build Status](https://api.travis-ci.org/simcap/fuzzgun.svg?branch=master)](https://travis-ci.org/simcap/fuzzgun)
[![Go Report Card](https://goreportcard.com/badge/github.com/simcap/fuzzgun)](https://goreportcard.com/report/github.com/simcap/fuzzgun)
[![GoDoc](https://godoc.org/github.com/simcap/fuzzgun?status.svg)](https://godoc.org/github.com/simcap/fuzzgun)

Fuzzgun generates mutated, invalid, random, unexpected data and exotic format from a given input string. 

As a fuzzer Fuzzgun:
* is black box fuzzer (unaware of internal program structure)
* is a mutation based and aware of input structure.
* takes a string layout as its input model

It can be used by developers, security tester and quality assurance teams alike.

## Usage

Fuzzgun only takes a string as input. This string can be anything you want and fuzzgun will try to work its best on it.

On Internet though, programs ingest structured inputs. So to harness a system the best inputs for fuzzgun will be string examples (or layouts) of what your systems is expecting.

### Library usage

For full usage, documentation and examples report to the [Godoc](https://godoc.org/github.com/simcap/fuzzgun)

```go
import ( 
    "http"
    "time"
    "github.com/simcap/fuzzgun"
)

func main() {
    for mutant := range fuzzgun.FuzzEvery("07/08/2018", 3 * time.Second) {
        url := fmt.Sprintf("http://example.com?date=%s", mutant)
        if resp, err := http.Get(url); err != nil {
            panic(err)
        } else if resp.StatusCode == 500 {
            panic("ouch")
        }
    }
}
```

### CLI usage

If you have [Golang](https://golang.org/dl/) (>= 1.10) installed, the following will fetch and install the CLI executable:
```sh
$ go get -u github.com/simcap/fuzzgun
```

Otherwise grab a [binary for Linux, Windows or Mac](https://github.com/simcap/fuzzgun/releases)

Then to get started run:
```sh
$ fuzzgun -h

 # start to mutate some stuff
$ fuzzgun bob@mail.net
$ fuzzgun http://example.com
$ fuzzgun 07/12/2016
```

## How it works

Fuzzgun takes as input a string layout. A **layout** is an string example of a structured input. Here is the basic algorithm (of my own cooking, i.e. feedback welcome) that will be applied to the input:

1. _Tokenizing_ separates the input string into either _alpha, numerical or separator_ tokens
2. _Labelizing_ 

    * known types: marks the input after a successful detection of a known type: URL, IP address, Date, e-mail, etc.
    * known encoding: marks the input after a successful detection of a known encoding: base64, URL encoding, etc.

3. _Grouping_ extracts set of tokens using various stategy: _arrangement_, _shifting_, _separators only_, etc...
4. _Fuzzing_ mutates the data in parallel given the different groups, labels, encoding, etc.
5. _Generating_ will finalizes the fuzzed output putting back groups to original input; encoding the result if the input was detected encoded

### Tokenizing 

The input is tokenized into either _alpha, numerical or separator_. For instance "bob@mail.net" would output: "bob" (alpha), "@" (separator), "mail" (alpha), "." (separator), "net" (alpha)

### Labelizing

Since structured input on the internet can easily have known format, fuzzgun will labelizes the input string according to detected format: _ip address, URL, date, unix timestamp_

This will allows to mutate data according to known issues or valid but exotic formats.

Examples:

* detecting an IP address we can generates output such as: IPv6, IP overflow values, octal/hexadecimal, etc.
* detecting an e-mail address we can generates *valid yet uncommon* e-mail addresses according to [RFC 5322](https://tools.ietf.org/html/rfc5322)

### Grouping 

Grouping allows to isolate array of tokens using various strategy to be fuzzed indenpendently of others.

Basically we extract some tokens to be fuzzed while letting others in their original form. Groups will go through fuzzing and will then be re-arranged with the original string to present the final fuzzed output.

We can think of the original input string as the main group and the grouping step will basically generates subgroups.

For instance given the string "bob@mail.net", some generated group will be:

```
# shifting
["bob"] (group 1), ["@"] (group 2), etc. 
["bob", "@"] (group 1), ["@", "mail"] (group 2), etc.
["bob", "@", "mail"] (group 1), ["@", "mail", ""] (group 2), etc.
...
# separators only
["@", "."]
```

For example given "bob@mail.net", the simplest group after tokenization could be ["bob"]: the group ["bob"] will then be mutated to be then re-arranged to "@mail.net".

## Notes

In future versions mutated values should be fed back as input string in fuzzgun itself!