package disjointSet

type disjointSetNode struct {
	parent *disjointSetNode
	rank   int
}

type DisjointSet struct {
	master map[int64]*disjointSetNode
}

func NewDisjointSet() *DisjointSet {
	return &DisjointSet{master: make(map[int64]*disjointSetNode, 0)}
}

func (ds *DisjointSet) MakeSet(e int64) {
	if _, ok := ds.master[e]; ok {
		return
	}

	dsNode := &disjointSetNode{rank: 0}
	dsNode.parent = dsNode
	ds.master[e] = dsNode
}

func (ds *DisjointSet) Union(x, y *disjointSetNode) {
	if x == nil || y == nil {
		panic("Disjoint set unit on nil sets")
	}

	xRoot := find(x)
	yRoot := find(y)
	if xRoot == nil || yRoot == nil ||  xRoot == yRoot{
		return
	}

	if xRoot.rank < yRoot.rank {
		xRoot.parent = yRoot
	}else if yRoot.rank < xRoot.rank{
		yRoot.parent = xRoot
	}else{
		yRoot.parent = xRoot
		xRoot.rank ++
	}
}

func (ds *DisjointSet) Find(e int64) *disjointSetNode {
	dsNode, ok := ds.master[e]
	if !ok {
		return nil
	}

	return find(dsNode)
}

func find(dsNode *disjointSetNode) *disjointSetNode {
	if dsNode.parent != dsNode {
		dsNode = find(dsNode.parent)
	}
	return dsNode.parent
}
