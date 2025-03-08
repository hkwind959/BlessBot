package model

type UserNode struct {
	NodeID     string `mapstructure:"node_id"`
	Proxy      string `mapstructure:"proxy"`
	HardwareID string `mapstructure:"hardware_id"`
}

// User represents a user with their nodes
type User struct {
	UserToken string     `mapstructure:"user_token"`
	Remark    string     `mapstructure:"remark"`
	Nodes     []UserNode `mapstructure:"nodes"`
}

// Account represents the entire configuration
type Account struct {
	Users []User `mapstructure:"users"`
}

type RegisterNoe struct {
	Token      string
	NodeID     string
	Proxy      string
	HardwareID string
	IpAddress  string
	Remark     string // 备注
}
