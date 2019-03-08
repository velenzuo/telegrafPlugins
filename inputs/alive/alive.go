package alive

import (
    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
    "time"
)

type AliveStats struct{}

func (_ *AliveStats) Description() string {
    return "Read metrics about network interface usage"
}

var netSampleConfig = ` `

func (_ *AliveStats) SampleConfig() string {
    return netSampleConfig
}

func (s *AliveStats) Gather(acc telegraf.Accumulator) error {
    tags := map[string]string{}
    currentTime:=time.Now()
    fields := map[string]interface{}{
        "alive": 1,
        "time": currentTime.Unix(),
    }

    acc.AddCounter("agent", fields,tags)

    return nil
}

func init() {
    inputs.Add("alive", func() telegraf.Input {
        return &AliveStats{}
    })
}

