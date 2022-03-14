package app

// appConfig is struct for parsing ENV configuration.
type appConfig struct {
	// WoWHost is host of Word of Wisdom server.
	WoWHost string `config:"WOW_HOST,required"`
	// MaxComplexity is byte-length of challenge.
	MaxComplexity int `config:"MAX_COMPLEXITY,required"`
}
