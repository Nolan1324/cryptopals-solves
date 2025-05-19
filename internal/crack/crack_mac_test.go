package crack

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"testing"
)

func TestExtendMac(t *testing.T) {
	const originalMessage = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis vitae facilisis nunc. Sed pulvinar imperdiet eros, eget tempus mauris sollicitudin ac. Etiam auctor tincidunt lorem quis lacinia. Aliquam erat volutpat. Aliquam erat volutpat. Nullam molestie, ipsum et tincidunt condimentum, diam dui dignissim tortor, eget consequat nulla orci non odio."
	const newMessage = "hello world"
	const key = "YELLOW SUBMARINE"

	m := cipherx.NewSha1Mac([]byte(key))

	originalMac := m.Sign([]byte(originalMessage))

	fullMessage, extendedMac := ExtendSha1Mac(originalMac, len(key), []byte(originalMessage), []byte(newMessage))

	if !bytes.Equal(extendedMac, m.Sign(fullMessage)) {
		t.Errorf("mac do not match")
	}
}
