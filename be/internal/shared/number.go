package shared

import "math"

func CheckOverFlowUintWithMulOperator(a uint, b uint) error {
	if a == 0 || b == 0 {
		return nil
	}
	// Công thức kiểm tra tràn số cho phép nhân:
	// Nếu a * b > MaxUint => b > MaxUint / a
	if b > math.MaxUint/a {
		return ErrIntegerOverflow
	}
	return nil
}
