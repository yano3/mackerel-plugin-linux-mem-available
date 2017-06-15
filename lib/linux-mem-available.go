package mplinuxmemavailable

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"regexp"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// MemAvailablePlugin mackerel plugin for meminfo
type MemAvailablePlugin struct {
	Tempfile string
}

var memReg = regexp.MustCompile(`^([A-Za-z]+):\s+([0-9]+)\s+kB`)

var memItems = map[string]string{
	"MemTotal":     "total",
	"MemAvailable": "available",
}

// FetchMetrics interface for mackerelplugin
func (p MemAvailablePlugin) FetchMetrics() (map[string]interface{}, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := memReg.FindStringSubmatch(line); len(matches) == 3 {
			k, ok := memItems[matches[1]]
			if !ok {
				continue
			}
			value, _ := strconv.ParseUint(matches[2], 10, 64)
			ret[k] = value * 1024
		}
	}
	return ret, nil
}

// GraphDefinition interface for mackerelplugin
func (p MemAvailablePlugin) GraphDefinition() map[string]mp.Graphs {
	var graphdef = map[string]mp.Graphs{
		"linux-mem-available.memory": {
			Label: "Linux Available Memory",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "total", Label: "Total Memory", Diff: false, Type: "uint64", Stacked: false},
				{Name: "available", Label: "Available Memory", Diff: false, Type: "uint64", Stacked: false},
			},
		},
	}
	return graphdef
}

// Do the plugin
func Do() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var linuxMemAvailable MemAvailablePlugin

	helper := mp.NewMackerelPlugin(linuxMemAvailable)
	helper.Tempfile = *optTempfile
	helper.Run()
}
