

<a/><a id="section-1"/>
# Simple Markdown Generator

Have you every written a larger set of interconnected markdown files in
your github project? Then you will for sure have been stumbled over some major
pain points in markdown:
- Links are based on section headers and file locations

  You cannot just move a section form one file into another or even
  change the section header without invalidating the links
  all over the document tree.

- Subsections of different sections may not use the same title, if they
  should be usable as link target.

- The different font sizes of section headers are nice but do not really
  help to understand the document structure, because the difference
  in the size is only recognizable in a direct comparison, but does
  no help to capture the document structure if there is text among
  the section titles.

- No more broken hyperlinks in your document tree?

- Are you missing section numbers? Manually maintaining numbers is possible,
  but then reorganizing of the structure is a hell. The numbers must be
  changed manually, and all the links are invalidated, again.

- Maintaining a common layout for description elements, which should
  look similar all over the document tree is a huge manual effort.
  For example describing structure fields or syntactical elements.

- You want to provide a highly interconnected document tree, where
  dedicated terms are always linked? Don't need to mention that this is a hell
  with pure markdown.

- You want to provide a glossary or a global table of contents?

If you feel comfortable with these problems, you can just forget about this
*Markdown Generator*.
 But if you want to get rid of this pain and want to concentrate on
the content instead of the form, then you should have a glimpse on this tool.

### Table Of Contents:
&nbsp;&nbsp;  <a href="usage.md#/usage">1 Usage</a> <br/>
&nbsp;&nbsp;  <a href="syntax.md#/syntax">2 The Document Graph</a> <br/>
&nbsp;&nbsp;  <a href="statements.md#/statements">3 Statements</a> <br/>
&nbsp;&nbsp;  <a href="examples.md#/examples">4 Examples</a> <br/>
&nbsp;&nbsp;  <a href="glossary.md#/glossary">Glossary</a> <br/>

  This tool is a markdown generator using markdown files enriched with new
  <a href="syntax.md#/directives">keywords</a> to create pure interlinked markdown files.


