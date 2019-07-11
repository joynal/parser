package core

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Message - Notification message
type Message struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string
	Message  string
	Language string
}

// Action - Browser push action configs
type Action struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Action  string
	Title   string
	URL     string
	IconURL string `bson:"icon"`
}

// Browser push configs
type Browser struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	BrowserName        string             `bson:"browserName"`
	IconURL            string             `bson:"iconURL"`
	ImageURL           string             `bson:"imageURL"`
	Badge              string
	Vibration          bool
	IsActive           bool `bson:"isActive"`
	IsEnabledCTAButton bool `bson:"isEnabledCTAButton"`
	Actions            []Action
}

// HideRule when notification popup will be closed.
type HideRule struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Type  string
	Value int
}

type SendTo struct {
	AllSubscriber bool                 `bson:"allSubscriber"`
	Segments      []primitive.ObjectID `bson:"segments"`
}

// Notification model
type Notification struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	SendTo             SendTo             `bson:"sendTo"`
	Priority           string             `bson:"priority"`
	TimeToLive         int                `bson:"timeToLive"`
	TotalSent          int                `bson:"totalSent"`
	TotalDeliver       int                `bson:"totalDeliver"`
	TotalShow          int                `bson:"totalShow"`
	TotalError         int                `bson:"totalError"`
	TotalClick         int                `bson:"totalClick"`
	TotalClose         int                `bson:"totalClose"`
	IsAtLocalTime      bool               `bson:"isAtLocalTime"`
	IsProcessed        string             `bson:"isProcessed"`
	IsSchedule         bool               `bson:"isSchedule"`
	TimezonesCompleted []string           `bson:"timezonesCompleted"`
	Timezone           string
	IsDeleted          bool               `bson:"isDeleted"`
	FromRSSFeed        bool               `bson:"fromRSSFeed"`
	SiteID             primitive.ObjectID `bson:"siteId"`
	Messages           []Message
	Browsers           []Browser
	Actions            []Action
	HideRules          HideRule           `bson:"hideRules"`
	LaunchURL          string             `bson:"launchUrl"`
	UserID             primitive.ObjectID `bson:"userId"`
	SentAt             time.Time          `bson:"sentAt"`
	CreatedAt          time.Time          `bson:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt"`
}

// ProcessedNotification model

type VAPIDOptions struct {
	Subject    string
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type ProcessedNotification struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SiteID        primitive.ObjectID `bson:"siteId"`
	TimeToLive    int                `bson:"timeToLive"`
	LaunchURL     string             `bson:"launchUrl"`
	Message       Message
	Browser       []Browser
	Actions       []Action
	HideRules     HideRule     `bson:"hideRules"`
	SendTo        SendTo       `bson:"sendTo"`
	IsAtLocalTime bool         `bson:"isAtLocalTime"`
	VapidDetails  VapidDetails `bson:"vapidDetails"`
	Timezone      string
}

type NotificationPayload struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	LaunchURL string             `bson:"launchUrl"`
	Message   Message
	Browser   []Browser
	HideRules HideRule `bson:"hideRules"`
	Actions   []Action
}

type VapidDetails struct {
	VapidPublicKeys  string `bson:"vapidPublicKeys"`
	VapidPrivateKeys string `bson:"vapidPrivateKeys"`
}

type DisplayCondition struct {
	Type  string `bson:"type"`
	Value int    `bson:"value"`
}

// Notification model
type NotificationAccount struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	DisplayCondition  DisplayCondition   `bson:"displayCondition"`
	VapidDetails      VapidDetails       `bson:"vapidDetails"`
	TotalSubscriber   int                `bson:"totalSubscriber"`
	TotalUnSubscriber int                `bson:"totalUnSubscriber"`
	Status            bool               `bson:"status"`
	HTTPSEnabled      bool               `bson:"httpsEnabled"`
	IsDeleted         bool               `bson:"isDeleted"`
	SiteID            primitive.ObjectID `bson:"siteId"`
	UserID            primitive.ObjectID `bson:"userId"`
	Domain            string             `bson:"domain"`
	SubDomain         string             `bson:"subDomain"`
	CreatedAt         time.Time          `bson:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt"`
}
