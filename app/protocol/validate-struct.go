package protocol

type ValidateStructProtocol interface {
	ValidateStruct(params interface{}) error
}
