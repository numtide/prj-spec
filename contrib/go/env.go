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
// SetAll sets all environment variables in the current environment
func SetAll() {
    Reload()
    os.Setenv("PRJ_ROOT", PrjRoot)
    os.Setenv("PRJ_ID", PrjId)
    os.Setenv("PRJ_CONFIG_HOME", PrjConfigHome)
    os.Setenv("PRJ_RUNTIME_DIR", PrjRuntimeDir)
    os.Setenv("PRJ_CACHE_HOME", PrjCacheHome)
    os.Setenv("PRJ_DATA_HOME", PrjDataHome)
    os.Setenv("PRJ_STATE_HOME", PrjStateHome)
    os.Setenv("PRJ_PATH", PrjPath)
}

// DataFile returns a suitable location for the specified data file.
// The relPath parameter must contain the name of the data file, and
// optionally, a set of parent directories (e.g. appname/app.data).
// If the specified directories do not exist, they will be created relative
// to the base data directory. On failure, an error containing the
// attempted paths is returned.
func DataFile(relPath string) (string, error) {
    	Reload()
	return pathutil.Create(relPath, []string{PrjDataHome})
}

// ConfigFile returns a suitable location for the specified config file.
// The relPath parameter must contain the name of the config file, and
// optionally, a set of parent directories (e.g. appname/app.yaml).
// If the specified directories do not exist, they will be created relative
// to the base config directory. On failure, an error containing the
// attempted paths is returned.
func ConfigFile(relPath string) (string, error) {
    	Reload()
	return pathutil.Create(relPath, []string{PrjConfigHome})
}

// StateFile returns a suitable location for the specified state file. State
// files are usually volatile data files, not suitable to be stored relative
// to the $XDG_DATA_HOME directory.
// The relPath parameter must contain the name of the state file, and
// optionally, a set of parent directories (e.g. appname/app.state).
// If the specified directories do not exist, they will be created relative
// to the base state directory. On failure, an error containing the
// attempted paths is returned.
func StateFile(relPath string) (string, error) {
    	Reload()
	return pathutil.Create(relPath, []string{PrjStateHome})
}

// CacheFile returns a suitable location for the specified cache file.
// The relPath parameter must contain the name of the cache file, and
// optionally, a set of parent directories (e.g. appname/app.cache).
// If the specified directories do not exist, they will be created relative
// to the base cache directory. On failure, an error containing the
// attempted paths is returned.
func CacheFile(relPath string) (string, error) {
    	Reload()
	return pathutil.Create(relPath, []string{PrjCacheHome})
}

// RuntimeFile returns a suitable location for the specified runtime file.
// The relPath parameter must contain the name of the runtime file, and
// optionally, a set of parent directories (e.g. appname/app.pid).
// If the specified directories do not exist, they will be created relative
// to the base runtime directory. On failure, an error containing the
// attempted paths is returned.
func RuntimeFile(relPath string) (string, error) {
    	Reload()
	return pathutil.Create(relPath, []string{PrjRuntimeDir})
}
