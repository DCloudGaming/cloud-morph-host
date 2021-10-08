package model

type RegisterBatchResponse struct {
	RowsAffected int64 `json:"rows_affected"`
}

type VerifyOTPResponse struct {
	WalletAddress string `json:"wallet_address"`
	Token string `json:"token"`
}

type UserDetailProfileResponse struct {
	WalletAddress string           `json:"wallet_address"`
	CurUnreleasedBalance int       `json:"cur_unreleased_balance"`
	HourlyRate int                 `json:"hourly_rate"`
	Location  string               `json:"location"`
	Machine string                 `json:"machine"`
	RegisteredApps []RegisteredApp `json:"registered_apps"`
	PlaySessions []StreamSession   `json:"play_sessions"`
	HostSessions []StreamSession   `json:"host_sessions"`
}

type DiscoverAppResponse struct {
	ID string `json:"id"'`
	HostWalletAddress string `json:"host_wallet_address"`
	AppName string `json:"app_name"`
	AppPath string `json:"app_path"`
	Machine string `json:"machine"`
	HourlyRate int `json:"hourly_rate"`
	MaxDuration int `json:"max_duration"`
	Rating int `json:"rating"`
	Image string `json:"image"`
}

type GetAllowAppResponse struct {
	AppName string `json:"app_name"`
	VoteCount int `json:"vote_count"`
}

type GetAdminConfigsResponse struct {
	HourlyRate int `json:"hourly_rate"`
	AllowedApps []string `json:"allowed_apps"`
}
