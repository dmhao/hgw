package modules

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Gw struct {
	ServerName	string		`json:"server_name"`
}

type Record struct {
	Count		float64			`json:"count"`
	Mean		float64			`json:"mean"`
	Max			float64			`json:"max"`
	Min			float64			`json:"min"`
	TimeStr		string			`json:"time_str"`
}

const (
	hour = 60
	halfHour = 30
	fifteenMinute = 15
	fiveMinute = 5
)

type RecordsData struct {
	Time			int64									`json:"time"`
	MetricsData		map[string]map[string]interface{}		`json:"metrics_data"`
}

func Gateways(c *gin.Context) {
	rsp, err := gatewaysData()
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	var gateways []*Gw
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			gateway := new(Gw)
			gateway.ServerName = strings.Replace(string(kv.Key), metricsGatewayActivePrefix, "", 1)
			gateways = append(gateways, gateway)
		}
		mContext{c}.SuccessOP(gateways)
		return
	}
	mContext{c}.SuccessOP(make([]string, 0))
}

func Gateway(c *gin.Context) {
	serName := c.Param("server_name")
	if serName == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	mcRsp,_ := gatewayMachineData(serName)
	machineData := string(mcRsp.Kvs[0].Value)

	rsp, err := gatewayData(serName, halfHour)
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	metricsMap := make(map[string][]Record)
	var timesData []string
	spanMap := make(map[string]string)

	for _, kv := range rsp.Kvs {
		recordData := RecordsData{}
		err := json.Unmarshal(kv.Value, &recordData)
		if err != nil {
			continue
		}
		fmt.Println(recordData)
		recordTime := recordData.Time
		timeData := time.Unix(recordTime, 0).Format("15:04")
		timesData = append(timesData, timeData)
		metricsData := recordData.MetricsData
		for oldSpan,_ := range metricsData {
			if _, ok := spanMap[oldSpan]; !ok {
				spans := strings.Split(oldSpan, "|-|")
				spanMap[oldSpan] = spans[0]+spans[1]
			}
			span := spanMap[oldSpan]
			if _, ok := metricsMap[span]; !ok {
				var data []Record
				metricsMap[span] = data
			}
		}
	}
	start := 0
	for _, kv := range rsp.Kvs {
		recordData := RecordsData{}
		err := json.Unmarshal(kv.Value, &recordData)
		if err != nil {
			continue
		}
		recordTime := recordData.Time
		timeData := time.Unix(recordTime, 0).Format("15:04")
		for k := range metricsMap {
			metricsMap[k] = append(metricsMap[k], Record{0,0,0, 0, timeData})
		}
		for oldSpan,metrics := range recordData.MetricsData {
			span := spanMap[oldSpan]
			metricsMap[span][start].Count = metrics["count"].(float64)
			metricsMap[span][start].Mean = Milliseconds(time.Duration(int64(metrics["mean"].(float64))))
			metricsMap[span][start].Max = Milliseconds(time.Duration(int64(metrics["max"].(float64))))
			metricsMap[span][start].Min = Milliseconds(time.Duration(int64(metrics["min"].(float64))))
		}
		start ++
	}
	allData := make(map[string]interface{})
	allData["metrics"] = metricsMap
	allData["times"] = timesData
	allData["machine"] = machineData
	mContext{c}.SuccessOP(allData)
}


func Milliseconds(d time.Duration) float64 {
	mill := d / time.Millisecond
	micrs := d % time.Microsecond
	return float64(mill) + float64(micrs)/1e6
}