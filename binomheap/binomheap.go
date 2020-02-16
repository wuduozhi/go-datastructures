package binomheap

import (
	"errors"
	"fmt"
)

type BinomialHeap struct {
	head *BinomialNode
	size int
}

func (heap *BinomialHeap) Union(other BinomialHeap) error {
	head := binomialLink(heap.head, other.head)
	if head == nil {
		return errors.New("head is nil")
	}
	var preNode *BinomialNode
	curNode, nextNode := head, head.next

	for nextNode != nil {
		if curNode.degree != nextNode.degree ||
			(nextNode.next != nil && nextNode.degree == nextNode.next.degree) {
			// Case 1: x->degree != next_x->degree
			// Case 2: x->degree == next_x->degree == next_x->next->degree
			preNode = curNode
			curNode = nextNode
		} else if curNode.key < nextNode.key {
			// Case 3: x->degree == next_x->degree != next_x->next->degree
			//      && x->key    <= next_x->key
			curNode.next = nextNode.next
			binomialMerge(curNode, nextNode)
		} else {
			// Case 4: x->degree == next_x->degree != next_x->next->degree
			//      && x->key    >  next_x->key
			if preNode == nil {
				head = nextNode
			} else {
				preNode.next = nextNode
			}
			binomialMerge(nextNode, curNode)
			curNode = nextNode
		}
		nextNode = curNode.next
	}

	heap.head = head
	return nil
}

func (heap *BinomialHeap) Search(key int) *BinomialNode {
	return binomialSearch(key, heap.head)
}

func (heap *BinomialHeap) Insert(key int) error {
	if heap.Search(key) != nil {
		return errors.New("key is existed")
	}

	node := newBinomialNode(key)
	return heap.Union(BinomialHeap{head: node})
}

func (heap *BinomialHeap) Delete(key int) error {
	node := heap.Search(key)
	if node == nil {
		return nil
	}

	// 将被删除的节点的数据数据上移到它所在的二项树的根节点
	parent := node.parent
	for parent != nil {
		tmpKey := parent.key
		parent.key = node.key
		node.key = tmpKey
		node = parent
		parent = parent.parent
	}

	// 找到node的前一个根节点(prev)
	var pre *BinomialNode
	pos := heap.head
	for pos != node {
		pre = pos
		pos = pos.next
	}

	// 移除node节点
	if pre != nil {
		pre.next = node.next
	} else {
		heap.head = node.next
	}

	reverseNode := binomialReverse(node.child)
	node = nil
	reverseHeap := BinomialHeap{head: reverseNode}
	heap.Union(reverseHeap)

	return nil
}

func (heap *BinomialHeap) ToString() {
	if heap.head == nil {
		fmt.Println("nothing!")
		return
	}
	head := heap.head
	cur := heap.head
	fmt.Print("== binomialheap(")

	for cur != nil {
		fmt.Printf("B%v ", cur.degree)
		cur = cur.next
	}

	fmt.Print(") detail：\n")

	i := 0
	for head != nil {
		i++
		fmt.Printf("%v. binomial node root:%v,degree:%v .\n", i, head.key, head.degree)
		printBinomial(head, head.child, 1)
		head = head.next
	}

	fmt.Println()
}
