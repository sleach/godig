package main

import (
    "log"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"strconv"
	"strings"
    "net"
)

var query = map[string]string{
	"server": "",
	"qname":  "",
	"qtype":  "",
	"qclass": "IN",
}

func invalidArg(message string) {
	log.Fatalf("ERROR: Invalid argument %s\n", message)
}

func doQuery(nameserver, qname, qtype, qclass string) {
    if len(nameserver) == 0 {
        conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(2)
        }
        nameserver = conf.Servers[0]
    }

    // if the nameserver is from /etc/resolv.conf the [ and ] are already
    // added, thereby breaking net.ParseIP. Check for this and don't
    // fully qualify such a name
    if nameserver[0] == '[' && nameserver[len(nameserver)-1] == ']' {
        nameserver = nameserver[1 : len(nameserver)-1]
    }
    if i := net.ParseIP(nameserver); i != nil {
        nameserver = net.JoinHostPort(nameserver, strconv.Itoa(IntOptions["port"]))
    } else {
        nameserver = dns.Fqdn(nameserver) + ":" + strconv.Itoa(IntOptions["port"])
    }

    client := new(dns.Client)
    if BoolOptions["tcp"] {
        client.Net = "tcp"
    } else { 
        client.Net = "udp"
    }
    msg := new(dns.Msg)
    msg.MsgHdr.Authoritative = BoolOptions["aaflag"]
    msg.MsgHdr.AuthenticatedData = BoolOptions["adflag"]
    msg.MsgHdr.CheckingDisabled = BoolOptions["cdflag"]
    msg.MsgHdr.RecursionDesired = BoolOptions["rec"]
    msg.Question = make([]dns.Question, 1)
    if BoolOptions["dnssec"] || BoolOptions["nsid"] || StringOptions["client"] != "" {
        o := new(dns.OPT)
        o.Hdr.Name = "."
        o.Hdr.Rrtype = dns.TypeOPT
        if BoolOptions["dnssec"] {
            o.SetDo()
            o.SetUDPSize(dns.DefaultMsgSize)
        }
        if BoolOptions["nsid"] {
            e := new(dns.EDNS0_NSID)
            e.Code = dns.EDNS0NSID
            o.Option = append(o.Option, e)
            // NSD will not return nsid when the udp message size is too small
            o.SetUDPSize(dns.DefaultMsgSize)
        }
        if StringOptions["client"] != "" {
            e := new(dns.EDNS0_SUBNET)
            e.Code = dns.EDNS0SUBNET
            e.SourceScope = 0
            e.Address = net.ParseIP(StringOptions["client"])
            if e.Address == nil {
                fmt.Fprintf(os.Stderr, "Failure to parse IP address: %s\n", StringOptions["client"])
                return
            }
            e.Family = 1 // IP4
            e.SourceNetmask = net.IPv4len * 8
            if e.Address.To4() == nil {
                e.Family = 2 // IP6
                e.SourceNetmask = net.IPv6len * 8
            }
            o.Option = append(o.Option, e)
        }
        msg.Extra = append(msg.Extra, o)
    }
    msg.Question[0] = dns.Question{dns.Fqdn(qname), dns.StringToType[qtype], dns.StringToClass[qclass]}
    msg.Id = dns.Id()
    resp, rtt, err := client.Exchange(msg, nameserver)
    if err != nil {
        fmt.Fprintf(os.Stderr, ";; %s\n", err.Error())
        return
    }
    if resp.Id != msg.Id {
        fmt.Fprintf(os.Stderr, "Id Mismatch\n")
        return
    }
    fmt.Printf("%v", resp)
    fmt.Printf("\n;; query time: %3d ms, server: %s(%s), size: %d bytes\n", rtt/1000000, nameserver, client.Net, resp.Len())
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
				} else if _, ok := StringOptions[items[0]]; ok {
                    StringOptions[items[0]] = items[1]
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
    if query["qname"] != "" {
        if query["qtype"] == "" {
            query["qtype"] = "A"
        }
    } else {
        query["qname"] = "."
        if query["qtype"] == "" {
            query["qtype"] = "NS"
        }
    }
}

func main() {
	args := os.Args
	HandleArgs(args)
    doQuery(query["server"], query["qname"], query["qtype"], query["qclass"])
    //result, err := CallOutputFunc(StringOptions["output"], resp)
    //fmt.Printf("Result: %s, err: %s\n", result, err)
}
