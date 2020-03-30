package plist

type Meta struct {
	Nodes []Node
}

type Node interface {
	AddAnnotation(val string)
	Annotations() []Annotation
	Value() string
	Nodes() []Node
	AddNode(node Node)
	SetValue(val string)
	//AddChildAnnotation(key, val string)
}

//type Node struct {
//	Annotations []Annotation
//}

type MetaNode struct {
	value string
	annotations []Annotation
	nodes []Node
	childAnnotations map[string]Annotation
}

func (n *MetaNode) Nodes() []Node {
	return n.nodes
}

func (n *MetaNode) Annotations() []Annotation {
	return n.annotations
}

func (n *MetaNode) AddAnnotation(val string) {
	if n.annotations == nil {
		n.annotations = []Annotation{}
	}
	n.annotations = append(n.annotations, Annotation{value:val})
}

func (n *MetaNode) AddChildAnnotation(key, value string) {
	if n.childAnnotations == nil {
		n.childAnnotations = make(map[string]Annotation)
	}
	n.childAnnotations[key] = Annotation{value:value}
}


func (n *MetaNode) Value() string {
	return n.value
}

func (n *MetaNode) SetValue(value string) {
	n.value = value
}

func (n *MetaNode) AddNode(node Node) {
	if n.nodes == nil {
		n.nodes = []Node{}
	}
	//if n.childAnnotations == nil {
	//	n.childAnnotations = make(map[string]Annotation)
	//}
	//if annotation, ok := n.childAnnotations[node.Value()]; ok {
	//	node.AddAnnotation(annotation.Value())
	//}
	n.nodes = append(n.nodes, node)
}

func (n *MetaNode) addAnnotation(str string) {
	if n.annotations == nil {
		n.annotations = []Annotation{}
	}
	n.annotations = append(n.annotations, Annotation{value:str})
}

type Annotation struct {
	value string
}

func (n *Annotation) Annotations() []Annotation {
	return nil
}

func (n *Annotation) Value() string {
	return n.value
}

func (n *Annotation) Nodes() []Node {
	return nil
}

func (*Annotation) AddNode(_ Node) {

}

func (*Annotation) AddAnnotation(_ string) {

}

func (n *Annotation) SetValue(value string) {
	n.value = value
}

func NewMeta() *Meta {
	return &Meta{Nodes: []Node{}}
}

func (m *Meta) addMetaNode(str string) {
	m.checkInitialized()
	m.addNode(&MetaNode{value: str})
}

func (m *Meta) addAnnotationNode(str string) {
	m.checkInitialized()
	m.addNode(&Annotation{value:str})
}

func (m *Meta) addNode(n Node) {
	m.checkInitialized()
	m.Nodes = append(m.Nodes, n)
}

func (m *Meta) checkInitialized() {
	if m.Nodes == nil {
		m.Nodes = []Node{}
	}
}
