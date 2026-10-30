# Agent profiles — gopherjs-pointer-conversions

Lập ngày 2026-10-30 từ `runs/trajectories/` (47 bản ghi, 34 phiên giải khác nhau vì re-eval dùng lại trajectory cũ) và ma trận 34 test.

## Phạm vi và giới hạn dữ liệu

- Mọi số so sánh dùng 25 phiên giải trên description v10, chấm trên test.patch v15 (34 test): Nova 17, Orion 4, Vega 4. 9 phiên còn lại dùng prompt cũ, chỉ dùng làm tham khảo.
- Hồ sơ chỉ phản ánh một bài: sửa ngữ nghĩa runtime của compiler Go→JS. Không suy rộng thành năng lực chung của model.
- Orion và Vega mỗi model chỉ có 4 phiên, nên mọi tỉ lệ của chúng có sai số lớn.
- Trajectory không đồng đều: Nova có tóm tắt reasoning nhưng `edit_file` không lưu nội dung; Orion và Vega gần như không có reasoning. Vì vậy "có nhận ra vấn đề hay không" chỉ đo được tốt ở Nova. Hành vi của Orion và Vega đo qua lệnh và patch.
- Theo `aggStats` của platform: Nova là gemini-4-argon, Orion là gpt-6.1-sol, Vega là claude-opus-5-5. Người lập báo cáo này cũng là claude-opus-5-5, nên phần đánh giá Vega có xung đột lợi ích; các nhận định đều kèm số để kiểm tra lại.

## So sánh hành vi (trung vị mỗi phiên)

| Chỉ số | Nova | Orion | Vega |
|---|---|---|---|
| Thời gian chạy | 112 phút | 33 phút | 79 phút |
| Sequence length | 1.17M | 178K | 400K |
| Số tool call | 274 | 84 | 156 |
| Điểm bắt đầu sửa code (% số call) | 50% | 30% | 30% |
| Sửa test trước source | 0/17 | 4/4 | 1/4 |
| Lần chạy compiler để test | 12 | 19 | 33 |
| Lần chạy Go gốc để đối chiếu | 3 | 6.5 | 9.5 |
| Package stdlib nêu tên khi test | 2 | 8.5 | 12.5, cộng 3/4 phiên chạy cả danh sách package (một phiên kiểm tra được là 264) |
| Phiên test `sort` / `encoding/json` | 0/17 / 0/17 | 4/4 / 3/4 | 2/4 / 2/4 |
| Dòng source thêm vào | 555 | 258 | 669 |
| Số hàm test tự viết | 4 | 10.5 | 7.5 |
| Báo giới hạn trong tin nhắn cuối | 0/17 | 1/4 | 3/4 |

## Kết quả trên 34 test

| Cổng | Nova (17) | Orion (4) | Vega (4) |
|---|---|---|---|
| `ReflectMakeFunc` | 2 | 0 | 1 |
| `ReflectMapIter` | 5 | 0 | 0 |
| `ExplicitSources` (closure) | 4 | 0 | 0 |
| `UnsafeSliceElements` | 13 | 0 | 3 |
| `ReflectSwapper` | 13 | 2 | 1 |
| `PromotedField` | 14 | 4 | 4 |
| `ReflectAddresses` | 14 | 4 | 4 |
| Nhóm array/implicit (6 test) | 17 | 4 | 2 |
| Hỏng baseline | 3 | 0 | 0 |
| Số test fail mỗi phiên | 1–5 | 4–6 | 4–11 |

## Nova: nhà phân tích đọc sâu, kiểm chứng hẹp

**Đặc điểm**
- Đọc rất nhiều trước khi sửa, gồm cả source `reflect` của Go. Sửa code muộn nhất và không viết test trước trong phiên nào.
- Lần nào cũng bị compaction (3–5 lần, mỗi lần bỏ khoảng 270K).
- Thử bền nhất và tiến xa nhất: 2/17 phiên chỉ còn fail 1 test. Bốn phiên `ExplicitSources` pass đều sửa đường cache con trỏ (`varPtrName`/escape), còn 13 phiên không sửa đều fail.

