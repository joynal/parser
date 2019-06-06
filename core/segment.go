package core

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Value struct {
	Value     string `bson:"value"`
	Key       string `bson:"key"`
	Condition string `bson:"condition"`
}

type Property struct {
	Type   string  `bson:"type"`
	Values []Value `bson:"values"`
}

type RuleSet struct {
	ActionType string     `bson:"actionType"`
	Properties []Property `bson:"properties"`
}

// Segment model
type Segment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      primitive.ObjectID `bson:"_id,omitempty"`
	SiteID    primitive.ObjectID `bson:"siteId"`
	RuleSets  []RuleSet          `bson:"ruleSets"`
	IsDeleted bool               `bson:"isDeleted"`
	DeletedAt time.Time          `bson:"deletedAt"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
