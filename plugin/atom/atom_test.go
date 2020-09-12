package atom

import (
	"fmt"
	"github.com/spf13/afero"
	"testing"

	"github.com/verless/verless/model"
	"github.com/verless/verless/test"
)

var (
	// testPages is a set of pages used for testing.
	testPages = []model.Page{
		{ID: "page-0", Route: "/route-0", Title: "Page 1"},
		{ID: "page-1", Route: "/route-1/route-22/route-333", Title: "Page 2"},
		{ID: "page-2", Route: "/route-2", Title: "Page 3"},
		{ID: "page-3", Route: "/route-3", Title: "Page 4"},
	}
)

// TestAtom_ProcessPage checks if the atom plugin creates a new
// RSS feed item for each processed page.
func TestAtom_ProcessPage(t *testing.T) {
	tests := map[string]struct {
		pages         []model.Page
		expectedError error
	}{
		"normal pages": {
			pages: testPages,
		},
	}

	for name, testCase := range tests {
		t.Log(name)

		a := New(&model.Meta{
			Base: "https://example.com",
		}, afero.NewOsFs(), "")

		for i, page := range testCase.pages {
			t.Logf("process page number %v, route '%v'", i, page.Route)
			err := a.ProcessPage(&page)
			if test.ExpectedError(t, testCase.expectedError, err) != test.IsCorrectNil {
				return
			}

			item := a.feed.Items[i]
			test.Equals(t, page.Title, item.Title)

			canonicalLink := fmt.Sprintf("%s%s/%s", a.meta.Base, page.Route, page.ID)
			test.Equals(t, canonicalLink, item.Link.Href)
		}

		test.Equals(t, len(testCase.pages), len(a.feed.Items))
	}
}
