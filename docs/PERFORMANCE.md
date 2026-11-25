# Performance baselines

These numbers are a reproducible development baseline, not a promise of application frame rate. Compare candidates against the same workload and environment rather than treating cached fast paths as whole-application results.

## Environment

- OS / architecture: Windows, amd64
- CPU: 12th Gen Intel Core i5-12600K
- Go: 1.26.0
- CGO: disabled
- Benchmark duration: 100 ms per benchmark, one run
- G3N and race results: not included in this environment

## Results

| Benchmark | Time/op | Bytes/op | Allocations/op |
| --- | ---: | ---: | ---: |
| `BenchmarkWrapUncached` | 7,074 ns | 416 | 11 |
| `BenchmarkWrapCached` | 3.014 ns | 0 | 0 |
| `BenchmarkGioFrameCacheStatic` | 6.076 ns | 0 | 0 |
| `BenchmarkLayoutCleanTree` | 13.66 ns | 0 | 0 |
| `BenchmarkMountClose` | 8,764 ns | 10,660 | 62 |
| `BenchmarkStartupComplexTree` | 8,222,350 ns | 4,473,086 | 16,036 |
| `BenchmarkIdleFrameComplexTree` | 20.19 ns | 0 | 0 |
| `BenchmarkImageBlur` | 6,786,478 ns | 3,146,075 | 4 |
| `BenchmarkImageGrayscale` | 1,536,438 ns | 1,048,641 | 2 |
| `BenchmarkImageGrayscaleYCbCr` | 2,480,489 ns | 1,048,640 | 2 |
| `BenchmarkImageGrayscaleGray` | 1,826,183 ns | 1,048,640 | 2 |
| `BenchmarkReconcileListUnkeyed` | 17,162 ns | 0 | 0 |
| `BenchmarkReconcileListUnkeyedWithText` | 29,039 ns | 0 | 0 |
| `BenchmarkYogaLayoutDirty` | 151,911 ns | 29 | 1 |
| `BenchmarkYogaLayoutCached` | 88.51 ns | 0 | 0 |
| `BenchmarkFullFrameComplexTree` | 38,397 ns | 25 | 0 |
| `BenchmarkVirtualList100kStartup` | 110,941 ns | 169,083 | 760 |
| `BenchmarkVirtualList100kScroll` | 69,130 ns | 49,707 | 301 |

`BenchmarkFullFrameComplexTree` drives a 500-row tree through `Harness.Step` and `Harness.PaintTo` with a discard painter. It measures the combined frame/layout/paint-command path without GPU submission, so it complements rather than replaces real-window profiling. An earlier 150,026 ns / 714,399 B / 59-allocation measurement was dominated by the test recording painter; the current result uses the discard painter to isolate the engine. This is a measurement correction, not a production performance claim.

`BenchmarkVirtualList100kStartup` and `BenchmarkVirtualList100kScroll` use 100,000 logical rows, a 200px viewport and 20px rows. A setup assertion verifies that exactly 17 visible/overscan rows mount and row 99,999 remains absent. Startup includes mount, initial layout, discard paint and close; scroll repeatedly changes the offset in both directions and includes reconciliation, layout and discard paint. Neither benchmark includes OS input delivery or GPU submission.

`BenchmarkMountClose` covers a complete Fiber/Yoga lifecycle with Store and Signal subscriptions and an Effect cleanup. `TestRepeatedMountCloseReleasesSubscriptions` additionally verifies 50 mount/close cycles: every cycle leaves zero Store and Signal subscriptions and exactly 50 effect cleanups. Allocation count is not itself a leak measurement; long-running memory profiles remain a separate release check.

`BenchmarkStartupComplexTree` builds the same immutable 500-row description once, then measures first mount, initial layout, discard-painter command generation and close. It excludes native window creation and GPU submission. `BenchmarkIdleFrameComplexTree` measures only `Harness.Step(0)` on that settled tree; it does not paint and does not include the OS event loop. Its zero allocations complement the separate static Gio macro-replay benchmark rather than replacing real-window idle measurements.

## Complex lifecycle memory trend

`TestLifecycleMemoryTrend` mounts and closes the 200-row benchmark tree 200 times after a 20-cycle warmup. In the documented local run, `HeapAlloc` was 37,810,416 bytes before and 35,685,496 bytes after the cycles and final GC, with no retained-heap growth. This is a useful regression signal, but a single GC trend is not proof that every possible leak is absent; long-running profiling and native backend checks remain release work.

