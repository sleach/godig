package main

var BoolOptions = map[string] bool {
    "debug": false,
    "tcp": false,
    "ignore": false,
    "aaflag": false,
    "adflag": false,
    "cdflag": false,
    "rec": true,
    "cl": true,
    "ttlid": true,
    "recurse": false,
    "nssearch": false,
    "trace": false,
    "cmd": true,
    "short": false,
    "identify": false,
    "comments": true,
    "rrcomments": true,
    "stats": true,
    "qr": false,
    "question": true,
    "answer": true,
    "authority": true,
    "additional": true,
    "all": false,
    "multiline": false,
    "onesoa": false,
    "fail": true,
    "besteffort": false,
    "dnssec": false,
    "sigchase": false,
    "topdown": false,
    "nsid": false,
}

var IntOptions = map[string] int {
    "time": 5,
    "tries": 3,
    "retry": 2,
    "ndots": 1,
    "bufsize": 0,
    "edns": -1,
}

var TrustedKey string = "trusted-key.key"
