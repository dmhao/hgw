package core

import (
	"encoding/json"
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"hgw/gateway/def"
	"time"
)

var hgwMetrics metrics.Registry

func init() {
	hgwMetrics = metrics.NewRegistry()
}

func GetMetrics() metrics.Registry {
	return hgwMetrics
}

type Metrics struct {
	Timer		metrics.Timer
	Histograms	metrics.Histogram
}

func NewDomainMetrics(domain *def.Domain) *Metrics {
	m := metrics.NewPrefixedChildRegistry(hgwMetrics, domain.DomainUrl + "|-|" + "/*")
	return initMetrics(m)
}

func NewDomainPathMetrics(domain *def.Domain, path *def.Path) *Metrics {
	m := metrics.NewPrefixedChildRegistry(hgwMetrics, domain.DomainUrl + "|-|" + path.ReqPath)
	return initMetrics(m)
}

func initMetrics(m metrics.Registry) *Metrics {
	his := metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))
	m.GetOrRegister("|-|Histograms", his)
	return &Metrics{Histograms: his}
}

type RecordsData struct {
	Time			int64									`json:"time"`
	MetricsData		map[string]map[string]interface{}		`json:"metrics_data"`
}

func getRecordsData() RecordsData {
	//获取统计数据， 并清理历史数据
	data := make(map[string]map[string]interface{})
	hgwMetrics.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case metrics.Counter:
		case metrics.Gauge:
		case metrics.GaugeFloat64:
		case metrics.Healthcheck:
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			metric.Clear()
		case metrics.Meter:
		case metrics.Timer:
		}
		data[name] = values
	})
	recordsData := RecordsData{Time: time.Now().Unix(), MetricsData: data}
	return recordsData
}

func getMachineData() map[string]interface{} {
	data := make(map[string]interface{})
	mem, _ := mem.VirtualMemory()
	data["Mem"] = mem.String()
	avg,_ := load.Avg()
	data["Avg"] = avg.String()
	return data
}

//定时更新监控数据
func RecordMetrics(serName string) {
	for {
		machineData := getMachineData()
		machineBytes, err := json.Marshal(machineData)
		if err == nil {
			putGatewayMachineData(serName, string(machineBytes))
		}

		recordsData := getRecordsData()
		recordsBytes, err := json.Marshal(recordsData)
		if err == nil {
			putGatewayActiveData(serName, string(recordsBytes))
		}
		time.Sleep(1 * time.Minute)
	}
}