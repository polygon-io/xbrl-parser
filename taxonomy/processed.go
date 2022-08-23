package main

import "path"

type ProcessedXBRLTaxonomy struct {
	TargetNamespace string

	ElementsByName        map[string]*XBRLTaxonomyElement
	ElementsByLocatorHRef map[string]*XBRLTaxonomyElement

	RoleTypesByHRef map[string]*XBRLRoleType

	Labels              ProcessedLabels
	PresentationsByRole map[string]*ProcessedPresentation
}

func NewProcessedXBRLTaxonomy(raw *XBRLTaxonomy) (*ProcessedXBRLTaxonomy, error) {
	processed := &ProcessedXBRLTaxonomy{
		TargetNamespace:       raw.TargetNamespace,
		ElementsByName:        make(map[string]*XBRLTaxonomyElement, len(raw.Elements)),
		ElementsByLocatorHRef: make(map[string]*XBRLTaxonomyElement, len(raw.Elements)),
		RoleTypesByHRef:       make(map[string]*XBRLRoleType, len(raw.AppInfo.RoleTypes)),
		Labels:                make(ProcessedLabels),
		PresentationsByRole:   make(map[string]*ProcessedPresentation),
	}

	// Populate ElementsByName map
	for _, element := range raw.Elements {
		processed.ElementsByName[element.Name] = element
		processed.ElementsByLocatorHRef[element.OriginalFile+"#"+element.ID] = element
	}

	// Populate Labels
	for _, linkbase := range raw.AppInfo.LinkBases {
		for _, labelLink := range linkbase.LabelLinks {
			labels := processLabelLinks(linkbase.OriginalFile, &labelLink, processed.ElementsByLocatorHRef)
			processed.Labels.CombineWith(labels)
		}
	}

	// Populate Roles
	for _, roleType := range raw.AppInfo.RoleTypes {
		processed.RoleTypesByHRef[roleType.OriginalFile+"#"+roleType.ID] = roleType
	}

	// Populate PresentationGraphs
	for _, linkbase := range raw.AppInfo.LinkBases {
		for _, presentationLink := range linkbase.PresentationLinks {
			var definition string
			if ref := linkbase.ResolveRoleRef(presentationLink.Role); ref != nil {
				absoluteHRef := path.Join(path.Dir(linkbase.OriginalFile), ref.HRef)
				if roleType, exists := processed.RoleTypesByHRef[absoluteHRef]; exists {
					definition = roleType.Definition
				}
			}

			graph := processPresentationLink(linkbase.OriginalFile, &presentationLink, processed.ElementsByLocatorHRef)

			processed.PresentationsByRole[presentationLink.Role] = &ProcessedPresentation{
				Definition: definition,
				Graph:      graph,
			}
		}
	}

	return processed, nil
}
