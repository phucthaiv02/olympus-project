# Checklist lỗi thường gặp (Olympus)

Tổng hợp từ **jsonnet `??`/`?.`** (`first_chall/`, accepted), **jedi overload** (`jedi-overload-arguments/`, draft), gopherjs, myst, otto, go-lua và xgo. Chi tiết từng vòng của jedi: `jedi-overload-arguments/v5-type-variadic/solution-quality-issues.md`. Report FP check: `first_chall/fp-check-passed.json` (pass), `jedi-overload-arguments/checks/fp-check-2026-10-20-failed.txt` (fail), `olympus-kh7d8qtdcfr0k5mt1rvftcvebn8gb2va-v0/checks/fp-check-2026-10-22.json` (pass with warnings). Auto Review và run của jedi ngày 2026-10-22: `olympus-kh7d8qtdcfr0k5mt1rvftcvebn8gb2va-v0/checks/`. Quality check của gopherjs pointer conversions: `gopherjs-pointer-conversions/SUBMISSION.md`.

## 1. Chọn ý tưởng / repo

- [ ] Phạm vi chỉ 1–2 nhóm hành vi. Rộng thì dính "too many requirements" (không khiếu nại được). Liệt kê nhiều dạng (list/Sequence/tuple/dict…) cũng bị đếm thành nhiều requirement (jedi bản đầu, v1, v4, H).
- [ ] Không quá hẹp hoặc trùng ý tưởng cũ, nếu không sẽ bị "derivative" (jedi v1 0.78, linkml uniqueness 56%). Issue cũ, nổi tiếng, có maintainer mời làm thì nhiều khả năng đã có người nộp.
- [ ] Scope Gate: tìm open PR theo **cơ chế và hàm sẽ sửa**, không chỉ theo từ khóa feature, và không cắt output. PR không link issue mà cover >50% test vẫn bị coi là "publicly-solved" (myst #607 / PR #1102).
- [ ] Scope Gate còn so với submission cũ trên platform và với code tham chiếu công khai:
  - otto ES5 strict mode bị "duplicate" vì 2 submission cũ đã làm cùng engine.
  - go-lua `debug.getinfo/getlocal/setlocal` bị "publicly-solved" dù PR công khai chỉ cover 25% test, vì reference port gần nguyên văn `ldblib.c`/`ldebug.c` của Lua 5.2.
  - Tránh task mà phần lõi port được từng dòng từ implementation tham chiếu (C Lua, source của spec).
- [ ] Plagiarism "derivative" với submission cũ cùng repo (xgo virtual fields, gopherjs evaluation order 0.81). Repo đã có submission khác thì liệt kê các ý tưởng gần đó trước khi build.
- [ ] Check repo cap / blocked trước (`.claude/skills/olympus-task-finder/scripts/repo_cap.mjs`).

## 2. Description

- [ ] Viết tiếng Việt rồi dịch qua API gtx. Description tự sinh bằng tiếng Anh bị nhận là AI (fractionAi = 1) hoặc quá prescriptive (jsonnet phải viết lại khoảng 6 lần).
  - Áp dụng cho **mọi lần sửa**, kể cả sửa một câu để thu hẹp phạm vi. Gopherjs v3–v8 viết thẳng tiếng Anh, bị descriptionQuality FAIL `ai_generated` ("too prescriptive… write it like how you would prompt an LLM", 2026-10-27). Viết tiếng Việt chưa đủ: jedi v10 viết tiếng Việt rồi dịch gtx, PASS trên v3 (2026-10-26) nhưng cùng text FAIL `ai_generated` trên v4 (2026-10-27). Nguyên nhân là bản tiếng Việt viết như spec: chuỗi quy tắc nối bằng dấu hai chấm và dấu phẩy, câu bị động kiểu "Annotation được kiểm tra theo…", liệt kê các construct. Hãy viết như người dùng báo bug: hiện tượng, ví dụ cụ thể, kết quả mong muốn. Lưu bản tiếng Việt ở `desc-vi/description.vi.md` cạnh bản dịch gtx thô, và chỉ sửa tay những chỗ gtx dịch sai nghĩa.
- [ ] Mô tả hành vi, không mô tả cách làm.
- [ ] Bỏ các câu thừa bị flag HIGH:
  - ví dụ chỉ lặp lại quy tắc (`1 ?? error "unreachable"`)
  - "các trường hợp khác giữ nguyên hành vi cũ"
  - đoạn mô tả hành vi lỗi hiện tại
  - câu "X không thuộc thay đổi này"
  - câu nhắc lại phần vẫn đúng ("struct pointers must keep working" của gopherjs). Phần này để test regression lo.
- [ ] Không viết câu rộng hơn những gì test và reference làm được. "in any call form that `dict` supports" khiến FP panel probe `dict(zip(...))`. Chỉ nêu các dạng mà test kiểm tra.
- [ ] Mọi điều test assert đều phải có trong description. jedi ghi "như type checker" nhưng test xử lý `Callable` riêng. jsonnet không nói rõ chuỗi `?.` dừng ở đâu.
- [ ] Sau khi đổi description, upload luôn patch tương ứng. Description mới đi với patch cũ thì check "Problem and tests are aligned" fail.
- [ ] Gợi ý cho từng edge case cũng bị đếm thành requirement. Jedi v2 thêm 7 mệnh đề gợi ý (411 từ) và descriptionQuality FAIL blocking "too many requirements". Muốn thêm gợi ý thì phải bỏ bớt một nhóm tính năng: jedi bỏ `dict(...)` và `functools.partial` khỏi description và test (290 từ) thì PASS. Sau mỗi lần đổi description, chạy riêng `olympus checks run <id> --check=descriptionQuality` trước khi làm việc khác.
- [ ] Mỗi câu trong description là một lời hứa mà solutionQuality soi lại trên reference. Description mới nhắc `{"k": C, **d}` và "biến chứa literal" thì reviewer AI thử ngay `d` có điều kiện. Đổi description xong phải viết probe cho từng câu mới (kể cả biến thể có điều kiện) và chạy trên reference.

- [ ] **Cổng bắt buộc trước mỗi lần `problems edit --description-file`: `wc -w task.md` ≤ 330.** Lỗi "too many requirements" đã lặp lại 3 lần dù đã có trong checklist.

  | Bản | Số từ | Kết quả |
  |---|---|---|
  | gopherjs | 273 (khoảng 18 mệnh đề) | FAIL |
  | jedi v2b | 290 | PASS |
  | gopherjs v4 | 291 (khoảng 15 mệnh đề) | PASS |
  | jedi v8 | 325 | PASS |
  | gopherjs | 405 | FAIL |
  | jedi v2 | 411 | FAIL |
  | jedi v9 | 435 | FAIL (2026-10-25) |

  Số từ chỉ là thước đo gần đúng. Thứ check thực sự đếm là số mệnh đề yêu cầu riêng biệt: đếm từng mệnh đề trước khi push và giữ ở mức ≤ 15.
- [ ] **Trả lời reviewer (người hoặc AI) không được bằng cách thêm một câu cho mỗi edge case.** Jedi v9 thêm 4 câu (key bằng nhau, entry có thể thiếu, comprehension rỗng, `*xs` nhiều độ dài) và bị FAIL. Làm thế này thay vào đó:
  - Gộp các edge case vào một quy tắc tổng quát mà mọi case đều suy ra được, ví dụ "nội dung là thứ Python tạo ra lúc chạy; mọi khả năng đều phải khớp".
  - Hoặc bỏ một nhóm tính năng khác (cả trong description lẫn test) để giữ ngân sách.
- [ ] `descriptionQuality` có giới hạn số yêu cầu và chặn cứng ("There are too many requirements in this task"). Khi làm khó task, đừng kể thêm từng trường hợp vào description. Viết lại thành ít quy tắc tổng quát. Mỗi hành vi được test vẫn phải suy ra được từ một quy tắc nào đó. Check này đếm số mệnh đề, không đếm số từ. Gopherjs: bản 405 từ FAIL, bản 273 từ (khoảng 18 mệnh đề) vẫn FAIL, bản v4 291 từ (khoảng 15 mệnh đề) PASS. Bản cuối gom còn 3 đoạn, mỗi đoạn một quy tắc, và bỏ các danh sách liệt kê (API reflect, cách kiểm tra kiểu, map key…).
- [ ] Mỗi lần solutionQuality FAIL, trước khi upload hãy tự rà các trường hợp "anh em" của lỗi được báo, rồi sửa cùng lúc. Gopherjs: reviewer báo field address qua struct view và nil pointer-to-array trong mảng. Tự rà thêm thì thấy phần tử slice đọc qua `*[]MyInt64` cũng sai (Y18), nên sửa luôn trước khi reviewer bắt ở vòng sau.
- [ ] Phạm vi theo loại kiểu phải là **danh sách đóng** ("pointer, mảng hoặc struct tạo từ các kiểu chỉ khác tên"), không viết "các kiểu tạo nên chúng". Gopherjs review lượt 5: cụm này kéo theo map, chan, func, những thứ test không kiểm tra (MEDIUM).
- [ ] Tránh cụm từ kỹ thuật có nghĩa rộng hơn phạm vi task. Gopherjs: "same memory layout" bị đọc thành phải đổi bit giữa `*uint32` và `*float32` (HIGH). Đổi thành "differs only in the names of the types it is made of" và sửa cả comment trong solution cho khớp.
- [ ] Viết lại description thì solutionQuality đọc nó như hợp đồng mới. Một ví dụ đặt sai chỗ cũng có thể bị hiểu là yêu cầu cho mọi nhánh (gopherjs: "`*[2]int64` → `*[2]MyInt64`" bị đọc thành đổi tường minh, dù trong Go phép đó không hợp lệ). Đặt ví dụ ngay trong câu về `unsafe`, và cho nhánh liên quan trong solution dùng chung một đường xử lý.

## 3. Test patch / `test.sh`

- [ ] Tên dir/file test mới có hash ngẫu nhiên (`openssl rand -hex 3`), không dùng tên dễ đoán như `using_test.go` hay `ptrconv_test.go`.
- [ ] `test.sh` không có `go install`, `go build`, `pip install`.
- [ ] Base và new phải chạy **hai tập test tách rời**. Nếu không, verifySolution FAIL `p2p_f2p_overlap`.
  - Mọi lệnh ở base dùng `./...` hoặc `go list ./...` phải loại package test mới ra, ví dụ `grep -v -e /tests/gorepo -e /tests/ptrconv_658733`.
  - Gopherjs bị FAIL vì `go test $(go list ./...)` chạy test mới bằng Go native (ở đó test pass), trong khi `./tests` không có lỗi gì.
  - Kiểm tra local: file junit của base không được chứa tên package mới (`grep -c <pkg> base.xml` = 0).
  - Warning "base_new_mode_support" của precheck là dấu hiệu sớm của lỗi này, nên đọc kỹ xem lệnh nào bị chỉ ra.
- [ ] Mỗi tên test (package + tên) chỉ được xuất hiện **một lần** trong junit của base. Nếu không, verifySolution FAIL `before_p2p_not_passing`.
  - Harness lấy kết quả sau cùng của tên bị trùng. Gopherjs chạy package `tests` cả bằng `go test` native lẫn `gopherjs test`, nên junit có 205 tên trùng. 11 test integration (gọi `gopherjs` qua os/exec) pass ở native nhưng skip dưới JS, nên bị báo là "skipped".
  - Cách sửa: chia mỗi test vào đúng một runner. Native chạy `-run '^(<11 test>)$'`, gopherjs chạy `--run` với danh sách lấy từ `func Test*` trong source, đã trừ 11 test đó. `go test -list` chạy native không thấy các file test chỉ build cho JS.
  - Kiểm tra local: parse junit của base, đếm tên trùng (phải bằng 0) và tên có cả pass lẫn skip.
- [ ] Pattern package không bao gồm thư mục con (như `./tests`) thì ghi comment giải thích, để check AI không đoán nhầm.
- [ ] `test.sh` dùng flag hợp lệ (jsonnet dùng `go test -skip` bị check báo là flag không hợp lệ).
- [ ] `test.sh` không nuốt lỗi:
  - không dùng `|| true`
  - không bỏ qua lỗi của `go test -list`
  - không dùng đường dẫn tương đối sau khi `cd`
- [ ] Mọi test mới phải **fail ở chế độ base**. `TestNullSafeSuperDoesNotSuppressOtherErrors` và `TestLintStillReportsOutsideSkippedChain` từng pass khi chưa có solution.
- [ ] Mọi hành vi mà reviewer, agent review hay quality check chỉ ra đều phải có **test**, không chỉ sửa reference. Jedi v14 bị Auto Review revision_requested (High: một constrained TypeVar dùng chung nhiều tham số; Medium: `f(**{'x': C})` vào tham số tên `x`). Lỗi thứ nhất đã được agent review vòng 1 tìm ra và chỉ sửa solution (probe8); lỗi thứ hai verifyFairness đã gợi ý trong `coverageSuggestions` từ v9 nhưng bị bỏ qua vì verdict PASS. Trước mỗi push: duyệt `coverageSuggestions` của verifyFairness và `gaps` của verifierIncompleteness (kể cả khi PASS), và mọi case BAD đã sửa trong probe; mỗi mục hoặc có assertion, hoặc ghi lý do bỏ (ví dụ run pass duy nhất fail -> chỉ thêm nếu chấp nhận mất run pass).
- [ ] Không assert điều description không nói (unfair). Các ví dụ đã bị bắt:
  - giữ đúng kiểu `RuntimeError`
  - formatter idempotent ở lần format thứ hai
  - `dict(unknown)` → combined
  - nội dung dict đi qua giá trị trả về của hàm helper (`wrap(d)` trả `{**d}`) khi description chỉ nói literal hoặc biến chứa literal. Jedi v13 verifyFairness FAIL. Mỗi assertion phải dùng đúng nguồn mà description liệt kê, không "suy rộng" sang call result.
- [ ] Coverage đủ để bắt mutant hợp lý:
  - đánh giá hai lần
  - lỗi xảy ra trước `?.`
  - slice step
  - chỉ số lớn
  - mọi call form có trong description
- [ ] Test không phụ thuộc thứ tự hash: chạy với nhiều `PYTHONHASHSEED` (0–3).
- [ ] Error test phải có phần positive để fail ở base. Go: ổn định với `-count=5` và `GOMAXPROCS=1`.
- [ ] `go vet` từ chối test có sẵn của repo thì dùng `go test -vet=off` (otto, go-lua).
- [ ] Verify giống sandbox của quality check: user `model`, `ulimit -s unlimited`. Gopherjs `gopherjs test` exit 255 không có output vì node nhận `--stack_size` khổng lồ, nên verifyTests/verifySolution fail. Chạy `docker run -u model ... bash -lc 'ulimit -s unlimited && ./test.sh base'`.
- [ ] Quy tắc hai chiều ("or the other way around") thì test cả hai chiều. verifyFairness bắt gopherjs thiếu chiều named → unnamed, thiếu `**T` non-nil, và thiếu case regression cho struct pointer.
- [ ] Test làm crash cả test binary ở base (gopherjs `(*Struct)(unsafe nil)` ném TypeError) thì đặt cuối, để các test khác vẫn chạy và báo fail.
- [ ] Verify trên clean clone với `--network none`. Test cần mạng thì skip ở base (xgo `TestTestxgo/domaintext-md`).
- [ ] Không thêm hidden test chỉ để chứng minh reference sửa được bug do solutionQuality tìm ra. Giữ fix trong solution, còn test chỉ kiểm tra điều description yêu cầu trực tiếp. Nếu không, mọi agent đều fail (jedi D: 0/23).
- [ ] Không sửa test có sẵn của repo.
- [ ] Chuẩn hoá `--output_path` thành đường dẫn tuyệt đối trước khi `cd`. Kiểm tra bằng cách gọi với đường dẫn tương đối từ thư mục khác (gopherjs: report nằm trong thư mục tạm, exit 0).
- [ ] Mở `<failure>` trong JUnit của một lần fail thật: Go example in got/want không thụt lề nên go-junit-report để `<failure>` rỗng, phải tự chèn từ raw log. Fallback không được chỉ giữ `tail`, vì lỗi gốc có thể nằm đầu output.

## 4. Solution patch

- [ ] Hành vi đã hứa phải đúng ở **mọi dạng tương đương**, không chỉ dạng được test.
  - jsonnet: `Slice.Step` khi duyệt cây, chỉ số âm dạng phân số, float bị tràn khi đổi sang int.
  - jedi: `dict(p=C)`, `{...}`, pairs, lồng nhau, alternatives, `**` trong literal, `[*[C]]`, partial, classmethod.
- [ ] Unpack lồng nhau phải mở ở mọi cấp. Auto Review đánh High S1 cho reference jedi vì `f(*[*[A], C])` mất `A` trước khi so overload: `SequenceLiteralValue.py__iter__` (`get_tree_entries`) lọc bỏ `star_expr`. Trước khi dùng helper sẵn có của repo, kiểm tra nó có bỏ bớt phần tử không.
- [ ] Không dùng giới hạn đệ quy cố định để biến giá trị vẫn suy ra được thành unknown/rỗng (jedi: `dict(...)` lồng sâu hơn 10 cấp, Medium). Chặn vòng lặp bằng tập đã thăm.
- [ ] Mọi giới hạn kích thước (số biến thể, độ sâu) phải **fail closed**: vượt giới hạn thì trả "không khớp", không được rơi về một xấp xỉ có thể cho ra khớp sai. Jedi v9: vượt `_MAX_ARGUMENT_VARIANTS` thì fallback về `[self]`, làm mất việc dịch vị trí đối số, và bị Medium.
- [ ] Mỗi câu mới trong description phải được probe ở **mọi biến thể cú pháp**. Jedi v9: câu "comprehension trên literal rỗng" nhưng chỉ xét `for` đầu tiên, nên `{k: C for k in [1] for _ in []}` bị High.
- [ ] Helper trả thêm thông tin (cờ `exact`, "có thể thiếu"...) thì **mọi caller** phải dùng thông tin đó. Jedi v11 FAIL High vì `_star_arg_sequences` bỏ cờ `exact`: `f(*xs)` với `xs: tuple[type[C], ...]` bị coi là đúng 1 đối số. Grep mọi chỗ gọi helper và kiểm tra từng chỗ.
- [ ] Với mỗi quy tắc unpack/container, probe đủ các **loại nguồn**: literal, biến chứa literal, giá trị annotation độ dài cố định (`tuple[X, Y]`), giá trị annotation không biết độ dài (`tuple[X, ...]`, `list[X]`), và giá trị không suy ra được. Mỗi loại nguồn đều phải đối chiếu với overload có số tham số cố định, overload có default và overload `*args`.
- [ ] Sau mỗi lần solutionQuality FAIL, giao cho một agent độc lập đóng vai reviewer: đọc description như hợp đồng, probe reference bằng biến thể của từng câu, chỉ báo lỗi đã xác nhận. Sửa hết rồi cho chạy thêm một vòng nữa trên patch mới, sau đó mới upload. Jedi v12: vòng 1 tìm ra 8 lỗi nằm trong phạm vi description mà probe của tôi đã bỏ sót.
- [ ] Sau mỗi lần sửa, chạy **toàn bộ** probe (của mình lẫn của reviewer) trên cả patch cũ và patch mới, rồi diff kết quả. Đừng chỉ xem các case vừa sửa. Jedi v13: diff kiểu này bắt được `RecursionError` mới phát sinh mà reviewer không báo. BAD nào còn lại thì tự xét lại kỳ vọng, vì reviewer cũng có thể đoán sai, ví dụ case "multiple values" mà cả hai overload thực ra đều lỗi.
- [ ] Không xấp xỉ lỏng, không để gap đã biết (Protocol "luôn accept", Union kiểm tra any-of thay vì all-of).
- [ ] Mọi chỗ "không suy ra được" đều phải dẫn tới "không khớp".
- [ ] Quy tắc chung trong description (ví dụ "đổi sang **mọi** kiểu pointer") thì rà **mọi nhánh** của hàm đã sửa, kể cả nhánh đặc biệt theo package hay platform mà test không chạm tới được. Gopherjs có hai nhánh riêng cho `syscall` (struct pointer ↔ `unsafe.Pointer`) không xử lý nil, và Solution Quality đánh HIGH.
- [ ] Đường đi chỉ biết kiểu lúc runtime (như `unsafe.Pointer` → `*T`) vẫn phải cho ra đúng kiểu đích. Không được trả nguyên object nguồn (gopherjs: cần proxy mang kiểu đích, có cache).
- [ ] Sau khi Solution Quality FAIL, **đừng chỉ sửa đúng điểm bị báo**. Lần chạy sau sẽ tìm ra điểm khác (gopherjs FAIL 3 lần liên tiếp: syscall nil → unsafe khác kiểu → syscall identity). Trước khi upload lại:
  - grep mọi nhánh đặc biệt trong hàm đã sửa: theo tên package (`== "syscall"`), theo cú pháp (`&x[i]`, `new(T)`, tên biến `_zero`), và các đường runtime song song (reflect `cvtDirect`, generics).
  - Viết probe cho **từng câu** của description qua từng đường đó, chạy trên reference và trên Go native, rồi so kết quả. Gopherjs: 56 probe trong `gopherjs-pointer-conversions/probes/`.
  - Nhánh mà code thường không chạy tới (chỉ áp dụng cho `syscall`) thì tạm sửa điều kiện package ngay trong container để chạy được trong runtime thật, rồi đọc JS sinh ra.
  - Sửa nhánh đặc biệt mà code stdlib có thể đang dựa vào thì chạy thêm các package stdlib ngoài base, và so kết quả khi có và không có solution.
- [ ] Lời hứa về identity hoặc kiểu đích ("bằng nhau sau khi đổi về cùng kiểu", "có kiểu đích") phải đúng ở **mọi chỗ tạo hoặc cache object**, không chỉ ở helper chuyển kiểu mới viết. Gopherjs FAIL 6 vòng liên tiếp, mỗi vòng reviewer tìm ra một cache khác:
  - field qua struct view (`$ptr_` trên struct);
  - phần tử mảng qua `$indexPtr`: cache theo vị trí, không theo kiểu, nên `&(*[2]MyInt32)(unsafe.Pointer(&a))[0]` trả lại `*int32` đã cache (HIGH, 2026-10-27);
  - reflect `Field(i).Addr()`, `Index(i).Addr()` và `UnsafeAddr`.

  Liệt kê **mọi** chỗ trong repo tạo pointer (grep `new .*ptr`, `$ptr`, `$indexPtr`, `JsPtrTo().New`) rồi probe từng chỗ.
  - Lời hứa "bằng nhau" phải probe ở cả hai dạng: round trip (`(*T)(unsafe.Pointer(p)) == p`) **và** đổi cùng một biến hai lần (`unsafe.Pointer(&x) == unsafe.Pointer(&x)`). Gopherjs vòng 10: nhánh `syscall` đã có round trip nhưng mỗi lần đổi vẫn tạo một `Uint8Array` mới (HIGH). Grep `new Uint8Array`, `.slice(`, `.subarray(` trong mọi đường chuyển kiểu.
- [ ] Nhánh mà runtime không bao giờ chạy tới (reflect `cvtDirect` nhánh `UnsafePointer` với nguồn là pointer, vì Go đã panic trước đó) vẫn bị reviewer đọc code rồi đánh HIGH, 2 vòng liên tiếp. Comment giải thích "không chạy tới được" không cứu được. Mọi nhánh song song phải gọi helper **với đúng các tham số** mà nhánh chính dùng (gopherjs: `$pointerToUnsafe(obj, type)` ở compiler, `Pointer`, `UnsafeAddr`, `cvtDirect`). Grep tên helper và so tham số ở từng chỗ gọi.
- [ ] Probe cả kiểu **không mang wrapper runtime**: alias của `int32`, `byte`, `float64`, `string`, `bool` (typed array và plain Array), không chỉ `int64` hay `complex` (có wrapper). Những kiểu không có wrapper đi qua đường "không cần view", và chính đường này bỏ sót việc đổi kiểu.
- [ ] Scratchpad bị xoá khi sang phiên hoặc sang ngày mới. Mọi patch, probe và script chạy phải lưu trong thư mục project. Worktree dựng lại được bằng `docker cp` từ `/app` của image cộng với `solution.patch` của project.
- [ ] Gộp các khả năng thay vì ghi đè (`dict.update`). Sentinel không được trùng giá trị thật (`'?'`).
- [ ] Khi gộp alternatives (biến có thể mang nhiều giá trị), giữ cả **sự có mặt**, không chỉ gộp giá trị. Jedi v2 bị solutionQuality FAIL vì quên điều này:
  - dict: key chỉ có ở một số alternative thì "có thể thiếu", không được ghi đè entry trước cùng key. `d = {"k": A} if c else {}`; `f({"k": C, **d})` phải xét cả `C` lẫn `A` (High).
  - sequence: vị trí chỉ có ở alternative dài hơn thì "có thể thiếu". `xs = (A,) if c else (A, A)`; `f(*xs)` không được coi là luôn có 2 đối số (Medium).
  - Rà mọi đường có alternatives: `**` tại lời gọi, `**` trong literal, `*` tại lời gọi, `dict(...)`, partial. Sửa một đường thì kiểm tra các đường còn lại.
- [ ] Unpack một tên không suy ra được (`f(*undefined)`, `f(**undefined)`) là "không khớp", không phải "không có đối số" (jedi từng chọn overload đầu tiên).
- [ ] Tính năng phụ mà solution thêm vào để không vỡ test cũ (jedi: so khớp Protocol) cũng bị reviewer và solutionQuality soi kỹ như phần chính. Mỗi vòng review jedi đều bắt một lỗ hổng ở đây: member không callable (`meth = 1`), thừa tham số bắt buộc, kiểu của tham số keyword-only. Hoặc làm đủ (callable, số tham số positional/keyword-only, kiểu từng tham số, kiểu trả về), hoặc giữ đúng hành vi của base.
- [ ] Tuân theo semantics của ngôn ngữ (Python binding, keyword override của `partial`), không dựa vào shortcut sẵn có của repo.
- [ ] Không sửa toàn cục cho một tính năng cục bộ. jsonnet giữ `Parens` ở mọi chỗ nên bị stack overflow sớm hơn. jedi mở rộng `**dict` ở mọi chỗ nên làm vỡ test cũ. Giới hạn thay đổi bằng cờ (`strict`/expand) chỉ trên đường mới.
- [ ] Thêm node hoặc field mới thì rà mọi chỗ duyệt cây (children, static analyzer, clone, linter, find dependencies).
- [ ] Build patch từ **mọi** thư mục đã sửa (`git diff -- jedi`, không phải `-- jedi/inference`).
- [ ] Lint sạch: flake8 E303, dòng không quá 100 ký tự. Không để comment chỉ kể lại code.
- [ ] Test **chính file patch** trên worktree sạch, ở cả `base` và `new`, rồi chạy full suite trước khi upload.

## 5. Dockerfile

- [ ] Python: pin chính xác (`==` hoặc từ lockfile), rồi `pip install --no-deps -e .`. Thiếu editable install thì check FAIL hẳn. Dùng `>=` thì bị WARN.
- [ ] Go: commit `go.mod`/`go.sum`, và pin version tool cài thêm.
- [ ] Không chạy `go test` trong Dockerfile (hard fail ở xgo go-expression). Chỉ build hoặc warm cache.
- [ ] `go build ./...` fail ở thư mục demo thì build từng package cần dùng (xgo: cl/cmd/parser/printer).
- [ ] Cài thêm toolchain thứ hai (GopherJS cần Go 1.21) sẽ bị WARN. Ghi lý do bằng comment ngay trong Dockerfile.
- [ ] Check Dockerfile (AI) có thể đoán sai là thiếu tool. Gopherjs bị báo base Go slim không có Node dù có `/usr/bin/node`. Ghi comment và assert ngay lúc build: `RUN test -x /usr/bin/node && /usr/bin/node --version`.
- [ ] Tool phụ thuộc stack limit (node của gopherjs) thì đặt wrapper trong Dockerfile, ví dụ `/usr/local/bin/node` dùng stack 8 MiB khi `ulimit -s` là unlimited. Wrapper này cũng áp dụng cho rollout agent.
- [ ] Môi trường phải tự đủ, kể cả submodule (jedi: typeshed, django-stubs). FP panel dựng worktree mới, nên thiếu submodule thì 48/53 test fail vì lý do không liên quan.

## 6. Rollout / độ khó

- [ ] Tiêu chí hiện tại (`runs-2026-10-22.json`, tiêu chí cũ là Nova ≥8 run và 256K):
  - Mỗi model Nova, Orion, Vega có ít nhất 4 run hoàn thành.
  - **Đúng một** model giải được, pass >0% và ≤25%; hai model còn lại 0%. Jedi đạt với Nova 0/4, Orion 0/4, Vega 1/4.
  - Run Nova trung bình trên **400K** token sequence length.
- [ ] Task quá dễ (gopherjs: Vega pass 3/3) thì đừng đoán. Thử nghiệm trước:
  - tải `workspaceDiff` của từng run pass, chỉ áp phần sửa source lên repo sạch;
  - chạy bộ probe tách riêng từng kiểm tra (mỗi kiểm tra có `recover`, để một crash không che các kiểm tra sau);
  - so từng run với reference và với Go native.
  - Kiểm tra mà reference đúng còn run sai, và có trong description (hoặc thêm được một câu), sẽ thành test mới. Chọn sao cho mỗi run pass hiện tại trượt ít nhất một test, nhưng mỗi run chỉ trượt 1–2 nhóm (gopherjs: nhóm `unsafe` khác kiểu và reflect với Vega #1 và #2, `unsafe.Pointer(new(T))` với Vega #3).
  - Lỗ hổng mà `verifierIncompleteness` báo (INCOMPLETE) là nguồn ứng viên tốt, vì nó đã chứng minh bỏ phần sửa thì test cũ vẫn pass.
  - Test mới có thể panic ở base thì bọc bằng `defer` + `recover` → `t.Errorf`, để không làm dừng cả test binary.
- [ ] Khi nhiều model cùng pass (gopherjs: Nova, Orion và Vega đều pass), probe các trường hợp biên thông thường ít khi phân biệt được, vì model mạnh đã làm đúng hết. Hãy tìm trục làm **sâu hơn** của cùng quy tắc, rồi xếp hạng các run theo trục đó:
  - đệ quy vào mảng, field struct, struct lồng nhau;
  - API phụ ít người nhớ (`reflect.NewAt`, `Pointer()`);
  - bất biến toàn cục (pointer khác kiểu cho ra cùng `unsafe.Pointer`).

  Chọn test sao cho run mạnh nhất chỉ trượt 1 trường hợp sâu nhất (Nova #5 chỉ trượt mảng lồng trong field). Reference thường phải nâng cấp theo. Chạy lại regression baseline và stdlib sau mỗi lần nâng cấp.
- [ ] Nova pass là xấu (jedi v2 bị Nova giải ngay lượt đầu). Quá khó cũng trượt: bản D có Orion 0/14, bản E có Vega 0/4.
- [ ] Median trên các run pass: ≥150 LOC và ≥2 file.
- [ ] Cheat rate dưới 20%, không có run bị chặn bởi môi trường.
- [ ] Run bị platform ngắt vẫn làm fail **cả** "Fair task" lẫn "No environment blockers", dù task không có lỗi. Jedi `rd76dwr3…`: `fault: problem`, `FAIL_EXECUTION_ERROR_EXTERNAL`, `environmentBlocker.type: mixed`, do idle cancellation 90 giây, 17 lần resume và 2 lần restart session. Sau mỗi đợt rollout, lọc `fault` khác `agent`/`pass` trong runs JSON và đọc `summary`. Lỗi do infra thì scratch run đó hoặc chạy run thay thế. Lỗi thật do môi trường thì sửa Dockerfile/test.sh trước.
- [ ] `FAIL_EARLY_TERMINATION` (`fault: infra`, agent dừng giữa chừng) không bị tính là environment blocker, nên không cần xử lý.
- [ ] Phân tích run fail bằng cách chấm lại patch của agent ở local (soft-assert), để biết rào cản nằm ở đâu trước khi chỉnh.
- [ ] Chấm local giống hệt platform. Harness `.<harness>-eval/` (`runner.sh`, `f2p.json`, `p2p.json`) nằm trong workspaceDiff của mỗi run:
  - apply diff của agent, trừ thư mục harness, `test.sh` và file hidden test;
  - apply test patch, copy harness vào (thêm test mới vào f2p);
  - chạy `runner.sh --junit-output-path`, và dedupe `classname == name` như wrapper.
- [ ] Môi trường local khác image (typeshed chưa pin, parso khác version) cho ra failure giả. Chấm trong image, hoặc so với các run khác.
- [ ] Sau khi đổi test, dùng re-eval (`reEvalRuns.triggerReEvalRuns`, cohort cố định, khoảng 1.8 token/run) thay vì rollout mới.

## 7. Auto Review

Auto Review đọc description, reference, test patch và patch của run pass, rồi trả `outcome` (`revision_requested` là warn ở readiness). Report nằm ở `dynamicChecks.autoReview.output` (`issues`, `authorGuidance`).

- [ ] Test có `*`/`**` trong literal phải **trộn kiểu**, để implementation bỏ qua phần unpack ra kết quả sai. Jedi chỉ có `{"k": C, **{"m": C}}`, nên run pass dùng `get_tree_entries()` (bỏ entry sau `**`) vẫn đạt 62/62. Cần thêm `f([C, *[A]])` và `f({"k": C, **{"m": A}})` với kết quả fallback (High T3/T4).
- [ ] Đọc patch của run pass và tìm helper của repo mà nó dùng, rồi kiểm tra helper đó có bỏ bớt dạng cú pháp nào description có nhắc. Mỗi chỗ bỏ sót là một test phân biệt cần thêm.
- [ ] Lỗi Auto Review tìm ra trong reference (S1) thì sửa trong solution. Chỉ thêm test khi description đã hứa hành vi đó (xem mục 3).

- [ ] Trước khi submit, chạy ma trận probe (mỗi case lấy từ một câu trong description) trên reference **và** trên patch của run pass, rồi so cả hai với hành vi thật của Python. Mỗi bên sai ở đâu thì rủi ro ở đó:
  - Reference sai: human reviewer báo S + T4 (jedi bị review 2 lần đúng kiểu này).
  - Run pass sai trong khi reference đúng: Auto Review đánh High T3/T4, check FP có thể đánh false positive.
  - Cả hai cùng sai: chọn một trong hai rủi ro trên. Ghi lại lựa chọn vào NOTES.
  Jedi v11: reference sai 3 case comprehension mà Orion làm đúng. Đã sửa reference trước khi re-eval.

## 8. FP check ("No false positives")

Cách panel làm việc (rút ra từ report jsonnet pass và jedi fail ngày 2026-10-20): có 2 judge và 1 adjudicator. Họ viết probe từ **từng câu trong description**, chạy trên patch của agent và trên reference, lặp lại với nhiều `PYTHONHASHSEED`, và chạy từng test riêng lẻ. Chỉ cần một run pass bị đánh `false_positive` là trượt cả tiêu chí. Không tính: probe mà reference cũng fail, probe không phân biệt được hai bên, và khác biệt ở chỗ description không nói (`UNDERSPECIFIED_EDGE`).

**Khi viết test**

- [ ] Mỗi quy tắc trong description cần một test **phân biệt được** cách làm sai phổ biến. Case mà implementation sai vẫn tình cờ ra đúng thì không tính.
  - jedi Vega #2: `test_dict_with_alternative_entries` chỉ có alternatives cùng key mà cả hai giá trị đều được `type[A]` nhận, hoặc khác key. Candidate ghi đè key (`non_matching_keys[key] = argument`) nên vẫn pass.
  - Probe của panel: `m = {"p": C} if random.random() else {"p": A}; f(**m)` phải ra RC, candidate ra RB.
- [ ] Kết quả sai của test phân biệt phải **ổn định** ở mọi seed. Chọn overload sao cho ghi đè bất kỳ alternative nào cũng ra kết quả sai (`type[A]`→RA, `type[C]`→RB, `type`→RC).
- [ ] Phủ mọi dạng cú pháp của cùng một quy tắc:
  - `f(**m)` lẫn `f(**dict(m))`;
  - tham số có tên lẫn `**kwargs`;
  - đảo thứ tự alternatives;
  - gọi qua `functools.partial`.
- [ ] **Không bỏ test để tăng pass rate khi description vẫn giữ quy tắc đó.** Version J bỏ case alt C/A để Vega có run pass, và FP check bắt đúng lỗ hổng đó. Nếu bỏ test thì thu hẹp description cùng lúc, hoặc giữ test và chấp nhận rate thấp.
- [ ] Chạy test mới **riêng lẻ từng test** và với nhiều seed, không chỉ chạy cả suite.
  - jedi Vega #8 nhận ra `tuple[X, ...]` qua tên class `'ellipsis'`. Typeshed có `ellipsis = EllipsisType`, nên case này fail khi chạy riêng ở seed 0–5, nhưng pass trong cả suite nhờ thứ tự test và cache.
  - Cân nhắc cho `test.sh` chạy test mới trong process riêng.
- [ ] Thêm biến thể cho những case có thể khớp nhờ tên hoặc cache: `Tuple[Type[A], ...]` với `(C,)`, `tuple[int, ...]` so với `tuple[str, ...]`.

**Trước khi trigger FP check** (8 token, cần ít nhất 12 working runs)

- [ ] Chỉ trigger khi run set đã chốt. Mọi edit đều làm FP check bị stale.
- [ ] Tự chạy FP check ở local cho từng run pass:
  - tải `olympus runs artifact <id> <run> --key workspaceDiff` và chấm bằng harness;
  - viết probe cho từng câu trong description và so kết quả với reference, ở nhiều seed và từng test riêng lẻ;
  - run nào khác reference ở một quy tắc có trong description thì bổ sung test trước.
- [ ] Lấy pattern sai từ các run **fail** để viết test phân biệt. Lỗi của run fail (ghi đè alternatives, bỏ qua `**`/`*` trong literal) thường cũng có ở run pass, chỉ là test chưa chạm tới.
- [ ] Đọc diff của run pass để xem agent có sửa test có sẵn không. `protected` rỗng và không có restore, nên expectation bị sửa (jedi `arrays:294`, `test_param_resolving_to_static`) vẫn giữ nguyên lúc chấm và che regression.
- [ ] Source nằm dưới thư mục tên `build`, `dist`, `target`, `node_modules`, `__pycache__` hoặc `.venv` thì khai báo ở Protected source folders. Nếu không, panel bỏ qua thư mục đó và thấy solution thiếu.

**Khi FP fail**

- [ ] Lưu ngay JSON report (`fpReview.getFpCheckForVersion {versionId}`) vào thư mục challenge. Sau khi edit, report thành stale và mất `resultSummary`. `olympus runs show` có thể trả Server Error.
- [ ] Chỉ xử lý probe mà adjudicator đánh `fairProbe` và `discriminates` là true. Judge-c bị loại hai lần: một probe dùng `if int():` mà jedi đánh giá tĩnh, một probe `dict(zip(...))` mà reference cũng fail.
- [ ] Sửa test (và description nếu cần). Test mới không được assert điều description không nói.
- [ ] Chấm lại mọi run cũ ở local để ước lượng rate mới còn >0% và ≤25%, rồi re-eval và chạy lại FP check.
- [ ] `passed_with_warnings` (`judge_dissent`) cho ra warn ở readiness, không chặn submit. Jedi 2026-10-22: judge-c lại replay với submodule typeshed rỗng (`ENV_SUBMODULE_MISSING_IN_WORKTREE`) nên báo 15 fail giả, và adjudicator loại judge này. Probe của judge-c `dict({"other": 0}, item=1)` cũng fail ở reference và sai kỳ vọng (keyword lạ thì không overload nào khớp, nên ra union).
- [ ] Run pass mà panic trên một chuyển đổi thuộc phạm vi là HIGH, kể cả kiểu hiếm (gopherjs: `type P *js.Object`). Reviewer đọc chính các probe này, không gọi là vùng xám.

## 9. Quy trình

- [ ] Kiểm tra `olympus runs presets` trước khi giả định agent nào có sẵn.
- [ ] Trước khi sửa, tải bản mới nhất từ platform (`olympus problems download`). Thư mục local `olympus-kh7d8...-v0` từng cũ hơn platform (thiếu test tuple). Sau khi upload, tải lại để xác nhận.
- [ ] `olympus runs show` (`jobs.get`) và `runPrecheckStage` có lúc trả Server Error. Khi đó đọc qua Convex query (xem memory `olympus-api-convex`).
- [ ] Mutation như `triggerReEvalRuns` bị auto mode chặn, nên để user tự trigger.
- [ ] Tính trước token: FP check cần 12 working runs, ví dụ 7 Nova + 2 Vega + FP ≈ 52 token.
- [ ] Không trigger solutionQuality khi verifySolution/taskQuality còn đang chạy (sẽ bị `Server Error`).
- [ ] Mỗi lần upload làm mọi check bị stale, nên kiểm tra kỹ ở local trước. Verdict là PASS thì không sửa thêm các note không chặn, **trừ khi** sắp phải upload lại phần đó vì lý do khác.
  - Solution Quality chấm lại từ đầu mỗi lần chạy. Note MEDIUM bị bỏ qua lần trước có thể cộng với một lỗi mới thành FAIL (gopherjs: MEDIUM về `unsafe.Pointer` khác kiểu, cộng HIGH về `syscall` → "Not Met").
  - Phải upload lại solution thì sửa luôn mọi note MEDIUM trở lên.
- [ ] Tự quyết theo lựa chọn khuyến nghị, không hỏi lại liên tục.
- [ ] Scratchpad bị xóa khi session khởi động lại. Repo làm việc, patch của run và script probe để trong thư mục project (jedi: `jedi-overload-arguments/work-v2/`).
- [ ] Chạy từng test riêng lẻ bằng script Python (`work-v2/tools/isolated.py`), không dùng `xargs`: node id có tham số chứa dấu cách bị cắt và báo FAIL giả.
- [ ] Thu hẹp phạm vi hoặc đổi description thì chạy lại solutionQuality, không chỉ descriptionQuality. Các check AI soi lại reference theo câu chữ mới.
- [ ] Đọc kết quả `verify.sh` theo dòng `FAIL`, không theo exit code của lệnh bọc (`| tail`, `; tail`). Jedi 2026-10-31: `verify.sh` báo 1 base test fail nhưng tác vụ nền vẫn exit 0, suýt push solution có regression.
- [ ] **Đọc lại file này trước mỗi lần push lên draft.** Lỗi jedi v9 (2026-10-25) đã nằm sẵn ở mục 2, nhưng vẫn bị lặp lại vì không đọc. Lặp lại lần nữa ngày 2026-10-29: jedi push 2 bản description (v11, v12) mà không đọc lại mục 2, nên viết dạng spec/câu mệnh lệnh thay vì dạng báo bug, và cả 2 đều FAIL `ai_generated`.
- [ ] Bản đã submit hoặc đã bị review thì bị khóa: `problems edit` trả `Server Error`. Chạy `olympus problems version create <id>` trước, rồi mới edit.
- [ ] Không lọc output của probe bằng `grep` khi chưa biết hết các dòng. Jedi v9 có một filter giấu mất case `{1: C, True: A}` sai, nên lỗi bool-for-int suýt bị bỏ sót. In toàn bộ output, hoặc in riêng các dòng không khớp.
- [ ] Kiểm tra mọi expectation có thể đạt được với số overload đã khai báo. Ví dụ hai overload thì không thể ra `RC`. Chạy test trên reference trước khi tin vào kỳ vọng.
- [ ] Task có ngữ nghĩa rộng (gộp alternatives, dict/sequence động) kéo reviewer người vào vô số edge case: jedi bị review 2 lần, lần sau có 8 lỗi mới. Khi chọn task mới, ưu tiên quy tắc hẹp và đóng.
