package repo

import "context"

type WalletRepoImpl struct {
}

func NewWallet() *WalletRepoImpl {
	return &WalletRepoImpl{}
}

func (wr *WalletRepoImpl) GetWalletTransactions(ctx context.Context) ([]string, error) {
	return []string{"1234", "5678"}, nil
}
