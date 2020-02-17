package skiplist

import (
	"fmt"
	"math/rand"
)

const DefaultMaxLevel = 32
const DefaultP = 0.25

type SkipList struct {
	LessThen func(l, r interface{}) bool
	header   *node // header 节点不存 key/value
	footer   *node //末尾节点
	length   int
	MaxLevel int
	P        float64
}

func (s *SkipList) Print() {
	fmt.Println("********Skip List**********")
	for i := 0; i <= s.level(); i++ {
		curNode := s.header.forward[i]
		fmt.Printf("Level %v:", i)
		for curNode != nil {
			fmt.Printf("%v ", curNode.key)
			curNode = curNode.forward[i]
		}
		fmt.Println()
	}
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func (s *SkipList) effectiveMaxLevel() int {
	return maxInt(s.level(), s.MaxLevel)
}

func (s *SkipList) Delete(key interface{}) (value interface{}, ok bool) {
	if key == nil {
		panic("skiplist:nil key are not supported")
	}
	update := make([]*node, s.level()+1, s.effectiveMaxLevel()+1)
	candidate := s.getPath(s.header, update, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	previous := candidate.backward
	if s.footer == candidate {
		s.footer = previous
	}

	next := candidate.next()
	if next != nil {
		next.backward = previous
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}
	s.length--

	return candidate.value, ok
}

func (s *SkipList) Set(key, value interface{}) {
	if key == nil {
		return
	}
	// s.level starts from 0, so we need to allocate one.
	update := make([]*node, s.level()+1, s.effectiveMaxLevel()+1)
	candidate := s.getPath(s.header, update, key)

	if candidate != nil && candidate.key == key {
		candidate.value = value
		return
	}

	newLevel := s.randomLevel()
	if currentLevel := s.level(); currentLevel < newLevel {
		for i := currentLevel + 1; i <= newLevel; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	newNode := &node{
		forward: make([]*node, newLevel+1, s.effectiveMaxLevel()+1),
		key:     key,
		value:   value,
	}

	// level 0 为双向链表
	if previous := update[0]; previous.key != nil {
		newNode.backward = previous
	}

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	s.length++

	// level 0 为双向链表
	if newNode.forward[0] != nil {
		if newNode.forward[0].backward != newNode {
			newNode.forward[0].backward = newNode
		}
	}

	if s.footer == nil || s.LessThen(s.footer.key, key) {
		s.footer = newNode
	}
}

func (s *SkipList) Get(key interface{}) (value interface{}, ok bool) {
	candidate := s.getPath(s.header, nil, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	return candidate.value, true
}

// Iterator returns an Iterator that will go through all elements s.
func (s *SkipList) Iterator() Iterator {
	return &iter{
		current: s.header,
		list:    s,
	}
}

func (s *SkipList) Len() int {
	return s.length
}

func (s *SkipList) randomLevel() (n int) {
	for n = 0; n < s.effectiveMaxLevel() && rand.Float64() < s.P; n++ {
	}

	return
}

func (s *SkipList) getPath(current *node, update []*node, key interface{}) *node {
	depth := len(current.forward) - 1

	for i := depth; i >= 0; i-- {
		for current.forward[i] != nil && s.LessThen(current.forward[i].key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}

	return current.next()
}

// Ordered is an interface which can be linearly ordered by the
// LessThan method, whereby this instance is deemed to be less than
// other. Additionally, Ordered instances should behave properly when
// compared using == and !=.
type Ordered interface {
	LessThan(other Ordered) bool
}

// New returns a new SkipList.
// Its keys must implement the Ordered interface.
func New() *SkipList {
	comparator := func(left, right interface{}) bool {
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewCustomMap(comparator)

}

func NewCustomMap(lessThan func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		LessThen: lessThan,
		header: &node{
			forward: []*node{nil},
		},
		MaxLevel: DefaultMaxLevel,
		P:        DefaultP,
	}
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func NewIntMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}
