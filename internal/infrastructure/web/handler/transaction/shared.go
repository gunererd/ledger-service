package transaction

type UserResponse struct {
	Id   string `json:"id" doc:"User ID"`
	Type string `json:"type" doc:"User type"`
}