package clean_arch

import (
	"reflect"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util"
	"go-hexagonal/util/log"
)

var layersAliases = map[string]Layer{
	// Domain
	"domain": LayerDomain,

	// Application
	"application": LayerApplication,

	// Interfaces
	"interface": LayerInterfaces,
	"api":       LayerInterfaces,

	// Infrastructure
	"infrastructure": LayerInfrastructure,
	"adapter":        LayerInfrastructure,
}

func TestValidator_Validate(t *testing.T) {
	config.Init("../../config", "config")
	log.Init()

	aliases := make(map[string]Layer)
	for alias, layer := range layersAliases {
		aliases[alias] = layer
	}

	ignoredPackages := []string{"cmd", "config", "util", "tests"}

	root := util.GetProjectRootPath()
	log.SugaredLogger.Infof("[Clean Arch] start checking, root: %s", root)

	validator := NewValidator(aliases)
	count, isValid, _, err := validator.Validate(root, true, ignoredPackages)
	if err != nil {
		panic(err)
	}

	log.SugaredLogger.Infof("[Clean Arch] scaned %d files", count)
	if isValid {
		log.Logger.Info("[Clean Arch] Good Job!")
	} else {
		log.Logger.Warn("[Clean Arch] your arch is not clean enough")
	}
}

func TestParseLayerMetadata(t *testing.T) {
	testCases := []struct {
		Path                 string
		ExpectedFileMetadata LayerMetadata
	}{
		// domain layer
		{
			Path: "/go-hexagonal/domain/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "domain",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/domain/sub-package/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "domain",
				Layer:  LayerDomain,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Path, func(t *testing.T) {
			metadata := ParseLayerMetadata(c.Path, layersAliases)

			if !reflect.DeepEqual(metadata, c.ExpectedFileMetadata) {
				t.Errorf("invalid metadata: %+v, expected %+v", metadata, c.ExpectedFileMetadata)
			}
		})
	}
}
