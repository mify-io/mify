package generate

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

type GeneratorLanguage string

const (
	GENERATOR_LANGUAGE_GO GeneratorLanguage = "go"
)

type OpenAPIGenerator struct {
	schemaDir string
	language GeneratorLanguage
}

func NewOpenAPIGenerator(schemaDir string, language GeneratorLanguage) OpenAPIGenerator {
	return OpenAPIGenerator{
		schemaDir: schemaDir,
		language: language,
	}
}

func (g *OpenAPIGenerator) GenerateServer(outputDir string) error {
	// TODO: maybe pass context from caller
	// ctx := context.Background()
	// loader := &openapi3.Loader{
		// Context: ctx,
		// IsExternalRefsAllowed: true,
	// }
	// doc, err := loader.LoadFromFile(g.schemaPath)
	// if err != nil {
		// return err
	// }
	// docCopy := doc
	// enrichedPaths := doc.Paths
	// for pathName, path := range doc.Paths {
		// fmt.Printf("path: %s %+v\n", pathName, path)
		// ops := path.Operations()
		// for method, op := range ops {
			// if op == nil {
				// continue
			// }
			// newOp := *op
			// newOp.Tags = []string {pathName}
			// enrichedPaths[pathName].SetOperation(method, &newOp)
		// }
		// // doc.Paths[k].Set
	// }
	// docCopy.Paths = enrichedPaths


	// // https://github.com/getkin/kin-openapi/issues/241
	// jsonData, err := docCopy.MarshalJSON()
	// if err != nil {
		// return fmt.Errorf("failed to create api yaml: %w", err)
	// }
	// tmp := map[string]interface{}{}
	// err = json.Unmarshal(jsonData, &tmp)
	// if err != nil {
		// return fmt.Errorf("failed to create api yaml: %w", err)
	// }

	// f, err := os.Create(fmt.Sprintf("/tmp/api.yaml"))
	// if err != nil {
		// return fmt.Errorf("failed to create api yaml: %w", err)
	// }

	// err = yaml.NewEncoder(f).Encode(tmp)
	// if err != nil {
		// return fmt.Errorf("failed to create api yaml: %w", err)
	// }

	// yaml.NewEncoder

	// data, err := yaml.Marshal(&docCopy)

	// err = ioutil.WriteFile(, data, 0644)
	// if err != nil {
		// return fmt.Errorf("failed to create api yaml: %w", err)
	// }

	schemaPath, err := g.makeEnrichedSchema()
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	err = g.doGenerate(schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// read generator config and templates
	// run subprocess with docker
	// perform language-specific post generation
	return nil

}

// private


func (g *OpenAPIGenerator) makeEnrichedSchema() (string, error) {
	schemaPath := g.schemaDir+"/api.yaml"


	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %s: %w", schemaPath, err)
	}
	// TODO: maybe pass context from caller
	ctx := context.Background()
	loader := &openapi3.Loader{
		Context: ctx,
		IsExternalRefsAllowed: true,
	}
	url, err := url.Parse(schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}

	openapiDoc, err := loader.LoadFromDataWithPath(data, url)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}
	err = openapiDoc.Validate(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}

	doc := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("failed to parse schema: %s: %w", schemaPath, err)
	}

	pathsIface, ok := doc["paths"]
	if !ok {
		return "", fmt.Errorf("missing paths in schema: %s", schemaPath)
	}
	// TODO mapstructure
	paths := pathsIface.(map[interface{}]interface{})
	for path, v := range paths {
		fmt.Printf("debug: processing path: %s\n", path)
		methods := v.(map[interface{}]interface{})
		if _, ok := methods["$ref"]; ok {
			return "", fmt.Errorf("paths with $ref are not supported yet")
		}
		for m, vv := range methods {
			fmt.Printf("debug: processing method: %s\n", m)
			method := vv.(map[interface{}]interface{})
			method["tags"] = []string{path.(string)}
			methods[m] = method
		}
	}

	cacheDir := config.GetCacheDirectory()
	fmt.Printf("debug: cache dir: %s\n", cacheDir)
	targetDir := cacheDir+"/"+g.schemaDir

	err = copy.Copy(g.schemaDir, targetDir, copy.Options{
		OnDirExists: func(src, dest string) copy.DirExistsAction {
			return copy.Replace
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to prepare temp api schema: %w", err)
	}


	targetPath := targetDir+"/api.yaml"
	f, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create api yaml: %w", err)
	}

	err = yaml.NewEncoder(f).Encode(doc)
	if err != nil {
		return "", fmt.Errorf("failed to create api yaml: %w", err)
	}

	return targetPath, nil
}

func (g *OpenAPIGenerator) doGenerate(schemaPath string, targetPath string) error {
	path, err := config.DumpAssets("openapi/server-template", "openapi")
	if err != nil {
		return err
	}
	fmt.Printf("debug: dumped path: %s\n", path)


	// TODO: maybe use provided path?
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(targetPath, "generated")

	err = runOpenapiGenerator(curDir, schemaPath, filepath.Join(path, "server-template"), generatedPath)
	if err != nil {
		return err
	}

//docker run --rm \
//  --user $(id -u ${USER}):$(id -g ${GROUP}) \
//  -v ${PWD}:/local openapitools/openapi-generator-cli "$@"


//mkdir -p out/go
//cp tpl/ignore-list.txt out/go/.openapi-generator-ignore
//./tool.sh generate -c /local/config.yaml -o /local/out/go
	return nil
}

func runOpenapiGenerator(basePath string, schemaPath string, templatePath string, targetDir string) error {
	curUser, err := user.Current()
	if err != nil {
		return err
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	err = copyFile(
		filepath.Join(templatePath, "ignore-list.txt"),
		filepath.Join(targetDir, ".openapi-generator-ignore"),
	)
	if err != nil {
		return err
	}

	args := [] string{
		"run",
		"--rm",
		"--user", curUser.Uid+":"+curUser.Gid,
		"-v", basePath+":/repo",
		"openapitools/openapi-generator-cli",
		"generate",
		"-c", filepath.Join("/repo", templatePath, "config.yaml"),
		"-i", filepath.Join("/repo", schemaPath),
		"-o", filepath.Join("/repo", targetDir),
	}
	fmt.Printf("debug: running docker %s\n", args)

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	// TODO only if verbose
	fmt.Printf("%s", output)
	if err != nil {
		return err
	}
	fmt.Printf("debug: successfully generated openapi\n")

	return nil
}

func copyFile(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0644)
}
