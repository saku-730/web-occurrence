package model

type CreateWorkstationRequest struct {
	WorkstationName string `json:"workstation_name" binding:"required"`
}

type CreateWorkstationResponse struct {
	WorkstationID   int64  `json:"workstation_id"`
	WorkstationName string `json:"workstation_name"`
}
