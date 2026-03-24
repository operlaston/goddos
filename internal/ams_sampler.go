package main

import (
	"fmt"
	"hash/maphash"
	"io"
	"math"
	"sort"
	"strings"
)

type AMSHasher struct {
	seeds              [][]maphash.Seed
	buckets            [][]int
	numBuckets         int
	numItemsPerBucket  int
}

func NewAMSHasher(numBuckets int, numItemsPerBucket int) *AMSHasher {
	buckets := make([][]int, numBuckets)
	for i := range buckets {
		buckets[i] = make([]int, numItemsPerBucket)
		for j := range buckets[i] {
			buckets[i][j] = 0
		}

	}

	seeds := make([][]maphash.Seed, numBuckets)
	for i := range seeds {
		seeds[i] = make([]maphash.Seed, numItemsPerBucket)
		for j := range buckets[i] {
			seeds[i][j] = maphash.MakeSeed()
		}

	}

	return &AMSHasher{
		seeds:             seeds,
		buckets:           buckets,
		numBuckets:        numBuckets,
		numItemsPerBucket: numItemsPerBucket,
	}
}

func (h *AMSHasher) Hash(item []byte, bucketIndex int, seedIndex int) int {
	var mh maphash.Hash
	mh.SetSeed(h.seeds[bucketIndex][seedIndex])
	mh.Write(item)

	hashVal := mh.Sum64()

	if hashVal&1 == 1 {
		return 1
	}
	return -1
}

func (ams *AMSHasher) Ams(b []byte, bucketIndex int) {
	for i := range ams.buckets[bucketIndex] {
		z := ams.Hash(b, bucketIndex, i)
		ams.buckets[bucketIndex][i] += z
	}
}

func (ams *AMSHasher) combine() float64 {
	bucketAverages := make([]float64, ams.numBuckets)

	for i := range ams.buckets {
		sumOfSquares := 0.0
		for j := range ams.buckets[i] {
			z := float64(ams.buckets[i][j])
			sumOfSquares += z * z
		}
		bucketAverages[i] = sumOfSquares / float64(ams.numItemsPerBucket)
	}

	sort.Float64s(bucketAverages)

	mid := ams.numBuckets / 2
	if ams.numBuckets%2 == 0 {
		return (bucketAverages[mid-1] + bucketAverages[mid]) / 2.0
	}
	return bucketAverages[mid]
}

func calculateF2() float64 {
	r := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	desiredEpsilon := .1
	desiredDelta := .05

	numBuckets := int(math.Ceil(math.Log(1.0 / desiredDelta)))
	numItemsPerBucket := int(math.Ceil(1 / (desiredEpsilon * desiredEpsilon)))

	hasher := NewAMSHasher(numBuckets, numItemsPerBucket)

	b := make([]byte, 1)
	for {
		_, err := r.Read(b) // might want to handle n later
		if err == io.EOF {
			break
		}
		for bucketIndex := range hasher.buckets {
			// fmt.Printf("Processing bucket %v\n", bucketIndex)
			hasher.Ams(b, bucketIndex)
		}
	}
	return hasher.combine()
	// return 0.0
}

func main() {
	fmt.Printf("Estimated f_2 is %v\n", calculateF2())
}
