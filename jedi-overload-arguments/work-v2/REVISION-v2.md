# v2: sửa theo review v1 (Ali Hashir, 2026-10-22) + Auto Review

Chưa upload. Patch mới: `out-v2/solution.patch`, `out-v2/test.patch` (test.sh không đổi, description không đổi).
Repo làm việc: `work-v2/jedi` (branch `v2`, commit đầu = v1). Probe: `tools/probe_review.py <worktree>`.

## Solution

| Mục review | Sửa |
|---|---|
| S1.1 TypeVar bound/constraint dạng chuỗi | `type_var.py`: `_infer_type` resolve forward reference bằng `infer_annotation` |
| S1.2 `f(*[A, *[C]])`, `f(*[*[C]])` | `arguments.py`: khi chọn overload, `*` dùng `iterate_elements` (đếm cả `*` lồng trong literal) |
| S1.3a `dict(M())` với `keys()`/`__getitem__` | `_keys_and_getitem_items`, dùng cho cả entries (`**`) và pairs (`dict[K, V]`) |
| S1.3b `{k: C for k in []}` | `is_empty_comprehension` (gradual/base.py): rỗng -> không entry; comprehension không rõ -> phần tử unknown |
| S1.4 `dict({"p": A} if c else {}, p=C)` | `_merge_alternative_pairs`: giữ key tĩnh khi gộp alternatives, entry sau ghi đè |
| S1.5 `tuple[int, ...]` vào `tuple[int]` | `iterate_elements`: chỉ literal, FakeTuple/FakeList, `Tuple[X, Y]` mới có độ dài chính xác |
| S1.6 partial: key không rõ / cùng có thể thiếu | `stdlib.py`: chỉ key str mới thay thế; giữ PossiblyMissing khi cả hai có thể thiếu. `param.py`: gộp các entry key không rõ thay vì ghi đè |
| S1.7 Protocol: `meth = 1`, thừa tham số bắt buộc | `klass.py`: member phải callable; `_can_be_called_like` kiểm tra positional/keyword-only |
| S1.8 zuban | Sửa kiểu; zuban còn đúng 9 lỗi như base |
| Auto Review depth=10 | Bỏ depth, phát hiện vòng lặp bằng tập node đang mở rộng (`expanded`) |
| FP panel `type[A] \| None` | `base_value.py`: annotation `None` nhận None, `object` nhận mọi giá trị |

## Test (68 case, đều fail ở base)

Thêm: `*` lồng (`f(*[*[A], C])`...), `*`/`**` trong list/tuple/dict literal, mapping có `keys`/`__getitem__`,
comprehension rỗng, TypeVar forward ref (`test_type_var_with_forward_references`), key không rõ trong partial
(`test_uncertain_keywords_in_partial`), ghi đè dict có điều kiện, tuple độ dài không rõ (`test_tuples_of_unknown_length`),
`dict(M())` (`test_dict_of_mapping_with_keys`), gợi ý Fairness (`f(str)` TypeVar không bound, partial positional, dict kết hợp).
Không test: Protocol (không có trong description), depth/vòng lặp, `X | None`, partial "cả hai có thể thiếu".

## Ảnh hưởng rollout

Chấm lại 13 run hiện tại với test mới: 0/13 pass (`regrade.txt`). Run Vega pass duy nhất fail 7 test.

## Đã kiểm tra (2026-10-22)

- Test mới: 68/68 fail ở base, 68/68 pass với solution; pass khi chạy riêng từng test với PYTHONHASHSEED 0–3.
- Full suite base mode: 3953 pass. 4 test fail (`test_environment` x2, `test_module_exists_only_as_stub`, `sys_path:13`) cũng fail ở base local, do môi trường.
- zuban: 9 lỗi, giống base. flake8 sạch.
- Patch apply sạch trên commit 2af1a0f3.

## v2b (đã push lên draft v2, 2026-10-22)

Description v2 (411 từ, liệt kê mọi dạng) bị descriptionQuality FAIL "too many requirements" (blocking).
Đã thu hẹp phạm vi: bỏ `dict(...)` call forms và `functools.partial` khỏi description và test (solution vẫn giữ).
Description `desc/description_v2b.md` (290 từ) giữ gợi ý: `*`/`**` trong literal, tuple độ dài không xác định,
`*` lồng trong literal được unpack khi gọi. descriptionQuality: PASS. Test: 63 (fail hết ở base, pass hết với ref,
pass riêng lẻ seed 0–3). Bản v2 trước (68 test) lưu ở `tools/test_v2a_9c3e7a.py`.

## v2c/v2d (đã push, 2026-10-22)

solutionQuality FAIL trên v2b: (High) `{"k": C, **d}` với `d = {"k": A} if c else {}` ghi đè `C`; (Medium) `xs = (A,) if c else (A, A)`; `f(*xs)`
coi như luôn có 2 đối số. Sửa: `_update_with_alternative_pairs` giữ trạng thái "có thể thiếu" (kể cả qua mapping lồng và `dict(...)`);
`*` khi chọn overload đánh dấu các vị trí chỉ có ở alternative dài hơn; `*undefined`/`**undefined` là không khớp.
Lần chạy sau còn note Medium: Protocol không so kiểu tham số keyword-only, đã sửa luôn.
Kết quả: verifySolution PASS, solutionQuality PASS ("implements the described overload matching and unpacking semantics").
Test không đổi (63). Full suite: chỉ 4 test môi trường fail, giống base.
