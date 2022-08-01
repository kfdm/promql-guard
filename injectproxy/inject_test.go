package injectproxy

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
)

func TestInject(t *testing.T) {
	var err error

	matchers := make([]*labels.Matcher, 0)
	matcher, _ := labels.NewMatcher(labels.MatchEqual, "service", "tenant")
	matchers = append(matchers, matcher)

	expr, err := promql.ParseExpr("a - b")
	require.NoError(t, err)

	err = InjectMatchers(expr, matchers)
	require.NoError(t, err)

	require.Equal(t, expr.String(), "a{service=\"tenant\"} - b{service=\"tenant\"}")
}
