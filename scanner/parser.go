/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"fmt"
	"io"
	"strings"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	Tokens.Register("end", parseEnd)
}

type Consumer interface {
	Inventory
}

type Node interface {
	Located
	GetDocument() Document
	Print(gap string)
}

type NodeBase struct {
	location
	document *document
}

func (n *NodeBase) GetDocument() Document {
	return n.document
}

func (n *NodeBase) EvaluateStatic(ctx ResolutionContext) error {
	return n.Errorf("no static evaluation possible")
}

func NewNodeBase(d Document, location Location) NodeBase {
	return NodeBase{location, d}
}

type InventoryContainer interface {
	Inventory
	SetLabelRule(loc *Location, typ string, abbrev, sep string, l labels.Rule, lvl int) error
	SetLabelMaster(typ string, master string, sep string, limit int) error
}

type NodeSequence interface {
	AddNode(n Node) error
	GetNodes() []Node

	Walk(f Resolver) error
	Register
	LabelResolver
	ValueResolver
	Emitter
	StaticEvaluator

	Print(gap string)
}

type NodeContainer interface {
	Node
	InventoryContainer
	NodeSequence
	Type() string
}

type inventoryScope struct {
	Inventory
}

func NewInventoryScope(i Inventory) InventoryContainer {
	return &inventoryScope{i}
}

func (s *inventoryScope) SetLabelRule(loc *Location, typ string, abbrev, sep string, l labels.Rule, lvl int) error {
	return fmt.Errorf("no label rule possible at block level")
}

func (s *inventoryScope) SetLabelMaster(typ string, master, sep string, lvl int) error {
	return fmt.Errorf("no label rule possible at block level")
}

type nodesequence struct {
	document []Node
}

var _ NodeSequence = (*nodesequence)(nil)

func NewNodeSequence() NodeSequence {
	return &nodesequence{}
}

func (c *nodesequence) AddNode(n Node) error {
	c.document = append(c.document, n)
	return nil
}

func (c *nodesequence) Print(gap string) {
	for _, d := range c.document {
		d.Print(gap)
	}
}

func (c *nodesequence) GetNodes() []Node {
	return c.document
}

