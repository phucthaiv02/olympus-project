# Ứng viên task Bug Fix: repo Python

Cập nhật: 2026-10-15. Đã quét 464 repo Python bằng `repo_cap.mjs` (tra theo tên repo chính thức). Mọi repo dưới đây **không bị chặn** và **all_contributors ≤ 5**, đồng thời qua `repo_check.sh` (≥500 sao, có commit trong 12 tháng, license MIT/Apache-2.0). Các issue đều đang mở, **không có PR nào liên kết**, và không bị maintainer từ chối.

| # | Issue | Repo (hạn mức) | Độ khó | Rủi ro chính |
|---|---|---|---|---|
| 1 | gel #3346: GraphQL nhiều inline fragment cùng tên field trả `null` | gel 2/50 | Cao | Repo gần như ngừng phát triển; Docker nặng |
| 2 | gel #6713: lỗi phạm vi shape lồng nhau trong `order by` | gel 2/50 | Cao | Như trên |
| 3 | gel #6373: `group by` trên `for` có `??` mất phần tử | gel 2/50 | Cao | Như trên |
| 4 | pip #4086: `pip check` không kiểm tra dependency của extras | pip 3/50 | Trung bình | Phạm vi có thể hơi hẹp |
| 5 | myst-parser #607: `{only}` bị bỏ qua khi chứa directive/heading lồng nhau | myst-parser 5/50 | Trung bình | Có thể dính giới hạn của Sphinx |

Commit đã ghim: gel `8519106` (2025-12-24), pip `41ce840`, myst-parser `723cffc`. Lấy lại SHA đầy đủ và mới nhất bằng `repo_check.sh` trước khi tạo task.

---

## 1. gel #3346: GraphQL, nhiều inline fragment có cùng tên field trả `null`

- **Repo:** `geldata/gel`. Python, Apache-2.0, 14k⭐.
- **Issue:** https://github.com/geldata/gel/issues/3346. Nhãn `bug,graphql`.
- **Bug:** truy vấn GraphQL trên kiểu abstract kèm `... on Post { deleted }` và `... on User { deleted }`: chỉ fragment đầu tiên có giá trị, fragment sau trả `null`. Nếu field là required thì truy vấn báo lỗi. Đổi thứ tự fragment thì kết quả đổi theo.
- **Nguyên nhân (theo contributor vpetrovykh):** bộ dịch GraphQL sang EdgeQL tạo shape đa hình kiểu `[is Post].deleted, [is User].deleted`, và hai mục này trùng tên nên một mục bị đè.
- **Hướng task:** mỗi đối tượng phải nhận giá trị field theo đúng kiểu cụ thể của nó, bất kể thứ tự fragment. Áp dụng cho cả property, link, field lồng nhau, alias GraphQL và field required.
- **Các tầng:** bộ dịch GraphQL (`edb/graphql`), sinh shape EdgeQL, có thể cả compiler.
- **Kiểm tra trùng:** đã tìm "inline fragment", "polymorphic graphql". Chỉ có PR cũ #2632 và #1606 đã merge, sửa các lỗi khác. Không có gì trùng.

## 2. gel #6713: lỗi phạm vi shape lồng nhau trong `order by`

- **Issue:** https://github.com/geldata/gel/issues/6713. Nhãn `bug,compiler`. Do maintainer msullivan tạo.
- **Bug:** `select User { deck: { name } order by { asdf := User.deck.name }.asdf } filter .name = 'Alice'` báo `InternalServerError: more than one row returned by a subquery`. Nguyên nhân là `User.deck` bên trong mệnh đề của shape lồng nhau bị tính sai phạm vi.
- **Hướng task:** tham chiếu tới đường dẫn của shape cha bên trong `order by`/`filter` của shape con phải được gắn đúng với phần tử đang xét, cho kết quả giống hệt khi viết bằng `.name`.
- **Các tầng:** compiler EdgeQL (scope tree, path binding), bước sinh SQL.
- **Kiểm tra trùng:** đã tìm "nested shape order by": không có kết quả.

