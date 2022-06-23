package attributepath_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/attributepath"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestToString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in  *tftypes.AttributePath
		exp string
	}

	testCases := map[string]testCase{
		"only-attribute-names": {
			in:  tftypes.NewAttributePath().WithAttributeName("foo").WithAttributeName("bar").WithAttributeName("baz"),
			exp: "foo.bar.baz",
		},
		"with-element-key-string": {
			in:  tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyString("bar").WithAttributeName("baz"),
			exp: `foo["bar"].baz`,
		},
		"with-element-key-int": {
			in:  tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyInt(10).WithAttributeName("baz"),
			exp: `foo[10].baz`,
		},
		"with-element-key-value": {
			in:  tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyInt(10).WithElementKeyValue(tftypes.NewValue(tftypes.Object{}, nil)),
			exp: `foo[10][tftypes.Object[]<null>]`,
		},
		"with-element-key-string-and-int": {
			in:  tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyString("bar").WithElementKeyInt(10),
			exp: `foo["bar"][10]`,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := attributepath.ToString(test.in)
			if test.exp != actual {
				t.Fatalf("expected %q, got %q", test.exp, actual)
			}
		})
	}
}

func TestJoinToString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in  []*tftypes.AttributePath
		exp string
	}

	testCases := map[string]testCase{
		"only-attribute-names": {
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo").WithAttributeName("bar").WithAttributeName("baz"),
				tftypes.NewAttributePath().WithAttributeName("bob").WithAttributeName("alice"),
			},
			exp: `foo.bar.baz,bob.alice`,
		},
		"with-element-key-string": {
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("bob").WithElementKeyString("alice"),
				tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyString("bar").WithAttributeName("baz"),
			},
			exp: `bob["alice"],foo["bar"].baz`,
		},
		"with-element-key-int": {
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
				tftypes.NewAttributePath().WithAttributeName("bar").WithElementKeyInt(10),
				tftypes.NewAttributePath().WithAttributeName("baz").WithElementKeyInt(100),
			},
			exp: `foo,bar[10],baz[100]`,
		},
		"with-element-key-value": {
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyValue(tftypes.NewValue(tftypes.Object{}, nil)),
				tftypes.NewAttributePath().WithAttributeName("bob").WithElementKeyString("alice"),
				tftypes.NewAttributePath().WithAttributeName("baz"),
			},
			exp: `foo[tftypes.Object[]<null>],bob["alice"],baz`,
		},
		"with-element-key-string-and-int": {
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("bob").WithAttributeName("alice"),
				tftypes.NewAttributePath().WithAttributeName("foo").WithElementKeyInt(10),
				tftypes.NewAttributePath().WithAttributeName("bar").WithElementKeyString("baz"),
			},
			exp: `bob.alice,foo[10],bar["baz"]`,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := attributepath.JoinToString(test.in...)
			if test.exp != actual {
				t.Fatalf("expected %q, got %q", test.exp, actual)
			}
		})
	}
}