func (c nodesequence) Walk(f Resolver) error {
	for _, n := range c.GetNodes() {
		err := f(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *nodesequence) Register(ctx ResolutionContext) error {
	return c.Walk(Resolve[Register](ctx))
}

func (c *nodesequence) ResolveLabels(ctx ResolutionContext) error {
	return c.Walk(Resolve[LabelResolver](ctx))
}

func (c *nodesequence) ResolveValues(ctx ResolutionContext) error {
	return c.Walk(Resolve[ValueResolver](ctx))
}

func (c *nodesequence) Emit(ctx ResolutionContext) error {
	return c.Walk(Resolve[Emitter](ctx))
}

func (c *nodesequence) EvaluateStatic(ctx ResolutionContext) error {
	return c.Walk(Resolve[StaticEvaluator](ctx))
}

type NodeContainerBase struct {
	NodeBase
	InventoryContainer
	NodeSequence
	typ string
}

func NewContainerBase(typ string, d Document, location Location, parent ...InventoryContainer) NodeContainerBase {
	return *NewNodeSequenceContainer(typ, d, location, NewNodeSequence(), parent...)
}

func NewNodeSequenceContainer(typ string, d Document, location Location, c NodeSequence, parent ...InventoryContainer) *NodeContainerBase {
	return &NodeContainerBase{
		NodeBase:           NewNodeBase(d, location),
		NodeSequence:       c,
		typ:                typ,
		InventoryContainer: utils.Optional(parent...),
	}
}

var _ NodeContainer = (*NodeContainerBase)(nil)

func (c *NodeContainerBase) Type() string {
	return c.typ
}

func (c *NodeContainerBase) Print(gap string) {
	fmt.Printf("%snodes: \n", gap)
	c.NodeSequence.Print(gap + "  ")
}

func (c *NodeContainerBase) EvaluateStatic(ctx ResolutionContext) error {
	return c.NodeBase.EvaluateStatic(ctx)
}

////////////////////////////////////////////////////////////////////////////////

type TaggedNodeBase struct {
	sid TaggedId
	tag string
}

func NewTaggedNodeBase(sid TaggedId, tag string) TaggedNodeBase {
	return TaggedNodeBase{
		sid: sid,
		tag: tag,
	}
}

func (n *TaggedNodeBase) Id() TaggedId {
	return n.sid
}

func (n *TaggedNodeBase) Tag() string {
	return n.tag
}

////////////////////////////////////////////////////////////////////////////////

type Token func(p Parser, e Element) (Element, error)

type tokens struct {
	tokens map[string]Token
}

func (t *tokens) Register(name string, tok Token, skipnl ...bool) {
	t.tokens[name] = tok
	if utils.Optional(skipnl...) {
		SkipNewline.Register(name)
	}
}

func (t *tokens) Get(name string) Token {
	return t.tokens[name]
}

var Tokens = &tokens{map[string]Token{}}

type keywords map[string]bool

func (k keywords) Register(name string, skipnl ...bool) {
	k[name] = true
	if utils.Optional(skipnl...) {
		SkipNewline.Register(name)
	}
}

var Keywords = keywords{}

////////////////////////////////////////////////////////////////////////////////

type Statement interface {
	Name() string
	Start(p Parser, e Element) (Element, error)
}

type Finisher interface {
	End(p Parser, e Element) (Element, error)
}

func (t *tokens) RegisterStatement(s Statement, skipnl ...bool) {
	start, end := false, false
	if len(skipnl) > 0 {
		start = skipnl[0]
	}
	if len(skipnl) > 1 {
		end = skipnl[1]
	}
	t.Register(s.Name(), s.Start, start)
	if e, ok := s.(Finisher); ok {
		t.Register("end"+s.Name(), e.End, end)
	}
}

type StatementBase struct {
	name string
}

func NewStatementBase(name string) StatementBase {
	return StatementBase{name: name}
}

func (s *StatementBase) Name() string {
	return s.name
}

type BracketedStatement[N Node] struct {
	StatementBase
	useAsNode bool
}

func NewBracketedStatement[N Node](name string, useAsNode bool) BracketedStatement[N] {
	return BracketedStatement[N]{
		StatementBase: NewStatementBase(name),
		useAsNode:     useAsNode,
	}
}

func (s *BracketedStatement[N]) End(p Parser, e Element) (Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("no tag possible for %s", e.Token())
	}
	if _, err := Assure[N](s.Name(), p, e); err != nil {
		return nil, err
	}
	return Pop(p, s.useAsNode)
}

////////////////////////////////////////////////////////////////////////////////

type Ids map[string]labels.Rule

func (s Ids) getId(typ string) labels.Rule {
	l := s[typ]
	if l == nil {
		l = labels.NewVoid(typ, 0)
	}
	return l
}

func (s Ids) NextId(typ string) labels.Rule {
	id := s.getId(typ)
	v := id.Next()
	s[typ] = v
	return v
}

func (s Ids) SubId(typ string) {
	id := s.getId(typ)
	s[typ] = id.Sub()
}

type Parser = *parser

type parser struct {
	doc       *document
	tokenizer Tokenizer
	State     *state
}

func (p *parser) Document() Document {
	return p.doc
}

func (p *parser) NextElement() (Element, error) {
	return p.tokenizer.NextElement()
}

type ids = Ids
type state struct {
	parent    *state
	Container NodeContainer
	scopename string
	ids
	lasttag string
}

func (s *state) ScopeName() string {
	return s.scopename
}

func (s *state) LastTag() string {
	return s.lasttag
}

func (s *state) Sub(c NodeContainer, scopename ...string) *state {
	if c == nil {
		c = s.Container
	}
	sn := s.scopename
	if len(scopename) > 0 && scopename[0] != "" {
		sn = scopename[0]
	}
	ids := Ids{}
	for k, l := range s.ids {
		ids[k] = l
	}
	return &state{
		parent:    s,
		scopename: sn,
		Container: c,
		ids:       ids,
	}
}

func (s *state) SetLastTag(tag string) {
	s.lasttag = tag
}

func NewParser(source, refpath string, r io.Reader) Parser {
	p := &parser{
		tokenizer: NewTokenizer(source, r),
		doc:       NewDocument(source, refpath),
	}
	p.State = &state{Container: p.doc, ids: Ids{}, scopename: p.doc.refpath + "#"}
	return p
}

func (p *parser) Parse() (Document, error) {
	_, err := ParseUntil(p, nil, nil)
	if err == nil {
		if p.State.parent != nil {
			return nil, p.Errorf("unfinished %s", p.State.Container.Type())
		}
	}
	return p.doc, err
}

func (p *parser) Errorf(msg string, args ...interface{}) error {
	return p.doc.Location().Errorf(msg, args...)
}

func Assure[T Node](t string, p Parser, e Element) (T, error) {
	var zero T
	if p.State.parent != nil {
		if node, ok := p.State.Container.(T); ok {
			return node, nil
		}
		s := p.State.parent
		for s != nil {
			if _, ok := s.Container.(T); ok {
				return zero, e.Errorf("unfinished token %q pending", p.State.Container.Type())
			}
			s = s.parent
		}
	}
	return zero, e.Errorf("no %q found for %q", t, e.Token())
}

func Lookup[T Node](p Parser) T {
	var zero T
	s := p.State
	for s != nil {
		if c, ok := s.Container.(T); ok {
			return c
		}
		s = s.parent
	}
	return zero
}

func ForbidNesting[T Node](t string, p Parser, e Element) error {
	if p.State.parent != nil {
		if _, ok := p.State.Container.(T); ok {
			return e.Errorf("%s not allowed in %s", e.Token(), t)
		}
	}
	return nil
}

func ForbidNestingInTypes(p Parser, e Element, types ...string) error {
	s := p.State
	for s.parent != nil {
		for _, t := range types {
			if s.Container.Type() == t {
				return e.Errorf("%s not allowed in %s", e.Token(), t)
			}
		}
		s = s.parent
	}
	return nil
}

func RequireNesting[T Node](t string, p Parser, e Element, fs ...func(T) (bool, error)) error {
	var zero T
	found := false
	s := p.State
	for s.parent != nil {
		if c, ok := s.Container.(T); ok {
			if len(fs) == 0 {
				return nil
			}
			found = true
			for _, f := range fs {
				if stop, err := f(c); stop || err != nil {
					return err
				}
			}
		}
		s = s.parent
	}
	if !found {
		return e.Errorf("%s may be used in %s, only", e.Token(), t)
	}
	for _, f := range fs {
		if _, err := f(zero); err != nil {
			return err
		}
	}
	return nil
}

func Pop(p Parser, isContent bool) (Element, error) {
	if isContent {
		p.State.parent.Container.AddNode(p.State.Container)
	}
	p.State = p.State.parent
	return p.tokenizer.NextElement()
}

func parseEnd(p Parser, e Element) (Element, error) {
	if &p.State.parent == nil {
		return nil, e.Errorf("no unfinished element found")
	}
	t := Tokens.Get("end" + p.State.Container.Type())
	if t == nil {
		return nil, e.Errorf("end%s not known", p.State.Container.Type())
	}
	return t(p, e)
}

func ParseUntil(p Parser, stop func(p Parser, e Element) bool, c NodeContainer) (Element, error) {
	e, err := p.tokenizer.NextElement()

	if err != nil {
		return nil, err
	}
	if c != nil {
		p.State = p.State.Sub(c)
		defer func() {
			p.State = p.State.parent
		}()
	}
	for e != nil {
		if stop != nil && stop(p, e) {
			return e, nil
		}
		if e.IsText() {
			err = p.State.Container.AddNode(NewTextNode(p.doc, e.Location(), e.Text()))
			if err == nil {
				e, err = p.tokenizer.NextElement()
			}
		} else {
			t := Tokens.Get(e.Token())
			if t == nil {
				if Keywords[e.Token()] {
					last := p.State.Container.Type()
					if last != "" {
						return nil, e.Errorf("unexpected token %q (last unfinished element is %q)", e.Token(), last)
					}
					return nil, e.Errorf("unexpected token %q (may be another element is not finished)", e.Token())
				}
				return nil, e.Errorf("unknown token %q", e.Token())
			}
			e, err = t(p, e)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func ParseSequence(p Parser, e Element) (Element, NodeSequence, error) {
	var err error
	seq := NewNodeSequence()
	c := NewNodeSequenceContainer(e.Token(), p.doc, e.Location(), seq, p.State.Container)
	tok := e.Token()
	end := "end" + e.Token()
	stop := func(p Parser, e Element) bool {
		if p.State.Container != c {
			return false
		}
		return e.Token() == "end" || e.Token() == end
	}
	e, err = ParseUntil(p, stop, c)
	if err != nil {
		return nil, nil, err
	}
	if e == nil {
		return nil, nil, p.Errorf("unfinished {{%s}} at end of document", tok)
	}
	if e.HasTags() {
		return nil, nil, e.Errorf("no tag possible for {{end%s}}", tok)
	}
	e, err = p.tokenizer.NextElement()
	return e, seq, err
}

func ParseSequenceUntil(p Parser, e Element, stop func(p Parser, e Element) bool) (Element, NodeSequence, error) {
	var err error
	seq := NewNodeSequence()
	c := NewNodeSequenceContainer(e.Token(), p.doc, e.Location(), seq, p.State.Container)
	end := "end" + e.Token()
	effstop := func(p Parser, e Element) bool {
		if p.State.Container != c {
			return false
		}
		if e.Token() == "end" || e.Token() == end {
			return true
		}
		return stop != nil && stop(p, e)
	}

	e, err = ParseUntil(p, effstop, c)
	if err != nil {
		return nil, nil, err
	}
	return e, seq, err
}

func ParseElementsUntil(p Parser, f Token) (Element, error) {
	var pendingText Element
	e, err := p.tokenizer.NextElement()
	if err != nil {
		return nil, err
	}
	for e != nil {
		if e.IsText() {
			if pendingText == nil {
				pendingText = e
			} else {
				pendingText.Append(e.Text())
			}
			if strings.TrimSpace(pendingText.Text()) != "" {
				e = nil
				break
			}
			e, err = p.tokenizer.NextElement()
		} else {
			last := e
			e, err = f(p, e)
			if last == e {
				break
			}
			pendingText = nil
		}
		if err != nil {
			return nil, err
		}
	}
	if pendingText != nil {
		if e != nil {
			p.tokenizer.Push(e)
		}
		return pendingText, nil
	}
	return e, nil
}
