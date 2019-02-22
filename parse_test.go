package main

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"testing"
)

var (
	expectedSample1 = Result{
		SysbenchVersion:       "1.0.9",
		LuajitVersion:         "2.0.4",
		Threads:               16,
		TotalRead:             24833060,
		TotalWrite:            7095097,
		TotalOther:            3547560,
		TotalTx:               1773770,
		Tps:                   5912.46,
		TotalQuery:            35475717,
		Qps:                   118250.34,
		IgnoredErrors:         20,
		Reconnects:            0,
		TotalTime:             300.0032,
		TotalEvents:           1773770,
		MinLatency:            2.24,
		AvgLatency:            2.7,
		MaxLatency:            79.97,
		P95thLatency:          3.25,
		SumLatency:            4.79663538e+06,
		ThreadsEventsAvg:      110860.625,
		ThreadsEventsStddev:   959.19,
		ThreadsExecTimeAvg:    299.7897,
		ThreadsExecTimeStddev: 0,
	}
)

func TestParseRows(t *testing.T) {
	var r Result
	content, err := ioutil.ReadFile("./output_samples/sample1.txt")
	if err != nil {
		log.Fatal(err)
	}
	ParseRows(&r, string(content))

	if diff := cmp.Diff(r, expectedSample1); diff != "" {
		t.Errorf("Result struct differs: (-got +want)\n%s", diff)
	}
}

func TestParseFile(t *testing.T) {
	var r Result
	ParseFile(&r, "output_samples/sample1.txt")

	if diff := cmp.Diff(r, expectedSample1); diff != "" {
		t.Errorf("Result struct differs: (-got +want)\n%s", diff)
	}
}

func TestOutputs(t *testing.T) {
	expected := "1.0.9 2.0.4 16 24833060 7095097 3547560 1773770 5912.460 35475717 118250.340 20 0 300.003 1773770 2.240 2.700 79.970 3.250 4796635.380 110860.625 959.190 299.790 0.000"
	var r Result
	ParseFile(&r, "output_samples/sample1.txt")
	str := r.toString()
	if str != expected {
		t.Errorf("output is different from expected.\n Expected: %s\n      Got: %s", expected, str)
	}
}

func TestCSVOutputs(t *testing.T) {
	expected := "1.0.9, 2.0.4, 16, 24833060, 7095097, 3547560, 1773770, 5912.460, 35475717, 118250.340, 20, 0, 300.003, 1773770, 2.240, 2.700, 79.970, 3.250, 4796635.380, 110860.625, 959.190, 299.790, 0.000"
	var r Result
	ParseFile(&r, "output_samples/sample1.txt")
	str := r.toCSVString()
	if str != expected {
		t.Errorf("output is different from expected.\n Expected: %s\n      Got: %s", expected, str)
	}
}
