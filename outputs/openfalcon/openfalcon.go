package openfalcon

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

type Openfalcon struct {
	Files []string

	Addr string


	writers []io.Writer
	closers []io.Closer

	serializer serializers.Serializer
	client  http.Client
}


type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}

type JsonMetaData struct {
	Metric      string      `json:"metric"`
	Endpoint    string      `json:"endpoint"`
	Timestamp   int64       `json:"timestamp"`
	Step        int64       `json:"step"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	Tags        string      `json:"tags"`
}

func (t *JsonMetaData) String() string {
	return fmt.Sprintf("<JsonMetaData Endpoint:%s, Metric:%s, Tags:%s, DsType:%s, Step:%d, Value:%v, Timestamp:%d>",
		t.Endpoint, t.Metric, t.Tags, t.CounterType, t.Step, t.Value, t.Timestamp)
}

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

func (t *MetaData) String() string {
	return fmt.Sprintf("<MetaData Endpoint:%s, Metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v>",
		t.Endpoint, t.Metric, t.Timestamp, t.Step, t.Value, t.Tags)
}

var sampleConfig = `
  addr = "127.0.0.1:8433"
`

func (f *Openfalcon) SetSerializer(serializer serializers.Serializer) {
	f.serializer = serializer
}

func (f *Openfalcon) Connect() error {
	timeout :=  time.Second * 10
	f.client = http.Client{
		Transport: &http.Transport{},
		Timeout: timeout,
	}

	return nil
}

func (f *Openfalcon) Close() error {
	fmt.Println("关闭链接")
	TransferClients[f.Addr].Close()
	return nil
}

func (f *Openfalcon) SampleConfig() string {
	return sampleConfig
}

func (f *Openfalcon) Description() string {
	return "Send telegraf metrics to open-falcon"
}



func (f *Openfalcon) Write(metrics []telegraf.Metric) error {
	var writeErr error = nil
	var  falconMetrics []*MetricValue
	for _, metric := range metrics {
		host, _ := metric.GetTag("host")
		for _, i := range (metric.FieldList()) {
			tags := ""
			for k,v := range metric.Tags() {
				tags = tags + fmt.Sprintf("%s=%s,",k,v)
			}
			tags = strings.Trim(tags,",")
			falconMetric := MetricValue{
				Endpoint: host,
				Metric:   metric.Name()+"."+i.Key,
				Value:    i.Value,
				Tags: tags,
				Timestamp: metric.Time().Unix(),
				Step: 10,
				Type: "GAUGE",
			}
			b := []byte(falconMetric.String())

			for _, writer := range f.writers {
				_, err := writer.Write(b)
				if err != nil && writer != os.Stdout {
					writeErr = fmt.Errorf("E! failed to write message: %s, %s", b, err)
				}
			}
			falconMetrics = append(falconMetrics, &falconMetric)
		}

	}
	f.SendToTransfer(falconMetrics)
	return writeErr
}

func init() {
	outputs.Add("openfalcon", func() telegraf.Output {
		return &Openfalcon{}
	})
}
