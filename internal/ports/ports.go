package ports

type BalanceProvider interface {
	Balance(address string) (float64, error)
}
