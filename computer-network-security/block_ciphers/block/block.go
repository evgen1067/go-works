package block

type Block []byte

type MessageBlocks struct {
	blocks []Block
	cipher CipherGost
}

type BlockEncoder interface {
	Encrypt() []byte
	Decrypt() []byte
	SetMessage([]byte)
}

func (g *MessageBlocks) SetMessage(message []byte) {
	var blocks = make([]Block, 0)
	for len(message)%8 != 0 {
		message = append(message, byte(32))
	}
	length := len(message)
	for i := 0; i < length; i += 8 {
		blocks = append(blocks, message[i:i+8])
	}
	g.blocks = blocks
}

func (g *MessageBlocks) Encrypt() []byte {
	result := make([]byte, 0)
	for i := range g.blocks {
		result = append(result, g.cipher.Encrypt(g.blocks[i])...)
	}
	return result
}

func (g *MessageBlocks) Decrypt() []byte {
	result := make([]byte, 0)
	for i := range g.blocks {
		result = append(result, g.cipher.Decrypt(g.blocks[i])...)
	}
	return result
}

func NewBlockEncoder(message []byte) BlockEncoder {
	g := &MessageBlocks{
		cipher: NewCipherGost(),
	}
	g.SetMessage(message)
	return g
}