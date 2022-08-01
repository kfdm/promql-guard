package injectproxy

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

type Injector struct {
	matchers []*labels.Matcher
}

func InjectMatchers(expr parser.Expr, matchers []*labels.Matcher) error {
	visitor := Injector{matchers: matchers}
	return parser.Walk(visitor, expr, nil)
}

func (i Injector) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	switch n := node.(type) {
	// case *parser.MatrixSelector:
	// 	// inject labelselector
	// 	n.LabelMatchers = append(n.LabelMatchers, i.matchers...)

	case *parser.VectorSelector:
		// inject labelselector
		n.LabelMatchers = append(n.LabelMatchers, i.matchers...)
	}
	return i, nil
}
