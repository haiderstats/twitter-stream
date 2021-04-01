package models

type GetMeta struct {
	Sent string
}

type Rule struct {
	Id    string
	Value string
	Tag   string
}

type GetResponse struct {
	Data []Rule
	Meta GetMeta
}

type DeleteSummary struct {
	Deleted    int
	NotDeleted int `json:"not_deleted"`
}

type DeleteMeta struct {
	Sent    string
	Summary DeleteSummary
}

type DeleteError struct {
	Message string
}

type DeleteErrors struct {
	Errors []DeleteError
	Title  string
	Detail string
	Type   string
}

type DeleteResponse struct {
	Meta   DeleteMeta
	Errors *[]DeleteErrors
}

type IdList struct {
	Ids []string `json:"ids"`
}

type DeleteRules struct {
	Delete IdList `json:"delete"`
}

type AddRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

type CreateRules struct {
	Add []AddRule `json:"add"`
}

type CreateError struct {
	Value   string
	Details []string
	Title   string
	Type    string
}

type CreateSummary struct {
	Created    int
	NotCreated int `json:"not_created"`
	Valid      int
	Invalid    int
}

type CreateMeta struct {
	Sent    string
	Summary CreateSummary
}

type CreateResponse struct {
	Data   []Rule
	Meta   CreateMeta
	Errors *[]CreateError
}

type Attachment struct {
	MediaKeys []string `json:"media_keys"`
}

type TweetData struct {
	Id          string
	Text        string
	Attachments Attachment
}

type TwitterImage struct {
	MediaKey string `json:"media_key"`
	Type     string
	Url      string
}

type TweetInclusions struct {
	Media []TwitterImage
}

type Match struct {
	Id  int
	Tag string
}

type Tweet struct {
	Data          TweetData
	Includes      TweetInclusions
	MatchingRules []Match `json:"matching_rules"`
}