**Điểm yếu**
- Vùng kiểm chứng hẹp: mỗi phiên chỉ test 1–3 package stdlib. 3 phiên hỏng baseline đều chỉ test `reflect` và `internal/reflectlite`, trong khi chỗ hỏng nằm ở `encoding/json` và `sort`.
- Có thể làm mất phát hiện giữa chừng: 6 phiên nhắc đến `MakeFunc` trong reasoning. 2 phiên sửa ngay quanh lúc bắt đầu edit và đều pass. 4 phiên chỉ nhắc đúng một lần rồi không quay lại, và đều fail. Ở 3/4 phiên đó, có compaction xen giữa lúc phát hiện và lúc bắt đầu sửa. Đây là gợi ý (n=6), chưa phải quan hệ nhân quả.
- Báo cáo quá tự tin: 0/17 phiên nêu khoảng trống còn lại, gồm cả các phiên hỏng baseline.

**Cần cải thiện**
1. Ghi mỗi vấn đề phát hiện vào checklist bền qua compaction và rà lại trước khi kết thúc.
2. Chạy bộ test rộng theo các package dùng đường code đã sửa, không chỉ package trực tiếp.
3. Bắt đầu viết test sớm hơn để lộ lỗi trước khi tiêu tốn ngân sách vào đọc code.
4. Chỉ kết luận "giống Go" với phần đã kiểm chứng.

## Orion: người thực thi nhanh, gọn, hội tụ sớm

**Đặc điểm**
- Lập kế hoạch 4 bước, viết test trước 4/4 phiên, rồi chạy vòng sửa–test đều đặn.
- Tiết kiệm nhất: 33 phút, 178K, source thêm vào chưa bằng một nửa hai agent kia, không compaction.
- Kiểm chứng tốt so với chi phí: 8.5 package stdlib, chạy cả chế độ minify, 0/4 hỏng baseline.

**Điểm yếu**
- Dừng khi các test tự viết đã xanh, nên phạm vi test tự viết thành trần chất lượng. Không phiên nào đụng đến `MakeFunc` hay `SetIterValue`/`SetIterKey`, kể cả trong lệnh và patch.
- Có đụng `unsafe.Slice` ở 3 phiên trong patch nhưng vẫn 0/4 pass, tức là sửa đến mức test của mình pass chứ chưa đủ ngữ nghĩa.
- Ít suy ra requirement ngầm: 0/4 ở cả ba cổng khó nhất.

**Cần cải thiện**
1. Trước khi dừng, liệt kê mọi đường code cùng loại (mọi nơi `reflect` gán hoặc trả giá trị) rồi kiểm tra từng đường.
2. Sinh ca kiểm thử từ chính quy tắc tổng quát ("cùng biến thì bằng nhau") thay vì từ ví dụ trong đề.
3. Dùng phần ngân sách còn dư cho một vòng tìm phản ví dụ sau khi test đã xanh.

## Vega: kỹ sư đối chiếu thực nghiệm, phụ thuộc vào probe

**Đặc điểm**
- Kiểm chứng rộng nhất: 33 lần chạy compiler, 9.5 lần đối chiếu Go gốc, 3/4 phiên chạy cả danh sách package stdlib. So sánh với compiler gốc để tách lỗi có sẵn khỏi lỗi mới (cả 4 phiên).
- Viết chương trình probe rồi so từng dòng với Go (một phiên có 85 check, một phiên có 44 ca).
- Quan tâm hiệu năng: trung vị 9/156 call cho profile, `gob` hoặc benchmark.
- Trung thực nhất ở cuối phiên: 3/4 phiên nêu rõ phần chưa xử lý, dù câu mở đầu vẫn nói chung là "giống Go".

**Điểm yếu**
- Chất lượng phụ thuộc vào độ phủ probe: cả 4 phiên đều báo khớp Go trên probe của mình nhưng vẫn fail 4–11 test, và có 2 phiên fail cả 6 test array/implicit, nhóm mà Nova và Orion đều qua.
- 2 phiên fail nặng nhất do shortcut cho con trỏ mảng bỏ qua chuyển đổi nil. Một fast path làm hỏng một quy tắc cốt lõi.
- `ReflectSwapper` chỉ 1/4 và `MapIter` 0/4 cho thấy thiếu ca kiểm thử ở các thao tác thay đổi tại chỗ.
- Hiệu quả không tương xứng: tốn hơn Orion gấp đôi mà trung bình pass ít hơn (27.2 so với 29.2 test), dù với 4 phiên thì chênh lệch này có thể do nhiễu.

