package list

type Node struct {
	Val  any
	next *Node
	prev *Node
}
type List struct {
	head *Node
	size int
}

// NewList 初始化一个 *List 对象
func NewList() *List {
	list := &List{}
	list.size = 0
	list.head = &Node{}
	list.head.next = list.head
	list.head.prev = list.head
	return list
}

func (list *List) Size() int {
	return list.size
}

func (list *List) Empty() bool {
	if list.head.next == list.head {
		return true
	}
	return false
}

// InsertFront 头插法
func (list *List) InsertFront(val any) {
	node := &Node{
		Val: val,
	}
	node.next = list.head.next
	node.prev = list.head
	list.head.next.prev = node // 最后一个节点的next会连接到 head 节点
	list.head.next = node
	list.size++
}

func (list *List) MoveToHead(node *Node) {
	if node == list.head {
		return
	}
	if node == nil {
		return
	}
	node.prev.next = node.next
	node.next.prev = node.prev

	node.next = list.head.next
	node.prev = list.head
	list.head.next.prev = node
	list.head.next = node
}

func (list *List) Erase(node *Node) {
	if node == list.head {
		return
	}
	if node == nil {
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	list.size--
}
