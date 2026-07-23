package model

import (
	"fmt"
	"icon_exchange/internal/shared"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Market struct {
	gorm.Model

	Name           string `gorm:"type:varchar(255);not null;unique"`
	Symbol         string `gorm:"type:varchar(255);not null;unique"`
	Decimals       int32  `gorm:"type:int;not null"`          // number of smallest unit which one symbol can device
	Precision      int32  `gorm:"type:int;not null"`          // maximum digit number after comma of quantity value
	SmallestUnit   string `gorm:"type:varchar(255);not null"` // name of smallest unit
	IsBaseCurrency bool   `gorm:"type:bool;not null"`
}

// ValidateQuantity kiểm tra số lượng có hợp lệ theo Decimals và Precision hay không.
// Trả về error dạng shared.Error nếu không hợp lệ, ngược lại trả về nil.
func (m Market) ValidateQuantity(floatQuantity float64) error {
	if floatQuantity <= 0 {
		return shared.ErrOrderQuantityInvalid
	}

	// 1. Chuyển float64 sang string định dạng tiêu chuẩn để tránh sai số dấu phẩy động
	quantityStr := fmt.Sprintf("%.18f", floatQuantity)
	quantityStr = strings.TrimRight(strings.TrimRight(quantityStr, "0"), ".")

	// Tách phần nguyên và phần thập phân để đếm số chữ số sau dấu phẩy
	parts := strings.Split(quantityStr, ".")
	decimalPlaces := 0
	if len(parts) > 1 {
		decimalPlaces = len(parts[1])
	}

	// 2. Validate với Precision (số chữ số tối đa sau dấu phẩy cho phép đặt lệnh)
	if int32(decimalPlaces) > m.Precision {
		return shared.ErrOrderQuantityPrecisionExceeded
	}

	// 3. Validate với Decimals (giới hạn chia nhỏ nhất gốc của blockchain/hệ thống)
	if int32(decimalPlaces) > m.Decimals {
		return shared.ErrOrderQuantityDecimalsExceeded
	}

	// 4. Validate giá trị tối thiểu cho phép (tương ứng với 10^-Precision)
	minQuantity := 1 / math.Pow(10, float64(m.Precision))
	epsilon := 1e-9
	if floatQuantity < (minQuantity - epsilon) {
		return shared.ErrOrderQuantityTooSmall
	}

	return nil
}

// ToSmallestUnit nhận vào floatQuantity, thực hiện validate trước.
// Nếu hợp lệ, trả về số lượng quy đổi ra đơn vị nhỏ nhất (Satoshi, Wei...) và nil.
// Nếu KHÔNG hợp lệ, trả về 0 và error chi tiết.
func (m Market) ToSmallestUnit(floatQuantity float64) (uint, error) {
	// Thực hiện validate trước khi chuyển đổi
	if err := m.ValidateQuantity(floatQuantity); err != nil {
		return 0, err
	}

	// Tính hệ số nhân: 10^Decimals
	multiplier := math.Pow(10, float64(m.Decimals))

	// Làm tròn để triệt tiêu sai số float64 trước khi ép kiểu
	smallestUnitQty := uint(math.Round(floatQuantity * multiplier))

	return smallestUnitQty, nil
}
