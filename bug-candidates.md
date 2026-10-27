# Ứng viên task Bug Fix

Cập nhật: 2026-10-15. Hạn mức lấy từ `repo_cap.mjs` (tra theo tên repo chính thức): mọi repo dưới đây **không bị chặn** và **all_contributors ≤ 5**. Các issue đều đang mở, **không có PR nào liên kết** (đã kiểm tra bằng `issue_check.sh show`), và không bị maintainer từ chối.

| # | Issue | Repo (hạn mức) | Độ khó | Rủi ro chính |
|---|---|---|---|---|
| 1 | typst #4359: tham chiếu tới figure nổi trỏ sai vị trí | typst 5/50 | Cao | Hạn mức chạm 5 |
| 2 | typst #6297: figure cột và figure parent-scope chồng lên nhau | typst 5/50 | Cao | Ít thảo luận |
| 3 | typst #5124: `place.flush` không áp dụng cho float đáy | typst 5/50 | Trung bình–cao | Có thể hơi hẹp |
| 4 | typst #5765: footnote mất trong float `scope: "parent"` | typst 5/50 | Trung bình–cao | Patch lộ trong gist |
| 5 | typst #8890: mất dòng của block khi float bị dời trang | typst 5/50 | Cao | Người báo lỗi có patch riêng |
| 6 | trino #17209: correlated `IN` sai với giá trị không xác định được | trino 2/50 | Trung bình–cao | Cùng nhóm với PR đã đóng |
| 7 | rustfmt #5464: comment giữa tên field và `:` bị xóa | rustfmt 4/50 | Trung bình (yếu) | Phạm vi hẹp, đã có người claim |

Commit đã ghim: typst `d7dfb34ae79689ed78857ab38a3f6a24a1747770`, trino `5706548dade9b1dc81144af6850c5dde9a356713`, rustfmt `c69af7786cdfef5da968f53eebd6c9b1af242b5f`. Lấy lại SHA mới nhất bằng `repo_check.sh` trước khi tạo task.

---

## 1. typst #4359: tham chiếu tới figure nổi trỏ sai vị trí (khuyên chọn)

- **Issue:** https://github.com/typst/typst/issues/4359. Nhãn `bug,introspection`, 15 reactions.
- **Bug:** figure có `placement: top/bottom/auto` được ghi nhận vị trí ở chỗ khai báo, không phải chỗ thực sự đặt. Hậu quả: `locate(<fig>).position()` sai, link của `@ref` trong PDF nhảy sai chỗ, outline danh sách hình in sai số trang.
- **Hướng task:** vị trí, link và số trang trong outline phải theo chỗ figure thực sự được đặt. Riêng số thứ tự figure và mọi counter vẫn giữ theo thứ tự trong luồng văn bản.
- **Các tầng:** introspection (Locator, Tag), layout (flow, float, `scope: "parent"`), outline, link trong PDF và HTML.
- **Vì sao khó:** maintainer (laurmaedje) nhận xét đây là bug khó sửa, vì vị trí phải nằm trong luồng thì counter mới tính đúng. Sửa ngây thơ sẽ làm hỏng cách đánh số.
- **Kiểm tra trùng:** đã tìm "floating figure reference", "float location": không có PR nào.
- **Rủi ro:** prompt phải chốt rõ quy tắc "counter theo luồng, vị trí theo chỗ đặt".

## 2. typst #6297: figure cột và figure parent-scope chồng lên nhau

- **Issue:** https://github.com/typst/typst/issues/6297. Nhãn `bug,layout`. Collaborator MDLC01 đã xác nhận bằng ảnh chụp.
- **Bug:** trong trang 2 cột, figure `placement: top` trong cột và figure `scope: "parent"` trên cùng trang có thể được đặt đè lên nhau.
- **Hướng task:** float trong cột và float trải ngang trang phải chia sẻ không gian. Float trong cột phải nằm dưới (hoặc trên) vùng đã bị float parent chiếm, không được chồng lên nhau, kể cả khi float bị dời sang trang sau.
- **Các tầng:** layout (distribute float theo cột và theo trang, tính lại vùng còn trống), có thể cả flush và pagebreak.
- **Vì sao khó:** thứ tự giữa layout cột và layout trang, float bị dời trang, nhiều float cùng lúc, kết hợp với `clearance`.
- **Rủi ro:** issue ít thảo luận. Test so sánh vị trí qua `locate().position()`, cách này ổn định hơn so ảnh.

