package payment

type RequestService struct {
	Method string            `json:"method,omitempty"`
	Params requestParameters `json:"params,omitempty"`
}
type requestParameters struct {
	Data string `json:"data"`
	Uid  string `json:"uid"`
}
type responsePayment struct {
	Result string `json:"result"`
}
