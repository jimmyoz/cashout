package _type

type Checkbook struct {
	Lastcheques []Lastcheques `json:"lastcheques"`
}
type Lastreceived struct {
	Beneficiary string `json:"beneficiary"`
	Chequebook  string `json:"chequebook"`
	Payout      int64  `json:"payout"`
}
type Lastsent struct {
	Beneficiary string `json:"beneficiary"`
	Chequebook  string `json:"chequebook"`
	Payout      int64  `json:"payout"`
}
type Lastcheques struct {
	Peer         string       `json:"peer"`
	Lastreceived Lastreceived `json:"lastreceived"`
	Lastsent     Lastsent     `json:"lastsent"`
}

type Balance struct {
	TotalBalance     int `json:"totalBalance"`
	AvailableBalance int `json:"availableBalance"`
}
