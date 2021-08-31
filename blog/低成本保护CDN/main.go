package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

var signName string
var signKey string
var domain string

func init() {
	signName = os.Getenv("signName")
	signKey = os.Getenv("signKey")
	domain = os.Getenv("domain")
}

func signUrl(ctx context.Context, event events.APIGatewayRequest) (interface{}, error) {

	path := event.Path
	resourcePath := path[5:]

	ts := time.Now().Unix()
	rand_str := uuid.NewString()[:8]
	signStr := fmt.Sprintf("%s-%d-%s-%d-%s", resourcePath, ts, rand_str, 0, signKey)
	sign := md5.Sum([]byte(signStr))
	queryParams := fmt.Sprintf("%s=%s", signName, fmt.Sprintf("%d-%s-%d-%x", ts, rand_str, 0, sign))

	resp := &events.APIGatewayResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": fmt.Sprintf("%s%s?%s", domain, resourcePath, queryParams),
		},
	}

	return resp, nil
}

func main() {

	cloudfunction.Start(signUrl)
}
