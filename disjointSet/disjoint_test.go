package disjointSet

import (
	"testing"
)

func TestDisjointSetMakeSet(t *testing.T) {
	ds := NewDisjointSet()
	if ds.master == nil {
		t.Fatal("Internal disjoint set map erroneously nil")
	} else if len(ds.master) != 0 {
		t.Error("Disjoint set master map of wrong size")
	}

	ds.MakeSet(3)
	if len(ds.master) != 1 {
		t.Error("Disjoint set master map of wrong size")
	}

	if node, ok := ds.master[3]; !ok {
		t.Error("Make set did not successfully add element")
	} else {
		if node == nil {
			t.Fatal("Disjoint set node from makeSet is nil")
		}

		if node.rank != 0 {
			t.Error("Node rank set incorrectly")
		}

		if node.parent != node {
			t.Error("Node parent set incorrectly")
		}
	}
}

func TestDisjointSetFind(t *testing.T) {
	ds := NewDisjointSet()

	ds.MakeSet(3)
	ds.MakeSet(5)

	if ds.Find(3) == ds.Find(5) {
		t.Error("Disjoint sets incorrectly found to be the same")
	}
}

func TestUnion(t *testing.T) {
	ds := NewDisjointSet()

	ds.MakeSet(3)
	ds.MakeSet(5)

	ds.Union(ds.Find(3), ds.Find(5))

	if ds.Find(3) != ds.Find(5) {
		t.Error("Sets found to be disjoint after union")
	}
}