package gos

import (
	"strconv"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	var ta *Table
	n := "new table"
	t.Run("new", func(t *testing.T) {
		ta = Get(n)
		if ta == nil {
			t.Fail()
		}
	})
	t.Run("get", func(t *testing.T) {
		table := Get(n)
		if table == nil {
			t.Fail()
		}
		if table != ta {
			t.Fail()
		}
	})
	t.Run("delete", func(t *testing.T) {
		Delete(ta)
	})
}

func TestPut(t *testing.T) {
	ta := Get("test")
	key := []byte("key")
	value := "value"
	ta.Put(key, value)
	v, ok := ta.Get(key)
	if !ok {
		t.Fail()
	}
	if v.(string) != value {
		t.Fail()
	}
	Delete(ta)
}

func TestDelete(t *testing.T) {
	ta := Get("test")
	key := []byte("key")
	value := "value"
	ta.Put(key, value)
	if ta.Size() != 1 {
		t.Fail()
	}
	v, ok := ta.Delete(key)
	if !ok {
		t.Fail()
	}
	if v.(string) != value {
		t.Fail()
	}
	if ta.Size() != 0 {
		t.Fail()
	}
	Delete(ta)
}

func TestCopy(t *testing.T) {
	ta := Get("test")
	key := []byte("key")
	value1 := "value1"
	value2 := "value2"
	ta.Put(key, value1)
	tb := ta.Copy()
	tb.Put(key, value2)
	va, _ := ta.Get(key)
	if va.(string) != value1 {
		t.Fail()
	}
	vb, _ := tb.Get(key)
	if vb.(string) != value2 {
		t.Fail()
	}
	Delete(ta)
}

func TestRange(t *testing.T) {
	ta := Get("test")
	expect := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}
	for i := 0; i < 10; i++ {
		ta.Put([]byte(strconv.Itoa(i)), i)
	}
	t.Logf("size %d", ta.Size())
	ta.Range(func(k []byte, v interface{}) bool {
		t.Logf("range %s", k)
		if expect[string(k)] != v {
			t.Fail()
			return true
		}
		if v == 5 {
			return true
		}
		return false
	})
	Delete(ta)
}

func BenchmarkPut(b *testing.B) {
	ta := Get("test")
	b.Run("map", func(b *testing.B) {
		m := &sync.Map{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Store(strconv.Itoa(i), i)
		}
	})
	b.Run("ctrie", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ta.Put([]byte(strconv.Itoa(i)), i)
		}
	})
	Delete(ta)
}
