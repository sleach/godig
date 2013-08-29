package main

import (
    "errors"
    "reflect"
    "fmt"
    "strings"
    "github.com/miekg/unbound"
)

func CallOutputFunc(name string, params ... interface{}) (result []reflect.Value, err error) {
    f := reflect.ValueOf(OutputMap[name])
    if len(params) != f.Type().NumIn() {
        err = errors.New("The number of params is not adapted.")
        return
    }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    result = f.Call(in)
    return
}

var OutputMap = map[string] interface{} {
    "text": OutputText,
    "txt": OutputText,
    "json": OutputJSON,
    "xml": OutputXML,
}

func OutputJSON(resp *unbound.Result) {
}

func OutputText(resp *unbound.Result) {
    fmt.Printf("\n; <<>> GoDig v0.1 <<>> %s %s %s\n", query["qname"], query["qtype"], query["qclass"])
    fmt.Printf(";; global options: TODO\n")
    fmt.Printf(";; Got answer:\n")
    fmt.Printf(";; ->>HEADER<<- opcode: TODO, status: TODO, id: TODO\n")
    fmt.Printf(";; flags: TODO; QUERY: TODO, ANSWER: TODO, AUTHORITY: TODO, ADDITIONAL: TODO\n\n")
    fmt.Printf(";; QUESTION SECTION: \n")
    fmt.Printf(";%-40s%-10s%-10s\n\n", strings.ToUpper(query["qname"]), strings.ToUpper(query["qclass"]), strings.ToUpper(query["qtype"]))
    if !resp.HaveData {
        fmt.Printf("Got no data\n")
    } else {
        fmt.Printf("; Answer:\n")
        for _, res := range resp.AnswerPacket.Answer {
            fmt.Printf("%s\n", res)
        }
        fmt.Printf("\n\nNs:\n")
        for _, res := range resp.AnswerPacket.Ns {
            fmt.Printf("%s\n", res)
        }
        fmt.Printf("\n\nExtra:\n")
        for _, res := range resp.AnswerPacket.Extra {
            fmt.Printf("%s\n", res)
        }
    }

    fmt.Printf(";; Query time: TODO\n")
    fmt.Printf(";; SERVER: TODO#53(TODO)\n")
    fmt.Printf(";; WHEN: TODO\n")
    fmt.Printf(";; MSG SIZE  rcvd: TODO\n\n")
}

func OutputXML(resp *unbound.Result) {
}

