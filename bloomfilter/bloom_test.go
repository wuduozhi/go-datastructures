package bloomfilter

import (
	"fmt"
	"testing"
)

func TestBloom(t *testing.T){
	bf := NewBloomFilter(3,1000)

	for i:=0;i<40;i++{
		bf.Add([]byte(fmt.Sprintf("%v%v%v",i,i+1,i+2)))
	}

	for i:=0;i<40;i++{
		check := bf.Check([]byte(fmt.Sprintf("%v%v%v",i,i+1,i+2)))
		if !check {
			t.Errorf("error happen")
		}
	}


}