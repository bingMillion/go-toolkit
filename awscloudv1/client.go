package awscloudv1

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Client struct {
	CeClient *costexplorer.CostExplorer
}

func NewCeClient(accessKeyID, secretAccessKey, region string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	cli := costexplorer.New(sess)
	return &Client{
		CeClient: cli,
	}, nil
}
