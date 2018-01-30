package litecoinpool

/* Container */
type PoolData struct {
	User    User    `json:"user"`
	Workers Workers `json:"workers"`
	Pool    Pool    `json:"pool"`
	Network Network `json:"network"`
	Market  Market  `json:"market"`
}

/* name:Worker */
type Workers map[string]Worker

/*
	"user": {
		"hash_rate": 087354349,
		"expected_24h_rewards": 2.01037317042,
		"total_rewards": 1.4898754336,
		"paid_rewards": 10,
		"unpaid_rewards": 1.4898754335637,
		"past_24h_rewards": 1.132397627773,
		"total_work": "151077177713885184",
		"blocks_found": 12
	},
*/
type User struct {
	HashRate           int64   `json:"hash_rate"`
	Expected24hRewards float64 `json:"expected_24h_rewards"`
	TotalRewards       float64 `json:"total_rewards"`
	PaidRewards        float64 `json:"paid_rewards"`
	UnpaidRewards      float64 `json:"unpaid_rewards"`
	Past24hRewards     float64 `json:"past_24h_rewards"`
	//TotalWork          string  `json:"total_work"`
	BlocksFound int64 `json:"blocks_found"`
}

/*
	"NAME.1": {
		"connected": true,
		"hash_rate": 499310.8,
		"hash_rate_24h": 533975.3,
		"valid_shares": "000765502976",
		"stale_shares": "209415168",
		"invalid_shares": "80247",
		"rewards": 1.196998671636,
		"rewards_24h": 0.07401437913088,
		"last_share_time": 1517124416,
		"reset_time": 1512630506
	},
*/
type Worker struct {
	Connected     bool    `json:"connected"`
	HashRate      float64 `json:"hash_rate"`
	HashRate24h   float64 `json:"hash_rate_24h"`
	ValidShares   int64   `json:"valid_shares,string"`
	StaleShares   int64   `json:"stale_shares,string"`
	InvalidShares int64   `json:"invalid_shares,string"`
	Rewards       float64 `json:"rewards"`
	Rewards24h    float64 `json:"rewards_24h"`
	LastShareTime int64   `json:"last_share_time"`
	ResetTime     int64   `json:"reset_time"`
}

/*
	"pool": {
		"hash_rate": 28840000000,
		"active_users": 9914,
		"total_work": "145016027186199134208",
		"pps_ratio": 1.03,
		"pps_rate": 1.05305e-10
	},
*/
type Pool struct {
	HashRate    int64 `json:"hash_rate"`
	ActiveUsers int64 `json:"active_users"`
	//TotalWork   string  `json:"total_work"`
	PPSRatio float64 `json:"pps_ratio"`
	PPSRate  float64 `json:"pps_rate"`
}

/*
	"network": {
	   "hash_rate": 95925913949,
	   "block_number": 1358312,
	   "time_per_block": 167,
	   "difficulty": 3731199.1833743,
	   "next_difficulty": 3854043.0221471,
	   "retarget_time": 78852
   },
*/
type Network struct {
	HashRate       int64   `json:"hash_rate"`
	BlockNumber    int64   `json:"block_number"`
	TimePerBlock   int64   `json:"time_per_block"`
	Difficulty     float64 `json:"difficulty"`
	NextDifficulty float64 `json:"next_difficulty"`
	RetargetTime   int64   `json:"retarget_time"`
}

/*
	"market": {
		"ltc_btc": 0.01576,
		"ltc_usd": 183.02,
		"ltc_eur": 148.64,
		"ltc_gbp": 131.13855,
		"ltc_rub": 10444.23225,
		"ltc_cny": 1172.85335,
		"btc_usd": 11761.7
	}
*/
type Market struct {
	LTC_BTC float64 `json:"ltc_btc"`
	LTC_USD float64 `json:"ltc_usd"`
	LTC_EUR float64 `json:"ltc_eur"`
	LTC_GBP float64 `json:"ltc_gbp"`
	LTC_RUB float64 `json:"ltc_rub"`
	LTC_CNY float64 `json:"ltc_cny"`
	BTC_USD float64 `json:"btc_usd"`
}
