package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Resource represents a Terraform resource.
type Resource struct {
	Type       string
	Name       string
	Labels     []string
	Attributes map[string]any
}

// Module represents a Terraform module.
type Module struct {
	Source     string
	Labels     []string
	Attributes map[string]any
}

type Local struct {
	Attributes map[string]any
}

// Config represents the Terraform configuration.
type Config struct {
	Resources []*Resource
	Modules   []*Module
	Locals    []*Local
}

func Parse(directory string) (Config, error) {
	config := Config{}

	parser := hclparse.NewParser()

	// Walk through all .tf files in the directory.
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			file, diags := parser.ParseHCLFile(path)
			if diags.HasErrors() {
				return fmt.Errorf("failed to load config file %s: %v", path, diags.Errs())
			}

			parsedConfig := parseConfig(file)

			config.Modules = append(config.Modules, parsedConfig.Modules...)
			config.Resources = append(config.Resources, parsedConfig.Resources...)
			config.Locals = append(config.Locals, parsedConfig.Locals...)
		}

		return nil
	})

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseConfig(file *hcl.File) Config {
	resources := make([]*Resource, 0)
	modules := make([]*Module, 0)
	locals := make([]*Local, 0)

	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		switch block.Type {
		case "module":
			modules = append(modules, parseModule(block))
		case "resource":
			resources = append(resources, parseResource(block))
		case "locals":
			locals = append(locals, parseLocals(block))
		}
	}

	return Config{Resources: resources, Modules: modules, Locals: locals}
}

func parseModule(block *hclsyntax.Block) *Module {
	module := &Module{
		Labels:     block.Labels,
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		module.Attributes[attribute.Name] = value

		if attribute.Name == "source" {
			module.Source = value.(string)
		}
	}

	return module
}

func parseResource(block *hclsyntax.Block) *Resource {
	resource := &Resource{
		Type:       block.Labels[0],
		Name:       block.Labels[1],
		Labels:     block.Labels,
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		resource.Attributes[attribute.Name] = value
	}

	return resource
}

func parseLocals(block *hclsyntax.Block) *Local {
	local := &Local{
		Attributes: map[string]any{},
	}

	for _, attribute := range block.Body.Attributes {
		value := evaluateExpression(attribute.Expr)
		local.Attributes[attribute.Name] = value
	}

	return local
}

func buildVarExpressions(traversal hcl.Traversal) string {
	varExp := make([]string, 0, len(traversal))

	for _, v := range traversal {
		switch v := v.(type) {
		case hcl.TraverseRoot:
			if v.Name != "" {
				varExp = append(varExp, v.Name)
			}
		case hcl.TraverseAttr:
			if v.Name != "" {
				varExp = append(varExp, v.Name)
			}
		}
	}

	return strings.Join(varExp, ".")
}

func convertValueToString(val cty.Value) string {
	switch val.Type() {
	case cty.Number:
		return val.AsBigFloat().String()
	case cty.String:
		return val.AsString()
	case cty.Bool:
		var v bool
		_ = gocty.FromCtyValue(val, &v)

		return fmt.Sprintf("%v", v)
	default:
		fmt.Println("Unsupported type:", val.Type().GoString())
		return ""
	}
}

// evaluateExpression evaluates the HCL expression and returns its value as a string or map[string]string.
func evaluateExpression(expr hcl.Expression) any {
	resultString := ""
	resultMap := map[string]any{}

	switch expr := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		resultString += buildVarExpressions(expr.Traversal)
	case *hclsyntax.LiteralValueExpr:
		resultString += convertValueToString(expr.Val)
	case *hclsyntax.TemplateExpr:
		parts := expr.Parts
		for _, part := range parts {
			resultString += evaluateExpression(part).(string)
		}
	case *hclsyntax.TupleConsExpr:
		for _, elem := range expr.Exprs {
			resultString += evaluateExpression(elem).(string) + ","
		}
	case *hclsyntax.ObjectConsKeyExpr:
		resultString += evaluateExpression(expr.Wrapped).(string)
	case *hclsyntax.ObjectConsExpr:
		for i := range expr.Items {
			item := expr.Items[i]

			resultMap[evaluateExpression(item.KeyExpr).(string)] = evaluateExpression(item.ValueExpr)
		}

		return resultMap
	case *hclsyntax.FunctionCallExpr: // TODO: Implement
	}

	return resultString
}
