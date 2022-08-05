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

func MakeHashEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hashReq)
		v, err := srv.Hash(ctx, req.Password)
		if err != nil {
			return hashRes{v, err.Error()}, nil
		}
		return hashRes{v, ""}, nil
	}
}

func MakeValidateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateReq)
		v, err := srv.Validate(ctx, req.Password, req.Hash)
		if err != nil {
			return validateRes{false, err.Error()}, nil
		}
		return validateRes{v, ""}, nil
	}
}

type Endpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
	req := hashReq{password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	hashResp := resp.(hashRes)
	if hashResp.Err != "" {
		return "", errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}

func (e Endpoints) Validate(ctx context.Context, password, hash string) (bool, error) {
	req := validateReq{password, hash}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(validateRes)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}
