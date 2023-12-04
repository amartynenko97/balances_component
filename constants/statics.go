package constants

const (
	ExNameBalances            = "e.balances.forward"
	CreateAccountRequestQueue = "q.balances.request.CreateAccountRequest"
	RkCreateAccountResponse   = "r.balances.#.CreateAccountResponse.#"

	GetAccountBalanceRequestQueue  = "q.balances.request.GetAccountBalanceRequest"
	GetAccountBalanceResponseQueue = "q.balances.request.GetAccountBalanceResponse"
	RkGetAccountBalanceRequest     = "r.balances.#.GetAccountBalanceRequest.#"
	RkGetAccountBalanceResponse    = "r.balances.#.GetAccountBalanceResponse.#"
)