## BiDi display-cache comparison

A 32-repeat mixed English/Hebrew/Arabic string showed why repeated visual-run computation needs caching. The same `BenchmarkBidiNodeDisplay` workload was run before and after adding the per-renderNode, single-entry display cache:

| Version | Time/op | Bytes/op | Allocations/op |
| --- | ---: | ---: | ---: |
| Before cache | 246,194 ns | 280,793 | 540 |
| After cache | 1.565 ns | 0 | 0 |

The cache stores only the node's most recent logical/display string pair. It does not change logical editing offsets, caret mapping or global state. Uncached `bidiDisplay` and `bidiCaretX` remain separate workloads; this comparison applies to repeated display of the same text, not to every BiDi operation.

A separate pure-LTR fast path scans bidi character classes before constructing a paragraph. `BenchmarkBidiDisplayLTR` changed from 33,623 ns / 90,072 B / 63 allocations to 2,327 ns / 64 B / 1 allocation on the same environment. Text containing RTL, Arabic-letter or right-to-left embedding/isolate classes still uses the full BiDi algorithm.

## Image filter conversion comparison

Both workloads use a 512×512 NRGBA source: blur radius 8 and grayscale amount 0.65. The initial implementation called `color.NRGBAModel.Convert` per pixel. Direct NRGBA/NRGBA64/RGBA/RGBA64 reads and allocation-free premultiplied conversion remove that per-pixel interface allocation; blur also reuses its input buffer for the second pass output. A follow-up JPEG-style YCbCr workload exposed the same generic-path allocation, so YCbCr now has a direct path too. Standard Gray already measured 2 allocations without a dedicated branch and retains the generic path.

| Version | Benchmark | Time/op | Bytes/op | Allocations/op |
| --- | --- | ---: | ---: | ---: |
| Before | `BenchmarkImageBlur` | 8,455,825 ns | 5,243,398 | 262,149 |
| After | `BenchmarkImageBlur` | 6,577,275 ns | 3,146,127 | 4 |
| Before | `BenchmarkImageGrayscale` | 4,238,742 ns | 2,097,233 | 262,146 |
| After | `BenchmarkImageGrayscale` | 1,525,565 ns | 1,048,645 | 2 |
| Before | `BenchmarkImageGrayscaleYCbCr` | 5,350,535 ns | 1,887,764 | 262,146 |
| After | `BenchmarkImageGrayscaleYCbCr` | 2,494,436 ns | 1,048,753 | 2 |

`TestImageNRGBAAtMatchesStandardConversion` compares the common fast paths and the generic path against Go's standard color conversion. These are CPU image-processing measurements, not frame-rate or GPU-upload results.

## Reproduce

```sh
go test ./pkg/ui -run '^$' \
  -bench 'Benchmark(MountClose|StartupComplexTree|IdleFrameComplexTree|ImageBlur|ImageGrayscale|ImageGrayscaleYCbCr|ImageGrayscaleGray|ReconcileListUnkeyed|ReconcileListUnkeyedWithText|FullFrameComplexTree|VirtualList100kStartup|VirtualList100kScroll|WrapCached|WrapUncached|GioFrameCacheStatic|LayoutCleanTree)$' \
  -benchtime=100ms -benchmem -count=1

go test ./yoga -run '^$' \
  -bench 'BenchmarkYogaLayout(Dirty|Cached)$' \
  -benchtime=100ms -benchmem -count=1

go test ./pkg/ui -run TestRepeatedMountCloseReleasesSubscriptions -count=1
```

For a lifecycle-specific quick run, use `-bench BenchmarkMountClose -benchtime=100x` to inspect one bounded batch. For the BiDi cache comparison, use `-bench BenchmarkBidiNodeDisplay -benchtime=100ms -benchmem`.

## Performance change policy

- Record hardware, OS, Go version, workload, duration and allocations with every comparison.
- Keep correctness and idle behavior checks green before accepting a performance change.
- Profile real component trees before replacing keyed maps, slices or event registries.
- The local child bitmap remains limited to hover state; broader bitmap/index adoption requires measured time and memory benefits.
- Add complex-application, long-running memory and native-backend benchmarks before the stable release gates.
