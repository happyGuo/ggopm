package action

import (
	"os"
	//"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/dependency"
	"github.com/Masterminds/glide/msg"
	gpath "github.com/Masterminds/glide/path"
	"github.com/Masterminds/glide/util"
	//"fmt"
)


func Create(base string, skipScan bool) {
	//文件路径
	config := gpath.GlideFile
	// 检查是否存在
	//guardYAML(config)

	//跳过扫描，生成空配置文件
	if skipScan {

	}
	// 检查依赖
	conf := guessDeps(base, skipScan)

	msg.Info("写入配置文件 (%s)", config)
	if err := conf.WriteFile(config); err != nil {
		msg.Die("写入失败 %s: %s", config, err)
	}

}

// guardYAML fails if the given file already exists.
//
// This prevents an important file from being overwritten.
func guardYAML(filename string) {
	if _, err := os.Stat(filename); err == nil {
		msg.Die("Cowardly refusing to overwrite existing YAML.")
	}
}

// guessDeps attempts to resolve all of the dependencies for a given project.
//
// base is the directory to start with.
// skipImport will skip running the automatic imports.
//
// FIXME: This function is likely a one-off that has a more standard alternative.
// It's also long and could use a refactor.
func guessDeps(base string, skipImport bool) *cfg.Config {
	buildContext, err := util.GetBuildContext()
	if err != nil {
		msg.Die("Failed to build an import context: %s", err)
	}
	name := buildContext.PackageName(base)


	msg.Info("Generating a YAML configuration file and guessing the dependencies")

	config := new(cfg.Config)

	// Get the name of the top level package
	config.Name = name

	// Import by looking at other package managers and looking over the
	// entire directory structure.

	// Attempt to import from other package managers.
	if !skipImport {
		guessImportDeps(base, config)
	}

	importLen := len(config.Imports)
	if importLen == 0 {
		msg.Info("Scanning code to look for dependencies")
	} else {
		msg.Info("Scanning code to look for dependencies not found in import")
	}

	// 返回依赖解析器
	r, err := dependency.NewResolver(base)


	if err != nil {
		msg.Die("Error creating a dependency resolver: %s", err)
	}

	// 初始化以测试模式
	r.ResolveTest = true

	h := &dependency.DefaultMissingPackageHandler{Missing: []string{}, Gopath: []string{}}
	r.Handler = h

	sortable, testSortable, err := r.ResolveLocal(false)
	if err != nil {
		msg.Die("分析本地依赖失败: %s", err)
	}

	sort.Strings(sortable)
	sort.Strings(testSortable)

	vpath := r.VendorDir
	if !strings.HasSuffix(vpath, "/") {
		vpath = vpath + string(os.PathSeparator)
	}

	for _, pa := range sortable {
		n := strings.TrimPrefix(pa, vpath)
		root, subpkg := util.NormalizeName(n)

		if !config.Imports.Has(root) && root != config.Name {
			msg.Info("--> Found reference to %s\n", n)
			d := &cfg.Dependency{
				Name: root,
			}
			if len(subpkg) > 0 {
				d.Subpackages = []string{subpkg}
			}
			config.Imports = append(config.Imports, d)
		} else if config.Imports.Has(root) {
			if len(subpkg) > 0 {
				subpkg = strings.TrimPrefix(subpkg, "/")
				d := config.Imports.Get(root)
				if !d.HasSubpackage(subpkg) {
					msg.Info("--> Adding sub-package %s to %s\n", subpkg, root)
					d.Subpackages = append(d.Subpackages, subpkg)
				}
			}
		}
	}

	for _, pa := range testSortable {
		n := strings.TrimPrefix(pa, vpath)
		root, subpkg := util.NormalizeName(n)

		if config.Imports.Has(root) && root != config.Name {
			msg.Debug("--> Found test reference to %s already listed as an import", n)
		} else if !config.DevImports.Has(root) && root != config.Name {
			msg.Info("--> Found test reference to %s", n)
			d := &cfg.Dependency{
				Name: root,
			}
			if len(subpkg) > 0 {
				d.Subpackages = []string{subpkg}
			}
			config.DevImports = append(config.DevImports, d)
		} else if config.DevImports.Has(root) {
			if len(subpkg) > 0 {
				subpkg = strings.TrimPrefix(subpkg, "/")
				d := config.DevImports.Get(root)
				if !d.HasSubpackage(subpkg) {
					msg.Info("--> Adding test sub-package %s to %s\n", subpkg, root)
					d.Subpackages = append(d.Subpackages, subpkg)
				}
			}
		}
	}

	if len(config.Imports) == importLen && importLen != 0 {
		msg.Info("--> Code scanning found no additional imports")
	}

	return config
}

func guessImportDeps(base string, config *cfg.Config) {
	msg.Info("Attempting to import from other package managers (use --skip-import to skip)")
	deps := []*cfg.Dependency{}
	//absBase, err := filepath.Abs(base)


	for _, i := range deps {
		if i.Reference == "" {
			msg.Info("--> Found imported reference to %s", i.Name)
		} else {
			msg.Info("--> Found imported reference to %s at revision %s", i.Name, i.Reference)
		}

		config.Imports = append(config.Imports, i)
	}
}


