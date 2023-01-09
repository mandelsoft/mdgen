[Simple Markdown Generator](README.md)

---


<a/><a id="/statements"/><a id="section-1"/>
## 3 Statements
### Table of Contents

&nbsp;&nbsp;&nbsp;&nbsp; [3.1 Referencable Document Elements](#/statements/referencables)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.1 Document Structure](#/statements/structure)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.1.1 Statement `section`](#/statement/section)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.1.2 Statement `sectionref`](#/statement/sectionref)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.2 Statement `anchor`](#/statement/anchor)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.3 Statement `figure`](#/statement/figure)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.1.4 Statement `labeled`](#/statement/labeled)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.2 Element Information](#/statements/info)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.2.1 Statement `label`](#/statement/label)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.2.2 Statement `title`](#/statement/title)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.3 Hyperlinks](#/statements/links)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.3.1 Statement `link`](#/statement/link)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.3.2 Statement `ref`](#/statement/ref)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.4 Terms](#/statements/terms)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.4.1 Statement `termdef`](#/statement/termdef)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.4.2 Statement `term`](#/statement/term)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.4.3 Statement `glossary`](#/statement/glossary)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.5 Text Modules](#/statements/textmodules)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.5.1 Statement `block`](#/statement/block)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.5.2 Statement `blockref`](#/statement/blockref)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.5.3 Statement `value`](#/statement/value)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.5.4 Statement `template`](#/statement/template)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.6 Formatting Hints](#/statements/formatting)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.6.1 Statement `center`](#/statement/center)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.7 Miscellaneous Statements](#/statements/misc)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.7.1 Statement `numberrange`](#/statement/numberrange)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.7.2 Statement `toc`](#/statement/toc)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.7.3 Statement `include`](#/statement/include)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.7.4 Statement `escape`](#/statement/escape)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [3.7.5 Statement `syntax`](#/statement/syntax)<br>
&nbsp;&nbsp;&nbsp;&nbsp; [3.8 Symbols](#/symbols)<br>

The <a href="README.md#section-1">*Markdown Generator*</a> uses special *statements* to control the generation of the markdown files.
A statement is a sequence of one or more <a href="syntax.md#/directives">directives</a>. If it consists of
multiple directives, there might be other content in between, which is used
by the statement to influence and structure the generation of a document tree.

If a statement is composition of multiple directives in mode cases it is finished
by an appropriate end directive, this uses the <a href="syntax.md#/directives">keyword</a> `end` or alternatively
the more expressice `end<statement-name>`, for example `{{link #tag}}Title{{endlink}}`.

The following statements are used:


<a/><a id="/statements/referencables"/><a id="section-1-1"/>
### 3.1 Referencable Document Elements

Referencable elements provide <a href="syntax.md#/anchors">anchors</a> which can be used
to establish hyperlinks with several other <a href="#/statements/links">statements</a>.


<a/><a id="/statements/structure"/><a id="section-1-1-1"/>
#### 3.1.1 Document Structure

As for regular markdown a single document may be hierarchically structured
in sections. But the <a href="README.md#section-1">*Markdown Generator*</a> additionally supports to span
a section hierarchy over multiple documents of the <a href="syntax.md#/filestructure">source tree</a>.
The <a href="syntax.md#/numberranges">numbering</a> of the sections will then
be consistently provided for the complete document tree.


<a/><a id="/statement/section"/><a id="section-1-1-1-1"/>
##### 3.1.1.1 Statement `section`
#### Synopsis
`{{section` [&lt;*anchor*&gt;]`}}` &lt;*title*&gt; *&lt;newline*&gt; ... `{{endsection}}`


#### Description
A section is the structural element of the document tree. It uses the <a href="syntax.md#/numberranges">number range</a> `section` to receive
a numbering. If sections are nested the sub-level elements get appropriate sub level labels.
In contrast to the native markdown heading feature sections can carry a logical (stable) <a href="syntax.md#/anchors">anchor</a>
independent of the heading text used by <a href="#/statements/links">other statements</a> to establish hyperlinks.

The section structure is not limited to a single document, it might span multiple <a href="syntax.md#/sourcedoc">source documents</a>.


<a/><a id="/statement/sectionref"/><a id="section-1-1-1-2"/>
##### 3.1.1.2 Statement `sectionref`
#### Synopsis
  `{{*sectionref` &lt;*ref*&gt;`}}` ... `{{endsectionref}}`\&lt;br&gt;
  `{{sectionref` &lt;*ref*&gt;`}}`


#### Description
This statement is used to link the section structure of another <a href="syntax.md#/sourcedoc">source document</a> at the actual
location of the actual <a href="syntax.md#/sourcedoc">source document</a>. The top level section of the document referred to by the
reference is added to the actual section hierarchy, the target document must have only one top level
section.

The flagged variant can be used to create an inline toc line in the actual document, see statements
<a href="#/statement/label">`label`</a> and <a href="#/statement/title">`title`</a>. The given reference is set as current
reference for the nested content.


<a/><a id="reftoc"/><a id="example-1"/>
<div align="center"><table><tr><td>


```
### Table Of Contents:
  {{*sectionref #/syntax    }}{{label}} {{title}}{{end}} <br/>
  {{*sectionref #/statements}}{{label}} {{title}}{{end}} <br/>
  {{*sectionref #/examples  }}{{label}} {{title}}{{end}} <br/>
  {{link #/glossary}}{{title #/glossary}}{{end}} <br/>
```
</td></tr></table>
 Example 3-a: Explicit Table of Contents for Linked Documents
</br></br>
</div>


<a/><a id="/statement/anchor"/><a id="section-1-1-2"/>
#### 3.1.2 Statement `anchor`
#### Synopsis
  `{{*anchor` [ &lt;*numberrange*&gt; &#39;`:`&#39;] [&#39;`!`&#39;] &lt;*anchor*&gt; `}} &lt;*caption text*&gt; {{endanchor}}`</br>
  `{{anchor` [ &lt;*numberrange*&gt; &#39;`:`&#39;] &lt;*anchor*&gt; `}}


#### Description
Define an <a href="syntax.md#/anchors">anchor</a> for the actual location. The optional caption text is used as title.
Additionally the anchor is labeled with the specified <a href="syntax.md#/numberranges">number range</a>, If no
number range is given the name `anchor` is used.
The title and label is avaiable as tag information, for example for the
<a href="#/statement/title">`title`</a> or <a href="#/statement/label">`label`</a> statement.


<a/><a id="/statement/figure"/><a id="section-1-1-3"/>
#### 3.1.3 Statement `figure`
#### Synopsis
`{{figure` [ &lt;*anchor arg*&gt; ] &lt;*filepath arg*&gt; { &lt;*attribute arg*&gt; } `}} &lt;*caption text*&gt; {{endfigure}}`


#### Description
Add a centered image to the output with a caption. This statement uses the <a href="syntax.md#/numberranges">number range</a> to label the
caprions. The caption prefixed with the label and potential number range name abbreviation
is placed below the image.
HTML image attributes may be given by additional attribute arguments of the form &lt;name&gt;`=`&lt;value&gt;,
for example `width=800`.


<a/><a id="/statement/labeled"/><a id="section-1-1-4"/>
#### 3.1.4 Statement `labeled`
#### Synopsis
`{{labeled` &lt;*numberrange*&gt; [ &#39;`:`&#39; &lt;*anchor*&gt; ]  &lt;*mode arg*&gt;`}} &lt;*caption text*&gt;  {{content}} &lt;*content*&gt; {{endlabeled}}`


#### Description
Add content to the output, which carries a caption and label according to the given <a href="syntax.md#/numberranges">number range</a>.
The output mode can be influence by a second argument:
- `box`: (default) the content is placed in a framed box, which is centered together with the caption to visually
  separate the complete element from the rest of the text flow.
- `float`: the content just follows the <a href="syntax.md#/anchors">anchor</a> followed by a centered caption.
   So, the content has complete control over its formatting.





<a/><a id="/statements/info"/><a id="section-1-2"/>
### 3.2 Element Information

<a/><a id="/statement/label"/><a id="section-1-2-1"/>
#### 3.2.1 Statement `label`
#### Synopsis
`{{label` [&lt;*ref*&gt;] `}}`


#### Description
The label of the referenced element is added to the document. All elements supporting a <a href="syntax.md#/numberranges">number range</a>
can be used, e.g. a <a href="#/statement/section">`section`</a>.
If no reference is given the current reference is used. The used reference is set as current reference.


<a/><a id="/statement/title"/><a id="section-1-2-2"/>
#### 3.2.2 Statement `title`
#### Synopsis
`{{title` [&lt;*ref*&gt;] `}}`


#### Description
The title of the referenced element is added to the document. All elements supporting a <a href="syntax.md#/numberranges">number range</a>
can be used, e.g. a <a href="#/statement/section">`section`</a>.
If no reference is given the current reference is used. The used reference is set as current reference.



<a/><a id="/statements/links"/><a id="section-1-3"/>
### 3.3 Hyperlinks

There are several statements used to generate hyperlinks.


<a/><a id="/statement/link"/><a id="section-1-3-1"/>
#### 3.3.1 Statement `link`
#### Synopsis
`{{link` &lt;*ref*&gt; `}}` &lt;*content*&gt; &#39;`{{endlink}}`&#39;


#### Description
Establish a hyperlink to some other part of the document tree on the embedded content. All elements
providing <a href="syntax.md#/anchors">anchors</a> can be used to link to.


<a/><a id="/statement/ref"/><a id="section-1-3-2"/>
#### 3.3.2 Statement `ref`
#### Synopsis
`{{ref` [&#39;`*`&#39; [&#39;`^`&#39;]] &lt;*ref*&gt; `}}`


#### Description
Establish a hyperlink on the label of a referenced element (see <a href="syntax.md#/anchors">→2.4.1</a>).
If the asterisk (`*`) is given a the label is preceded with the abbreviation text of the
<a href="syntax.md#/numberranges">label type</a>. If additionally the `^` prefix is given, the
abbreviation text will be converted to upper case first.



<a/><a id="/statements/terms"/><a id="section-1-4"/>
### 3.4 Terms

<a href="syntax.md#/terms">Terms</a> can be used to formally define a <a href="syntax.md#/tags">tag</a> and assign it to  a <a href="syntax.md#/tags">tag</a>.
This tag can then be used all over the document to substitute the centrally
defined text and to automatically establish a hyperlink to the section of its definition.

Additionally, the defined terms with their definition texts can then be
used to automatically generate a glossary for the document tree. For more details see
<a href="syntax.md#/terms">→2.6</a>.


<a/><a id="/statement/termdef"/><a id="section-1-4-1"/>
#### 3.4.1 Statement `termdef`
#### Synopsis
`{{termdef` [&#39;`*`&#39; | &#39;`-`&#39;] &lt;*tag*&gt;`}}` &lt;*term name as content*&gt; `{{description}}` &lt;*glossary content*&gt; `{{endtermdef}}`


#### Description
The given text is defined as <a href="syntax.md#/terms">term</a> with the given logical <a href="syntax.md#/tags">tag</a>.
The defined tag can be used all over the document tree with the <a href="#/statement/term">`term`</a>.
The definition of a term should be placed at the location, where the term is explained.
The usage of the term will automatically be linked to the section containing the term definition.

To support the various usage scenarios of the term, the term text should always be
specified in its singular form. To provides explicit highlighting the term may use
the markdown formatting characters `*`, `_` or <code>`</code>.

If a term is used the specification of the term tag determines the output mode.
If an asterisk (`*`) is used the plural form is substituted. The term tag is always
defined in its lower case form. If the specified tag uses a first upper case letter,
the upper case form is substituted.

By default the `termdef` statement outputs the defined term according to the output
mode derived from the tag in a highlighted format to indicate its definition location.
If the `-` prefix is given, the output is omitted, and the statement just defines
the term for usage elsewhere.
its definition.


<a/><a id="/statement/term"/><a id="section-1-4-2"/>
#### 3.4.2 Statement `term`
#### Synopsis
`{{term` [&#39;`!`&#39;] ([&#39;`#`&#39;] | [&#39;`*`&#39;]) &lt;*tag*&gt;`}}`


#### Description

This <a href="#/statements">statement</a> outputs a <a href="syntax.md#/terms">term</a> previously defined with the <a href="#/statement/termdef">`termdef`</a>
statement. Unless the prefix `!` a hyperlink to the section defining the term will be added to the term.

The used tag specification determines the output mode. The term tag is always defined
in its lower case form. If an asterisk (`*`) is used the plural form is substituted,
if the tag uses a first upper case character, the upper case form is substituted.
If the prefix `#` is given, instead of the term the label of the section containing
the term definition is used.


<a/><a id="/statement/glossary"/><a id="section-1-4-3"/>
#### 3.4.3 Statement `glossary`
#### Synopsis
`{{glossary` [&lt;*prefix*&gt;]`}}`


#### Description

This <a href="#/statements">statement</a> outputs a glossary, an alphabetical index of the used <a href="syntax.md#/terms">terms</a>
in the document tree. For an example, please refer to our own <a href="glossary.md">glossary</a>.

The optional prefix can be used to restrict the glossary to a dedicated term tag prefix.


<a/><a id="/statements/textmodules"/><a id="section-1-5"/>
### 3.5 Text Modules

<a href="syntax.md#/textmodules">text modules</a> can be used to define reusable and parameterized block content
to
- share the same content at different places in the document tree
- provide parameterized text patterns for descriptions of the same kind, which
  should look identical for every incarnation. For example, this complete document is
  based on a statement pattern.


<a/><a id="/statement/block"/><a id="section-1-5-1"/>
#### 3.5.1 Statement `block`
#### Synopsis
`{{block` &lt;*tag*&gt;`}}` { &lt;*parameter*&gt; } &lt;*block content*&gt; `{{endblock}}`

  A parameter might optionally be defaulted. Without a default:

  `{{param &lt;*name*&gt; { &#39;`,`&#39; &lt;*name*&gt; } `}}`

  and with a default:

  `{{*param` &lt;*name*&gt; { &#39;`,`&#39; &lt;*name*&gt; } `}}` &lt;*content*&gt; `{{endparam}}`
  


#### Description

This <a href="#/statements">statement</a> defines a <a href="syntax.md#/textmodules">text module</a>, which can be referred to using the
specified tag, it might be a <a href="syntax.md#/anchors">global anchor</a> or a <a href="syntax.md#/anchors">local anchor</a>. Optionally
the block definition may declare parameters. A parameter may be defaulted using the
flagged syntax for the `param` directive.



<a/><a id="/statement/blockref"/><a id="section-1-5-2"/>
#### 3.5.2 Statement `blockref`
#### Synopsis
`{{blockref` [ &lt;*name*&gt; &#39;`:`&#39;] &lt;*ref*&gt;`}}` { `{{arg` &lt;*name*&gt; `}}` &lt;*content*&gt; `{{endarg}} }`


#### Description

This <a href="#/statements">statement</a> instantiates a <a href="syntax.md#/textmodules">text module</a> defined with the
statement <a href="#/statement/block">`block`</a>, which can be referred to using the
specified tag, it might be a <a href="syntax.md#/anchors">global anchor</a> or a <a href="syntax.md#/anchors">local anchor</a>. Optionally
a name for the instantiation can be set. It is used to label the scope
provided by the instantiated text module. If no explicit name is given,
an impicit name is generated, numbered according to the instantiation count of
the referenced block in the actual scope.




<a/><a id="/statement/value"/><a id="section-1-5-3"/>
#### 3.5.3 Statement `value`
#### Synopsis
`{{value` [&#39;`*`&#39;]&lt;*parameter name*&gt; `}}`


#### Description
Inside a <a href="syntax.md#/textmodules">text module</a> body this <a href="#/statements">statement</a> is used to
access the argument value of a parameter. Parameter names are resolved
up the static scope chain, this means an inner text module may access
argument values of outer text modules.

If the name is prefixed with a asterisk (`*`) the value of the appropriate
<a href="syntax.md#/scoped">scope</a> attribute is substitute.


<a/><a id="/statement/template"/><a id="section-1-5-4"/>
#### 3.5.4 Statement `template`
#### Synopsis
`{{template}}`


#### Description
This statement can be used to flag a <a href="syntax.md#/sourcedoc">source document</a> to be omitted
from the generation of a markdown file. Neverthess, the content is interpreted
and the <a href="syntax.md#/textmodules">text modules</a> defined in this file are available to be used
in other source documents.





<a/><a id="/statements/formatting"/><a id="section-1-6"/>
### 3.6 Formatting Hints


<a/><a id="/statement/center"/><a id="section-1-6-1"/>
#### 3.6.1 Statement `center`
#### Synopsis
`{{center}}` &lt;*content*&gt; `{{endcenter}}`


#### Description
This <a href="#/statements">statement</a> centers the lines of the embedded content.


<a/><a id="/statements/misc"/><a id="section-1-7"/>
### 3.7 Miscellaneous Statements


<a/><a id="/statement/numberrange"/><a id="section-1-7-1"/>
#### 3.7.1 Statement `numberrange`
#### Synopsis
`{{numberrange` &lt;*name*&gt; [&#39;`:`&#39; &lt;*format*&gt;] [ &#39;`:#`&#39; &lt;*starting heading level*&gt;] { &lt;*attribute arg*&gt; } }}

  


#### Description
This <a href="#/statements">statement</a> declares and/or configures a <a href="syntax.md#/numberranges">number range</a>. It can
only be used at the top-level <a href="syntax.md#/scoped">scope</a> outside of any other statement
block.

The format can be specified as described in <a href="syntax.md#/numberranges">→2.5</a>. Optionally the
heading level can be specified, which should be used for the first hierarch level (default is 1).

The optional atribute arguments are of the form &lt;*attr*&gt; `=` &lt;*value*&gt;.
The following attributes are supported:
- `master=`&lt;*name*&gt;[`:`&lt;*level*&gt;]: the name of the number range to be used as <a href="syntax.md#/numberranges">master</a>
- `abbrev=`&lt;*text*&gt;: the abbreviation name of the number range used to prefix a label.



<a/><a id="/statement/toc"/><a id="section-1-7-2"/>
#### 3.7.2 Statement `toc`
#### Synopsis
`{{toc` [&lt;*ref*&gt;] `}}`


#### Description
This <a href="#/statements">statement</a> outputs a table of contents. If a reference is specified
the table is limited to the given section.



<a/><a id="/statement/include"/><a id="section-1-7-3"/>
#### 3.7.3 Statement `include`
#### Synopsis
`{{include` &lt;*path argument*&gt; `}}`


#### Description
This statement can be used to include the content of a file. The content is
not interpreted, it is just forwarded to the generated output.

If interpreted content should be provided in a reusable manner a
<a href="syntax.md#/textmodules">text module</a> has to be used. Using the <a href="#/statement/template">`template`</a> statement
the generation of a markdown document for a <a href="syntax.md#/sourcedoc">source document</a> can be omitted.



<a/><a id="/statement/escape"/><a id="section-1-7-4"/>
#### 3.7.4 Statement `escape`
#### Synopsis
`{{escape}}` &lt;*content*&gt; `{{endescape}}`


#### Description
The content of this <a href="#/statements">statement</a> is HTML-escaped. Breaking rules (`</br>`)
are not escaped.



<a/><a id="/statement/syntax"/><a id="section-1-7-5"/>
#### 3.7.5 Statement `syntax`
#### Synopsis
`{{syntax}}` &lt;*expression*&gt; `{{endsyntax}}`


#### Description
  It might be a line-based list of EBNF-like rules or simple EBNF-like expressions.

  <table><tr><td>
  <div>

&lt;*expression*&gt; = [ &lt;*identifier*&gt; '`=`' ] &lt;*syntaxexpr*&gt;</br>
</div>
</td><td>
  for a rule
  </td></tr><tr><td>
    <div>

'`{`' &lt;*syntaxexpr*&gt; '`}`' [ '`+`' ]</br>
</div>
</td><td> an arbitray number of occurrences,
    if a (`+`) is added, at least one occurrence is required.
  </td></tr><tr><td>
  <div>

'`[`' &lt;*syntaxexpr*&gt; '`]`'</br>
</div>
 </td><td>an optional expression.
   </td></tr><tr><td>
  <div>

'`<`' &lt;*identifier*&gt; '`>`'</br>
</div>
</td><td>an identifier for a rule.
   </td></tr><tr><td>
   <div>

&lt;*literal*&gt;</br>
</div>
</td><td>any other character sequence is taken a literal.
    </td></tr><tr><td>
      <div>

&lt;*syntaxexpr*&gt; '`|`' &lt;*syntaxexpr*&gt;</br>
</div>
</td><td>the left or the right expression.
   </td></tr><tr><td>
   one space</td><td> a single space.
   </td></tr><tr><td>
      two spaces</td><td> an arbitrary number of spaces</td></tr><tr><td>
      three spaces</td><td>at least one space.
  </td></tr></table>


  
<a/><a id="syntaxexpr"/><a id="example-2"/>
<div align="center"><table><tr><td>


  <div align="center">

  ```
  "<argument>{   <argument>}"
  ```
  </div>

  is rendered as
  <div align="center">

  <div>

&lt;*argument*&gt; { {'` `'}+ &lt;*argument*&gt; }</br>
</div>

  </div>

  </td></tr></table>
 Example 3-b: An example for a syntax expression
</br></br>
</div>

  



<a/><a id="/symbols"/><a id="section-1-8"/>
### 3.8 Symbols

The following <a href="syntax.md#/directives">directives</a> provide fixed symbols:
- `nl`: The newline character.
- `cs`: The start sequence (`/#`) of a comment.
