# concurrency

```go
savarin@esa-e-sf-a4aba0 concurrency % go test . -bench=.

goos: darwin
goarch: amd64
pkg: advanced-programming/concurrency
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkBasic-12       862995500            1.290 ns/op
BenchmarkAtomic-12      197048946            6.003 ns/op
BenchmarkMutex-12       97199598            12.13 ns/op
BenchmarkChannel-12      6192446           190.5 ns/op
PASS
ok      advanced-programming/concurrency    6.175s

```
