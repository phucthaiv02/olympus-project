When a function has multiple `@overload` signatures, Jedi should return the result of the overload that matches the arguments in the call. For example, with the following code, when typing `get_obj(Foo).`, Jedi should only suggest methods of `FooBuilder`.

```python
@overload
def get_obj(input: type[Foo]) -> FooBuilder: ...
@overload
def get_obj(input: type[Bar]) -> BarBuilder: ...
```

Choose the first overload that accepts all arguments, including arguments passed by name. When comparing types, use normal typing rules. If no overload matches, combine the results of all overloads. Suppose `A` and `C` are two unrelated classes. In the `f(A if c else C)` call, the argument can be `A` or `C`, so the overload with the `x: type[A]` parameter does not match. Similarly, if `xs` has type `list[type[C]]`, then `f(*xs)` can pass any number of arguments, even no arguments, so Jedi does not choose an overload with exactly one `x: type[C]` parameter.

With tuples, lists and dicts written directly in the code, Jedi should look at the actual content that Python creates at runtime. This is still true when they are in a variable, unpacked with `*` or `**`, or nested like `f(*[A, *b])`. For example, in `{1: A, True: C}`, the two keys are equal, so the value written later overwrites the value written before, and the dict only has the value `C`. In `{"k": C, **d}`, if `d` can be `{"k": A}` or `{}`, then the value of `"k"` can be `A` or `C`. A comprehension on an empty literal, for example `[C for _ in []]`, produces no elements. With `xs = [A] if c else []`, when calling `f(*xs, C)`, `C` can be the first or second argument.
