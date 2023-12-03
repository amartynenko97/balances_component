package constants

const (
	ExNameBalances                   = "e.balances.forward"
	CreateAccountBalanceRequestQueue = "q.balances.request.CreateAccountBalanceRequest"
	RkCreateAccountBalanceResponse   = "r.balances.#.CreateAccountBalanceResponse.#"

	GetAccountBalanceRequestQueue  = "q.balances.request.GetAccountBalanceRequest"
	GetAccountBalanceResponseQueue = "q.balances.request.GetAccountBalanceResponse"
	RkGetAccountBalanceRequest     = "r.balances.#.GetAccountBalanceRequest.#"
	RkGetAccountBalanceResponse    = "r.balances.#.GetAccountBalanceResponse.#"
)
