package clean_arch

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

type Layer string

const (
	// LayerDomain represents domain layer.
	LayerDomain Layer = "domain"

	// LayerApplication represents application layer.
	LayerApplication Layer = "application"

	// LayerInfrastructure represents infrastructure layer. aka adapters
	LayerInfrastructure Layer = "infrastructure"

	// LayerInterfaces represents interfaces layer. aka api
	LayerInterfaces Layer = "interfaces"
)

const (
	LayerDomainWeight         = 1
	LayerApplicationWeight    = 2
	LayerInterfacesWeight     = 3
	LayerInfrastructureWeight = 4
)

var layersHierarchy = map[Layer]int{
	LayerDomain:         LayerDomainWeight,
	LayerApplication:    LayerApplicationWeight,
	LayerInterfaces:     LayerInterfacesWeight,
	LayerInfrastructure: LayerInfrastructureWeight,
}

// NewValidator creates new Validator.
func NewValidator(alias map[string]Layer) *Validator {
	filesMetadata := make(map[string]LayerMetadata, 0)

	return &Validator{
		filesMetadata: filesMetadata,
		alias:         alias,
	}
}

// ValidationError represents an error when Clean Architecture rule is not keep.
type ValidationError error

// Validator is responsible for Clean Architecture validation.
type Validator struct {
	filesMetadata map[string]LayerMetadata
	alias         map[string]Layer
}

// Validate validates provided a path for Clean Architecture rules.
func (v *Validator) Validate(root string, ignoreTests bool, ignoredPackages []string) (int, bool, []ValidationError, error) {
	errors := make([]ValidationError, 0)
	count := 0

	err := filepath.Walk(root, func(path string, fi os.FileInfo, walkErr error) error {
		for _, ignored := range ignoredPackages {
			if strings.Contains(path, ignored) {
				return nil
			}
		}
		if fi.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if ignoreTests && strings.HasSuffix(path, "_test.go") {
			return nil
		}

		if strings.Contains(path, "/vendor/") {
			// Skip vendor directory - contains third-party dependencies
			return nil
		}

		if strings.Contains(path, "/.") {
			return nil
		}

		fset := token.NewFileSet()

		f, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if parseErr != nil {
			panic(parseErr)
		}

		importerMeta := v.fileMetadata(path)
		log.SugaredLogger.Infof("file: %s, metadata: %+v", path, importerMeta)
		count++

		if importerMeta.Layer == "" || importerMeta.Module == "" {
			// Unable to parse layer metadata - skip validation for this file
			log.SugaredLogger.Warnf("cannot parse metadata for file %s, meta: %+v", path, importerMeta)

			return nil
		}

	ImportsLoop:
		for _, imp := range f.Imports {
			for _, ignoredPackage := range ignoredPackages {
				if strings.Contains(imp.Path.Value, ignoredPackage) {
					continue ImportsLoop
				}
			}

			validationErrors := v.validateImport(imp, importerMeta, path)
			errors = append(errors, validationErrors...)
		}

		return nil
	})
	if err != nil {
		return 0, false, nil, err
	}

	return count, len(errors) == 0, errors, nil
}

func (v *Validator) validateImport(imp *ast.ImportSpec, importerMeta LayerMetadata, path string) []ValidationError {
	errors := make([]ValidationError, 0)

	importPath := imp.Path.Value
	importPath = strings.TrimSuffix(importPath, `"`)
	importPath = strings.TrimPrefix(importPath, `"`)
	importMeta := v.fileMetadata(importPath)

	if !strings.Contains(importPath, config.GlobalConfig.App.Name) {
		log.SugaredLogger.Debugf("[%s] filtered due to third part dependency", importPath)
		return nil
	}

	if importMeta.Layer == importerMeta.Layer {
		// pass
	} else {
		importHierarchy := layersHierarchy[importMeta.Layer]
		importerHierarchy := layersHierarchy[importerMeta.Layer]
		// log.SugaredLogger.Infof("import hierarchy: %d, importer hierarchy: %d", importHierarchy, importerHierarchy)

		if importHierarchy > importerHierarchy {
			err := fmt.Errorf(
				"anti-clean [hit-0]: %s import %s(%s) to %s",
				path, importMeta.Layer, importPath,
				importerMeta.Layer,
			)
			errors = append(errors, err)
		}
	}
	// if importMeta.Module == importerMeta.Module {
	// 	importHierarchy := layersHierarchy[importMeta.Layer]
	// 	importerHierarchy := layersHierarchy[importerMeta.Layer]
	// 	// log.SugaredLogger.Infof("import hierarchy: %d, importer hierarchy: %d", importHierarchy, importerHierarchy)
	//
	// 	if importHierarchy > importerHierarchy {
	// 		err := fmt.Errorf(
	// 			"anti-clean [hit-1]: %s, import %s (%s) to %s",
	// 			path, importMeta.Layer, importPath,
	// 			importerMeta.Layer,
	// 		)
	// 		errors = append(errors, err)
	// 	}
	// } else if importMeta.Layer != "" {
	// 	if importMeta.Layer != LayerInterfaces || importerMeta.Layer != LayerInfrastructure {
	// 		err := fmt.Errorf(
	// 			"anti-clean [hit-2]: %s imported %s, between %s and %s",
	// 			path, importPath,
	// 			importMeta.Module, importerMeta.Module,
	// 		)
	// 		errors = append(errors, err)
	// 	} else {
	// 		panic("exists unhandled case")
	// 	}
	// } else {
	// 	panic("exists unhandled case")
	// }

	if len(errors) == 0 {
		log.SugaredLogger.Infof("%s imported: %s passed ✅ (%s import %s)", path, importPath, importMeta.Layer, importerMeta.Layer)
	} else {
		for _, err := range errors {
			log.Logger.Warn(err.Error())
		}
	}

	return errors
}

func (v *Validator) fileMetadata(path string) LayerMetadata {
	if metadata, ok := v.filesMetadata[path]; ok {
		return metadata
	}

	v.filesMetadata[path] = ParseLayerMetadata(path, v.alias)

	return v.filesMetadata[path]
}

// LayerMetadata contains information about directory module and software layer.
type LayerMetadata struct {
	Module string
	Layer  Layer
}

// ParseLayerMetadata parses metadata of provided path.
func ParseLayerMetadata(path string, alias map[string]Layer) LayerMetadata {
	metadata := LayerMetadata{}

	for alia, layer := range alias {
		if strings.Contains(path, alia) {
			if metadata.Module != "" && len(layer) < len(metadata.Module) {
				continue
			}
			metadata.Layer = layer
			metadata.Module = alia
			break // we assume that one file belongs to one module
		}
	}

	return metadata
}
