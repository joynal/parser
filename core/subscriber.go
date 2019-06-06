package core

import (
	"github.com/SherClockHolmes/webpush-go"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// PUsh keys
type Keys struct {
	P256Dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

// push endpoint object
type PushEndPoint struct {
	Endpoint       string      `json:"endpoint"`
	ExpirationTime interface{} `json:"expirationTime"`
	Keys           Keys        `json:"keys"`
}

// Subscriber model
type Subscriber struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SiteID        primitive.ObjectID `bson:"siteId"`
	VisitorID     primitive.ObjectID `bson:"visitorId"`
	SubscriberID  string             `bson:"subscriberId"`
	PushEndpoint  string             `bson:"pushEndPoint"`
	Status        string
	LastActive    time.Time `bson:"lastActive"`
	FirstSession  time.Time `bson:"firstSession"`
	Platform      string
	Device        string
	Os            string
	Version       string
	IsMobile      bool   `bson:"isMobile"`
	IsDesktop     bool   `bson:"isDesktop"`
	CookieID      string `bson:"cookieId"`
	AppVersion    string `bson:"appVersion"`
	Browser       string
	UserAgent     string `bson:"userAgent"`
	Timezone      string
	Country       string
	Language      string
	UsageDuration int                  `bson:"usageDuration"`
	Segments      []primitive.ObjectID `bson:"segments"`
	IsDeleted     bool                 `bson:"isDeleted"`
	CreatedAt     time.Time            `bson:"createdAt"`
	UpdatedAt     time.Time            `bson:"updatedAt"`
}

type SubscriberPayload struct {
	PushEndpoint string `bson:"pushEndPoint"`
	Data         string
	Options      webpush.Options
	SubscriberID primitive.ObjectID `bson:"subscriberId"`
}
