package main

import (
	"path"
	"sort"
)

type PresentationNode struct {
	Element        *XBRLTaxonomyElement // TODO make this work
	PreferredLabel string
	Order          int

	Children []*PresentationNode
}

// PresentationGraph consists of one or more disjoint presentation trees.
// Each PresentationNode in the slice is the root of a presentation tree.
// Each presentation tree in the graph MUST be a directed acyclic graph.
type PresentationGraph []*PresentationNode

// Traverse traverses the PresentationGraph in depth first order and calls the given 'action' func on each node it encounters.
// If the 'action' returns false, traversal stops early and 'action' will not be called on any more nodes.
func (g PresentationGraph) Traverse(action func(node *PresentationNode, depth int) bool) {
	// Traverse each tree in the graph
	for _, root := range g {
		if !traverseTree(root, 0, action) {
			return
		}
	}
}

func traverseTree(node *PresentationNode, depth int, action func(node *PresentationNode, depth int) bool) bool {
	if !action(node, depth) {
		return false
	}

	for _, child := range node.Children {
		if !traverseTree(child, depth+1, action) {
			return false
		}
	}

	return true
}

func (g *PresentationGraph) FindNodeByElement(elementName string) *PresentationNode {
	var node *PresentationNode

	g.Traverse(func(n *PresentationNode, _ int) bool {
		if n.Element.Name == elementName {
			node = n
			return false
		}

		return true
	})

	return node
}

func processPresentationLink(originalFile string, presentationLink *XLinkPresentationExtendedLink, elementsByLocationHRef map[string]*XBRLTaxonomyElement) PresentationGraph {
	var graph PresentationGraph

	// Build the presentation graph
	for _, arc := range presentationLink.PresentationArcs {
		if arc.ArcRole != "http://www.xbrl.org/2003/arcrole/parent-child" {
			// TODO: log? error? idk do something
			continue
		}

		// Find the parent and child locators
		var parentLocatorHRef, childLocatorHRef string
		for _, loc := range presentationLink.Locators {
			switch loc.Label {
			case arc.From:
				parentLocatorHRef = path.Join(path.Dir(originalFile), loc.HRef)
			case arc.To:
				childLocatorHRef = path.Join(path.Dir(originalFile), loc.HRef)
			}
		}

		if parentLocatorHRef == "" || childLocatorHRef == "" {
			// TODO: log? error? idk do something
			continue
		}

		parentElement, exists := elementsByLocationHRef[parentLocatorHRef]
		if !exists {
			// TODO: log? error? idk do something
			continue
		}

		parentNode := graph.FindNodeByElement(parentElement.Name)
		if parentNode == nil { // If we haven't seen this parent yet, we need to make the node
			parentNode = &PresentationNode{Element: parentElement}
			// TODO: this might not work depending on the order of presentationArcs in the linkbase
			graph = append(graph, parentNode)
		}

		childElement, exists := elementsByLocationHRef[childLocatorHRef]
		if !exists {
			// TODO: log? error? idk do something
			continue
		}

		child := &PresentationNode{
			Element:        childElement,
			PreferredLabel: arc.PreferredLabel,
			Order:          arc.Order,
		}

		parentNode.Children = append(parentNode.Children, child)
	}

	// Ensure children are sorted by their order attribute
	graph.Traverse(func(node *PresentationNode, _ int) bool {
		sort.SliceStable(node.Children, func(i, j int) bool {
			return node.Children[i].Order < node.Children[j].Order
		})

		return true
	})

	return graph
}
