package resource

type xmlMeta struct {
	Infos []xmlInfo `xml:"info"`
}

type xmlInfo struct {
	Name        string `xml:"name,attr"`
	Version     string `xml:"version,attr"`
	Description string `xml:"description,attr"`
	Type        string `xml:"type,attr"`
}
