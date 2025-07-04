package dashboard_test

import (
	"github.com/NoamFav/Zvezda/src/dashboard"
	"testing"
)

func TestInfoModel_View(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repo dashboard.RepoModel
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := dashboard.NewInfoModel(tt.repo)
			got := m.View()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("View() = %v, want %v", got, tt.want)
			}
		})
	}
}
