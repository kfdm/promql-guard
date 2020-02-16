package injectproxy

import (
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

type Injector struct {
	matchers []*labels.Matcher
}

func InjectMatchers(expr promql.Expr, matchers []*labels.Matcher) error {
	visitor := Injector{matchers: matchers}
	return promql.Walk(visitor, expr, nil)
}

func (i Injector) Visit(node promql.Node, path []promql.Node) (promql.Visitor, error) {
	switch n := node.(type) {
	case *promql.MatrixSelector:
		// inject labelselector
		n.LabelMatchers = append(n.LabelMatchers, i.matchers...)

	case *promql.VectorSelector:
		// inject labelselector
		n.LabelMatchers = append(n.LabelMatchers, i.matchers...)
	}
	return i, nil
}
