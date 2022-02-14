package valid

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Man struct {
	Name string
	Age  int
}

type User struct {
	Man   Man
	Hobby []int32
	// Class *Class
	Class Class
}

type Class struct {
	Name         string
	TeatherNames []string
	Testhers     []*Testher
}

type Testher struct {
	Man Man
}

func TestDump0(t *testing.T) {
	m := Man{
		Name: "XUE",
		Age:  20,
	}
	t.Logf("%+v", m)
	t.Log(GetDumpStructStr(m))
}

func TestDump1(t *testing.T) {
	type SliceDemo struct {
		Name  string
		Hobby []int32
	}

	d := SliceDemo{
		Name:  "xue",
		Hobby: []int32{1, 2, 4},
	}
	t.Logf("%+v", d)
	t.Log(GetDumpStructStr(d))
}

func TestDump2(t *testing.T) {
	type SliceDemo struct {
		Name     string
		Testhers []*Testher
	}

	d := SliceDemo{
		Name:     "xue",
		Testhers: []*Testher{{Man{Name: "test1", Age: 11}}, {Man{Name: "test2", Age: 11}}},
	}
	t.Logf("%+v", d)
	t.Log(GetDumpStructStr(d))
}

func TestDump3(t *testing.T) {
	type Demo struct {
		Name    string
		Testher *Testher
	}

	d := Demo{
		Name: "xue",
		Testher: &Testher{
			Man: Man{
				Name: "xue1",
				Age:  10,
			},
		},
	}
	t.Logf("%+v", d)
	t.Log(GetDumpStructStr(d))
}

func TestDump4(t *testing.T) {
	type Demo struct {
		Name    string
		Testher *Testher
		Map     map[int32]string
	}

	d := Demo{
		Name: "xue",
		Testher: &Testher{
			Man: Man{
				Name: "xue1",
				Age:  10,
			},
		},
		Map: map[int32]string{1: "TEST"},
	}
	t.Logf("%+v", d)
	t.Log(GetDumpStructStr(d))
}

func TestDump5(t *testing.T) {
	type Demo struct {
		Name    string
		Testher *Testher
		Map     map[int32]string
		Fn      func() int
		Ch      chan int
	}
	d := &Demo{
		Name: "xue",
		Testher: &Testher{
			Man: Man{
				Name: "xue1",
				Age:  10,
			},
		},
		Map: map[int32]string{1: "TEST", 2: "test"},
		Fn:  func() int { return 1 },
		Ch:  make(chan int),
	}
	t.Logf("%+v", d)

	t.Log(GetDumpStructStr(d))

	b, err := json.Marshal(d)
	t.Log(err, string(b))
}

func TestDump(t *testing.T) {
	u := &User{
		Man: Man{
			Name: "xuesongtao",
			Age:  20,
		},
		Hobby: []int32{1},
		Class: Class{
			Name:         "社会大学1",
			TeatherNames: []string{"社佬"},
			Testhers: []*Testher{
				{
					Man{
						Name: "社佬",
						Age:  11,
					},
				},
				{
					Man{
						Name: "社佬1",
						Age:  11,
					},
				},
			},
		},
	}
	t.Logf("%+v", u)
	d := NewDumpStruct()
	t.Log(d.GetDumpStructStr(u))
}

func BenchmarkDump0(b *testing.B) {
	u := &User{
		Man: Man{
			Name: "xuesongtao",
			Age:  20,
		},
		Hobby: []int32{1},
		Class: Class{
			Name:         "社会大学1",
			TeatherNames: []string{"社佬"},
			Testhers: []*Testher{
				{
					Man{
						Name: "社佬",
						Age:  11,
					},
				},
				{
					Man{
						Name: "社佬1",
						Age:  11,
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		_ = GetDumpStructStr(u)
	}
}

func BenchmarkDump1(b *testing.B) {
	u := &User{
		Man: Man{
			Name: "xuesongtao",
			Age:  20,
		},
		Hobby: []int32{1},
		Class: Class{
			Name:         "社会大学1",
			TeatherNames: []string{"社佬"},
			Testhers: []*Testher{
				{
					Man{
						Name: "社佬",
						Age:  11,
					},
				},
				{
					Man{
						Name: "社佬1",
						Age:  11,
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%+v", u)
	}
}

func BenchmarkDump2(b *testing.B) {
	u := &User{
		Man: Man{
			Name: "xuesongtao",
			Age:  20,
		},
		Hobby: []int32{1},
		Class: Class{
			Name:         "社会大学1",
			TeatherNames: []string{"社佬"},
			Testhers: []*Testher{
				{
					Man{
						Name: "社佬",
						Age:  11,
					},
				},
				{
					Man{
						Name: "社佬1",
						Age:  11,
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		b, _ := json.Marshal(u)
		_ = string(b)
	}
}
