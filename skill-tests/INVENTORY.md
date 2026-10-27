# Danh sách ca thử cho quality-check và fp-check (jedi-overload-arguments)

Nguồn: `jedi-overload-arguments/` (các version), `check-errors/`, `passed-runs/`. Kết quả đúng ("ground truth") lấy từ ghi chú của tác giả khi các bản này bị fail thật.
Repo cố định: davidhalter/jedi @ 2af1a0f394. Dockerfile: `v9-review2/Dockerfile`.

## A. quality-check (prompt + test + solution)

| ID | Version | Prompt | Test | Solution | Kết quả đúng cần bắt |
|----|---------|--------|------|----------|----------------------|
| Q1 | v1-three-groups | description.md | test.patch | solution.patch | Prompt: quá nhiều nhóm yêu cầu (type[], None, Callable) |
| Q2 | v9 | task_v9.md | test_v10.patch | solution_v9.patch | Prompt 435 từ vượt ngân sách; Solution: comprehension rỗng chỉ xét `for` đầu (High), vượt giới hạn variant thì xấp xỉ (Medium), Protocol checker ngoài prompt (Medium) |
| Q3 | v10 | task_v10.md | test_v10.patch | solution_v10.patch | Đối chứng: các check qua (trừ ghi chú) |
| Q4 | v11 | task_v10.md | test_v11.patch | solution_v11.patch | Solution: helper bỏ cờ `exact`, `f(*xs)` với `tuple[X, ...]` bị coi là 1 đối số (High) |
| Q5 | v13 | task.md | test_v12.patch | solution_v13.patch | Test không fair: 3 assertion `wrap(...)` (dict đi qua kết quả helper, prompt không nêu) |
| Q6 | v14 | task.md | test_v14.patch | solution.patch | Đối chứng: đã bỏ assertion wrap, các check qua |

## B. fp-check (prompt + test + solution + run của agent)

| ID | Test | Run | Kết quả đúng cần bắt |
|----|------|-----|----------------------|
| F1 | test_v10.patch + solution_v10.patch | `v9-review2/runs-v10/*.patch` (8 run) | Vega x4 chỉ fail `test_keys_that_are_not_strings`; rd74gg0m, rd7000n9 fail; rd749, rd771 pass. rd749 khác reference ở `{k: C for k in "ab"}` (ra RC, đúng là RB) = FP cần thêm test |
| F2 | test_v14.patch + solution.patch | `runs-v10/*.patch` | Chấm lại trên test mới: rd749 69/69, rd771 fail phần unknown-length, các run còn lại fail |
| F3 | `olympus-kh7...-v0/test.patch` + solution.patch | `passed-runs/jedi-overload-arguments/*/workspace-diff.patch` (9 run) | Các run "pass" cũ được chấm trên bộ test nhỏ hơn (19/42 test): chấm lại trên test hiện tại, Nova_Orion_1 fail gần hết, Orion_Nova_1 fail ~6 hàm |

## Kết quả chạy thử (2026-10-27)

Công cụ chạy: `quality-check/scripts/{static_gate.sh,verify.sh,junit_check.py}`, `fp-check/scripts/{grade_run.sh,diff_audit.py,run_probe.sh,compare_probes.py}`. Review bằng subagent dùng skill (chỉ đọc + probe trong container).

### fp-check (công cụ, tái hiện được)
| ID | Kết quả | So với ground truth |
|----|---------|--------------------|
| F1 | rd749 67/67, rd771 67/67, rd7000n9 62/67, rd74gg0m 65/67, Vega 65-66/67 | Khớp ghi chú tác giả. Probe `{k: C for k in "ab"}`: rd749 ra RC, reference và truth là RB, đúng run lọt FP mà tác giả đã thêm test (v11) |
| F2 | rd749 69/69, rd771 68/69 (chỉ fail unknown-length), còn lại 63-67/69 | Khớp |
| F3 | 9 run "pass cũ" nay 25/62 (Nova_Orion_1) và 43-44/62 (8 run còn lại); fail chủ yếu partial, container, tuple | Khớp dự đoán: pass cũ là trên bộ test yếu hơn |

### quality-check
| ID | Lỗi đã biết | Lần 1 | Sau khi sửa skill |
|----|-------------|-------|-------------------|
| Q1 | quá nhiều nhóm yêu cầu | BỎ SÓT (check 7 PASS) | BẮT (3 feature group trong 200 từ) |
| Q2 | prompt 435 từ; comprehension rỗng chỉ xét `for` đầu; xấp xỉ khi vượt giới hạn variant; Protocol ngoài prompt | bắt prompt dài và Protocol; bỏ sót 2 lỗi solution | chưa chạy lại |
| Q3 | đối chứng (v10) | verify.sh: contract PASS, flaky base 0/3981 | |
| Q4 | cờ `exact` bị bỏ, `f(*xs)` với `tuple[X, ...]` | BỎ SÓT | vẫn BỎ SÓT; thêm báo Literal/defaultdict (chưa kiểm chứng) và báo check 7 FAIL cho prompt mà bên kia PASS |
| Q5 | assertion `wrap(...)` không fair | BẮT | chưa chạy lại |
| Q6 | đối chứng (v14) | báo động nhầm check 4, 5, 7 | còn báo check 6 (Literal) và check 7 borderline; check 4, 5 PASS |
