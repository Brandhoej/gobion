package hasher

type Hashable64 interface {
	Sum64() uint64
}