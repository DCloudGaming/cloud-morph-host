package model

type SignUpReq struct {
	WalletAddress string `json:"wallet_address"`
}

type GetOrCreateUserReq struct {
	WalletAddress string `json:"wallet_address"`
}

type UpdateUserReq struct {
	WalletAddress string `json:"wallet_address"`
	Machine string `json:"machine"`
	Location string `json:"location"`
	Name string `json:"name"`
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
	MaxDuration int64 `json:"max_duration"`
	ClientWalletAddress string `json:"client_wallet_address"`
	HostWalletAddress string `json:"host_wallet_address"`
	AppName string `json:"app_name"`
}

type UpdateSessionReq struct {
	SessionID int `json:"session_id"`
	TotalDuration int64 `json:"total_duration"`
	AccumCharge int64 `json:"accum_charge"`
	StreamStatus int `json:"stream_status"`
}

