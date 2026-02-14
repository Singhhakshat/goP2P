package p2p

const (
	HandShakeMessage = "P2P-SHARE-v1.0"
)

func validateHandshake(msg string) bool {
	return msg == HandShakeMessage
}
