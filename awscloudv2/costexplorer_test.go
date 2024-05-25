package awscloudv2

import (
	"context"
	"math"
	"strconv"
	"testing"
)

func TestGetCostAndUsage(t *testing.T) {
	// test case 结构体定义每个测试用例的参数和预期结果
	type testCase struct {
		dateRange       string
		granularity     string
		groupBy         []string
		filterService   []string
		filterAccount   []string
		filterRegion    []string
		filterUsageType []string
		expectedCost    string
	}
	// 测试集合
	tests := []testCase{
		{
			dateRange:       "2024-04-01~2024-04-30",
			granularity:     "MONTHLY",
			groupBy:         []string{},
			filterService:   []string{"Amazon Elastic Load Balancing"},
			filterAccount:   []string{},
			filterRegion:    []string{},
			filterUsageType: []string{},
			expectedCost:    case1Cost,
		},
	}

	client, err := NewCeClient(ak, sk, region)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	for _, tc := range tests {
		results, err := client.GetCostAndUsage(context.TODO(), tc.dateRange, tc.granularity, tc.groupBy, tc.filterService, tc.filterAccount, tc.filterRegion, tc.filterUsageType)
		if err != nil {
			t.Errorf("Failed to get cost and usage: %v", err)
			continue
		}
		// 检查预期结果
		// 注意：这里我们将实际的检查逻辑简化为对比预期 cost 和第一个结果的 cost
		costStr, err := roundStringToTwoDecimals(*results[0].Total["UnblendedCost"].Amount)
		if err != nil {
			t.Errorf("Failed to get cost and usage: %v", err)
			continue
		}

		if results == nil || len(results) == 0 || costStr != tc.expectedCost {
			t.Errorf("Unexpected cost: got %v want %v", results, tc.expectedCost)
		}
	}
}

func roundStringToTwoDecimals(s string) (string, error) {
	// 将字符串转换为 float64
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", err
	}

	// 四舍五入到两位小数
	rounded := math.Round(f*100) / 100

	// 将 float64 转换回字符串，并格式化为两位小数
	result := strconv.FormatFloat(rounded, 'f', 2, 64)

	return result, nil
}