## 3. typst #5124: `place.flush` không áp dụng cho float đáy

- **Issue:** https://github.com/typst/typst/issues/5124. Nhãn `bug,layout`. Người báo là EpicEricEE (contributor).
- **Bug:** `#place.flush()` lẽ ra phải đẩy nội dung phía sau xuống dưới mọi float đang chờ, nhưng float `bottom` vẫn bị đặt sau nội dung đó.
- **Hướng task:** `place.flush` đảm bảo mọi float đang chờ (top, bottom, auto, cả trong cột lẫn parent-scope) được đặt xong trước nội dung theo sau.
- **Liên quan:** #6821 "`place.flush` doesn't work for some use cases" (cũng chưa có PR) có thể gộp làm cùng một ngữ nghĩa nếu đúng một lỗi gốc. Cần đọc kỹ trước khi gộp.
- **Rủi ro:** nếu chỉ sửa trường hợp bottom thì phạm vi hẹp, nên bao quát cả float trong cột và parent-scope.

## 4. typst #5765: footnote mất trong `place(float: true, scope: "parent")`

- **Issue:** https://github.com/typst/typst/issues/5765. Nhãn `bug,layout`, 9 reactions.
- **Bug:** trong trang nhiều cột, footnote nằm trong khối float parent-scope chỉ hiện dấu, phần nội dung không được in ra và không có cảnh báo. Đây lại là cách tutorial chính thức hướng dẫn làm khối tiêu đề.
- **Các tầng:** layout (flow, footnote, float, cột), introspection.
- **Rủi ro:** ① chưa chốt footnote nên nằm ở cột trái hay trải ngang trang, prompt phải tự quyết. ② Trong issue có link gist chứa **patch do Claude viết**, tức lời giải lộ công khai.

## 5. typst #8890: mất dòng của block khi float bị dời sang trang sau

- **Issue:** https://github.com/typst/typst/issues/8890. Nhãn `bug,layout`. Maintainer laurmaedje xác nhận và đóng #8487 vì trùng.
- **Bug:** block có thể ngắt trang, khi tràn sang trang có float bị dời tới, sẽ **mất một dòng** (ví dụ dòng "D" biến mất). Nếu dòng bị mất chứa link thì xuất PDF bị lỗi internal error.
- **Hướng task:** không được mất hay lặp nội dung khi block ngắt trang gặp float bị dời. Xuất PDF phải không lỗi.
- **Các tầng:** layout flow (MultiSpill, deferred float), xuất PDF (link, tagging).
- **Rủi ro:** người báo lỗi nói họ đã có patch riêng và đề nghị mở PR, nên có thể sớm có PR. Kiểm tra lại trước khi làm.

## 6. trino #17209: correlated `IN` trả kết quả sai với giá trị không xác định được

- **Issue:** https://github.com/trinodb/trino/issues/17209. Nhãn `correctness`, do maintainer kasiafi tạo.
- **Bug:** `SELECT a IN (SELECT ARRAY[null] WHERE b = 2) FROM (SELECT ARRAY[null], 2) t(a, b)` trả `false`, đúng ra phải là `NULL`. Lỗi nằm ở rule optimizer `TransformCorrelatedInPredicateToJoin`.
- **Các tầng:** optimizer (rule decorrelation), cách join so sánh, ngữ nghĩa so sánh các giá trị cấu trúc (array/row) có chứa null.
- **Rủi ro:** issue này thuộc issue tổng #17213. Issue tổng có 2 PR đã đóng (#20340, #30735), nhưng hai PR đó sửa đường value-list và SemiJoin, không phải đường correlated. Repo Java rất lớn nên build Docker nặng.

## 7. rustfmt #5464: comment giữa tên field và `:` bị xóa (ứng viên yếu)

- **Issue:** https://github.com/rust-lang/rustfmt/issues/5464. Nhãn `C-bug,A-comments,P-medium`.
- **Bug:** `field /* comment */ : i32` trong struct và enum variant bị xóa comment mà không có cảnh báo.
- **Rủi ro:** phạm vi hẹp nên dễ bị Nova giải. Có người claim từ 2024 nhưng chưa có PR. Maintainer muốn tách riêng trường hợp `pub /* c */ fn`, nên không được gộp vào.
