package binomheap

import "fmt"

type BinomialNode struct {
	key    int
	degree int
	child  *BinomialNode
	parent *BinomialNode
	next   *BinomialNode
}

func newBinomialNode(key int) *BinomialNode {
	return &BinomialNode{key: key, degree: 0}
}

/*
 * 将h1, h2中的根表合并成一个按度数递增的链表，返回合并后的根节点
 */
func binomialLink(h1, h2 *BinomialNode) *BinomialNode {
	head := &BinomialNode{}
	pos := &head

	for h1 != nil && h2 != nil {
		if h1.degree < h2.degree {
			*pos = h1
			h1 = h1.next
		} else {
			*pos = h2
			h2 = h2.next
		}
		pos = &((*pos).next)
	}

	if h1 != nil {
		*pos = h1
	}

	if h2 != nil {
		*pos = h2
	}
	return head
}

// 合并两个二项树：将 child 合并到 parent 中
func binomialMerge(parent, child *BinomialNode) {
	child.parent = parent
	child.next = parent.child
	parent.child = child
	parent.degree += 1
}

func binomialSearch(key int, node *BinomialNode) *BinomialNode {
	var child *BinomialNode = nil
	parent := node

	for parent != nil {
		if parent.key == key {
			return parent
		} else {
			child = binomialSearch(key, parent.child)
			if child != nil {
				return child
			}
			parent = parent.next
		}
	}

	return nil
}

/*
 * 打印"二项堆"
 *
 * 参数说明：
 *     node       -- 当前节点
 *     prev       -- 当前节点的前一个节点(父节点or兄弟节点)
 *     direction  --  1，表示当前节点是一个左孩子;
 *                    2，表示当前节点是一个兄弟节点。
 */
func printBinomial(pre, node *BinomialNode, direction int) {
	for node != nil {
		if direction == 1 {
			fmt.Printf("\t%v(%v) is %v's child", node.key, node.degree, pre.key)
		} else {
			fmt.Printf("\t%v(%v) is %v's next", node.key, node.degree, pre.key)
		}

		if node.child != nil {
			printBinomial(node, node.child, 1)
		}

		pre = node
		node = node.next
		direction = 2
	}
	fmt.Println()
}

/*
 * 反转二项堆heap
 */
func binomialReverse(node *BinomialNode) *BinomialNode {
	if node == nil {
		return node
	}

	var pre *BinomialNode
	var next *BinomialNode
	cur := node
	for cur != nil {
		cur.parent = nil

		next = cur.next
		cur.next = pre
		pre = cur
		cur = next
	}

	return pre
}
