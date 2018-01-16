package util

import "testing"

func TestFldToCol(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"Name", "name"},
		{"FooBarBaz", "foo_bar_baz"},
		{"fooBar", "foo_bar"},
		{"API", "api"},
		{"UserID", "user_id"},
		{"FooAPIDoc", "foo_api_doc"},
		{"Account1", "account1"},
		{"A1B2", "a1_b2"},
	}

	for _, test := range tests {
		got := FldToCol(test.name)
		if got != test.want {
			t.Errorf("[%s]: want: %s, got: %s", test.name, test.want, got)
		}
	}
}

func TestColToFld(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"name", "Name"},
		{"foo_bar_baz", "FooBarBaz"},
		{"account1", "Account1"},
		{"a1_b2", "A1B2"},

		// We can't restore initialisms...
		{"api", "Api"},
		{"user_id", "UserId"},
	}

	for _, test := range tests {
		got := ColToFld(test.name)
		if got != test.want {
			t.Errorf("[%s]: want: %s, got: %s", test.name, test.want, got)
		}
	}
}
