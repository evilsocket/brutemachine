# BruteMachine

This is a Go library which main purpose is giving an interface to loop over a dictionary and use those words/lines as input for some 
custom logic such as HTTP file bruteforcing, DNS bruteforcing, etc.

[![baby-gopher](https://raw.githubusercontent.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org) 

## Example

The following is an example of how to use `brutemachine` to perform HTTP files bruteforcing.


    package main

    import (
        "fmt"
        "net/http"
        "strings"

        "github.com/evilsocket/brutemachine"
    )

    const base = "http://nas:8181/"

    func DoRequest(page string) interface{} {
        url := strings.Replace(fmt.Sprintf("%s%s", base, page), "%EXT%", "php", -1)
        resp, err := http.Head(url)
        // Only pass valid responses to the handler.
        if err == nil && resp.StatusCode == 200 {
            return url
        }

        return nil
    }

    func OnResult(res interface{}) {
        fmt.Printf("@ Found '%s'\n", res)
    }

    func main() {
        m := brutemachine.New( -1, "dictionary.txt", DoRequest, OnResult)
        if err := m.Start(); err != nil {
            panic(err)
        }

        m.Wait()

        fmt.Printf("\nDONE:\n")
        fmt.Printf("%+v\n", m.Stats)
    }

## Installation

    go get github.com/evilsocket/brutemachine

## License

This project is copyleft of [Simone Margaritelli](http://www.evilsocket.net/) and released under the GPL 3 license.

