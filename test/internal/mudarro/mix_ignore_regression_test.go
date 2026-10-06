package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"testing"
)

func TestMixArtifactIgnoresDoNotHideOtherProjects(t *testing.T) {
	for _, dir := range []string{"deps", "_build"} {
		t.Run(dir, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": `{"workspaces":["` + dir + `/*"]}`, dir + "/app/package.json": `{"scripts":{"test":"never"}}`})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(r.Config.Services) != 2 {
				t.Fatalf("unrelated workspace member lost: %+v", r)
			}
			if jsWorkspaceService(t, r, dir+"/app").WorkspaceRoot != "." {
				t.Fatal("ownership lost")
			}
		})
	}
}
