package common

type TableName string

const (
	Blacklist      TableName = "blacklist"
	Whitelist      TableName = "whitelist"
	LoginBucketKey string    = "login"
	PassBucketKey  string    = "pass"
	IPBucketKey    string    = "ip"
	BlacklistURL   string    = "/list/blacklist"
	WhitelistURL   string    = "/list/whitelist"
	ResetURL       string    = "/reset/bucket"
)
