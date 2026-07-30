package ratelimit

import (
	"sync"
	"time"
)

// bucket representa un balde de fichas para un cliente
type bucket struct {
	tokens    float64
	lastRefill time.Time
}

// Limiter es un Token Bucket rate limiter en memoria
type Limiter struct {
	mu          sync.Mutex
	buckets     map[string]*bucket
	capacity    float64   
	refillRate  float64   
	cleanupTick time.Duration 
}


// capacity = máx requests en una ráfaga
// refillRate = requests por segundo que se permiten sostenidos
func New(capacity int, refillRate float64) *Limiter {
	l := &Limiter{
		buckets:     make(map[string]*bucket),
		capacity:    float64(capacity),
		refillRate:  refillRate,
		cleanupTick: 5 * time.Minute,
	}

	// periodic clean up to avoid memory leak
	go l.cleanup()

	return l
}


func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[key]
	if !exists {
		b = &bucket{
			tokens:     l.capacity,
			lastRefill: time.Now(),
		}
		l.buckets[key] = b
	}

	// Refill: agregar fichas según el tiempo transcurrido
	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * l.refillRate
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.lastRefill = now

	// Consumir una ficha
	if b.tokens < 1 {
		return false 
	}

	b.tokens--
	return true 
}


func (l *Limiter) cleanup() {
	for {
		time.Sleep(l.cleanupTick)
		l.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for key, b := range l.buckets {
			if b.lastRefill.Before(cutoff) {
				delete(l.buckets, key)
			}
		}
		l.mu.Unlock()
	}
}
