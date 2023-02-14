package dto

type LoadDocumentsReqDTO struct {
	OthertmsLoads []string `json:"othertmsLoads" validate:"required"`
	BtmsLoads     []string `json:"btmsLoads" validate:"required"`
}
