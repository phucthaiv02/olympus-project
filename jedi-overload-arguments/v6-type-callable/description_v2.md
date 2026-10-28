# Tiếng Việt

**Title:** Overload resolution ignores type[...] and Callable annotations

Khi một hàm có nhiều overload và tham số được annotate bằng `type[...]` hoặc `Callable`, Jedi hiện đang gợi ý sai phương thức. Với hai overload `get_obj(input: type[Foo]) -> FooBuilder` và `get_obj(input: type[Bar]) -> BarBuilder`, khi gõ `get_obj(Foo).` Jedi gợi ý phương thức của cả `FooBuilder` lẫn `BarBuilder`, trong khi chỉ nên gợi ý của `FooBuilder`.

Jedi nên chọn overload giống một type checker: overload đầu tiên nhận được các đối số của lời gọi, dù chúng được truyền theo vị trí hay theo tên, sẽ được dùng. Chỉ khi không overload nào khớp mới trả về kết quả gộp của tất cả overload.

Một tham số `type[X]` (hoặc `Type[X]`) nhận class `X` và các subclass của nó, nhưng không nhận instance. Ngược lại, tham số annotate `X` không nhận chính class `X`. `type` hoặc `Type` không có tham số, và `type[Any]`, nhận mọi class. Union có thể nằm trong hoặc ngoài `type[...]`. Với `type[T]` mà `T` là TypeVar, class phải thỏa bound hoặc một trong các constraint của `T`; nếu `T` không có cả hai thì mọi class đều được nhận. Đối số có thể là tên class viết trực tiếp, hoặc một biến hay tham số đã được annotate `type[...]`.

Tham số `Callable` nhận mọi thứ gọi được (hàm, lambda, hàm builtin, class, bound method, và instance của class có `__call__`; như trong Python, `__call__` gán riêng trên instance thì không tính), và từ chối giá trị thông thường. Chỉ kiểu trả về của `Callable[...]` được kiểm tra, không kiểm tra danh sách tham số: mọi giá trị mà callable có thể trả về phải khớp kiểu đó. Gọi một class thì trả về instance của class đó. Nếu Jedi không suy ra được kết quả của callable thì coi như khớp.

# English (để submit)

**Title:** Overload resolution ignores type[...] and Callable annotations

Fix how Jedi picks an overload when its parameters are annotated with `type[...]` or `Callable`. Given overloads `get_obj(input: type[Foo]) -> FooBuilder` and `get_obj(input: type[Bar]) -> BarBuilder`, completing `get_obj(Foo).` offers the methods of both `FooBuilder` and `BarBuilder`, when only `FooBuilder`'s belong there.

Jedi should choose the overload the way a type checker does: take the first overload that accepts the call's arguments, whether they are passed by position or by keyword, and fall back to the combined results of all overloads only when none of them fits.

A `type[X]` (or `Type[X]`) parameter accepts the class `X` and its subclasses, but not instances; in the other direction, a parameter annotated `X` does not accept the class `X` itself. Bare `type` or `Type`, and `type[Any]`, accept any class. A union may appear inside or outside `type[...]`. For `type[T]` where `T` is a TypeVar, the class has to satisfy `T`'s bound or one of its constraints, and any class is accepted when `T` has neither. The argument can be a class named directly, or a variable or parameter annotated `type[...]`.

A `Callable` parameter accepts anything callable (functions, lambdas, builtins, classes, bound methods, and instances of classes that define `__call__`; as in Python, a `__call__` assigned on the instance doesn't count) and rejects plain values. Only the return type of a `Callable[...]` annotation is checked, not its parameter list: every value the callable may return has to fit it, calling a class returns an instance of it, and a callable whose result Jedi can't infer is treated as fitting.
