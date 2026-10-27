# Repository

**GitHub repo**

```
https://github.com/google/go-jsonnet
```

**Commit SHA**

```
567b61ac4a6c23546a62d79324bb4aaed6bdc941
```

# Challenge Details

**Title**

```
Add null-coalescing (??) and null-safe access (?.) operators
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Added the ?? and ?. operators to Jsonnet, implemented across the evaluator, formatter, and linter.

a ?? b returns b only if a is null; if a is not null, b is not evaluated. The ?? operator has lower precedence than all other binary operators, including ||. For example, 5 ?? 1 + 2 evaluates to 5, while false ?? true || true evaluates to false.

e?.name and e?.[index] work with any expression, including self and super. They return null if e is null, if the field is missing, or if the index is out of bounds (negative indices are considered out of bounds).

In short, a ?. chain yields the same result as standard access, except that it returns null immediately if any ?. step fails to find a value.

When this happens, the remainder of the chain is skipped and not evaluated.  The chain extends to the end of the expression or to a closing bracket, encompassing indices and function call arguments. However, only the ?. step itself suppresses errors: given o = {a: null}, o?.a.b still triggers an error.

The formatter prints a ?? b with a space on each side and a?.b without spaces. The linter does not report errors for valid code using these operators.
```
