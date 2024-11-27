package p2p

type HandshakeFunc func(any) error

func NOPHandShakeFunc(any) error {
	return nil
}
