package vault

import "context"

type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password, hash string) (bool, error)
}

type vaultService struct{}

func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (vaultService) Validate(ctx context.Context, password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
type hashReq struct {
	Password string `json:"password"`
}

func decodeHashReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req hashReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

type hashRes struct {
	Hash string `json:"hash"`
	Err  string `json:"err,omitempty"`
}

type validateReq struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type validateRes struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

func decodeValidateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req validateReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRes(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}
