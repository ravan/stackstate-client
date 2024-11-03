package stackstate

type StackState struct {
	ApiUrl    string `mapstructure:"api_url" validate:"required"`
	ApiKey    string `mapstructure:"api_key" validate:"required"`
	ApiToken  string `mapstructure:"api_token" validate:"required"`
	LegacyApi bool   `mapstructure:"legacy_api"`
}
