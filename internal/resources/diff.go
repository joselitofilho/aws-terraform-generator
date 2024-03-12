package resources

import (
	"fmt"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
)

// FindDifferences finds the differences between two resource collections.
func FindDifferences(
	rc1, rc2 *ResourceCollection) (addedResources, removedResources []Resource,
	addedRelationships, removedRelationships []Relationship,
) {
	// Find added and removed resources.
	rc1Resources := make(map[string]struct{})
	for _, res := range rc1.Resources {
		rc1Resources[res.Value()] = struct{}{}
	}

	rc2Resources := make(map[string]struct{})
	for _, res := range rc2.Resources {
		rc2Resources[res.Value()] = struct{}{}
	}

	for _, res := range rc1.Resources {
		if _, exists := rc2Resources[res.Value()]; !exists {
			removedResources = append(removedResources, res)
		}
	}

	for _, res := range rc2.Resources {
		if _, exists := rc1Resources[res.Value()]; !exists {
			addedResources = append(addedResources, res)
		}
	}

	// Find added and removed relationships.
	for _, rel := range rc2.Relationships {
		if !containsRelationship(rc1.Relationships, rel) {
			addedRelationships = append(addedRelationships, rel)
		}
	}

	for _, rel := range rc1.Relationships {
		if !containsRelationship(rc2.Relationships, rel) {
			removedRelationships = append(removedRelationships, rel)
		}
	}

	return addedResources, removedResources, addedRelationships, removedRelationships
}

// PrintDiff prints the differences between two resource collections.
func PrintDiff(rc1, rc2 *ResourceCollection) {
	addedResources, removedResources, addedRelationships, removedRelationships := FindDifferences(rc1, rc2)

	addedResourcesByType := map[ResourceType][]Resource{}

	for i := range addedResources {
		r := addedResources[i]
		addedResourcesByType[r.ResourceType()] = append(addedResourcesByType[r.ResourceType()], r)
	}

	removedResourcesByType := map[ResourceType][]Resource{}

	for i := range removedResources {
		r := removedResources[i]
		removedResourcesByType[r.ResourceType()] = append(removedResourcesByType[r.ResourceType()], r)
	}

	for _, k := range AvailableTypes {
		if len(addedResourcesByType[k]) > 0 || len(removedResourcesByType[k]) > 0 {
			fmtcolor.White.Printf("[%s]:\n", k)
			printResources(addedResourcesByType[k], "+")
			printResources(removedResourcesByType[k], "-")
			fmt.Println()
		}
	}

	fmtcolor.White.Println("[Relationships]:")
	printRelationships(addedRelationships, "+")
	printRelationships(removedRelationships, "-")
}

// containsRelationship checks if a relationship is present in a slice of relationships.
func containsRelationship(relationships []Relationship, rel Relationship) bool {
	for _, r := range relationships {
		if r.Source.Value() == rel.Source.Value() && r.Source.ResourceType() == rel.Source.ResourceType() &&
			r.Target.Value() == rel.Target.Value() && r.Target.ResourceType() == rel.Target.ResourceType() {
			return true
		}
	}

	return false
}

// printResources prints the resources.
func printResources(resources []Resource, simbol string) {
	for _, res := range resources {
		if simbol == "+" {
			fmtcolor.Green.Printf("%s %s\n", simbol, res.Value())
		} else {
			fmtcolor.Red.Printf("%s %s\n", simbol, res.Value())
		}
	}
}

// printRelationships prints the relationships.
func printRelationships(relationships []Relationship, simbol string) {
	for _, rel := range relationships {
		if simbol == "+" {
			fmtcolor.Green.Printf("%s Source: %s (%s)\n", simbol, rel.Source.Value(), rel.Source.ResourceType())
			fmtcolor.Green.Printf("  Target: %s (%s)\n", rel.Target.Value(), rel.Target.ResourceType())
		} else {
			fmtcolor.Red.Printf("%s Source: %s (%s)\n", simbol, rel.Source.Value(), rel.Source.ResourceType())
			fmtcolor.Red.Printf("  Target: %s (%s)\n", rel.Target.Value(), rel.Target.ResourceType())
		}
	}
}
