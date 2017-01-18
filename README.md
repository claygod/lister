# Lister

Library for work with lists of structures in a functional style.
Library works effectively only with small lists of up to one or two thousand.

# Usage

An example of a library can be seen in the tests

# Perfomance

- BenchmarkAdd-4            	2000000000	         0.05 ns/op
- BenchmarkMainParallel-4   	   10000	    103710 ns/op
- BenchmarkDel-4            	2000000000	         0.07 ns/op
- BenchmarkAddaDel-4        	 2000000	       573 ns/op
- BenchmarkSelect-4         	2000000000	         0.17 ns/op
