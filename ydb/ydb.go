package main

import (
	"context"
	"fmt"
	"log"

	//  "os"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := ydb.Open(ctx,
		"grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1gatc4m3hv1ldldhljp/etng5tv897avcvrror1v",
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("c:/go/key.json"),
	)
	log.Println(db)

	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
		whoAmI, err := db.Discovery().WhoAmI(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(whoAmI.String())

}
