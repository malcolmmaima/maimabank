package util

// All supported currencies in our banking service
const (
	USD = "USD"
	EUR = "EUR"
	KES = "KES"
	GBP = "GBP"
)

// Check if a currency is supported by our banking service
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, KES, GBP:
		return true
	}
	return false
}