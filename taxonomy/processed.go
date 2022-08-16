package main

type ProcessedXBRLTaxonomy struct {
	TargetNamespace string

	ElementsByName        map[string]*XBRLTaxonomyElement
	ElementsByLocatorHRef map[string]*XBRLTaxonomyElement

	Labels                  ProcessedLabels
	PresentationTreesByRole map[string]PresentationGraph
}

func NewProcessedXBRLTaxonomy(raw *XBRLTaxonomy) (*ProcessedXBRLTaxonomy, error) {
	processed := &ProcessedXBRLTaxonomy{
		TargetNamespace:         raw.TargetNamespace,
		ElementsByName:          make(map[string]*XBRLTaxonomyElement, len(raw.Elements)),
		ElementsByLocatorHRef:   make(map[string]*XBRLTaxonomyElement, len(raw.Elements)),
		Labels:                  make(ProcessedLabels),
		PresentationTreesByRole: make(map[string]PresentationGraph),
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

	// Populate PresentationGraphs
	for _, linkbase := range raw.AppInfo.LinkBases {
		for _, presentationLink := range linkbase.PresentationLinks {
			graph := processPresentationLink(linkbase.OriginalFile, &presentationLink, processed.ElementsByLocatorHRef)
			processed.PresentationTreesByRole[presentationLink.Role] = graph
		}
	}

	return processed, nil
}
