package main

import (
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	{
		builder := NewBuilder()
		q, a := builder.Where(EqID(1), EqName("name")).Build()
		expect := q == "select * from posts where id = ? and name = ?" && reflect.DeepEqual(a, []any{1, "name"})
		if !expect {
			t.Errorf(q, a)
		}
	}
	{
		builder := NewBuilder()
		q, a := builder.Where(EqID(1), And(EqID(2), EqName("name"))).Build()
		expect := q == "select * from posts where id = ? and (id = ? and name = ?)" && reflect.DeepEqual(a, []any{1, 2, "name"})
		if !expect {
			t.Errorf(q, a)
		}
	}
	// {
	// 	builder := NewBuilder()
	// 	q, a := builder.Where(And(EqID(1), EqName("name")), Or(EqID(2))).And(EqID(3)).Build()
	// 	expect := q == "select * from posts where ((id = ? and name = ?) or id = ?) and id = ?" && reflect.DeepEqual(a, []any{1, "name", 2, 3})
	// 	if !expect {
	// 		t.Errorf("")
	// 	}
	// }
}
