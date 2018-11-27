package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

const (
	authDataPath = hgwPrefix + "auth-data/"
	authInitDataPath = hgwPrefix + "auth-data/init"
	adminUserDataPathFormat = authDataPath + "auth-data/user/%s"
)

func adminUserPath(username string) string {
	return fmt.Sprintf(adminUserDataPathFormat, username)
}

func authInitData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, authInitDataPath)
	cancel()
	return rsp,err
}


func putAdminUser(userId string, userJson string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	txn := cli.Txn(ctx)
	rsp, err := txn.Then(clientv3.OpPut(adminUserPath(userId), userJson),
		clientv3.OpPut(authInitDataPath, time.Now().Format("2006-01-02 15:04"))).Commit()
	cancel()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}


func adminUser(userId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, adminUserPath(userId))
	cancel()
	return rsp, err
}