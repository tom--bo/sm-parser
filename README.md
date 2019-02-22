# sm-parser

sm-parser(Sysbench-MySQL output parser) just extract items from an output of sysbench for MySQL.
This package assume to be used as small library, but can be used as command line tool to extract output.

## How to use

```
go run parse.go -f output.txt
```

or 

```
cat output.txt | go run parse.go
```


## What items are extracted

sm-parser basically extract all information from output of sysbench for MySQL.  

- version of sysbench
- version of LuaJIT
- Number of threads
- Number of total read/write/other/Tx
- TPS/QPS
- Ignored errors
- Reconnects
- Total Time/Events
- MIN/AVG/MAX/95th/SUM Latency
- ThreadsEvent AVG/STDDEV
- ThreadsExecTime AVG/STDDEV

You can see these fields in Result struct.


