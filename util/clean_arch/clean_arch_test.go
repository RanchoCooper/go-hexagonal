package clean_arch

import (
	"reflect"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/12/30
 */

var layersAliases = map[string]Layer{
	// Domain
	"domain":                          LayerDomain,
	"domain/aggregate":                LayerDomain,
	"domain/event":                    LayerDomain,
	"domain/model":                    LayerDomain,
	"domain/repo":                     LayerDomain,
	"domain/service":                  LayerDomain,
	"domain/vo":                       LayerDomain,
	"adapter/repository/mysql/entity": LayerInfrastructure,

	// Application
	"application": LayerApplication,
	"app":         LayerApplication,

	// Interfaces
	"interfaces": LayerInterfaces,
	"interface":  LayerInterfaces,
	"api":        LayerInterfaces,
	"api/dto":    LayerInterfaces,
	"api/grpc":   LayerInterfaces,
	"api/http":   LayerInterfaces,

	// Infrastructure
	"infrastructure":           LayerInfrastructure,
	"adapters":                 LayerInfrastructure,
	"adapter":                  LayerInfrastructure,
	"adapter/repository":       LayerInfrastructure,
	"adapter/repository/mysql": LayerInfrastructure,
	"adapter/repository/redis": LayerInfrastructure,
}

func TestValidator_Validate(t *testing.T) {
	config.Init()
	log.Init()

	aliases := make(map[string]Layer)
	for alias, layer := range layersAliases {
		aliases[alias] = layer
	}

	ignoredPackages := []string{"cmd", "config", "util"}

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
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/service/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/model/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/aggregate/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/repo/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/event/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},
		{
			Path: "/go-hexagonal/vo/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerDomain,
			},
		},

		// application layer
		{
			Path: "/go-hexagonal/application/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerApplication,
			},
		},
		{
			Path: "/go-hexagonal/app/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerApplication,
			},
		},

		// interface layer
		{
			Path: "/go-hexagonal/interfaces/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},
		{
			Path: "/go-hexagonal/interface/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},
		{
			Path: "/go-hexagonal/dto/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},
		{
			Path: "/go-hexagonal/dto/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},
		{
			Path: "/go-hexagonal/grpc/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},
		{
			Path: "/go-hexagonal/http/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInterfaces,
			},
		},

		// infrastructure layer
		{
			Path: "/go-hexagonal/infrastructure/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
			},
		},
		{
			Path: "/go-hexagonal/adapters/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
			},
		},
		{
			Path: "/go-hexagonal/adapter/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
			},
		},
		{
			Path: "/go-hexagonal/repository/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
			},
		},
		{
			Path: "/go-hexagonal/mysql/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
			},
		},
		{
			Path: "/go-hexagonal/entity/file.go",
			ExpectedFileMetadata: LayerMetadata{
				Module: "go-hexagonal",
				Layer:  LayerInfrastructure,
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
