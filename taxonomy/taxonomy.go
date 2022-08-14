package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"
)

type XBRLTaxonomy struct {
	LocalDirectory string `xml:"-"`

	XMLName xml.Name `xml:"schema"` // TODO: specifiy namespace?

	TargetNamespace string `xml:"targetNamespace,attr"`

	AppInfo XBRLAppInfo `xml:"annotation>appinfo"`

	Imports []XBRLTaxonomyImport `xml:"import"`

	Documentation string `xml:"documentation"` // TODO: can there be more than one of these? What happens?

	Elements []*XBRLTaxonomyElement `xml:"element"`
}

type XBRLAppInfo struct {
	XMLName xml.Name `xml:"appinfo"`

	Refs []XLinkLinkBaseRef `xml:"linkbaseRef"`

	LinkBases []*XLinkLinkBase `xml:"linkbase"`
}

type XBRLTaxonomyImport struct {
	XMLName        xml.Name `xml:"import"` // TODO: specify namespace?
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
}

type XBRLTaxonomyElement struct {
	OriginalFile string `xml:"-"`

	XMLName xml.Name `xml:"element"` // TODO: namespace?

	ID                string `xml:"id,attr"` // Technically optional?
	Name              string `xml:"name,attr"`
	Type              string `xml:"type,attr"`
	SubstitutionGroup string `xml:"substitutionGroup,attr"`

	PeriodType string `xml:"periodType"` // TODO: xbrli namespace (http://www.xbrl.org/2003/instance
	Balance    string `xml:"balance"`    // TODO: xbrli namespace (http://www.xbrl.org/2003/instance

	Abstract bool `xml:"abstract,attr"`
	Nillable bool `xml:"nillable,attr"`
}

var ResolvedImportsSet = make(map[string]bool)  // TODO: This shouldn't be a global variable...obviously
var ResolvedLinkBaseSet = make(map[string]bool) //  TODO: This shouldn't be a global variable...obviously

func (t *XBRLTaxonomy) ResolveLocalImports(rootDir string, recursive bool) error {
	for _, i := range t.Imports {
		if strings.HasPrefix(i.SchemaLocation, "http") {
			fmt.Println("skipping remote import: " + i.SchemaLocation)
			continue
		}

		if err := t.importLocalTaxonomy(rootDir, i, recursive); err != nil {
			return err
		}
	}

	return nil
}

func (t *XBRLTaxonomy) importLocalTaxonomy(rootDir string, i XBRLTaxonomyImport, recursive bool) error {
	fileLocation := path.Join(rootDir, i.SchemaLocation)
	fmt.Printf("Resolving local import: %s -> %s\n", i.SchemaLocation, fileLocation)

	if _, exists := ResolvedImportsSet[fileLocation]; exists {
		return nil
	}

	ResolvedImportsSet[fileLocation] = true

	importedFileBytes, err := os.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	taxonomy := &XBRLTaxonomy{}
	if err := xml.Unmarshal(importedFileBytes, taxonomy); err != nil {
		return err
	}

	for _, elem := range taxonomy.Elements {
		elem.OriginalFile = fileLocation
	}

	for _, linkbase := range taxonomy.AppInfo.LinkBases {
		linkbase.OriginalFile = fileLocation
	}

	if err := taxonomy.ResolveLocalLinkBases(path.Dir(fileLocation), true); err != nil {
		return err
	}

	if recursive {
		if err := taxonomy.ResolveLocalImports(path.Dir(fileLocation), recursive); err != nil {
			return fmt.Errorf("resolve recursive imports ('%s' -> '%s'): %w", i.SchemaLocation, fileLocation, err)
		}
	}

	t.CombineWith(taxonomy)
	return nil
}

func (t *XBRLTaxonomy) CombineWith(other *XBRLTaxonomy) {
	t.Elements = append(t.Elements, other.Elements...)

	//t.AppInfo.Refs = append(t.AppInfo.Refs, other.AppInfo.Refs...)
	t.AppInfo.LinkBases = append(t.AppInfo.LinkBases, other.AppInfo.LinkBases...)
}

func (t *XBRLTaxonomy) LookupLabels(element string) ([]XLinkLabel, error) {
	var relevantLabels []XLinkLabel

	// wtf this is wild
	for _, linkbase := range t.AppInfo.LinkBases {
		for _, labelLink := range linkbase.LabelLinks {
			var assetsLocator *XLinkLocator

			// Find assets locator
			for _, loc := range labelLink.Locators {
				elem, err := t.ResolveLocator(linkbase.OriginalFile, loc)
				if err != nil {
					return nil, err // TODO: err msg?
				}

				if elem.Name == element {
					assetsLocator = &loc
					break
				}
			}

			if assetsLocator != nil {
				// Find labelArc referencing assets locator
				for _, arc := range labelLink.LabelArcs {
					if arc.From == assetsLocator.Label {
						// Find label referencing arc.To
						for _, label := range labelLink.Labels {
							if label.Label == arc.To {
								relevantLabels = append(relevantLabels, label)
							}
						}
					}
				}
			}
		}
	}

	return relevantLabels, nil
}