package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

const (
	metricsGatewayActivePrefix = hgwPrefix + "gateway-active/"
	metricsGatewayActiveFormat = metricsGatewayActivePrefix + "%s"
	metricsGatewayActiveDataPrefix = hgwPrefix + "gateway-active-data/"
	metricsGatewayActivesDataFormat = metricsGatewayActiveDataPrefix + "%s/"
	metricsGatewayActiveDataFormat = metricsGatewayActivesDataFormat + "%s"
)

func gatewayActivePath(serName string) string {
	return fmt.Sprintf(metricsGatewayActiveFormat, serName)
}

func gatewayActivesDataPath(serName string,) string {
	return fmt.Sprintf(metricsGatewayActivesDataFormat, serName)
}

func gatewaysData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, metricsGatewayActivePrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}

func gatewayMachineData(serName string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, gatewayActivePath(serName))
	cancel()
	return rsp, err
}

func gatewayData(serName string, limit int64) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	op := clientv3.WithLastKey()
	op = append(op, clientv3.WithLimit(limit))
	rsp, err := cli.Get(ctx, gatewayActivesDataPath(serName), op...)
	cancel()
	return rsp, err
}
