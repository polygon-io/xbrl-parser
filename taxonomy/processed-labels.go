package main

import "path"

const DefaultLabelRole = "http://www.xbrl.org/2003/arcrole/concept-label"

// ProcessedLabels is a map of element name to a map of label roles to labels.
// See methods on this type for ease of use
type ProcessedLabels map[string]map[string]string

func (p ProcessedLabels) CombineWith(other ProcessedLabels) {
	for elementName, labelsByRole := range other {
		p.AddLabelsForElement(elementName, labelsByRole)
	}
}

func (p ProcessedLabels) AddLabelsForElement(elementName string, labelsByRole map[string]string) {
	if _, exists := p[elementName]; !exists {
		p[elementName] = make(map[string]string, len(labelsByRole))
	}

	for role, label := range labelsByRole {
		p[elementName][role] = label
	}
}

func (p ProcessedLabels) ForElement(name string) map[string]string {
	return p[name]
}

func (p ProcessedLabels) DefaultForElement(name string) string {
	return p.ForElementAndRole(name, DefaultLabelRole)
}

func (p ProcessedLabels) ForElementAndRole(name, role string) string {
	labelsByRole, exists := p[name]
	if !exists {
		return ""
	}

	return labelsByRole[role]
}

func processLabelLinks(originalFile string, labelLink *XLinkLabelExtendedLink, elementsByLocatorHRef map[string]*XBRLTaxonomyElement) ProcessedLabels {
	results := make(ProcessedLabels)

	for _, arc := range labelLink.LabelArcs {
		if arc.ArcRole != DefaultLabelRole {
			// TODO: log? warn? count? idk do something
			continue // Ignore unknown roles
		}

		var locator *XLinkLocator

		// Find 'from' locator
		for _, loc := range labelLink.Locators {
			if loc.Label == arc.From {
				loc := loc // TODO: not sure if this is necessary
				locator = &loc
				break
			}
		}

		if locator == nil {
			// TODO: log? warn? count? idk do something
			continue
		}

		element, exists := elementsByLocatorHRef[path.Join(path.Dir(originalFile), locator.HRef)]
		if !exists {
			// TODO: log? warn? count? idk do something
			continue
		}

		// Find the 'to' label. There can be one or more labels with different roles here.
		for _, lab := range labelLink.Labels {
			if lab.Label == arc.To {
				if _, exists := results[element.Name]; !exists {
					results[element.Name] = make(map[string]string)
				}

				role := lab.Role
				if role == "" {
					role = "http://www.xbrl.org/2003/role/label" // Default role
				}

				results[element.Name][role] = lab.Value
			}
		}
	}

	return results
}
