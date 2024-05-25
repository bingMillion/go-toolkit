package awscloudv1

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"strings"
	"time"
)

func (c *Client) GetCostAndUsage(dateRange, granularity *string, groupBy, filterService, filterAccount, filterRegion, filterUsageType []*string) ([]*costexplorer.ResultByTime, error) {
	startDate, endDate, err := parseDateRange(*dateRange)
	if err != nil {
		return nil, err
	}
	// 构建过滤器
	var filters []*costexplorer.Expression
	if len(filterService) > 0 {
		filters = append(filters, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String(costexplorer.DimensionService),
				Values: filterService,
			},
		})
	}
	if len(filterAccount) > 0 {
		filters = append(filters, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String(costexplorer.DimensionLinkedAccount),
				Values: filterAccount,
			},
		})
	}
	if len(filterRegion) > 0 {
		filters = append(filters, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String(costexplorer.DimensionRegion),
				Values: filterRegion,
			},
		})
	}
	if len(filterUsageType) > 0 {
		filters = append(filters, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String(costexplorer.DimensionUsageType),
				Values: filterUsageType,
			},
		})
	}
	// 构建最终的过滤器
	var finalFilter *costexplorer.Expression
	if len(filters) > 0 {
		if len(filters) > 1 {
			finalFilter = &costexplorer.Expression{
				And: filters,
			}
		} else { // =1时不能用and
			finalFilter = &costexplorer.Expression{
				Dimensions: filters[0].Dimensions,
			}
		}

	}

	// 构建分组
	var groupDefinitions []costexplorer.GroupDefinition
	for _, g := range groupBy {
		key := g
		groupDefinitions = append(groupDefinitions, costexplorer.GroupDefinition{
			Key:  key, // 需要确保g是有效的分组键
			Type: aws.String(costexplorer.GroupDefinitionTypeDimension),
		})
	}

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: &startDate,
			End:   &endDate,
		},
		Granularity: granularity,
		Metrics:     []*string{aws.String("UsageQuantity"), aws.String("UnblendedCost")},
	}

	// 如果存在过滤器，将其添加到请求中
	if finalFilter != nil {
		input.Filter = finalFilter
	}

	var startToken *string
	var allResults []*costexplorer.ResultByTime
	for {
		input.NextPageToken = startToken
		// 发送查询请求
		output, err := c.CeClient.GetCostAndUsage(input)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, output.ResultsByTime...)
		if output.NextPageToken == nil {
			break
		}
		startToken = output.NextPageToken
	}

	return allResults, nil

}

// parseDateRange 校验输入时间是否符合2021-01-01~2021-01-31。 由于costexplorer api查询的日期不包含最后一天，这样月数据会有问题。 而aws页面上是包含最后一天的，故在接口中修正
func parseDateRange(dateRange string) (string, string, error) {
	dates := strings.Split(dateRange, "~")
	if len(dates) != 2 {
		return "", "", fmt.Errorf("invalid date range format, expected 'YYYY-MM-DD~YYYY-MM-DD'")
	}
	startDate := strings.TrimSpace(dates[0])
	endDate := strings.TrimSpace(dates[1])

	// Parse the endDate and add one day
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", "", fmt.Errorf("invalid end date format, expected 'YYYY-MM-DD'")
	}
	endTime = endTime.AddDate(0, 0, 1) // Add one day to include the end date
	endDate = endTime.Format("2006-01-02")

	return startDate, endDate, nil
}
