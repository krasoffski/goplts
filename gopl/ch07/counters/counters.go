package counters

type Bytes int

func (b *Bytes) Write(p []byte) (int, error) {
	*b += Bytes(len(p))
	return len(p), nil
}
