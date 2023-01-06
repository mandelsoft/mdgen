
<a/><a id="/syntax"/><a id="section-1"/>
## 1 The Document Graph
Table of Contents

 [1.1 File Structure](#/filestructure)<br>
 [1.2 Source Document](#/sourcedoc)<br>
 [1.3 Directives](#/directives)<br>
 [1.4 Tags](#/tags)<br>
&nbsp;&nbsp; [1.4.1 Anchors for Referencable Elements](#/anchors)<br>
&nbsp;&nbsp; [1.4.2 Tags and Anchors in Scoped Environments](#/scoped)<br>
 [1.5 Number ranges](#/numberranges)<br>
&nbsp;&nbsp; [1.5.1 Label Formats](#/labelformats)<br>
&nbsp;&nbsp; [1.5.2 Cascading Number Ranges](#/master)<br>
 [1.6 Formalized Terminology](#/terms)<br>
 [1.7 Text Modules](#/textmodules)<br>

The *Markdown Generator* maintains a documentation consisting of a set of
markdown files, which are highly interconnected via hyperlinks.
Therefore a set of source files will be parsed and transferred into a
document graph describing all the those connections among parts of the
involved documents. This graph is then transformed again into a plain
set of markdown files by generating the appropriate document
anchors and links.


<a/><a id="/filestructure"/><a id="section-1-1"/>
### 1.1 File Structure

The generator works on a *source tree*.

This is a tree of <a href="#/sourcedoc">source documents</a> and other files under a common root
folder containing source files for the generator, which build a closed
documentation with hyperlinks among different parts.
Non-source documents are ignored as long as they are not used
by `figure` or `include` statements.


The outcome of the generation process is a tree of interconnected markdown
files following the same structure.


<a/><a id="tree"/><a id="example-1"/>
<div align="center"><table><tr><td>


```
/
├── chapters
│   ├── figure.png
│   ├── intro.mdg
│   ├── spec.mdg
│   └── tree.txt
└── README.mdg


```

It contains three source documents denoted by the names
`/README`, `/chapters/intro` and `/chapters/spec`.
</td></tr></table>
 Example 1-a: A typical source tree
</br></br>
</div>



<a/><a id="/sourcedoc"/><a id="section-1-2"/>
### 1.2 Source Document

A *source document* is a file in the <a href="#/filestructure">source tree</a> with the file suffix `.mdg`.
The *name* of the source document is the path from the root of the
<a href="#/filestructure">source tree</a> to the file with the suffix omitted.
It will be used to generate an appropriate markdown file with the suffix `.md`.

A source document consists of a sequence of <a href="#/directives">directives</a> and regular
*markdown text*.

Line content is omitted after a `/#` character sequence. Lines starting
with the comment sequence are completely omitted (see <a href="#example">→example 1-b</a>).

Different source documents can be embedded in a common chapter hierarchy
regardless of their names and locations in the source tree by using
the `sectionref` statement.



<a/><a id="/directives"/><a id="section-1-3"/>
### 1.3 Directives

The markdown generator uses special *directives*. A <a href="#/sourcedoc">source document</a> consists if a sequence of directives and
regular markdown text.

A directive consists of a *keyword* and argument strings. It is described by the following syntax:
<div align=center>

  &#39;`{{`&#39; &#39;` `&#39; [&#39;`*`&#39;] &lt;keyword&gt; {&#39;` `&#39; {&#39;` `&#39;} &lt;arg&gt;} {&#39;` `&#39;} &#39;`}}`&#39;
</div>

An argument is a string which may be quoted with double quotes (`"`) to include
spaces. Additionally the escape character (`\ `) can be use to escape single
characters. The start sequence `{{` scan be escaped by a preceding (`\ `)
character.

A directive can be flagged with the asterisk (`*`) character. This is used by <a href="statements.md#/statements">statements</a>
to indicate an optional nested structure (see <a href="statements.md#/statements">→2</a>).

A newline after a directive is removed from the following text.

Directives are used to formulate <a href="statements.md#/statements">statements</a>, which are used by the generator
to incluence and structure the generated document tree.

A typical <a href="#/sourcedoc">source document</a> could like in <a href="#example">→example 1-b</a>.


<a/><a id="example"/><a id="example-2"/>
<div align="center"><table><tr><td>


**Source File:** README.mdg
```
/#############################################################
/# This is a sample source document for the markdown generator
/#############################################################
{{numberrange section:V1.}}
{{numberrange figure:-a master=section:#2 abbrev=fig}}

{{section}}A Game Changing Tool

{{section introduction}}Introduction
For further details see {{link #specification}}specification section{{endlink}}.
{{endsection}}

{{section specification}}Specification
Here are the details. For a common overview, please see {{ref #introduction}}
{{endsection}}

{{endsection}}
```
</td></tr></table>
 Example 1-b: A typical source document
</br></br>
</div>


<a/><a id="/tags"/><a id="section-1-4"/>
### 1.4 Tags

A *tag* is short identifier used to identify a
dedicated *taggable element* in the <a href="#/filestructure">source tree</a>. It is described by a <a href="statements.md#/statements">statement</a>
of the <a href="README.md#section-1">*Markdown Generator*</a> in a <a href="#/sourcedoc">source document</a>.

There are three kinds of taggable elements:
- <a href="#/anchors">referencable elements</a>: Elements, which might be targets for hyperlinks.
- <a href="#/terms">terms</a>: A formal part of the terminology used in the <a href="#/filestructure">source tree</a>.
- <a href="#/textmodules">text modules</a>: A parameterized block of reusable content.

The tags are maintained in separate namespaces, therefore different
kinds of elements may carry the same tag.

A tag is a sequence of identifiers separated by a slash ('/') to formulate
hierarchical tag names. An identifier should start with a lower case letter
followed by a sequence of digits or lower case letters, for example `statements/referencables`.

There are two flavors of tags:
- *Global tags* are used to identify locations in a <a href="#/filestructure">source tree</a> by a globally
unique name. This name can be used to refer to the tagged element without
requiring to know the document, the tag is located in. They always start with
the slash character (`/`).
- *Local tags* are used to identify locations relative to the
actually generated document. They do not start with a slash character (`/`).
When used to refer to an element a local tag can pre prefixed by the name of the
document defining the tag separated by a '#' to refer to elements defined in
another document.


<a/><a id="/anchors"/><a id="section-1-4-1"/>
#### 1.4.1 Anchors for Referencable Elements

*Anchors* are a central concept of the <a href="README.md#section-1">*Markdown Generator*</a>. They are used to
establish a stable hyperlink structure among various parts of the
<a href="#/filestructure">source tree</a>.

Anchors can be attached to various types of <a href="statements.md#/statements/referencables">*referencable elements*</a>.

Anchors are <a href="#/tags">tags</a>, therefore there are again two flavors:
- *Global anchors* are used to identify locations in a document by a globally
unique name. This name can be used to establish links to this location without
requiring to know the document, the anchor is located in. They always
start with the slash character (`/`).
- *Local anchors* are used to identify referencable elements relative to the actually
generated document. They do not start with a slash character (`/`).

When used, anchors are always preceded by the number sign (`#`).
Local anchors may be preceded with the name of the source document to refer
to anchors in other documents. <a href="#/sourcedoc">Source documents</a> can be denoted by their absolute
names, or names relative to the actual document (similar to regular file path names).


<a/><a id="/scoped"/><a id="section-1-4-2"/>
#### 1.4.2 Tags and Anchors in Scoped Environments

<a href="#/tags">Taggable elements</a> might be defined in <a href="#/textmodules">text modules</a>, also. Because
such modules are intended to be used multiple times, the locally specified
<a href="#/tags">tags</a> are typically neither unique to the global <a href="#/filestructure">source tree</a> nor
for the actually generated document.

Therefore three concepts are used to resolve this problem:
- *Scopes*
- *Tag extension*
- *Tag composition*

Module instantiations provide a scope for looking up local reference names
when refering to <a href="#/anchors">referencable elements</a>, <a href="#/terms">terms</a> and <a href="#/textmodules">text modules</a>.

For text modules this is the only way to resolve name conflicts. Only
top-level definitions are visible for global reuse. Nested <a href="#/textmodules">text modules</a> are
visible for the scopes, they are defined in, only.

For referencable elements and <a href="#/terms">terms</a> additionally the second
mechanism is used. The definition is propagated up the
dynamic scope chain by prefixing the originally specified <a href="#/tags">tag</a>
step-by-step with the name of the <a href="#/scoped">scope</a> as additional namespace component.
Following the dynamic scope chain makes them visible at the document level
of the initially processed <a href="#/sourcedoc">source document</a> (<a href="#/textmodules">text modules</a> defined in
a source document might be used
from other source documents, but the tags defined in such a
<a href="#/textmodules">text module</a> are not defined in the scope of the document containing
the definition, but in the one containing the first use of a top-level
text module.

- <a href="#/tags">local tags</a> are prefixed by the scope names up the dynamic scope chain
  separated by a slash (`/`).
- <a href="#/tags">global tags</a> are prefixed by the scope names up the dynamic scope chain
   separated by a slash (`/`) and adding a leading slash to keep a global name.
   At the document level scope the name of the document is added to provide
   a name globally unique in the document tree.

In nested scopes <a href="#/tags">local tags</a> are resolved up the static nesting by using
the unextended names as specified in the source document.


<a/><a id="anchors"/><a id="example-3"/>
<div align="center"><table><tr><td>


**Source: blocks.mdg**

```
{{template}}

{{block /a}}{{param pa}}
{{block b}}{{param pb}}
{{*anchor /b1}}This is instance {{value pb}} of block b in {{value pa}}{{endanchor}}
{{endblock}}

{{*anchor a1}}This is instance {{value pa}} of block a{{endanchor}}
{{blockref bref1:#b}}{{arg pb}}first call{{endarg}}
{{blockref bref2:#b}}{{arg pb}}second call{{endarg}}


{{endblock}}
```

**Source: README.mdg**

```
{{blockref ref1:#/a}}{{arg pa}}first call{{endarg}}
{{blockref ref2:#/a}}{{arg pa}}second call{{endarg}}
```

generates anchors global anchors `/README/ref1/bref1/b1`,
`/README/ref1/bref2/b1`, `/README/ref2/bref1/b1` and `/README/ref2/bref2/b1`
and local anchors `ref1/a1` and `ref2/a1` for the generated file `README.md`.
The block tag `b` used in block `a` is resolved in the static scope instance
of block `a` using the names `ref1` and `ref2`.
</td></tr></table>
 Example 1-c: Anchors in scoped environments
</br></br>
</div>

<a href="#/tags">Tag</a> definitions inside a <a href="#/textmodules">text module</a> may be composed using scope
attributes using the syntax `{`&lt;*attribute*&gt;`}`, for example
`/statement/{scope}`. It is possible to compose <a href="#/tags">local tags</a> as well
as <a href="#/tags">global tags</a>.

The following attributes are supported:
- `{scope}`: the name of the actual scope
- `{namespace}`: the complete scope name path up to document level (separated
   by a slash (`/`).

If <a href="#/scoped">tag composition</a> is used, the implicit name extension is omitted and the
composed names must be unqiue, either globally for global tag or locally
for a generated document for local tags.



<a/><a id="/numberranges"/><a id="section-1-5"/>
### 1.5 Number ranges

*Number ranges* are another central concept. They are used to provide a potentially
hierarchical labeling of elements. The standard use case is to provide
hierarchical section labels as used in this document. Number ranges have names,
which can be used to refer to them in various <a href="statements.md#/statements">statements</a>. The
number range implicitly use for sections is `section`. There is another one
used for figures called `figure`.

Besides these implicitly defined ones, it is possible to define an arbitrary
number of other ones used for dedicated purposes in a document tree using the
`numberrange` statement.


<a/><a id="/labelformats"/><a id="section-1-5-1"/>
#### 1.5.1 Label Formats
The look and feel of number ranges can be configured. There are several ways
to express the number format used for the different hierarchy levels:

- Arabic Numbers (type `1`)
- Roman numbers  (type `i`)
- Letters (type `a`)
- Void (type `V`).

Using an upper case type name switches to an upper case format. Additionally
the level separators can be chosen among the characters `-.+*~/_#°^`.

A complete format specification looks like this:

<div align=center>

 [&lt;separator] &lt;type&gt; [&lt;separator&gt;] {&lt;type&gt; [&lt;separator]} 
</div>

A number range format could, for example look like:
<div align=center>

`V1.`
</div>

This leaves the top level sections without a number and use arabic numbers
for the deeper levels separated by a dot.

If a separator is missing, the level names are not separated by a separator.
A training separator is used for the deeper levels not specified. A leading
separator is used for the number range composition, if a number range is concatenated
to another one.



<a/><a id="/master"/><a id="section-1-5-2"/>
#### 1.5.2 Cascading Number Ranges

It is possible to control the labeling of a number range from another one.
The controlled number range is a *slave number range*, it
racts on the actual state of its *master number range*.

When generating a new label the actual state of its master is observed.
The numbering is resetted if the master state is changing. The generated label
is finally a composition of the label of the master suffixed by the label of the
slave. The concatenation is optionally separated by the leading separator
character defined for the slave (see <a href="#/labelformats">→1.5.1</a>).
The condition to reset the slave numbering can be limits to a dedicated hierarchy
level for the master.

This can be used for example for numbering figures or examples relative to a (sub)section.


<a/><a id="nr"/><a id="example-4"/>
<div align="center"><table><tr><td>


```
{{numberrange section:V1.}}
{{numberrange example:-a master=section:#2 abbrev=example}}
```

It is configured to be a slave of the section numbering using letters as number format
and resetting its numbers
whenever at least section level 2 changes (the first level is unnumbered (type `V`)).
Additionally the number range is configured to the name `example` in the label.
</td></tr></table>
 Example 1-d: The Figure Number Range used in this Documentation
</br></br>
</div>


Additionally a number range can be enriched by an abbreviation text.
This text is, for example, used to output the title for titled elements
(not sections) or for <a href="statements.md#/statement/ref">short references</a>.




<a/><a id="/terms"/><a id="section-1-6"/>
### 1.6 Formalized Terminology
The <a href="README.md#section-1">*Markdown Generator*</a> supports keeping the terminology used in a document tree consistent.
Therefore dedicated *terms*
can explicitly be defined with their name and a
short definition text. Such a term is then denoted by an own <a href="#/tags">tag</a>, which can be used
all over the document tree to use this term in the document content. If the term
should be changed, just its definition statement has to be adapted.

The term substitution with the <a href="statements.md#/statement/term">`term`</a> statement may specify
the usages of an upper case and/or plural form of the term.
The usages of the term will also automatically provide a hyperlink to the section
defining the term with the <a href="statements.md#/statement/termdef">`termdef`</a> statement.

Additionally, the defined terms with their definition texts can then be used to automatically
generate a glossary for the document tree. The index entries and all usages of
the term automatically provide hyperlinks to the section containing the definition
of the term.

<a href="#/terms">Terms</a> tags are <a href="#/scoped">scoped</a> to support their
definition in <a href="#/textmodules">text modules</a>.



<a/><a id="/textmodules"/><a id="section-1-7"/>
### 1.7 Text Modules

A *text module* is a parameterized reusable block of content.

A block definition may define parameters, whose values can be used in the
block content. A block is defined using the <a href="statements.md#/statement/block">`block`</a> statement.

Block definitions may be nested, a block definition may again declare blocks.
Instantiating blocks with the <a href="statements.md#/statement/blockref">`blockref`</a> statement resolves
the block reference up the static instantiation chain.

Top-level blocks may use <a href="#/anchors">global anchors</a>, nested blocks should use
<a href="#/anchors">local anchors</a>. <a href="#/tags">Tag</a> definitions for
<a href="statements.md#/statements/referencables">referencable elements</a> in nested
blocks may use <a href="#/scoped">tag composition</a> or are subject to <a href="#/scoped">tag extension</a> according
 to <a href="#/scoped">→1.4.2</a>.

In a block argument values for named parameters can be accessed via the
<a href="statements.md#/statement/value">`value`</a> statement. Values for parameters of outer <a href="#/scoped">scopes</a>
are resolvable as long as different parameter names are used.

