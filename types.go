package client

import "time"

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type WalletResponse struct {
	ID             string `json:"id"`
	AddressPoolGap int    `json:"address_pool_gap"`
	Balance        struct {
		Available struct {
			Quantity int    `json:"quantity"`
			Unit     string `json:"unit"`
		} `json:"available"`
		Reward struct {
			Quantity int    `json:"quantity"`
			Unit     string `json:"unit"`
		} `json:"reward"`
		Total struct {
			Quantity int    `json:"quantity"`
			Unit     string `json:"unit"`
		} `json:"total"`
	} `json:"balance"`
	Assets struct {
		Available []struct {
			PolicyID  string `json:"policy_id"`
			AssetName string `json:"asset_name"`
			Quantity  int    `json:"quantity"`
		} `json:"available"`
		Total []struct {
			PolicyID  string `json:"policy_id"`
			AssetName string `json:"asset_name"`
			Quantity  int    `json:"quantity"`
		} `json:"total"`
	} `json:"assets"`
	Delegation struct {
		Active struct {
			Status string `json:"status"`
			Target string `json:"target"`
		} `json:"active"`
		Next []struct {
			Status    string `json:"status"`
			ChangesAt struct {
				EpochNumber    int       `json:"epoch_number"`
				EpochStartTime time.Time `json:"epoch_start_time"`
			} `json:"changes_at"`
		} `json:"next"`
	} `json:"delegation"`
	Name       string `json:"name"`
	Passphrase struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"passphrase"`
	State struct {
		Status string `json:"status"`
	} `json:"state"`
	Tip struct {
		AbsoluteSlotNumber int       `json:"absolute_slot_number"`
		SlotNumber         int       `json:"slot_number"`
		EpochNumber        int       `json:"epoch_number"`
		Time               time.Time `json:"time"`
		Height             struct {
			Quantity int    `json:"quantity"`
			Unit     string `json:"unit"`
		} `json:"height"`
	} `json:"tip"`
}
