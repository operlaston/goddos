package main

import (
	"fmt"
	"hash/maphash"
	"io"
	"sort"
	"strings"
)

type AMSHasher struct {
	seeds [][]maphash.Seed
	buckets [][]int
	num_buckets int
	num_items_per_bucket int
}

func NewAMSHasher(num_buckets int, num_items_per_bucket int) *AMSHasher {
	buckets := make([][]int, num_buckets)
	for i := range buckets {
		buckets[i] = make([]int, num_items_per_bucket)
		for j := range buckets[i] {
			buckets[i][j] = 0
		}

	}

	seeds := make([][]maphash.Seed, num_buckets)
	for i := range seeds {
		seeds[i] = make([]maphash.Seed, num_items_per_bucket)
		for j := range buckets[i] {
			seeds[i][j] = maphash.MakeSeed()
		}

	}
	
	return &AMSHasher{
		seeds: seeds,
		buckets: buckets,
		num_buckets: num_buckets,
		num_items_per_bucket: num_items_per_bucket,
	}
}

func (h *AMSHasher) Hash(item []byte, bucket_index int, seed_index int) int {
	var mh maphash.Hash
	mh.SetSeed(h.seeds[bucket_index][seed_index])
	mh.Write(item)

	hashVal := mh.Sum64()

	if hashVal&1 == 1 {
		return 1
	}
	return -1
}

func (ams* AMSHasher) Ams(b []byte, bucket_index int) {
	for i := range ams.buckets[bucket_index] {
		z := ams.Hash(b, bucket_index, i)
		ams.buckets[bucket_index][i] += z
	}
}

func (ams* AMSHasher) combine() float64 {
	bucket_averages := make([]float64, ams.num_buckets)

	for i := range ams.buckets {
		sum_of_squares := 0.0
		for j := range ams.buckets[i] {
			z := float64(ams.buckets[i][j])
			sum_of_squares += z * z
		}
		bucket_averages[i] = sum_of_squares / float64(ams.num_items_per_bucket)
	}

	sort.Float64s(bucket_averages)
	
	mid := ams.num_buckets / 2
	if ams.num_buckets%2 == 0 {
		return (bucket_averages[mid-1] + bucket_averages[mid]) / 2.0
	}
	return bucket_averages[mid]
}

func Calculate_f2() float64 {
	r := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	hasher := NewAMSHasher(21, 30)

	b := make([]byte, 1)
	for {
		_, err := r.Read(b) // might want to handle n later
		if err == io.EOF {
			break
		}
		for bucket_index := range hasher.buckets {
			// fmt.Printf("Processing bucket %v\n", bucket_index)
			hasher.Ams(b, bucket_index)
		}
	}
	return hasher.combine()
	// return 0.0
}

func main() {
	fmt.Printf("Estimated f_2 is %v\n", Calculate_f2())
}
