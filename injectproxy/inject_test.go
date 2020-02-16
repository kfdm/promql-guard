package injectproxy

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/testutil"
)

func TestInject(t *testing.T) {
	var err error

	matchers := make([]*labels.Matcher, 0)
	matcher, _ := labels.NewMatcher(labels.MatchEqual, "service", "tenant")
	matchers = append(matchers, matcher)

	expr, err := promql.ParseExpr("a - b")
	testutil.Ok(t, err)

	err = InjectMatchers(expr, matchers)
	testutil.Ok(t, err)

	testutil.Equals(t, expr.String(), "a{service=\"tenant\"} - b{service=\"tenant\"}")
}
