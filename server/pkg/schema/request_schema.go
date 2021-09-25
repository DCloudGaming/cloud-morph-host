package schema

type SignUpReq struct {
	WalletAddress string `json:"wallet_address"`
}

type AuthReq struct {
	WalletAddress string `json:"wallet_address"`
	Signature string `json:"signature"`
}

type MockAuthReq struct {
	WalletAddress string `json:"wallet_address"`
}

type RegisterAppReq struct {
	WalletAddress string `json:"wallet_address"`
	AppPaths []string `json:"app_paths"`
	AppNames []string `json:"app_names"`
}

type StartSessionReq struct {
	MaxDuration int `json:"max_duration"`
	ClientWalletAddress string `json:"client_wallet_address"`
	HostWalletAddress string `json:"host_wallet_address"`
	AppName string `json:"app_name"`
}

type UpdateSessionReq struct {
	SessionID int `json:"session_id"`
	TotalDuration int `json:"total_duration"`
	AccumCharge int `json:"accum_charge"`
	StreamStatus int `json:"stream_status"`
}