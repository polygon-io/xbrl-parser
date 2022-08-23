package main

import (
	"encoding/xml"
	"io"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	entryPointXSD := "../../_random_stuff/xbrl/us-gaap-2022/us-gaap-2022/entire/us-gaap-entryPoint-all-2022.xsd"
	entryPointRootDir := path.Dir(entryPointXSD)

	reader, err := os.Open(entryPointXSD)
	require.NoError(t, err)

	xmlBytes, err := io.ReadAll(reader)
	require.NoError(t, err)

	var taxonomy XBRLTaxonomy
	require.NoError(t, xml.Unmarshal(xmlBytes, &taxonomy))

	for _, elem := range taxonomy.Elements {
		elem.OriginalFile = entryPointXSD
	}

	for _, linkbase := range taxonomy.AppInfo.LinkBases {
		linkbase.OriginalFile = entryPointXSD
	}

	for _, roleType := range taxonomy.AppInfo.RoleTypes {
		roleType.OriginalFile = entryPointXSD
	}

	require.NoError(t, taxonomy.ResolveLocalImports(entryPointRootDir, true))

	t.Logf("total elements: %d", len(taxonomy.Elements))
	t.Logf("files imported: %d", len(ResolvedImportsSet))

	for _, elem := range taxonomy.Elements {
		if elem.Name == "Assets" {
			t.Log(elem)
		}
	}

	//require.NoError(t, taxonomy.ResolveLocalLinkBases(entryPointRootDir, true))

	t.Logf("total linkbases: %d", len(taxonomy.AppInfo.LinkBases))
	t.Logf("linkbases imported: %d", len(ResolvedLinkBaseSet))

	start := time.Now()
	totalLocators := 0

	labels, err := taxonomy.LookupLabels("Assets")
	require.NoError(t, err)

	for _, label := range labels {
		t.Logf("label for Assets: %s (%s)", label.Value, label.Role)
	}

	t.Logf("finding labels took %s", time.Since(start).String())
	t.Logf("total of %d locators", totalLocators)

	//require.NoError(t, FormatLegendForPresentationLink(&taxonomy, "http://fasb.org/us-gaap/role/statement/StatementOfFinancialPositionClassified"))

	processed, err := NewProcessedXBRLTaxonomy(&taxonomy)
	require.NoError(t, err)

	for role, label := range processed.Labels.ForElement("Assets") {
		t.Logf("label from processed: %s (%s)\n", label, role)
	}

	presentation, exists := processed.PresentationsByRole["http://fasb.org/us-gaap/role/statement/StatementOfFinancialPositionClassified"]
	require.True(t, exists)

	t.Log(ResolvedImportsSet[path.Join(path.Dir(entryPointXSD), "../elts/us-roles-2022.xsd")])

	t.Logf("Presentation Definition: %s", presentation.Definition)
	presentation.Graph.Traverse(func(node *PresentationNode, depth int) bool {
		label := processed.Labels.ForElementAndRole(node.Element.Name, node.PreferredLabel)
		if label == "" {
			label = node.Element.Name
		}

		t.Logf("%s-> %s", strings.Repeat(" ", depth), label)
		return true
	})
}

func TestLabelThing(t *testing.T) {
	l1 := ProcessedLabels{
		"Assets": map[string]string{
			"role1": "role1value",
		},
	}

	l2 := ProcessedLabels{
		"Assets": map[string]string{
			"role2": "role2value",
		},
	}

	l1.CombineWith(l2)
	t.Log(l1["Assets"]["role1"])
	t.Log(l1["Assets"]["role2"])
}
