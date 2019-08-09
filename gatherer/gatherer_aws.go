package gatherer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/logikone/gensops/config"
)

type awsGatherer struct {
	config *config.AWSConfig
	sess   *session.Session
}

func NewAWSGatherer(c *config.AWSConfig) AWSGatherer {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	return awsGatherer{
		config: c,
		sess:   sess,
	}
}

func (a awsGatherer) Gather() *AWSGathererResults {
	results := AWSGathererResults{}

	for _, account := range a.config.Accounts {
		roleName := account.IAMRole.Name

		if account.IAMRole.Prefix != "" {
			roleName = fmt.Sprintf("%s/%s", account.IAMRole.Prefix, account.IAMRole.Name)
		}

		roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", account.ID, roleName)

		creds := stscreds.NewCredentials(a.sess, roleArn)

		svc := kms.New(a.sess, &aws.Config{
			Credentials: creds,
		})

		out, err := svc.ListKeys(&kms.ListKeysInput{})

		if err != nil {
			fmt.Println("error listing keys:", err)
			continue
		}

		for _, key := range out.Keys {
			tagsOut, tagsErr := svc.ListResourceTags(&kms.ListResourceTagsInput{
				KeyId: key.KeyId,
			})

			if tagsErr != nil {
				continue
			}

			addKey := false

			for _, tag := range tagsOut.Tags {
				for _, match := range a.config.IncludeTags {
					if strings.ToLower(*tag.TagKey) == match {
						addKey = true
						break
					}
				}
			}

			if addKey {
				results.AddKey(&AWSKMSKey{
					Arn:  key.KeyArn,
					Tags: tagsOut.Tags,
				})
			}
		}

	}

	return &results
}
