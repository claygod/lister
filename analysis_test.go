package lister

// Lister
// Tests and benchmarks
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "fmt"
import "testing"
import "strconv"

func Test100(t *testing.T) {
	x := Article{}
	d := genDb()

	res := d.Start().
		Filter(x.FilterPub).
		SortDown(x.SortDate)

	if len(res) != 3 {
		t.Error("Incorrect number in the returned slice")
	}

	res = d.Start().
		Filter(x.FilterPub).
		Filter(x.FilterDate(234234234238))

	if len(res) != 2 {
		t.Error("Incorrect number in the returned slice")
	}
}

func TestSortUp(t *testing.T) {
	x := Article{}
	d := genDb()

	res := d.Start().
		SortUp(x.SortDate)
	a := res[0].(Article)

	if a.Date != 234234234231 {
		t.Error("Error sorting ascending")
	}
}

func TestSortDown(t *testing.T) {
	x := Article{}
	d := genDb()

	res := d.Start().
		SortDown(x.SortDate)
	a := res[0].(Article)

	if a.Date != 234234234239 {
		t.Error("Error sorting descending")
	}
}

func TestDel(t *testing.T) {
	x := Article{}
	d := genDb()
	d.Del(x.FilterButById("a1"))

	res := d.Start().
		Filter(x.FilterPub)
	if len(res) != 2 {
		t.Error("Error deleting the element")
	}
}

func TestMapDateMinus(t *testing.T) {
	x := Article{}
	d := genDb()

	d.Update(d.Start().Map(x.MapDateMinus(1e11)))

	res := d.Start().
		Filter(x.FilterPub)
	//fmt.Print(res, "---+++--------\r\n")
	for _, a := range res {
		if a.(Article).Date > 2e11 {
			t.Error("Failed to amend by means `Map` and `MapDateMinus`")
			break
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	x := Article{}
	fu := x.FilterById("a1")
	d := NewDb()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	b.StartTimer()
	for n := 0; n < 3000; n++ {
		a.Id = strconv.Itoa(n)
		d.Add(a, fu)
	}
}

func BenchmarkMainParallel(b *testing.B) {
	b.StopTimer()
	x := Article{}
	fu := x.FilterById("a1")
	d := NewDb()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body ",
		Tags:  []string{"news", "april"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	n := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Id = strconv.Itoa(n)
			d.Add(a, fu)
			n++
		}
	})
}

func BenchmarkDel(b *testing.B) {
	b.StopTimer()
	x := Article{}
	fu := x.FilterById("a1")
	d := NewDb()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	for n := 0; n < 3000; n++ {
		a.Id = strconv.Itoa(n)
		d.Add(a, fu)
	}

	b.StartTimer()
	for n := 3000; n > 0; n-- {
		d.Del(x.FilterButById(strconv.Itoa(n)))
	}
}

func BenchmarkAddaDel(b *testing.B) {
	b.StopTimer()
	x := Article{}
	fu := x.FilterById("a1")
	d := NewDb()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	for n := 0; n < 2000; n++ {
		a.Id = strconv.Itoa(n)
		d.Add(a, fu)
	}
	a.Id = "abc"
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		d.Add(a, fu)
		d.Del(x.FilterButById("abc"))
	}
}

func BenchmarkSelect(b *testing.B) {
	b.StopTimer()
	x := Article{}
	fu := x.FilterById("a1")
	d := NewDb()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	for n := 0; n < 3000; n++ {
		a.Id = strconv.Itoa(n)
		d.Add(a, fu)
	}

	b.StartTimer()
	for n := 3000; n > 0; n-- {
		d.Start().
			Filter(x.FilterById(strconv.Itoa(n)))
	}
}

// ----------------

func genDb() *Db {
	d := NewDb()
	x := Article{}
	fu := x.FilterById("abc")
	d.Add(Article{
		Id:    "a1",
		Pub:   true,
		Date:  234234234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}, fu)

	d.Add(Article{
		Id:    "a2",
		Pub:   false,
		Date:  234234234232,
		Title: "Two article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"blog", "april"},
	}, fu)

	d.Add(Article{
		Id:    "a3",
		Pub:   true,
		Date:  234234234239,
		Title: "3_article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"blog", "april"},
	}, fu)

	d.Add(Article{
		Id:    "a4",
		Pub:   true,
		Date:  234234234234,
		Title: "4_article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"blog", "may"},
	}, fu)

	return d
}

// --------------

type Article struct {
	Id    string
	Pub   bool
	Date  int
	Title string
	Desc  string
	Text  string
	Tags  []string
}

func (_ *Article) SortDate(a interface{}) int {
	return a.(Article).Date
}

func (_ *Article) FilterPub(a interface{}) bool {
	return a.(Article).Pub
}

func (_ *Article) FilterButById(id string) func(a interface{}) bool {
	fu := func(a interface{}) bool {
		return a.(Article).Id != id
	}
	return fu
}

func (_ *Article) FilterById(id string) func(a interface{}) bool {
	fu := func(a interface{}) bool {
		return a.(Article).Id == id
	}
	return fu
}

func (_ *Article) FilterDate(timeLimit int) func(a interface{}) bool {
	fu := func(a interface{}) bool {
		return a.(Article).Date < timeLimit
	}
	return fu
}

func (_ *Article) MapDateMinus(decr int) func(a interface{}) interface{} {
	fu := func(a interface{}) interface{} {
		a1 := a.(Article)
		a1.Date -= decr
		return a1
	}
	return fu
}
