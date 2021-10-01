package model

type RegisterBatchResponse struct {
	RowsAffected int64 `json:"rows_affected"`
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