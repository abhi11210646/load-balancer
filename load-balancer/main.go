package main

func main() {
	lb := NewLoadBalancer()
	lb.ListenAndServe()
}
