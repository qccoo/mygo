package lib

import (
    "fmt"
    "sync"
    "time"
)

// Simple rolling number implementation in go, referred to Hystrix, 
// with sum and max.
type RollingNumber struct {
    buckets           map[int64]*bucket
    numBuckets        int
    bucketSizeMillis  int64
    totalSizeMillis   int64
    lastBucketTimeKey int64 // Milliseconds since epoch
    mu                *sync.RWMutex
}

type bucket struct {
    adder     addBucket
    maxGetter maxBucket
}

type addBucket struct {
    value float64
}

type maxBucket struct {
    value float64
}

// New Rolling Number.
// totalSizeMillis: total window size for all the buckets in milliseconds, e.g. 10000 => 10s.
// numBuckets: Number of buckets. e.g. 10.
// NewRollingNumber(10000, 10) => a rolling number for 10s window, with each bucket 1s.
func NewRollingNumber(totalSizeMillis int64, numBuckets int) (*RollingNumber, error) {
    if totalSizeMillis <= 0 {
        return nil, fmt.Errorf(
            "Invalid totalSizeMillis %v, should be non-negative.", totalSizeMillis)
    }
    if numBuckets <= 0 {
        return nil, fmt.Errorf(
            "Invalid numBuckets %v, should be non-negative.", numBuckets)
    }
    if totalSizeMillis % int64(numBuckets) != 0 {
        return nil, fmt.Errorf(
            "Invalid totalSizeMillis %v and numBuckets %v, mod(totalSizeMillis, numBuckets) must be zero.",
             totalSizeMillis, numBuckets)
    }
    return &RollingNumber{
        buckets:           make(map[int64]*bucket, numBuckets),
        numBuckets:        numBuckets,
        bucketSizeMillis:  totalSizeMillis / int64(numBuckets),
        totalSizeMillis:   totalSizeMillis,
        lastBucketTimeKey: 0,
        mu:                &sync.RWMutex{},
    }, nil
}
// TODO: Consider time injection for easier testing.

// Add value to the current bucket.
func (r *RollingNumber) Add(num float64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    currBucket := r.getCurrentBucket()
    currBucket.add(num)
    return nil
}

// Update max value to the current bucket.
// Does not support negative values.
func (r *RollingNumber) UpdateMax(num float64) error {
    if num < 0 {
        return fmt.Errorf(
            "Max does not support negative value %v.", num)
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    currBucket := r.getCurrentBucket()
    currBucket.updateMax(num)
    return nil
}

// Get rolling sum for the current window.
func (r *RollingNumber) GetRollingSum() float64 {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.getCurrentBucket() // Force clear the old buckets
    sum := float64(0)
    for _, b := range r.buckets {
        sum += b.getSum()
    }
    return sum
}

// Get rolling max for the current window.
func (r *RollingNumber) GetRollingMax() float64 {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.getCurrentBucket() // Force clear the old buckets
    max := float64(0)
    for _, b := range r.buckets {
        if b.getMax() > max {
            max = b.getMax()
        }
    }
    return max
}

// Reset the rolling number to clear all values.
func (r *RollingNumber) Reset() {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.reset()
}

// Reset the rolling number to clear all values.
func (r *RollingNumber) reset() {
    r.buckets = make(map[int64]*bucket, r.numBuckets)
    r.lastBucketTimeKey = 0
}

// Calculate the current bucket based on current time.
// Private with synchronization handled by caller.
func (r *RollingNumber) getCurrentBucket() *bucket {
    return r.getCurrentBucketWithMillis(time.Now().UnixMilli())
}

func (r *RollingNumber) getCurrentBucketWithMillis(currMillis int64) *bucket {
    if r.lastBucketTimeKey == 0 {
        r.buckets[currMillis] = &bucket{}
        r.lastBucketTimeKey = currMillis
        return r.buckets[currMillis]
    }
    if currMillis < r.lastBucketTimeKey + r.bucketSizeMillis {
        return r.buckets[r.lastBucketTimeKey]
    }
    if currMillis > r.lastBucketTimeKey + r.totalSizeMillis {
        r.reset()
        return r.getCurrentBucketWithMillis(currMillis)
    }
    for key := range r.buckets {
        // Remove old buckets
        if currMillis > key + r.totalSizeMillis {
            delete(r.buckets, key)
        }
    }
    // Calculate the new key, intervals in between with empty values will be skipped.
    mod := (currMillis - r.lastBucketTimeKey) % r.bucketSizeMillis
    newKey := currMillis - mod
    r.buckets[newKey] = &bucket{}
    r.lastBucketTimeKey = newKey
    return r.buckets[newKey]
}

// Add value to the bucket.
// Private with synchronization handled by caller.
func (b *bucket) add(num float64) {
    b.adder.add(num)
}

// Update max value to the bucket.
// Private with synchronization handled by caller.
func (b *bucket) updateMax(num float64) {
    b.maxGetter.updateMax(num)
}

func (b *bucket) getSum() float64 {
    return b.adder.value
}

func (b *bucket) getMax() float64 {
    return b.maxGetter.value
}

func (a *addBucket) add(num float64) {
    a.value += num
}

func (m *maxBucket) updateMax(num float64) {
    if num > m.value {
        m.value = num
    }
}
