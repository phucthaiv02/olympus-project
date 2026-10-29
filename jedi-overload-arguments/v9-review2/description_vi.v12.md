Khi một hàm có nhiều chữ ký `@overload`, Jedi nên trả về kết quả của overload phù hợp với các đối số trong lời gọi. Ví dụ với đoạn code sau, khi gõ `get_obj(Foo).` thì Jedi chỉ nên gợi ý các method của `FooBuilder`.

```python
@overload
def get_obj(input: type[Foo]) -> FooBuilder: ...
@overload
def get_obj(input: type[Bar]) -> BarBuilder: ...
```

Hãy chọn overload đầu tiên chấp nhận tất cả các đối số, kể cả các đối số được truyền theo tên. Khi so sánh kiểu, hãy dùng các quy tắc typing thông thường. Nếu không có overload nào phù hợp thì gộp kết quả của tất cả các overload. Giả sử `A` và `C` là hai class không liên quan. Trong lời gọi `f(A if c else C)`, đối số có thể là `A` hoặc `C`, nên overload có tham số `x: type[A]` không khớp. Tương tự, nếu `xs` có kiểu `list[type[C]]` thì `f(*xs)` có thể truyền vào một số lượng đối số bất kỳ, kể cả không có đối số nào, nên Jedi cũng không chọn overload có đúng một tham số `x: type[C]`.

Với tuple, list và dict viết trực tiếp trong code, Jedi nên nhìn vào nội dung thật mà Python tạo ra lúc chạy. Điều này vẫn đúng khi chúng nằm trong một biến, được unpack bằng `*` hoặc `**`, hay lồng nhau như `f(*[A, *b])`. Ví dụ trong `{1: A, True: C}`, hai key bằng nhau nên value viết sau ghi đè value viết trước, và dict chỉ còn value `C`. Trong `{"k": C, **d}`, nếu `d` có thể là `{"k": A}` hoặc `{}` thì value của `"k"` có thể là `A` hoặc `C`. Một comprehension trên literal rỗng, ví dụ `[C for _ in []]`, không tạo ra phần tử nào. Với `xs = [A] if c else []`, khi gọi `f(*xs, C)` thì `C` có thể là đối số thứ nhất hoặc thứ hai.
