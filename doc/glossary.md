[Simple Markdown Generator](README.md)

---


<a/><a id="/glossary"/><a id="section-1"/>
# Glossary
[A](a) &nbsp;[B](b) &nbsp;[C](c) &nbsp;[D](d) &nbsp;[E](e) &nbsp;[F](f) &nbsp;[G](g) &nbsp;H &nbsp;[I](i) &nbsp;J &nbsp;[K](k) &nbsp;[L](l) &nbsp;[M](m) &nbsp;[N](n) &nbsp;O &nbsp;P &nbsp;Q &nbsp;[R](r) &nbsp;[S](s) &nbsp;[T](t) &nbsp;U &nbsp;[V](v) &nbsp;W &nbsp;X &nbsp;Y &nbsp;Z &nbsp;

## A

### [`anchor`](statements.md#/statement/anchor)<a id="glossary/statement/anchor"/>
A <a href="#glossary/statement">statement</a> used define a titled anchor.
### [Anchor](syntax.md#/anchors)<a id="glossary/anchor"/>
A <a href="#glossary/tag">tag</a> used to identify <a href="#glossary/refelem">referencable elements</a>. They can be
used to establish hyperlinks among documents. There are <a href="#glossary/globa">global anchors</a>
and <a href="#glossary/loca">local anchors</a>.

## B

### [`block`](statements.md#/statement/block)<a id="glossary/statement/block"/>
A <a href="#glossary/statement">statement</a> used to define a <a href="#glossary/textmodule">text module</a>.
### [`blockref`](statements.md#/statement/blockref)<a id="glossary/statement/blockref"/>
A <a href="#glossary/statement">statement</a> used to instantiate a <a href="#glossary/textmodule">text module</a>.
## C

### [`center`](statements.md#/statement/center)<a id="glossary/statement/center"/>
A <a href="#glossary/statement">statement</a> used to center the embedded content lines.
### [`cs`](statements.md#/symbols)<a id="glossary/statement/cs"/>
A <a href="#glossary/statement">statement</a> emitting the (c)omment (s)tart sequence (`/#`) comment. 
## D

### [Directive](syntax.md#/directives)<a id="glossary/directive"/>
Directives are dedicated markers used by the generator to influence and structure the
generation of a document tree.

## E

### [`escape`](statements.md#/statement/escape)<a id="glossary/statement/escape"/>
A <a href="#glossary/statement">statement</a> used to apply HTML escaping on its content.
## F

### [`figure`](statements.md#/statement/figure)<a id="glossary/statement/figure"/>
A <a href="#glossary/statement">statement</a> used add an image to the output.
## G

### [`glossary`](statements.md#/statement/glossary)<a id="glossary/statement/glossary"/>
A <a href="#glossary/statement">statement</a> used to generate a glossary for the defined <a href="#glossary/term">terms</a>.
### [Global Anchor](syntax.md#/anchors)<a id="glossary/globa"/>
Location independent anchor globally unique for the <a href="#glossary/sourcetree">source tree</a>.

### [Global Tag](syntax.md#/tags)<a id="glossary/globtag"/>
Location independent <a href="#glossary/tag">tag</a> globally unique for the complete <a href="#glossary/sourcetree">source tree</a>.

## I

### [`include`](statements.md#/statement/include)<a id="glossary/statement/include"/>
A <a href="#glossary/statement">statement</a> used to include the content of a file.
## K

### [Keyword](syntax.md#/directives)<a id="glossary/keyword"/>
The name of a <a href="#glossary/directive">directive</a>.

## L

### [`label`](statements.md#/statement/label)<a id="glossary/statement/label"/>
A <a href="#glossary/statement">statement</a> used to add the label of the referenced element to the document.
### [`labeled`](statements.md#/statement/labeled)<a id="glossary/statement/labeled"/>
A <a href="#glossary/statement">statement</a> used add a tagged element with a caption to the output.
### [`link`](statements.md#/statement/link)<a id="glossary/statement/link"/>
A <a href="#glossary/statement">statement</a> used to add a hyperlink to some embedded text.
### [Local Anchor](syntax.md#/anchors)<a id="glossary/loca"/>
An <a href="#glossary/anchor">anchor</a> locally valid in a generated document.

### [Local Tag](syntax.md#/tags)<a id="glossary/loctag"/>
A <a href="#glossary/tag">tag</a> locally valid in a
generated document.
## M

### [*Markdown Generator*](README.md#section-1)<a id="glossary/mdgen"/>
The tool this documentation is for. It maintains consistent hyperlinks
among a set of markdown documents.

### [Markdown Text](syntax.md#/sourcedoc)<a id="glossary/markdown"/>
Part of a <a href="#glossary/sourcedoc">source document</a> used as plain markdown text for the generated
markdown files.

### [Master Number Range](syntax.md#/master)<a id="glossary/masterrange"/>
A <a href="#glossary/numberrange">number range</a> controlling a <a href="#glossary/slaverange">slave number range</a>.
## N

### [`nl`](statements.md#/symbols)<a id="glossary/statement/nl"/>
A <a href="#glossary/statement">statement</a> emitting a newline character
### [`numberrange`](statements.md#/statement/numberrange)<a id="glossary/statement/numberrange"/>
A <a href="#glossary/statement">statement</a> used to declare and configure <a href="#glossary/numberrange">number ranges</a>.
### [Number Range](syntax.md#/numberranges)<a id="glossary/numberrange"/>
A hierarchical labeling mechanism, e.g. used to label sections.

