package main

// Mock Random Generator
type mockRandomGenerator struct {
	float64Result float64
	intnResult    int
}

func (m *mockRandomGenerator) Float64() float64 {
	return m.float64Result
}

func (m *mockRandomGenerator) Intn(n int) int {
	return m.intnResult
}
