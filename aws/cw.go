package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func NewCW(region string) *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(region),
	)
}

func ListLogGroup(ctx context.Context, client *cloudwatchlogs.CloudWatchLogs, prefix string) ([]string, error) {

	var names []string

	param := cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &prefix,
	}
	var getGroupNames func() error
	getGroupNames = func() error {
		groups, err := client.DescribeLogGroups(&param)
		if err != nil {
			return err
		}
		for _, g := range groups.LogGroups {
			names = append(names, *g.LogGroupName)
		}
		if groups.NextToken != nil {
			getGroupNames()
		}
		return nil
	}
	if err := getGroupNames(); err != nil {
		return nil, err
	}
	return names, nil
}
