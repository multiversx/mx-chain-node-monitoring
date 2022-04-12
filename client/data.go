package client

// APIResponse defines an api reponse holding a list of nodes
type APIResponse struct {
	Nodes []APINode
}

// APINode defines an api node structure
type APINode struct {
	Bls                        string  `json:"bls"`
	Name                       string  `json:"name"`
	Version                    string  `json:"version"`
	Identity                   string  `json:"identity"`
	Rating                     int     `json:"rating"`
	TempRating                 float64 `json:"tempRating"`
	RatingModifier             float64 `json:"ratingModifier"`
	Shard                      int     `json:"shard"`
	Type                       string  `json:"type"`
	Status                     string  `json:"status"`
	Online                     bool    `json:"online"`
	Nonce                      int     `json:"nonce"`
	Instances                  int     `json:"instances"`
	Owner                      string  `json:"owner"`
	Provider                   string  `json:"provider"`
	Stake                      string  `json:"stake"`
	TopUp                      string  `json:"topUp"`
	Locked                     string  `json:"locked"`
	LeaderFailure              int     `json:"leaderFailure"`
	LeaderSuccess              int     `json:"leaderSuccess"`
	ValidatorFailure           int     `json:"validatorFailure"`
	ValidatorIgnoredSignatures int     `json:"validatorIgnoredSignatures"`
	ValidatorSuccess           int     `json:"validatorSuccess"`
	Position                   int     `json:"position"`
}
