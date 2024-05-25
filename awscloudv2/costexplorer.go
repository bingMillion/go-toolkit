package awscloudv2

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"strings"
	"time"
)

func (c *Client) GetCostAndUsage(ctx context.Context, dateRange, granularity string, groupBy, filterService, filterAccount, filterRegion, filterUsageType []string) ([]types.ResultByTime, error) {
	startDate, endDate, err := parseDateRange(dateRange)
	if err != nil {
		return nil, err
	}
	// 构建过滤器
	var filters []types.Expression
	if len(filterService) > 0 {
		filters = append(filters, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionService,
				Values: filterService,
			},
		})
	}
	if len(filterAccount) > 0 {
		filters = append(filters, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionLinkedAccount,
				Values: filterAccount,
			},
		})
	}
	if len(filterRegion) > 0 {
		filters = append(filters, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionRegion,
				Values: filterRegion,
			},
		})
	}
	if len(filterUsageType) > 0 {
		filters = append(filters, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionUsageType,
				Values: filterUsageType,
			},
		})
	}

	// 构建最终的过滤器
	var finalFilter *types.Expression
	if len(filters) > 0 {
		if len(filters) == 1 {
			finalFilter = &types.Expression{
				Dimensions: filters[0].Dimensions,
			}
		}
		if len(filters) == 2 {
			finalFilter = &types.Expression{
				And: filters,
			}
		}
	}

	// 构建分组
	var groupDefinitions []types.GroupDefinition
	for _, g := range groupBy {
		key := g
		groupDefinitions = append(groupDefinitions, types.GroupDefinition{
			Key:  &key, // 需要确保g是有效的分组键
			Type: types.GroupDefinitionTypeDimension,
		})
	}

	req := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: &startDate,
			End:   &endDate,
		},
		Granularity: types.Granularity(granularity),
		GroupBy:     groupDefinitions,
		Metrics:     []string{"UnblendedCost", "UsageQuantity"},
	}
	// 如果存在过滤器，将其添加到请求中
	if finalFilter != nil {
		req.Filter = finalFilter
	}

	var allResults []types.ResultByTime
	for {
		output, err := c.CeClient.GetCostAndUsage(ctx, req)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, output.ResultsByTime...)
		if output.NextPageToken == nil {
			break
		}
		req.NextPageToken = output.NextPageToken
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
