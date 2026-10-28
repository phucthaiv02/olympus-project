# Tiếng Việt

**Title:** Overload resolution ignores type[...] annotations

Sửa cách Jedi chọn overload khi tham số được annotate bằng `type[...]`. Ví dụ có hai overload `get_obj(input: type[Foo]) -> FooBuilder` và `get_obj(input: type[Bar]) -> BarBuilder`. Khi gõ `get_obj(Foo).`, Jedi gợi ý phương thức của cả `FooBuilder` lẫn `BarBuilder`, trong khi lẽ ra chỉ gợi ý của `FooBuilder`.

Jedi cần chọn overload giống một type checker: lấy overload đầu tiên chấp nhận các đối số của lời gọi, dù đối số được truyền theo vị trí hay theo tên. Chỉ khi không overload nào khớp thì mới trả về kết quả gộp của tất cả các overload.

Tham số `type[X]` (hoặc `Type[X]`) nhận chính class `X` và các class con của nó, nhưng không nhận instance; ngược lại, tham số annotate bằng `X` không nhận class `X`. Bản `type` hoặc `Type` không có tham số, cũng như `type[Any]`, nhận mọi class. Union có thể nằm bên trong hoặc bên ngoài `type[...]`. Với `type[T]` mà `T` là một TypeVar, class phải thỏa bound của `T` hoặc một trong các constraint của nó; nếu `T` không có bound lẫn constraint thì mọi class đều được nhận. Đối số có thể là tên class viết trực tiếp, hoặc một biến hay tham số đã được annotate `type[...]`.

# English (để submit)

**Title:** Overload resolution ignores type[...] annotations

Fix how Jedi picks an overload when its parameters are annotated with `type[...]`. Given overloads `get_obj(input: type[Foo]) -> FooBuilder` and `get_obj(input: type[Bar]) -> BarBuilder`, completing `get_obj(Foo).` offers the methods of both `FooBuilder` and `BarBuilder`, when only `FooBuilder`'s belong there.

Jedi should choose the overload the way a type checker does: take the first overload that accepts the call's arguments, whether they are passed by position or by keyword, and fall back to the combined results of all overloads only when none of them fits.

A `type[X]` (or `Type[X]`) parameter accepts the class `X` and its subclasses, but not instances; in the other direction, a parameter annotated `X` does not accept the class `X` itself. Bare `type` or `Type`, and `type[Any]`, accept any class. A union may appear inside or outside `type[...]`. For `type[T]` where `T` is a TypeVar, the class has to satisfy `T`'s bound or one of its constraints, and any class is accepted when `T` has neither. The argument can be a class named directly, or a variable or parameter annotated `type[...]`.
