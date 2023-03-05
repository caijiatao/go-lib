test:
	# -benchmem 查看每次分配内存的次数
	# -benchtime=3s 执行的时间指定为3s，一般1s、3s、5s得到的结果差不多的，如果性能较差，执行时间越长得到的性能平均值越准确
	# -count=3 执行次数，多次执行能保证准确性
	# -cpu n 指定cpu的核数，一般情况下CPU核数增加会提升性能，但也不是正相关的关系，因为核数多了之后上下文切换会带来影响，需要看是IO密集型还是CPU密集型的应用，多协程的测试中可以进行对比
	go test -bench="." -benchmem -cpuprofile=cpu_profile.out -memprofile=mem_profile.out -benchtime=3s -count=3 .
	 go tool pprof -http="127.0.0.1:8080" cpu_profile.out
	#go tool pprof -http="127.0.0.1:8080" mem_profile.out
