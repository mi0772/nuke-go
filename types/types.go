package types

type PushItemRequest struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type PushItemStringRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PushItemFileRequest struct {
	Key string `json:"key"`
}

type PopRequest struct {
	Key string `json:"key"`
}

type PartitionDetailsResponse struct {
	Partition uint8
	Entries   uint
}

type KeysResponse struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Key       string `json:"key"`
	Partition uint8  `json:"partition"`
}
