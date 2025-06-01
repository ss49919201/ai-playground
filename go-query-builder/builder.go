package main

type builder struct {
	expressions []Expr
}

type ExprEqID struct {
	arg int
}

func (e ExprEqID) Query() string {
	return "id = ?"
}

func (e ExprEqID) Args() []any {
	return []any{e.arg}
}

func EqID(v int) ExprEqID {
	return ExprEqID{v}
}

type ExprEqName struct {
	arg string
}

func (e ExprEqName) Query() string {
	return "name = ?"
}

func (e ExprEqName) Args() []any {
	return []any{e.arg}
}

func EqName(v string) ExprEqName {
	return ExprEqName{v}
}

type ExprAnd struct {
	expressions []Expr
}

func (e ExprAnd) Query() string {
	var query string
	for i, expression := range e.expressions {
		if i > 0 {
			query += " and "
		}

		query += expression.Query()

		if i == len(e.expressions)-1 {
			query = "(" + query + ")"
		}
	}
	return query
}

func (e ExprAnd) Args() []any {
	args := make([]any, 0)
	for _, expression := range e.expressions {
		args = append(args, expression.Args()...)
	}
	return args
}

func And(expressions ...Expr) *ExprAnd {
	return &ExprAnd{
		expressions: expressions,
	}
}

type Expr interface {
	Query() string
	Args() []any
}

func (b *builder) Where(expressions ...Expr) *builder {
	b.expressions = append(b.expressions, expressions...)
	return b
}

func (b *builder) Build() (string, []any) {
	var (
		query string = "select * from posts"
		args  []any
	)

	if len(b.expressions) > 0 {
		query += " where "

		for i, expr := range b.expressions {
			if i > 0 {
				query += " and "
			}

			query += expr.Query()
			args = append(args, expr.Args()...)
		}
	}

	return query, args
}

func NewBuilder() *builder {
	return &builder{}
}
