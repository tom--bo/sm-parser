package smparser

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
	SysbenchVersion       string
	LuajitVersion         string
	Threads               int
	TotalRead             int
	TotalWrite            int
	TotalOther            int
	TotalTx               int
	Tps                   float64
	TotalQuery            int
	Qps                   float64
	IgnoredErrors         int
	Reconnects            int
	TotalTime             float64
	TotalEvents           int
	MinLatency            float64
	AvgLatency            float64
	MaxLatency            float64
	P95thLatency          float64
	SumLatency            float64
	ThreadsEventsAvg      float64
	ThreadsEventsStddev   float64
	ThreadsExecTimeAvg    float64
	ThreadsExecTimeStddev float64
}

func (r *Result) toCSVString() string {
	return fmt.Sprintf("%s, %s, %d, %d, %d, %d, %d, %.3f, %d, %.3f, %d, %d, %.3f, %d, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f",
		r.SysbenchVersion,
		r.LuajitVersion,
		r.Threads,
		r.TotalRead,
		r.TotalWrite,
		r.TotalOther,
		r.TotalTx,
		r.Tps,
		r.TotalQuery,
		r.Qps,
		r.IgnoredErrors,
		r.Reconnects,
		r.TotalTime,
		r.TotalEvents,
		r.MinLatency,
		r.AvgLatency,
		r.MaxLatency,
		r.P95thLatency,
		r.SumLatency,
		r.ThreadsEventsAvg,
		r.ThreadsEventsStddev,
		r.ThreadsExecTimeAvg,
		r.ThreadsExecTimeStddev,
	)
}

func (r *Result) toString() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d %.3f %d %.3f %d %d %.3f %d %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f",
		r.SysbenchVersion,
		r.LuajitVersion,
		r.Threads,
		r.TotalRead,
		r.TotalWrite,
		r.TotalOther,
		r.TotalTx,
		r.Tps,
		r.TotalQuery,
		r.Qps,
		r.IgnoredErrors,
		r.Reconnects,
		r.TotalTime,
		r.TotalEvents,
		r.MinLatency,
		r.AvgLatency,
		r.MaxLatency,
		r.P95thLatency,
		r.SumLatency,
		r.ThreadsEventsAvg,
		r.ThreadsEventsStddev,
		r.ThreadsExecTimeAvg,
		r.ThreadsExecTimeStddev,
	)
}

func parseRow(r *Result, row string) {
	re := regexp.MustCompile("[ \t]+")
	row = re.ReplaceAllLiteralString(row, " ")
	if strings.Index(row, "sysbench") != -1 {
		str := strings.Split(row, " ")
		s := strings.Replace(str[5], ")", "", -1)
		r.SysbenchVersion = str[1]
		r.LuajitVersion = s
	} else if strings.Index(row, "Number of threads") != -1 {
		str := strings.Split(row, " ")
		r.Threads, _ = strconv.Atoi(str[3])
	} else if strings.Index(row, "read:") != -1 {
		str := strings.Split(row, " ")
		r.TotalRead, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "write:") != -1 {
		str := strings.Split(row, " ")
		r.TotalWrite, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "other:") != -1 {
		str := strings.Split(row, " ")
		r.TotalOther, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "total:") != -1 {
		str := strings.Split(row, " ")
		r.TotalTx, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "transactions:") != -1 {
		str := strings.Split(row, " ")
		r.TotalTx, _ = strconv.Atoi(str[2])
		r.Tps, _ = strconv.ParseFloat(strings.Replace(str[3], "(", "", -1), 64)
	} else if strings.Index(row, "queries:") != -1 {
		str := strings.Split(row, " ")
		r.TotalQuery, _ = strconv.Atoi(str[2])
		r.Qps, _ = strconv.ParseFloat(strings.Replace(str[3], "(", "", -1), 64)
	} else if strings.Index(row, "errors:") != -1 {
		str := strings.Split(row, " ")
		r.IgnoredErrors, _ = strconv.Atoi(str[3])
	} else if strings.Index(row, "reconnects:") != -1 {
		str := strings.Split(row, " ")
		r.Reconnects, _ = strconv.Atoi(str[2])
	} else if strings.Index(row, "time:") != -1 {
		str := strings.Split(row, " ")
		r.TotalTime, _ = strconv.ParseFloat(strings.Replace(str[3], "s", "", -1), 64)
	} else if strings.Index(row, "events:") != -1 {
		str := strings.Split(row, " ")
		r.TotalEvents, _ = strconv.Atoi(str[5])
	} else if strings.Index(row, "min:") != -1 {
		str := strings.Split(row, " ")
		r.MinLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "avg:") != -1 {
		str := strings.Split(row, " ")
		r.AvgLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "max:") != -1 {
		str := strings.Split(row, " ")
		r.MaxLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "percentile:") != -1 {
		str := strings.Split(row, " ")
		r.P95thLatency, _ = strconv.ParseFloat(str[3], 64)
	} else if strings.Index(row, "sum:") != -1 {
		str := strings.Split(row, " ")
		r.SumLatency, _ = strconv.ParseFloat(str[2], 64)
	} else if strings.Index(row, "events (avg/stddev):") != -1 {
		str := strings.Split(row, " ")
		d := strings.Split(str[3], "/")
		r.ThreadsEventsAvg, _ = strconv.ParseFloat(d[0], 64)
		r.ThreadsEventsStddev, _ = strconv.ParseFloat(d[1], 64)
	} else if strings.Index(row, "time (avg/stddev):") != -1 {
		str := strings.Split(row, " ")
		d := strings.Split(str[4], "/")
		r.ThreadsExecTimeAvg, _ = strconv.ParseFloat(d[0], 64)
		r.ThreadsExecTimeStddev, _ = strconv.ParseFloat(d[1], 64)
	}
}

func ParseFile(r *Result, f string) error {
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

func ParseOutput(r *Result, s string) {
	rows := strings.Split(s, "\n")
	for i := 0; i < len(rows); i++ {
		parseRow(r, rows[i])
	}
}

func main() {
	var r Result
	filename := ""
	flag.BoolVar(&isCSV, "c", false, "csv?")
	flag.StringVar(&filename, "f", "", "read from file")
	flag.Parse()

	if filename != "" {
		err := ParseFile(&r, filename)
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
