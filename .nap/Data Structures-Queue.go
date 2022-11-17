package main

type Queue struct {
	len        int
	head, tail int
	q          []int
}

func New(n int) *Queue {
	return &Queue{n, 0, 0, make([]int, n)}
}

func (p *Queue) Enqueue(x int) bool {
	p.q[p.tail] = x
	ntail := (p.tail + 1) % p.len
	ok := false
	if ntail != p.head {
		p.tail = ntail
		ok = true
	}
	return ok
}

func (p *Queue) Dequeue() (int, bool) {
	if p.head == p.tail {
		return 0, false
	}
	x := p.q[p.head]
	p.head = (p.head + 1) % p.len
	return x, true
}
