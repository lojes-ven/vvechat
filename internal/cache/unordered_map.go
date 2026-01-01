package cache

type node struct {
	key   any
	value any
	next  *node
}
type HashMap struct {
	size       int
	capacity   int
	loadFactor float64
	bucket     []*node
}

// NewHashMap 初始化哈希表
func NewHashMap(capacity int) *HashMap {
	return &HashMap{
		capacity: capacity,
		bucket:   make([]*node, capacity),
	}
}

func hash(key any) uint64 {

}
