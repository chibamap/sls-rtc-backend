package connection

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
)

// Connection dynamodb record structure
type Participant struct {
	// ConnectionID request.RequestContext.ConnectionID
	PK            string  `json:"pk"`
	ParticipantID *string `json:"participantID"`
	P256dh        *string `json:"wpushP256dh"`
	Auth          *string `json:"wpushAuth"`
	Endpoint      *string `json:"endpoint"`
	ExpiredAt     int64   `json:"expiredAt"`
}

const (
	pkPrefixParticipantId = "participantID:"
	participantIdPrefix   = "p-"
	ttl                   = "24h"
)

func NewParticipant() *Participant {
	participantId := participantIdPrefix + xid.New().String()
	pk := pkPrefixParticipantId + participantId
	limit, _ := time.ParseDuration(ttl)
	expiredAt := time.Now().Add(limit)
	return &Participant{
		PK:            pk,
		ParticipantID: &participantId,
		ExpiredAt:     expiredAt.Unix(),
	}
}

func FindParticipant(participantId string) (*Participant, error) {
	pk := pkPrefixParticipantId + participantId
	table := GetDefaultTable()
	item, err := table.find(pk)
	if err != nil {
		return nil, err
	}
	p := &Participant{}
	err = dynamodbattribute.UnmarshalMap(item, p)
	return p, err
}

func (p *Participant) Save() error {
	attributeValues, _ := dynamodbattribute.MarshalMap(p)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: ddbtable.tableName,
	}
	table := GetDefaultTable()
	_, err := table.ddb.PutItem(input)
	return err
}
