package api

type AuthInfo struct {
	Environment string
	Role        string
	Username    string
	ValidFor    int
}
