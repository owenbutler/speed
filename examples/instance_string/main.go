package main

import (
	"flag"
	"time"

	"github.com/performancecopilot/speed"
)

var timelimit = flag.Int("time", 60, "number of seconds to run for")

func main() {
	flag.Parse()

	w, err := speed.NewPCPWriter("strings", speed.ProcessFlag)
	if err != nil {
		panic(err)
	}

	m, err := w.RegisterString("language[go, javascript, php].users", speed.Instances{
		"go":         1,
		"javascript": 100,
		"php":        10,
	}, speed.CounterSemantics, speed.Uint64Type, speed.OneUnit)
	if err != nil {
		panic(err)
	}

	err = w.Start()
	if err != nil {
		panic(err)
	}

	metric := m.(speed.InstanceMetric)
	for i := 0; i < *timelimit; i++ {
		v, _ := metric.ValInstance("go")
		metric.SetInstance("go", v.(uint64)*2)

		v, _ = metric.ValInstance("javascript")
		metric.SetInstance("javascript", v.(uint64)+10)

		v, _ = metric.ValInstance("php")
		metric.SetInstance("php", v.(uint64)+1)

		time.Sleep(time.Second)
	}

	err = w.Stop()
	if err != nil {
		panic(err)
	}
}
