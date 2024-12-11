package zmq

type ZMQTxInfo struct {
	TxID                        string        `json:"txid"`
	FromBlock                   bool          `json:"fromBlock"`
	Source                      string        `json:"source"`
	Address                     string        `json:"address"`
	NodeId                      int           `json:"nodeId"`
	Size                        int           `json:"size"`
	Hex                         string        `json:"hex"`
	IsInvalid                   bool          `json:"isInvalid"`
	IsValidationError           bool          `json:"isValidationError"`
	IsMissingInputs             bool          `json:"isMissingInputs"`
	IsDoubleSpendDetected       bool          `json:"isDoubleSpendDetected"`
	IsMempoolConflictDetected   bool          `json:"isMempoolConflictDetected"`
	IsNonFinal                  bool          `json:"isNonFinal"`
	IsValidationTimeoutExceeded bool          `json:"isValidationTimeoutExceeded"`
	IsStandardTx                bool          `json:"isStandardTx"`
	RejectionCode               int           `json:"rejectionCode"`
	Reason                      string        `json:"reason"`
	RejectionReason             string        `json:"rejectionReason"`
	CollidedWith                []interface{} `json:"collidedWith"`
	RejectionTime               string        `json:"rejectionTime"`
}

type ZMQDiscardFromMempool struct {
	TxID         string `json:"txid"`
	Reason       string `json:"reason"`
	CollidedWith struct {
		TxID string `json:"txid"`
		Size int    `json:"size"`
		Hex  string `json:"hex"`
	} `json:"collidedWith"`
	BlockHash string `json:"blockhash"`
}
