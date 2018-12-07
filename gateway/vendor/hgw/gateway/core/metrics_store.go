package core

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

const (
	metricsGatewayActivePrefix = hgwPrefix + "gateway-active/"
	metricsGatewayActiveFormat = metricsGatewayActivePrefix + "%s"
	metricsGatewayActiveDataPrefix = hgwPrefix + "gateway-active-data/"
	metricsGatewayActiveDataFormat = metricsGatewayActiveDataPrefix + "%s/%s"
)

const (
	activeTTL = 300
	activeDataTTL = 21600
)

func gatewayActiveK(serName string) string {
	return fmt.Sprintf(metricsGatewayActiveFormat, serName)
}

func gatewayActiveDataK(serName string,) string {
	return fmt.Sprintf(metricsGatewayActiveDataFormat, serName, time.Now().Format("2006-1-2-15-04"))
}

//设置网关机器系统数据
func putGatewayMachineData(serName string, jsonData string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	lease, err := cli.Grant(ctx, activeTTL)
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, gatewayActiveK(serName), jsonData, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	return nil
}

//设置网关统计数据
func putGatewayActiveData(serName string, jsonData string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	lease, err := cli.Grant(ctx, activeDataTTL)
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, gatewayActiveDataK(serName), jsonData, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	return nil
}