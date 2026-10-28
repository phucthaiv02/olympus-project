Completion is wrong after calling an overloaded function that is overloaded on classes, `None` or callables. With overloads `get_obj(input: type[Foo]) -> FooBuilder` and `get_obj(input: type[Bar]) -> BarBuilder`, `get_obj(Foo).` offers the methods of both `FooBuilder` and `BarBuilder`, because Jedi returns the results of every overload instead of picking one.

Make Jedi choose the overload the way a type checker would: take the first overload that accepts the call's arguments, and only fall back to the combined results when no overload fits.

The cases that need to work:

- `type[X]` / `Type[X]` parameters accept the class `X` or its subclasses, not instances. This includes unions, `Any`, bare `type`, and TypeVars with a bound or constraints, and the class may come from a variable or parameter annotated `type[...]`.
- `None` is accepted only where the annotation allows it (`None`, `Optional[...]`, a union containing `None`).
- `Callable` parameters accept anything callable (functions, lambdas, builtins, classes, bound methods, instances of classes that define `__call__`) and reject plain values. If the annotation states a return type, the callable's return value has to fit it; if Jedi can't infer what the callable returns, treat it as fitting. The parameter list inside `Callable[...]` doesn't need to be checked.
