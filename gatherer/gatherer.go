package gatherer

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/kms"
)

var lock sync.Mutex

type AWSGatherer interface {
	Gather() *AWSGathererResults
}

type AWSGathererResults struct {
	Keys []*AWSKMSKey
}

type AWSKMSKey struct {
	Arn      *string
	RoleArns []*string
	Tags     []*kms.Tag
}

func (r *AWSGathererResults) AddKey(key *AWSKMSKey) {
	lock.Lock()
	r.Keys = append(r.Keys, key)
	lock.Unlock()
}
