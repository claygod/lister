package lister

// Lister
// Library for work with lists of structures in a functional style
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "reflect"
//import "fmt"
import "sort"
import "sync"
import "sync/atomic"
import "time"

// import "errors"

// NewDb - create a new Db
func NewDb() *Db { //s interface{}
	d := &Db{
		arr:    make(Lister, 0, 100),
		update: time.Now().UnixNano(),
	}
	return d
}

// Db structure
type Db struct {
	sync.Mutex
	update int64 // time.Now().UnixNano()
	arr    Lister
}

func (d *Db) Updated() int64 {
	return atomic.LoadInt64(&d.update)
}

func (d *Db) Start() Lister {
	d.Lock()
	out := make(Lister, len(d.arr))
	copy(out, d.arr)
	d.Unlock()
	//for _, item := range db {
	//	out = append(out, item)
	//}
	return out
}

func (d *Db) Add(item interface{}, fu func(a interface{}) bool) bool {
	d.Lock()
	defer d.Unlock()
	for _, item := range d.arr {
		if fu(item) {
			return false
		}
	}
	// тут возможна проверка на тип
	d.arr = append(d.arr, item)
	return true

}

func (d *Db) Del(fu func(interface{}) bool) bool {
	d.Lock()
	defer d.Unlock()
	d.arr = d.arr.Filter(fu)
	for k, item := range d.arr {
		if fu(item) {
			if k == 0 {
				d.arr = d.arr[1:]
			} else {
				d.arr = append(d.arr[:(k-1)], d.arr[k:]...)
			}
			return true
		}
	}
	return false
}

func (d *Db) Update(newArr Lister) {
	d.Lock()
	d.arr = newArr
	d.update = time.Now().UnixNano()
	d.Unlock()
}

type Lister []interface{}

func (d Lister) Map(fu func(interface{}) interface{}) Lister { // трансформация
	var out Lister
	for _, item := range d {
		out = append(out, fu(item))
	}
	d = out
	return d
}

func (d Lister) Filter(fu func(interface{}) bool) Lister { // фильтр

	var out Lister
	for _, item := range d {
		if fu(item) {
			out = append(out, item)
		}
	}
	d = out
	return d
}

func (d Lister) SortUp(fu func(interface{}) int) Lister { // сортировка
	out := make(Lister, len(d))

	a := make(map[int]int)
	for k, item := range d {
		a[fu(item)] = k
	}

	var keys []int
	for k, _ := range a {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k, key := range keys {
		pos := a[key]
		out[k] = d[pos]
	}

	d = out
	return d
}

func (d Lister) SortDown(fu func(interface{}) int) Lister { // сортировка
	out := make(Lister, len(d))

	a := make(map[int]int)
	for k, item := range d {
		a[fu(item)] = k
	}

	var keys []int
	for k, _ := range a {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	//fmt.Print(keys, "-------------\r\n")
	for k, key := range keys {
		pos := a[key]
		out[k] = d[pos]
	}
	d = out
	return d
}