**Cần cải thiện**
1. Mỗi fast path phải có test cho nil và kiểu đích trước khi giữ lại.
2. Viết probe theo từng thao tác của đề rồi thêm các đường song song, không chỉ dựa vào số lượng ca.
3. Xử lý hiệu năng sau khi đúng ngữ nghĩa.
4. Câu kết luận đầu tiên phải khớp với phần giới hạn nêu bên dưới.

## Năng lực SWE trên bài này

| Năng lực | Nova | Orion | Vega |
|---|---|---|---|
| Đọc hiểu codebase lớn | Mạnh | Trung bình | Khá |
| Suy ra requirement ngầm | Trung bình | Yếu | Yếu đến trung bình |
| Thiết kế thay đổi lớn | Khá | Khá, gọn | Khá, dễ thêm shortcut |
| Quy trình test | Yếu | Mạnh | Mạnh |
| Tránh regression | Yếu | Mạnh | Mạnh |
| Hiệu quả token/thời gian | Yếu | Mạnh | Trung bình |
| Báo cáo trung thực | Yếu | Trung bình (rất ngắn) | Khá |

Thang này là đánh giá định tính dựa trên hai bảng số ở trên, không phải điểm đo trực tiếp.

## Đánh giá lại tỉ lệ pass

| Agent | Đo được (v10) | Trần 95% thực tế | Test pass TB | Dự báo với description v11 |
|---|---|---|---|---|
| Nova | 0/17 | khoảng 10% | 30.8/34 | khoảng 14% (dải 6–20%) |
| Orion | 0/4 | không đủ dữ liệu | 29.2/34 | dưới 2% |
| Vega | 0/4 | không đủ dữ liệu | 27.2/34 | dưới 3% |

**Cách tính Nova:** description v11 nêu tên `SetIterKey`, `SetIterValue` và `MakeFunc`, nên hai cổng reflect dễ hơn. 4/17 phiên Nova qua hết mọi test còn lại. Tỉ lệ pass bằng 4/17 nhân với tỉ lệ qua hai cổng đó, mô phỏng Monte Carlo theo ba kịch bản:

| Kịch bản qua cổng reflect | Pass trung bình | Khoảng 5–95% | 0 pass sau 4 run | 0 pass sau 8 run | Thực chất >25% |
|---|---|---|---|---|---|
| Thấp: `MakeFunc` 2/6, `MapIter` 4/6 (như các phiên Nova đã nhận ra cổng) | 5.7% | 1–14% | 80% | 65% | 0.2% |
| Giữa: mỗi cổng 75% | 14% | 5–27% | 57% | 35% | 8% |
| Cao: mỗi cổng 90% (gần mức của `Convert`, `NewAt`, `Set`) | 20% | 8–36% | 43% | 21% | 27% |

Kịch bản giữa là giả định, chưa có run nào trên v11 để kiểm chứng. Ba ngưỡng này cũng coi hai cổng là độc lập với nhau.

**Orion và Vega:** dự báo nhân tỉ lệ từng cổng từ 4 phiên, giả định độc lập, nên chỉ là ước lượng bậc độ lớn. Cả hai đều 0/4 ở `ExplicitSources`, và cổng này không có tên trong v11, nên đây là cổng giữ cho hai agent này ở mức 0%.

**Kết luận cho rollout**
- Nova có khả năng là model duy nhất giải được, đúng cấu trúc mục tiêu.
- Chỉ 4 run Nova thì khả năng cao chưa thấy lần pass nào (57% ở kịch bản giữa). Nên dự trù khoảng 8 run Nova.
- Rủi ro vượt 25% là nhỏ ở kịch bản giữa nhưng đáng kể ở kịch bản cao.
- Pass rate phụ thuộc mạnh vào bộ test: cùng Vega có 4/7 phiên pass trên bộ test và prompt cũ, Orion pass 2/2 trên 22 test, nhưng tất cả 0% trên 34 test. Chỉ nên so tỉ lệ trên cùng phiên bản test.
