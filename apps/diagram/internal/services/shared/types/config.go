package types

type StoreWebConfigs struct {
	EntitiesAddress     string `json:"entities_address"`
	RelationsAddress    string `json:"relations_address"`
	TimeoutMilliseconds int    `json:"timeout_milliseconds"`
}