## 3. gel #6373: `group by` trên kết quả `for … union (… ?? …)` làm mất phần tử

- **Issue:** https://github.com/geldata/gel/issues/6373. Nhãn `bug`.
- **Bug:** một vòng `for` "upsert" dùng `(update …) ?? (insert …)` cho kết quả đúng khi `select`. Nhưng khi đổi sang `group upserts by .action` thì nhóm `insert` bị rỗng. Vẫn lỗi ở bản 4.3.
- **Hướng task:** `group` phải nhìn thấy đúng tập phần tử mà `select` trả về, kể cả khi phần tử đến từ DML bên trong `??`/`for`.
- **Các tầng:** compiler (group, DML trong biểu thức, coalesce), sinh SQL.
- **Kiểm tra trùng:** đã tìm "group by coalesce": không có kết quả.

**Rủi ro chung của gel:**
- Commit gần nhất là 2025-12-24 ("Gel Cloud to self-hosted migration guide"), có vẻ dự án đang dừng phát triển. Repo vẫn đạt mốc "có commit trong 12 tháng" **đến 2026-12-24**, nên cần nộp trước ngày đó.
- Build Docker rất nặng: phải build Postgres kèm phần mở rộng Rust, và test cần bootstrap server (`edb test`). Nên làm thử Dockerfile trước khi đầu tư viết task.

## 4. pip #4086: `pip check` không kiểm tra dependency của extras

- **Repo:** `pypa/pip`. Python, MIT, 10k⭐.
- **Issue:** https://github.com/pypa/pip/issues/4086. Nhãn `type: bug, C: extras, C: check`. Các maintainer uranusjr, pradyunsg và sbidoul đều ủng hộ sửa.
- **Bug:** nếu gói A yêu cầu `B[x]` và extra `x` của B cần C, `pip check` vẫn báo không có lỗi gì khi C chưa được cài hoặc sai phiên bản.
- **Hướng task:** `pip check` phải tính cả các extras mà các gói đã cài yêu cầu (theo marker `extra == …` trong metadata), rồi báo thiếu hoặc xung đột phiên bản theo cùng định dạng hiện tại. Cần xử lý cả marker môi trường, extra lồng extra, tên extra viết hoa thường khác nhau, và nhiều gói cùng yêu cầu các extra khác nhau.
- **Các tầng:** `operations/check`, `metadata` (các backend importlib và pkg_resources), xử lý requirement và marker.
- **Kiểm tra trùng:** đã tìm "check extras", "pip check extra". Thấy #9299 (đã đóng, cùng chủ đề) và PR #14141 (đã đóng, sửa việc đánh giá marker lúc install, không phải `check`). Cần đọc #9299 để chắc nó được đóng vì trùng với #4086.
- **Rủi ro:** phạm vi có thể chưa đủ lớn để Nova thất bại. Nên mở rộng sang toàn bộ ngữ nghĩa extras trong `check`.

## 5. myst-parser #607: directive `{only}` bị bỏ qua khi chứa directive hoặc heading lồng nhau

- **Repo:** `executablebooks/myst-parser`. Python, MIT.
- **Issue:** https://github.com/executablebooks/MyST-Parser/issues/607. Nhãn `bug`. Nhiều người xác nhận, kể cả webknjaz khi tích hợp vào pip-tools.
- **Bug:** `{only} not html` chứa `{include}` (với file Markdown), directive khác, hoặc heading lồng bên trong: nội dung vẫn lọt vào bản build html. Theo webknjaz, có vẻ thuộc tính `expr` không được truyền cho các node lồng nhau.
- **Hướng task:** nội dung trong `{only}` phải được giữ hoặc bỏ đúng theo biểu thức, bất kể bên trong là include, directive lồng hay heading, giống hành vi của `.. only::` trong RST.
- **Các tầng:** chuyển đổi từ Markdown sang docutils (directive, section), mock directive của Sphinx, include.
- **Rủi ro:** với heading, ngay cả RST cũng có giới hạn về section bên trong `only`, nên prompt phải nói rõ hành vi mong muốn bằng cách so sánh với RST. Phạm vi ở mức trung bình.
