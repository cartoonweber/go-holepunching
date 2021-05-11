package model

type Action string

const (
	Get Action = "GET"
	Reg Action = "REG"
)

// type GETRequest struct {
// 	Action Action
// 	PeerID string
// 	IP     string
// 	LAddr  *net.UDPAddr `json:"-"`
// }

type Request struct {
	Action Action
	PeerID string
	IP     string
}

type Response struct {
	RAddr *string
}
