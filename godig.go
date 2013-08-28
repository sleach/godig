package main

import (
    "log"
	"fmt"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"os"
	"strconv"
	"strings"
)

var query = map[string]string{
	"server": "",
	"qname":  ".",
	"qtype":  "ns",
	"qclass": "in",
}

var option_output string = "text"

func invalidArg(message string) {
	log.Fatalf("ERROR: Invalid argument %s\n", message)
}

func HandleArgs(args []string) {
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "+") {
			set := true
			option := arg[1:]
			if strings.HasPrefix(option, "no") {
				option = option[2:]
				set = false
			}
			if _, ok := BoolOptions[option]; ok {
				BoolOptions[option] = set
			} else if strings.Count(option, "=") >= 0 {
				items := strings.Split(option, "=")
				if _, ok := IntOptions[items[0]]; ok {
					int_val, err := strconv.Atoi(items[1])
					if err != nil {
						invalidArg(fmt.Sprintf("%s %s", items[0], items[1]))
					} else {
						IntOptions[items[0]] = int_val
					}
				} else if items[0] == "trusted-key" {
					TrustedKey = items[1]
				} else {
					invalidArg(option)
				}
			} else {
				invalidArg(option)
			}
		} else if strings.HasPrefix(arg, "@") {
			query["server"] = arg[1:]
		} else if _, ok := dns.StringToType[strings.ToUpper(arg)]; ok {
			query["qtype"] = strings.ToUpper(arg)
		} else if _, ok := dns.StringToClass[strings.ToUpper(arg)]; ok {
			query["qclass"] = strings.ToUpper(arg)
		} else {
			query["qname"] = arg
		}
	}
}

func main() {
	args := os.Args
	HandleArgs(args)
    u := unbound.New()
    defer u.Destroy()
    if query["server"] != "" {
        fmt.Printf("Will use nameserver %s\n", query["server"])
        u.SetFwd(query["server"])
    } else {
        u.ResolvConf("/etc/resolv.conf")
    }
    if option_output == "text" {
        fmt.Printf("\n; <<>> GoDig v0.1 <<>> %s %s %s\n", query["qname"], query["qtype"], query["qclass"])
        fmt.Printf(";; global options: TODO\n")
        fmt.Printf(";; Got answer:\n")
        fmt.Printf(";; ->>HEADER<<- opcode: TODO, status: TODO, id: TODO\n")
        fmt.Printf(";; flags: TODO; QUERY: TODO, ANSWER: TODO, AUTHORITY: TODO, ADDITIONAL: TODO\n\n")
        fmt.Printf(";; QUESTION SECTION: \n")
        fmt.Printf(";%-40s%-10s%-10s\n\n", strings.ToUpper(query["qname"]), strings.ToUpper(query["qclass"]), strings.ToUpper(query["qtype"]))
    }
    resp, err := u.Resolve(query["qname"], dns.StringToType[query["qtype"]], dns.StringToClass[query["qclass"]])
    if err != nil {
        log.Fatalf("ERROR: query failed: %s\n", err) 
    } else {
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
    }
    if option_output == "text" {
        fmt.Printf(";; Query time: TODO\n")
        fmt.Printf(";; SERVER: TODO#53(TODO)\n")
        fmt.Printf(";; WHEN: TODO\n")
        fmt.Printf(";; MSG SIZE  rcvd: TODO\n\n")
    }
}
