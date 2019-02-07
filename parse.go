package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	isCSV bool
)

type Result struct {
	sysbenchVersion       string
	luajitVersion         string
	threads               int
	totalRead             int
	totalWrite            int
	totalOther            int
	totalTx               int
	tps                   float64
	totalQuery            int
	qps                   float64
	ignoredErrors         int
	reconnects            int
	totalTime             float64
	totalEvents           int
	minLatency            float64
	avgLatency            float64
	maxLatency            float64
	p95thLatency          float64
	sumLatency            float64
	threadsEventsAvg      float64
	threadsEventsStddev   float64
	threadsExecTimeAvg    float64
	threadsExecTimeStddev float64
}

func (r *Result) toCSVString() string {
	return fmt.Sprintf("%s, %s, %d, %d, %d, %d, %d, %.3f, %d, %.3f, %d, %d, %.3f, %d, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f",
		r.sysbenchVersion,
		r.luajitVersion,
		r.threads,
		r.totalRead,
		r.totalWrite,
		r.totalOther,
		r.totalTx,
		r.tps,
		r.totalQuery,
		r.qps,
		r.ignoredErrors,
		r.reconnects,
		r.totalTime,
		r.totalEvents,
		r.minLatency,
		r.avgLatency,
		r.maxLatency,
		r.p95thLatency,
		r.sumLatency,
		r.threadsEventsAvg,
		r.threadsEventsStddev,
		r.threadsExecTimeAvg,
		r.threadsExecTimeStddev,
	)
}

func (r *Result) toString() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d %.3f %d %.3f %d %d %.3f %d %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f",
		r.sysbenchVersion,
		r.luajitVersion,
		r.threads,
		r.totalRead,
		r.totalWrite,
		r.totalOther,
		r.totalTx,
		r.tps,
		r.totalQuery,
		r.qps,
		r.ignoredErrors,
		r.reconnects,
		r.totalTime,
		r.totalEvents,
		r.minLatency,
		r.avgLatency,
		r.maxLatency,
		r.p95thLatency,
		r.sumLatency,
		r.threadsEventsAvg,
		r.threadsEventsStddev,
		r.threadsExecTimeAvg,
		r.threadsExecTimeStddev,
	)
}

func parseRow(r *Result, row string) {
	re := regexp.MustCompile("[ \t]+")
	row = re.ReplaceAllLiteralString(row, " ")
	if strings.Index(row, "sysbench") != -1 {
		str := strings.Split(row, " ")
		s := strings.Replace(str[5], ")", "", -1)
		r.sysbenchVersion = str[1]
		r.luajitVersion = s
	} else if strings.Index(row, "Number of threads") != -1 {
		str := strings.Split(row, " ")
		r.threads, _ = strconv.Atoi(str[3])
	} else if strings.Index(row, "read:") != -1 {
		str := strings.Split(row, " ")
		r.totalRead, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "write:") != -1 {
		str := strings.Split(row, " ")
		r.totalWrite, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "other:") != -1 {
		str := strings.Split(row, " ")
		r.totalOther, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "total:") != -1 {
		str := strings.Split(row, " ")
		r.totalTx, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "transactions:") != -1 {
		str := strings.Split(row, " ")
		r.totalTx, _ = strconv.Atoi(str[2])
		r.tps, _ = strconv.ParseFloat(strings.Replace(str[3], "(", "", -1), 64)
	} else if strings.Index(row, "queries:") != -1 {
		str := strings.Split(row, " ")
		r.totalQuery, _ = strconv.Atoi(str[2])
		r.qps, _ = strconv.ParseFloat(strings.Replace(str[3], "(", "", -1), 64)
	} else if strings.Index(row, "errors:") != -1 {
		str := strings.Split(row, " ")
		r.ignoredErrors, _ = strconv.Atoi(str[3])
	} else if strings.Index(row, "reconnects:") != -1 {
		str := strings.Split(row, " ")
		r.reconnects, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "time:") != -1 {
		str := strings.Split(row, " ")
		r.totalTime, _ = strconv.ParseFloat(strings.Replace(str[3], "s", "", -1), 64)
	} else if strings.Index(row, "events:") != -1 {
		str := strings.Split(row, " ")
		r.totalEvents, _ = strconv.Atoi(str[5])
	} else if strings.Index(row, "min:") != -1 {
		str := strings.Split(row, " ")
		r.minLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "avg:") != -1 {
		str := strings.Split(row, " ")
		r.avgLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "max:") != -1 {
		str := strings.Split(row, " ")
		r.maxLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "percentile:") != -1 {
		str := strings.Split(row, " ")
		r.p95thLatency, _ = strconv.ParseFloat(str[3], 64)
	} else if strings.Index(row, "sum:") != -1 {
		str := strings.Split(row, " ")
		r.sumLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "events (avg/stddev):") != -1 {
		str := strings.Split(row, " ")
		d := strings.Split(str[3], "/")
		r.threadsEventsAvg, _ = strconv.ParseFloat(d[0], 64)
		r.threadsEventsStddev, _ = strconv.ParseFloat(d[1], 64)
	} else if strings.Index(row, "time (avg/stddev):") != -1 {
		str := strings.Split(row, " ")
		d := strings.Split(str[4], "/")
		r.threadsExecTimeAvg, _ = strconv.ParseFloat(d[0], 64)
		r.threadsExecTimeStddev, _ = strconv.ParseFloat(d[1], 64)
	}
}

func parseFile(r *Result, f string) error {
	fp, err := os.Open(f)
	if err != nil {
		fmt.Println("Error: read file!")
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		parseRow(r, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Error: read file!")
		return err
	}
	return nil
}

func ParseRows(r *Result, f string) {
	var r Result
	rows := strings.Split(f, "\n")
	for i := 0; i < len(rows); i++ {
		parseRow(&r, rows[i])
	}

	return r
}

func main() {
	var r Result
	filename := ""
	flag.BoolVar(&isCSV, "c", false, "csv?")
	flag.StringVar(&filename, "f", "", "read from file")
	flag.Parse()

	if filename != "" {
		err := parseFile(&r, filename)
		if err != nil {
			return
		}
	} else if terminal.IsTerminal(0) {
		fmt.Println("Error: no input from pipe neither not specified any file with -f option")
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error: read input from pipe")
			return
		}
		input := string(b)
		rows := strings.Split(input, "\n")
		for i := 0; i < len(rows); i++ {
			parseRow(&r, rows[i])
		}
	}

	if isCSV {
		fmt.Println(r.toCSVString())
	} else {
		fmt.Println(r.toString())
	}
}
