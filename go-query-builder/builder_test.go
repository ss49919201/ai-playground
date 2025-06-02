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
	{
		builder := NewBuilder()
		q, a := builder.Where(EqID(1), And(And(EqID(2), EqID(3)), EqName("name"))).Build()
		expect := q == "select * from posts where id = ? and ((id = ? and id = ?) and name = ?)" && reflect.DeepEqual(a, []any{1, 2, 3, "name"})
		if !expect {
			t.Errorf(q, a)
		}
	}
}
