package srcmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
)

type SourceMapFS struct {
	fs fs.FS
}

func NewSourceMapFS(fs fs.FS) *SourceMapFS {
	return &SourceMapFS{fs: fs}
}

type SourceID uint64

func (id *SourceID) UnmarshalText(data []byte) error {
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return err
	}
	*id = SourceID(v)
	return nil
}

type ForgeBuild struct {
	ID             string              `json:"id"`                // ID of the build itself
	SourceIDToPath map[SourceID]string `json:"source_id_to_path"` // srcmap ID to source filepath
}

func (s *SourceMapFS) readBuild(id string) (*ForgeBuild, error) {
	buildPath := path.Join("artifacts", "build-info", id+".json")
	f, err := s.fs.Open(buildPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open build: %w", err)
	}
	defer f.Close()
	var build ForgeBuild
	if err := json.NewDecoder(f).Decode(&build); err != nil {
		return nil, fmt.Errorf("failed to read build: %w", err)
	}
	return &build, nil
}

type BuildEntry struct {
	Path    string `json:"path"`
	BuildID string `json:"build_id"`
}

type BuildInfo struct {
	// contract name -> solidity version -> build entry
	Artifacts map[string]map[string]BuildEntry `json:"artifacts"`
}

type SolidityBuildCache struct {
	Files map[string]BuildInfo
}

func (s *SourceMapFS) readBuildCache() (*SolidityBuildCache, error) {
	cachePath := path.Join("cache", "solidity-files-cache.json")
	f, err := s.fs.Open(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open build cache: %w", err)
	}
	defer f.Close()
	var buildCache SolidityBuildCache
	if err := json.NewDecoder(f).Decode(&buildCache); err != nil {
		return nil, fmt.Errorf("failed to read build cache: %w", err)
	}
	return &buildCache, nil
}

func (s *SourceMapFS) ReadSourceIDs(path string, contract string, compilerVersion string) (map[SourceID]string, error) {
	buildCache, err := s.readBuildCache()
	if err != nil {
		return nil, err
	}
	artifactBuilds, ok := buildCache.Files[path]
	if !ok {
		return nil, fmt.Errorf("no known builds for path %q", path)
	}
	byCompilerVersion, ok := artifactBuilds.Artifacts[contract]
	if !ok {
		return nil, fmt.Errorf("contract not found in artifact: %q", contract)
	}
	var buildEntry BuildEntry
	if compilerVersion != "" {
		entry, ok := byCompilerVersion[compilerVersion]
		if !ok {
			return nil, fmt.Errorf("no known build for compiler version: %q", compilerVersion)
		}
		buildEntry = entry
	} else {
		if len(byCompilerVersion) == 0 {
			return nil, errors.New("no known build, unspecified compiler version")
		}
		if len(byCompilerVersion) > 1 {
			return nil, fmt.Errorf("no compiler version specified, and more than one option: %s", strings.Join(maps.Keys(byCompilerVersion), ", "))
		}
		for _, entry := range byCompilerVersion {
			buildEntry = entry
		}
	}
	build, err := s.readBuild(buildEntry.BuildID)
	if err != nil {
		return nil, fmt.Errorf("failed to read build %q of contract %q: %w", buildEntry.BuildID, contract, err)
	}
	return build.SourceIDToPath, nil
}

func (s *SourceMapFS) SourceMap(artifact *foundry.Artifact, contract string) (*SourceMap, error) {
	srcPath := ""
	for path, name := range artifact.Metadata.Settings.CompilationTarget {
		if name == contract {
			srcPath = path
			break
		}
	}
	if srcPath == "" {
		return nil, fmt.Errorf("no known source path for contract %s in artifact", contract)
	}
	basicCompilerVersion := strings.SplitN(artifact.Metadata.Compiler.Version, "+", 2)[0]
	ids, err := s.ReadSourceIDs(srcPath, contract, basicCompilerVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to read source IDs of %q: %w", srcPath, err)
	}
	return ParseSourceMap(s.fs, ids, artifact.DeployedBytecode.Object, artifact.DeployedBytecode.SourceMap)
}
