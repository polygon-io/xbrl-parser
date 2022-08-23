package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"
)

type XLinkLinkBaseRef struct {
	XMLName xml.Name `xml:"linkbaseRef"`

	ArcRole string `xml:"arcrole,attr"`
	HRef    string `xml:"href,attr"`
	Role    string `xml:"role,attr"`
	Type    string `xml:"type,attr"`
}

type XLinkLinkBase struct {
	OriginalFile string `xml:"-"`

	XMLName xml.Name `xml:"linkbase"`

	Base string `xml:"base,attr"` // See spec section 3.5.2.2

	RoleRefs []*XLinkRoleRef `xml:"roleRef"`

	LabelLinks        []XLinkLabelExtendedLink        `xml:"labelLink"`
	PresentationLinks []XLinkPresentationExtendedLink `xml:"presentationLink"`
}

type XLinkRoleRef struct {
	XMLName xml.Name `xml:"roleRef"`

	RoleURI string `xml:"roleURI,attr"`
	HRef    string `xml:"href,attr"`
	Type    string `xml:"type,attr"`
}

type XLinkPresentationExtendedLink struct {
	XMLName xml.Name `xml:"presentationLink"`

	Role string `xml:"role,attr"`
	Type string `xml:"type,attr"`

	Locators         []XLinkLocator         `xml:"loc"`
	PresentationArcs []XLinkPresentationArc `xml:"presentationArc"`
}

type XLinkPresentationArc struct {
	XMLName xml.Name `xml:"presentationArc"`

	Order   int    `xml:"order,attr"`
	ArcRole string `xml:"arcrole,attr"`

	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
	Type string `xml:"type,attr"`

	PreferredLabel string `xml:"preferredLabel,attr"`
}

type XLinkLabelExtendedLink struct {
	XMLName xml.Name `xml:"labelLink"`

	Locators  []XLinkLocator  `xml:"loc"`
	Labels    []XLinkLabel    `xml:"label"`
	LabelArcs []XLinkLabelArc `xml:"labelArc"`
}

type XLinkLocator struct {
	XMLName xml.Name `xml:"loc"`

	Label string `xml:"label,attr"`
	HRef  string `xml:"href,attr"`
	Type  string `xml:"type,attr"`
}

// Section 5.2.2.2
type XLinkLabel struct {
	XMLName xml.Name `xml:"label"`

	ID string `xml:"id,attr"`

	Language string `xml:"lang,attr"`
	Role     string `xml:"role,attr"`
	Label    string `xml:"label,attr"`
	Type     string `xml:"type,attr"`

	Value string `xml:",chardata"`
}

type XLinkLabelArc struct {
	XMLName xml.Name `xml:"labelArc"`

	Order   string `xml:"order,attr"`
	ArcRole string `xml:"arcrole,attr"`

	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
	Type string `xml:"type,attr"`
}

func (t *XBRLTaxonomy) ResolveLocalLinkBases(rootDir string, recursive bool) error {
	for _, ref := range t.AppInfo.Refs {
		if strings.HasPrefix(ref.HRef, "http") {
			fmt.Println("skipping remote linkbase: " + ref.HRef)
			continue
		}

		// TODO: process roleURI hrefs?

		if err := t.processLocalLinkBase(rootDir, ref); err != nil {
			return err

		}
	}

	return nil
}

func (t *XBRLTaxonomy) processLocalLinkBase(rootDir string, ref XLinkLinkBaseRef) error {
	fileLocation := path.Join(rootDir, ref.HRef)
	fmt.Printf("Resolving local linkbase: %s -> %s\n", ref.HRef, fileLocation)

	if _, exists := ResolvedLinkBaseSet[fileLocation]; exists {
		return nil
	}

	ResolvedLinkBaseSet[fileLocation] = true

	importedFileBytes, err := os.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	var linkbase XLinkLinkBase
	if err := xml.Unmarshal(importedFileBytes, &linkbase); err != nil {
		return err
	}

	linkbase.OriginalFile = fileLocation

	t.AppInfo.LinkBases = append(t.AppInfo.LinkBases, &linkbase)
	return nil
}

func (t *XBRLTaxonomy) ResolveLocator(locatorFile string, loc XLinkLocator) (*XBRLTaxonomyElement, error) {
	var file, id string

	split := strings.Split(loc.HRef, "#")
	switch len(split) {
	case 1: // xlink:href='#us-gaap_TypeOfAdoptionMember'
		file = locatorFile
		id = split[0]
	case 2: // xlink:href='us-gaap-2022.xsd#us-gaap_AOCIAttributableToParentAbstract'
		file = path.Join(path.Dir(locatorFile), split[0])
		id = split[1]
	default:
		return nil, fmt.Errorf("wtf is this href? %s", loc.HRef)
	}

	for _, elem := range t.Elements {
		if elem.OriginalFile == file && elem.ID == id {
			return elem, nil
		}
	}

	return nil, fmt.Errorf("could not locate element for locator: %s", loc.HRef)
}

func (l *XLinkLinkBase) ResolveRoleRef(roleURI string) *XLinkRoleRef {
	for _, ref := range l.RoleRefs {
		if ref.RoleURI == roleURI {
			return ref
		}
	}

	return nil
}
