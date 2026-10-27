# Repository

**GitHub repo**

```
https://github.com/executablebooks/MyST-Parser
```

**Commit SHA**

```
723cffcf84213f0cb58695b27eec9ad72052b53a
```

**GitHub issue**

```
https://github.com/executablebooks/MyST-Parser/issues/607
```

# Challenge Details

**Title**

```
Headings inside only and ifconfig directive content escape the directive
```

**Language:** Python

**Category:** Bug Fix

**Description**

```
Fix Markdown headings written inside the content of the Sphinx `only` and `ifconfig` directives. Today such a heading escapes the directive: its section is attached to the document's own section structure instead of to the directive, so it is rendered no matter what the `only` expression or the `ifconfig` condition says, it shows up in the toctree for every builder, and whatever follows it can end up in the wrong place. The same happens when the heading comes from a Markdown file pulled in with `include` inside one of these directives. Plain paragraphs and nested directives in the same content are already filtered correctly; only headings are affected.

A heading inside these directives must belong to the directive. Its section, together with everything under it, is kept or dropped with the directive's content, so with `{only} not html` it is absent from HTML output but present in, say, the text builder, and with `{ifconfig}` it follows the configuration value. Headings inside the content nest among themselves by level exactly as they would anywhere else in the document, and that includes headings from an included file, with any `heading-offset` applied. Content that comes after the directive continues in the section that contained the directive, as if the headings inside it were not there.

The resulting document structure must be the same as what Sphinx produces for the equivalent reStructuredText document, where `.. only::` and `.. ifconfig::` contain section titles whose adornments correspond to the Markdown heading levels. Concretely, for `only`: when the headings inside are deeper than the section containing the directive, the directive stays where it is and its sections nest inside it; when the first heading inside is at the same level as, or a higher level than, the section containing the directive, the whole `only` block moves up so that it sits next to the sections of that level, directly after the section it was written in, while the content around it keeps its place. For `ifconfig`, the directive always stays where it is written and its sections nest inside it. Nested `only` blocks follow the same rules. Because the sections are now part of the `only` block, the toctree lists them only for builders where the expression holds.

Headings inside directives that do not accept section titles, such as admonitions, are unchanged and keep turning into rubrics.
```

# Research notes

Olympus draft: `kh78a7rc9f9gpet5y0wecydq7h8ga4yp` (version `kd73r3acmzcnqtrr8zmj97pj458gb5sc`).

- Prechecks (2026-10-15): GitHub PASS; Problem & Tests WARNING (coverage: rubric case for admonitions and same-level `ifconfig` not tested; conciseness suggestions); Dockerfile WARNING (editable install, known rubric conflict); Solution patch PASS.
- Scope Gate (2026-10-15): **DROP**, publicly-solved. Open PR #1102 "Support headings in nested parse" (https://github.com/executablebooks/MyST-Parser/pull/1102) puts nested headings under the temp root node. That covers 9/16 hidden tests (56%). The parts it does not cover are moving same/higher-level `only` blocks up (title styles in the memo) and isolating sibling sections. Platform overlap: 0% with other submissions.
