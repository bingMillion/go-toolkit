package awscloudv2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

/*
1. 开发时，需明确指定下载v2的包：
* go get github.com/aws/aws-sdk-go-v2
* go get github.com/aws/aws-sdk-go-v2/config

2. 还需要什么服务，就下载对应service下的包。因为这些其实是独立的子项目，有自己的gomod。这里比如
* go get github.com/aws/aws-sdk-go-v2/service/costexplorer

*/

type Client struct {
	CeClient *costexplorer.Client
}

func NewCeClient(accessKeyID, secretAccessKey, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		),
	)
	if err != nil {
		// 处理配置加载错误
		return nil, err
	}

	// 创建 Cost Explorer 服务客户端实例
	ceClient := costexplorer.NewFromConfig(cfg)

	// 返回 Client 结构体实例
	return &Client{
		CeClient: ceClient,
	}, nil
}
