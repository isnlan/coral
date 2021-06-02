package gateway

type Producer interface {
	APIUpload(api *API) error
	APICallRecord(entity *APICallEntity) error
	ContractCallRecord(entity *ContractCallEntity) error
}

type Consumer interface {
	APIHandler(api *API) error
	APICallHandler(entity *APICallEntity) error
	ContractCallHandler(entity *ContractCallEntity) error
}
