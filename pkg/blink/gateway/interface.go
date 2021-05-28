package gateway

type Producer interface {
	ApiUpload(api *Api) error
	ApiCallRecord(entity *ApiCallEntity) error
	ContractCallRecord(entity *ContractCallEntity) error
}

type Consumer interface {
	ApiHandler(api *Api) error
	ApiCallHandler(entity *ApiCallEntity) error
	ContractCallHandler(entity *ContractCallEntity) error
}
