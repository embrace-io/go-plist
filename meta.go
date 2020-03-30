package plist

type Meta struct {
	Nodes []Node
}

type Node interface {
	AddAnnotation(annotation *Annotation)
	Annotations() []*Annotation
	Value() string
	Nodes() []Node
	AddNode(node Node)
	SetValue(val string)
}

type MetaNode struct {
	value string
	annotations []*Annotation
	nodes []Node
}

func NewMetaNode() *MetaNode {
	return &MetaNode{
		annotations: []*Annotation{},
		nodes:       []Node{},
	}
}

func (n *MetaNode) Nodes() []Node {
	return n.nodes
}

func (n *MetaNode) Annotations() []*Annotation {
	return n.annotations
}

func (n *MetaNode) AddAnnotation(annotation *Annotation) {
	if n.annotations == nil {
		n.annotations = []*Annotation{}
	}
	n.annotations = append(n.annotations, annotation)
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
	n.nodes = append(n.nodes, node)
}

type Annotation struct {
	value string
}

func NewAnnotation(value string) *Annotation {
	return &Annotation{value:value}
}

func (n *Annotation) Annotations() []*Annotation {
	return nil
}

func (n *Annotation) Value() string {
	return n.value
}

func (n *Annotation) Nodes() []Node {
	return nil
}

func (*Annotation) AddNode(_ Node) {}

func (*Annotation) AddAnnotation(_ *Annotation) {}

func (n *Annotation) SetValue(value string) {
	n.value = value
}

func NewMeta() *Meta {
	return &Meta{Nodes: []Node{}}
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

func nodeListToMap(nodes []Node) map[string]Node {
	m := make(map[string]Node)
	for _, node := range nodes {
		m[node.Value()] = node
	}
	return m
}
