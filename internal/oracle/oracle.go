package oracle

type Oracle interface {
	Encrypt([]byte) []byte
}
