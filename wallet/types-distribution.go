package wallet

// Distribution represents a distribution of UTXO sizes for a given address.
// This type is manually added here, because oapi-codegen fails to generate it.
type Distribution struct {
	Total struct {
		Quantity int    `json:"quantity"`
		Unit     string `json:"unit"`
	} `json:"total"`
	Scale        string          `json:"scale"` // Expected enum value: "log10"
	Distribution map[uint64]uint `json:"distribution"`
}
