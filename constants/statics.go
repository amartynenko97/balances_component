package constants

const (
	ExNameBalances            = "e.balances.forward"
	CreateAccountRequestQueue = "q.balances.request.CreateAccountRequest"
	RkCreateAccountResponse   = "r.balances.#.CreateAccountResponse.#"

	GetAccountBalanceRequestQueue = "q.balances.request.GetAccountBalanceRequest"
	RkGetAccountBalanceResponse   = "r.balances.#.GetAccountBalanceResponse.#"
)

const (
	QueueTypeCreateAccount      = "CreateAccount"
	QueueTypeGetAccountBalances = ""
)
