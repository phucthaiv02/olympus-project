# Language notes (TypeScript, JavaScript, Python, Go, Rust, C++, Java)

Apply the generic rules in `checks.md`; this file only lists what differs per language. Identify the language from the Dockerfile and the repo, and read only its section.

| Language | Flakiness knobs | Wildcards to exclude new tests from base | Environment pitfalls | Probe-vs-reality reference |
|---|---|---|---|---|
| Python | `PYTHONHASHSEED=0..3`, `-p no:randomly`, run each test alone | `pytest .`, `pytest test/` dirs, `-k`/`--ignore` | pin exact versions or lockfile, `pip install --no-deps -e .`; unpinned submodules/stubs differ from the image | run the same snippet in CPython |
| Go | `-count=5`, `GOMAXPROCS=1`, `-race` optional, `-shuffle=on` | `./...`, `go list ./...` (filter with `grep -v`) | commit `go.mod`/`go.sum`, no `go test` in the image, `-vet=off` if vet rejects existing tests, dual runners (native vs transpiled) duplicate names | run the same snippet with plain `go run` |
| TypeScript / JavaScript | run with `--runInBand`/single worker, shuffle if supported, repeat with different `TZ` | `jest`/`vitest` globs, `--testPathIgnorePatterns`, workspaces | lockfile + `npm ci`, node version pin, no network at test time, stack size/`--max-old-space-size` | run the same snippet in node |
| Rust | `cargo test -- --test-threads=1`, repeat runs | `cargo test --workspace` with `--exclude`, `--test` targets | `Cargo.lock` committed, `--locked --offline`, prebuild dependencies in the image | `cargo run` of the snippet |
| C++ | repeat with `--gtest_shuffle --gtest_repeat=N` or ctest `--repeat until-fail:N` | ctest regex (`-E`), globbed test dirs | pin compiler/CMake versions, build in the image not in `test.sh`, sanitizer flags change results | compile the snippet with the same compiler |
| Java | surefire `-DforkCount=1`, `-Dsurefire.runOrder=alphabetical` vs `random` | `-Dtest=` patterns, module lists in multi-module builds | pin JDK and build-tool versions, `--offline` after warm-up in the image | `jshell`/small main class |

Junit-based verification works the same in every language: each mode needs one junit file, each `classname + name` pair appearing exactly once. If the runner can't emit junit natively, convert its output (a failing converter must not hide the exit code), and check the duplicate-name and overlap rules on the converted file.

When the runner dedupes or renames tests differently (parametrized tests, subtests, table-driven cases), compare the ids in the junit file rather than the names in the source.

If the language isn't listed, derive the same four things from first principles: how to randomize and repeat runs, how a wildcard picks up new tests, what the image must pin, and how to run a snippet in the real language to get ground truth.
