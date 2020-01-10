package models

type Profile struct {
	Confirmed   bool     `json:"confirmed" json:"confirmed"`
	Visits      int64    `json:"visits" bson:"visits"`
	Likes       int64    `json:"likes" bson:"likes"`
	Fame        float64  `json:"fame" bson:"fame"`
	Propic      string   `json:"propic" bson:"propic"`
	Images      []string `json:"images" bson:"images"`
	Orientation string   `json:"orientation" bson:"orientation"`
	Interests   []string `json:"interests" bson:"interests"`
	Range       int64    `json:"range" bson:"range"`
	Index       int64    `Json:"index" bson:"index"`
}