## R

### [`ref`](statements.md#/statement/ref)<a id="glossary/statement/ref"/>
A <a href="#glossary/statement">statement</a> used to add a linked label to the document.
### [Referencable Element](syntax.md#/anchors)<a id="glossary/refelem"/>
Part of a document, which can be target of
a hyperlink (see <a href="statements.md#/statements/referencables">→3.1</a>).
## S

### [`section`](statements.md#/statement/section)<a id="glossary/statement/section"/>
A <a href="#glossary/statement">statement</a> used to describe a structural element in the final document tree.
### [`sectionref`](statements.md#/statement/sectionref)<a id="glossary/statement/sectionref"/>
A <a href="#glossary/statement">statement</a> used to link the section structure of another <a href="#glossary/sourcedoc">source document</a> into the
  own section structure. This statement is related to statement <a href="#glossary/statement/section">`section`</a>.
### [`subrange`](statements.md#/statement/subrange)<a id="glossary/statement/subrange"/>
A <a href="#glossary/statement">statement</a> used to open a new sub level for a <a href="#glossary/numberrange">number ranges</a>.
### [`syntax`](statements.md#/statement/syntax)<a id="glossary/statement/syntax"/>
A <a href="#glossary/statement">statement</a> used render simple syntax expressions.
### [Scope](syntax.md#/scoped)<a id="glossary/scope"/>
  A scope is used as namespace to resolve local <a href="#glossary/tag">tags</a> to <a href="#glossary/tagelem">taggable elements</a>.
  
### [Slave Number Range](syntax.md#/master)<a id="glossary/slaverange"/>
A <a href="#glossary/numberrange">number range</a> controlled by a <a href="#glossary/masterrange">master number range</a>.
### [Source Document](syntax.md#/sourcedoc)<a id="glossary/sourcedoc"/>
A file used a source for the markdown generator. It consists of
a sequence of <a href="#glossary/directive">directives</a> and markdown text and used the file suffix `.mdg`.

### [Source Tree](syntax.md#/filestructure)<a id="glossary/sourcetree"/>
A tree of source files under a common root folder containing source files for
the generator, which build an interconnected documentation.

### [Statement](statements.md#/statements)<a id="glossary/statement"/>
A sequence of one or more <a href="#glossary/directive">directives</a>.

## T

### [`template`](statements.md#/statement/template)<a id="glossary/statement/template"/>
A <a href="#glossary/statement">statement</a> flagging a <a href="#glossary/sourcedoc">source document</a> to be omitted
  from the generation of a markdown file.
### [`term`](statements.md#/statement/term)<a id="glossary/statement/term"/>
A <a href="#glossary/statement">statement</a> used to output a previously defined <a href="#glossary/term">term</a>.
### [`termdef`](statements.md#/statement/termdef)<a id="glossary/statement/termdef"/>
A <a href="#glossary/statement">statement</a> used to define a <a href="#glossary/term">term</a> used in the document tree.
### [`title`](statements.md#/statement/title)<a id="glossary/statement/title"/>
A <a href="#glossary/statement">statement</a> used to add the title of the referenced element to the document.
### [`toc`](statements.md#/statement/toc)<a id="glossary/statement/toc"/>
A <a href="#glossary/statement">statement</a> used to add a table of contents.
### [Tag](syntax.md#/tags)<a id="glossary/tag"/>
A short identifier used identify a
dedicated element in the <a href="#glossary/sourcetree">source tree</a>. There are three kinds
of <a href="#glossary/tagelem">taggable elements</a>: <a href="#glossary/refelem">referencable elements</a>, <a href="#glossary/term">terms</a>, and <a href="#glossary/textmodule">text modules</a>.

### [Tag Composition](syntax.md#/scoped)<a id="glossary/tagcomp"/>
Composing a <a href="#glossary/tag">tag</a> by substitution
  of various <a href="#glossary/scope">scope</a> related attributes to explicitly provide unique names for
  <a href="#glossary/tagelem">taggable elements</a> in <a href="#glossary/textmodule">text modules</a>.
### [Tag Extension](syntax.md#/scoped)<a id="glossary/tagext"/>
Adding additional namespace parts in
  front of a a <a href="#glossary/tag">tag</a> name.
### [Taggable Element](syntax.md#/tags)<a id="glossary/tagelem"/>
An identifiable element
descibed by a <a href="#glossary/statement">statement</a> of the <a href="#glossary/mdgen">*Markdown Generator*</a> in the <a href="#glossary/sourcetree">source tree</a>.

### [Term](syntax.md#/terms)<a id="glossary/term"/>
A part of the terminology of the
<a href="#glossary/sourcetree">source tree</a> with a dedicated meaning.
### [Text Module](syntax.md#/textmodules)<a id="glossary/textmodule"/>
Tagged and parameterized reusable block of content, which can be instantiated
multiple times all over the document tree.

## V

### [`value`](statements.md#/statement/value)<a id="glossary/statement/value"/>
A <a href="#glossary/statement">statement</a> used access the argument value of a
  {term textmodule}} parameter.
