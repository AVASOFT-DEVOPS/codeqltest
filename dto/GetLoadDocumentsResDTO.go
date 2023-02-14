package dto

type ImagingAPIResDTO struct {
	Error     bool   `json:"error"`
	Search    string `json:"search"`
	Message   string `json:"message"`
	Documents []struct {
		Id            string      `json:"id"`
		Load          string      `json:"load"`
		Path          string      `json:"path"`
		Type          *string     `json:"type"`
		TypeId        *string     `json:"type_id"`
		UsedIn        string      `json:"used_in"`
		Description   string      `json:"description"`
		CreatedAt     string      `json:"created_at"`
		UpdatedAt     string      `json:"updated_at"`
		User          string      `json:"user"`
		Carrier       *string     `json:"carrier"`
		ParentId      interface{} `json:"parent_id"`
		Connector     interface{} `json:"connector"`
		Truck         interface{} `json:"truck"`
		ReceivedAt    interface{} `json:"received_at"`
		Customer      interface{} `json:"customer"`
		CustomerAgent interface{} `json:"customer_agent"`
		IsLoad        string      `json:"is_load"`
		NumPages      *string     `json:"num_pages"`
		Settlement    interface{} `json:"settlement"`
		Body          string      `json:"body"`
		Categories    []string    `json:"categories"`
	} `json:"documents"`
	Field5 int `json:"0"`
}

type LoadDocumentsResDTO struct {
	LoadDocResult []LoadDocsResult `json:"loadDocResult"`
}

type LoadDocsResult struct {
	LoadId    string          `json:"loadId"`
	Documents []LoadDocuments `json:"documents"`
}

type LoadDocuments struct {
	Id      string `json:"id"`
	LoadId  string `json:"loadId"`
	Path    string `json:"path"`
	DocName string `json:"docName"`
	TypeId  string `json:"typeId"`
	DocUrl  string `json:"docUrl"`
}
