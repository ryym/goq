package gen

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseTag(t *testing.T) {
	type Attrs map[string]tagType
	type Result map[string]interface{}

	cases := []struct {
		tag   string
		attrs Attrs
		want  Result
		err   string
	}{
		{
			tag:   "flg",
			attrs: Attrs{"flg": tagTypeBool},
			want:  Result{"flg": true},
		},
		{
			tag: "name:bob;married;age:30",
			attrs: Attrs{
				"name":    tagTypeStr,
				"married": tagTypeBool,
				"age":     tagTypeStr,
			},
			want: Result{
				"name":    "bob",
				"married": true,
				"age":     "30",
			},
		},
		{
			tag:   "",
			attrs: Attrs{"s": tagTypeStr, "b": tagTypeBool},
			want:  Result{},
		},
		{
			tag:   "b:val",
			attrs: Attrs{"b": tagTypeBool},
			want:  nil,
			err:   "tagType unmatch",
		},
		{
			tag:   "foo",
			attrs: Attrs{"bar": tagTypeBool},
			want:  nil,
			err:   "Unknown attribute",
		},
	}

	for _, c := range cases {
		got, err := parseTag(c.tag, c.attrs)

		if err != nil {
			if c.err == "" {
				t.Errorf("Unexpected error occurred\n%s", err)
			}
		} else if c.err != "" {
			t.Errorf("Expected error did not occur\n%s", c.err)
		} else if diff := deep.Equal(Result(got), c.want); diff != nil {
			t.Errorf("Unexpected result\n%s", diff)
		}
	}
}

func TestParseColumnTag(t *testing.T) {
	got, _ := parseColumnTag("pk;name:foo_bar")
	want := columnTag{
		IsPK:    true,
		ColName: "foo_bar",
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Errorf("Unexpected model tag\n%s", diff)
	}
}

func TestParseTableTag(t *testing.T) {
	got, _ := parseTableTag("helper:Prefs")
	want := tableTag{
		HelperName: "Prefs",
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Errorf("Unexpected table tag\n%s", diff)
	}
}
