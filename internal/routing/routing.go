package routing

const (
	ArmyMovesPrefix = "army_moves"

	WarRecognitionsPrefix = "war"

	PauseKey    = "pause"
	ArmyMoveKey = "army_moves"
	GameLogKey  = "game_logs.*"
	DLX_Key     = "peril_dlq"

	GameLogSlug = "game_logs"
)

const (
	ExchangePerilDirect = "peril_direct"
	ExchangePerilTopic  = "peril_topic"
	Exchange_DLX        = "peril_dlx"
)
