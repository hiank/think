package sys

var (
	Export_newLoader = func(path string, lt int) *export_Loader {
		return &export_Loader{Loader: &Loader{path: path, lt: lt}}
	}
)

type export_Loader struct {
	*Loader
}

func (el *export_Loader) Match() []string {
	return el.match()
}

func (el *export_Loader) ListPaths(dpath string) ([]string, []string) {
	return el.listPaths(dpath)
}