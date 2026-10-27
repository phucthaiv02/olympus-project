# Project Memory

Instructions here apply to this project and are shared with team members.

## Context

## Workflow

Mỗi task đi qua pipeline theo thứ tự, mỗi stage phải PASS mới sang stage sau:

| # | Stage | Ai chạy | PASS khi | Ghi vào log |
|---|---|---|---|---|
| 1 | Task generator | skill `task-generator` | có repo, commit, issue, title, description (tiếng Việt rồi dịch sang tiếng Anh) | mode (A tìm mới / B sinh lại), findings đã xử lý, version mới, file đã đổi |
| 2 | Quality checks | skill `quality-check` | không còn finding HIGH | entry đầy đủ của skill |
| 3 | Platform quality checks | chạy trên platform | mọi check của platform qua | từng check: PASS/FAIL, nguyên văn lý do khi fail, mức độ |
| 4 | Platform agent runs | chạy trên platform | bộ run đạt mục tiêu rollout (xem `PROCESS.md`) | từng run: model, kết quả, số test pass, tỉ lệ pass, đường dẫn thư mục run |
| 5 | FP check | skill `fp-check` | không còn run lọt FP mức HIGH | entry đầy đủ của skill |
| 6 | Platform FP check | chạy trên platform | FP check của platform qua (pass with warnings tính là PASS, ghi cảnh báo thành finding MEDIUM/LOW) | verdict, nguyên văn feedback, probe nào bị coi là hợp lệ |
| 7 | Auto review | platform | outcome không yêu cầu sửa | issues theo mức độ |
| 8 | Expert review | reviewer người | được chấp nhận | từng nhận xét kèm mức độ |
| 9 | Submit | người dùng | đã nộp | version id, ngày |

Quy tắc chạy:
- Không bắt đầu stage sau khi stage trước chưa PASS và còn mới. Stage 4 tốn nhiều nhất, nên chỉ chạy khi stage 2 và 3 đã PASS trên đúng version hiện tại.
- **Fail ở bất kỳ stage nào thì ghi entry (Result FAIL, findings có mức độ và bằng chứng) rồi quay về stage 1.** Task generator đọc log, tự quyết Mode A (task bị reject, tìm repo/issue mới) hay Mode B (sinh lại description theo findings), rồi ghi entry stage 1 của vòng mới (Round + 1).
- Sau khi sinh lại, mọi stage từ 2 trở đi chạy lại từ đầu trên version mới; không dùng lại kết quả PASS cũ, vì mọi chỉnh sửa làm kết quả cũ hết hiệu lực.
- Phần test và solution bị finding yêu cầu sửa thì người làm task sửa theo kết quả của task generator (ví dụ bỏ một nhóm tính năng thì bỏ luôn test và solution tương ứng) trước khi sang stage 2.

## Check log

- Mỗi task có một file log riêng, đặt trong thư mục của task: `<thư-mục-task>/CHECK-LOG.md` (ví dụ `jedi-overload-arguments/CHECK-LOG.md`). Không gộp nhiều task vào một file. Chưa có file thì tạo từ phần header trong `.claude/skills/quality-check/references/severity-and-log.md`.
- **Mọi stage đều ghi log**, kể cả stage do người hoặc platform chạy, và kể cả lượt bị bỏ dở (ghi rõ phần chưa chạy). Mỗi lượt là một entry ở cuối mục "Entries" kèm một dòng trong bảng tóm tắt; entry có Round (vòng), tên stage, version, Result (PASS / FAIL / PARTIAL), Next (sang stage nào hoặc quay về task generator vì lý do gì). Agent chạy skill hoặc stage nào thì agent đó ghi entry; stage do người chạy thì agent ghi hộ từ kết quả người dùng đưa.
- Mỗi finding có đúng một mức độ HIGH / MEDIUM / LOW và một dòng bằng chứng (lệnh, file:dòng, probe id, câu trích hoặc nguyên văn từ platform). Không ghi finding không có bằng chứng. Với kết quả platform, quy đổi mức độ của platform sang HIGH/MEDIUM/LOW và giữ nguyên văn gốc trong phần bằng chứng.
- Không sửa entry cũ. Chỉ điền hai dòng `Reference result` và `Agreement` (và cột "Đối chiếu" trong bảng) khi có kết quả độc lập để đối chiếu: stage 3 đối chiếu stage 2, stage 6 đối chiếu stage 5 (và review ở stage 7–8 đối chiếu cả hai).
- Khi đối chiếu, gắn từng finding là: match / skill-only (báo nhầm) / reference-only (skill bỏ sót). Các mục skill-only và reference-only là đầu vào để cải thiện skill; cuối mỗi đợt, ghi lại mẫu lỗi lặp lại vào cuối file log của task.
- Skill giữ độc lập với platform: chỉ các file `CHECK-LOG.md` trong thư mục task và `CLAUDE.md` nhắc đến stage và kết quả platform.
