package spec

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"

	"github.com/numtide/prj-spec/contrib/go/internal/pathutil"
)

var (
	// PrjRoot ...
	PrjRoot string
	// PrjId ...
	PrjId string
	// PrjConfigHome ...
	PrjConfigHome string
	// PrjRuntimeDir ...
	PrjRuntimeDir string
	// PrjCacheHome ...
	PrjCacheHome string
	// PrjDataHome ...
	PrjDataHome string
	// PrjStateHome ...
	PrjStateHome string
	// PrjPath ...
	PrjPath string
)

func Load() {
	Reload()
}

// Reload refreshes PRJ by reading the environment.
// Defaults are applied for XDG variables which are empty or not present
// in the environment.
func Reload() {
	// Initialize Project Root.
	PrjRoot = getEnvOr("PRJ_ROOT", prjRoot)

	// Initialize Project Config Home (PrjId depends on it).
	PrjConfigHome = getEnvOr("PRJ_CONFIG_HOME", func() string { return filepath.Join(PrjRoot, ".config") })

	// Initialize Project Id.
	PrjId = getEnvOr("PRJ_ID", func() string {
		prjIdFile := filepath.Join(PrjConfigHome, "prj_id")
		if pathutil.Exists(prjIdFile) {
			buf, err := os.ReadFile(prjIdFile)
			if err != nil {
				log.Fatalln(err)
			}
			return string(buf)
		}
		return ""
	})

	// Set other directories.
	PrjRuntimeDir = getEnvOr("PRJ_RUNTIME_DIR", func() string { return filepath.Join(PrjRoot, ".run") })
	PrjCacheHome = getEnvOr("PRJ_CACHE_HOME", func() string {
		if PrjId != "" {
			return filepath.Join(xdg.CacheHome, "prj", PrjId)
		}
		return filepath.Join(PrjRoot, ".cache")
	})
	PrjDataHome = getEnvOr("PRJ_DATA_HOME", func() string {
		if PrjId != "" {
			return filepath.Join(xdg.DataHome, "prj", PrjId)
		}
		return filepath.Join(PrjRoot, ".local", "share")
	})
	PrjStateHome = getEnvOr("PRJ_STATE_HOME", func() string {
		if PrjId != "" {
			return filepath.Join(xdg.StateHome, "prj", PrjId)
		}
		return filepath.Join(PrjRoot, ".local", "state")
	})
	PrjPath = getEnvOr("PRJ_PATH", func() string { return filepath.Join(PrjRoot, ".local", "bin") })

}

func getEnvOr(v string, f func() string) string {
	if path := os.Getenv(v); path != "" {
		return path
	}
	return f()
}

func prjRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	var searchPaths []string
	searchPaths = append(searchPaths, dir)

	for {
    		old := dir
		dir = filepath.Dir(dir)
		if dir == old {
			break
		}
		searchPaths = append(searchPaths, dir)
	}
	p, err := pathutil.Search(".config", searchPaths)
	if err != nil {
		log.Fatalln(err)
	}
	return p
}
