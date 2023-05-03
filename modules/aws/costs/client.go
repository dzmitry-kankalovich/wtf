package awscosts

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"sort"
	"strconv"
	"time"
)

func getConfig(settings *Settings) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(settings.profile),
		config.WithDefaultRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

/*
Keep in mind that each API call has a cost of $0.01, so it's essential to construct your query
effectively using GroupBy, Filter, and other parameters to avoid unnecessary costs.
*/

func getMTD(client *costexplorer.Client) (string, string, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	timePeriod := &types.DateInterval{
		Start: aws.String(startOfMonth.Format("2006-01-02")),
		End:   aws.String(endOfMonth.Format("2006-01-02")),
	}

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod:  timePeriod,
		Granularity: types.GranularityMonthly,
		Metrics:     []string{"UnblendedCost"},
	}

	res, err := client.GetCostAndUsage(context.Background(), input)
	if err != nil {
		return "", "", err
	}

	if len(res.ResultsByTime) == 1 {
		unit := mapUnitToChar(aws.ToString(res.ResultsByTime[0].Total["UnblendedCost"].Unit))
		amount := trunkateAmount(aws.ToString(res.ResultsByTime[0].Total["UnblendedCost"].Amount))
		return unit, amount, nil
	} else {
		return "", "", fmt.Errorf("the MTD response size is %d, while expected 1", len(res.ResultsByTime))
	}
}

func getForecast(client *costexplorer.Client) (string, string, error) {
	now := time.Now()
	endOfMonth := now.AddDate(0, 1, 0)

	timePeriod := &types.DateInterval{
		Start: aws.String(now.Format("2006-01-02")),
		End:   aws.String(endOfMonth.Format("2006-01-02")),
	}

	input := &costexplorer.GetCostForecastInput{
		TimePeriod:  timePeriod,
		Granularity: types.GranularityMonthly,
		Metric:      "UNBLENDED_COST",
	}

	res, err := client.GetCostForecast(context.Background(), input)
	if err != nil {
		return "", "", err
	}

	unit := mapUnitToChar(aws.ToString(res.Total.Unit))
	amount := trunkateAmount(aws.ToString(res.Total.Amount))
	return unit, amount, nil
}

type ServiceCost struct {
	Name        string
	Amount      string
	AmountFloat float64
	Unit        string
}

func getTopN(client *costexplorer.Client, limit int) ([]ServiceCost, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	timePeriod := &types.DateInterval{
		Start: aws.String(startOfMonth.Format("2006-01-02")),
		End:   aws.String(endOfMonth.Format("2006-01-02")),
	}

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod:  timePeriod,
		Granularity: types.GranularityMonthly,
		Metrics:     []string{"UnblendedCost"},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	}

	res, err := client.GetCostAndUsage(context.Background(), input)
	if err != nil {
		return nil, err
	}

	costs := make([]ServiceCost, 0)

	for _, group := range res.ResultsByTime[0].Groups {
		name := group.Keys[0]
		amount := trunkateAmount(aws.ToString(group.Metrics["UnblendedCost"].Amount))
		amountFloat, _ := strconv.ParseFloat(amount, 32)
		unit := mapUnitToChar(aws.ToString(group.Metrics["UnblendedCost"].Unit))
		costs = append(costs, ServiceCost{Name: name, Amount: amount, AmountFloat: amountFloat, Unit: unit})
	}

	sort.Slice(costs, func(i, j int) bool {
		return costs[i].Amount > costs[j].Amount
	})

	var sliceLimit int
	if len(costs) < limit {
		sliceLimit = len(costs)
	} else {
		sliceLimit = limit
	}

	return costs[:sliceLimit], nil
}
