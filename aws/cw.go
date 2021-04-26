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
		LogGroupNamePrefix: aws.String(prefix),
	}
	var getGroupNames func(nt *string) error
	getGroupNames = func(nt *string) error {
		if nt != nil {
			param.NextToken = nt
		}
		groups, err := client.DescribeLogGroupsWithContext(ctx, &param)
		if err != nil {
			return err
		}
		for _, g := range groups.LogGroups {
			names = append(names, *g.LogGroupName)
		}
		if groups.NextToken != nil {
			getGroupNames(groups.NextToken)
		}
		return nil
	}
	if err := getGroupNames(nil); err != nil {
		return nil, err
	}
	return names, nil
}

func DescLogStreams(ctx context.Context, client *cloudwatchlogs.CloudWatchLogs, groupName string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	param := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(groupName),
		Descending:   aws.Bool(true),
		Limit:        aws.Int64(1),
	}
	return client.DescribeLogStreamsWithContext(ctx, &param)
}
