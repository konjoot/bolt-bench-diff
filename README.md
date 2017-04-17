```
➔ cat native.out
Benchmark_CreatePerTx-4                  	     500	   3006573 ns/op
Benchmark_PutDataByBatch-4               	  200000	    331578 ns/op
Benchmark_PutDataConcurrentlyByBatch-4   	  200000	    327767 ns/op
--- BENCH: Benchmark_PutDataConcurrentlyByBatch-4
	bench_test.go:186: count keys = 5
	bench_test.go:186: count keys = 104
	bench_test.go:186: count keys = 10004
	bench_test.go:186: count keys = 200004
Benchmark_GetData-4                      	 1000000	      1021 ns/op
Benchmark_Cursor-4                       	20000000	        54.5 ns/op
PASS
ok  	command-line-arguments	231.569s


➔ cat encrypt.out
Benchmark_CreatePerTx-4                  	     500	   3344007 ns/op
Benchmark_PutDataByBatch-4               	  100000	    160170 ns/op
Benchmark_PutDataConcurrentlyByBatch-4   	  100000	    160146 ns/op
--- BENCH: Benchmark_PutDataConcurrentlyByBatch-4
	bolt_encrypt_bench_test.go:192: count keys = 5
	bolt_encrypt_bench_test.go:192: count keys = 104
	bolt_encrypt_bench_test.go:192: count keys = 10004
	bolt_encrypt_bench_test.go:192: count keys = 100004
Benchmark_GetData-4                      	   50000	     32058 ns/op
Benchmark_Cursor-4                       	 2000000	       918 ns/op
PASS
ok  	command-line-arguments	62.841s

➔ benchcmp native.out encrypt.out
benchmark                                  old ns/op     new ns/op     delta
Benchmark_CreatePerTx-4                    3006573       3344007       +11.22%
Benchmark_PutDataByBatch-4                 331578        160170        -51.69%
Benchmark_PutDataConcurrentlyByBatch-4     327767        160146        -51.14%
Benchmark_GetData-4                        1021          32058         +3039.86%
Benchmark_Cursor-4                         54.5          918           +1584.40%

```
