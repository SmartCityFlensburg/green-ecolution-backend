package entities

import "time"

type WateringPlanStatus string // @Name WateringPlanStatus

const (
	WateringPlanStatusPlanned     WateringPlanStatus = "planned"
	WateringPlanStatusActive      WateringPlanStatus = "active"
	WateringPlanStatusCanceled    WateringPlanStatus = "canceled"
	WateringPlanStatusFinished    WateringPlanStatus = "finished"
	WateringPlanStatusNotCompeted WateringPlanStatus = "not competed"
	WateringPlanStatusUnknown     WateringPlanStatus = "unknown"
)

type WateringPlanResponse struct {
	ID                 int32                        `json:"id"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	Date               time.Time                    `json:"date"`
	Description        string                       `json:"description"`
	Status             WateringPlanStatus           `json:"status"`
	Distance           *float64                     `json:"distance"`
	TotalWaterRequired *float64                     `json:"total_water_required"`
	Users              []*UserResponse              `json:"users"`
	TreeClusters       []*TreeClusterInListResponse `json:"treeclusters"`
	Transporter        *VehicleResponse             `json:"transporter"`
	Trailer            *VehicleResponse             `json:"trailer" validate:"optional"`
	CancellationNote   string                       `json:"cancellation_note"`
	ConsumedWaterList  []*ConsumedWaterItem         `json:"consumed_water_list"`
} // @Name WateringPlan

type WateringPlanInListResponse struct {
	ID                 int32                        `json:"id"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	Date               time.Time                    `json:"date"`
	Description        string                       `json:"description"`
	Status             WateringPlanStatus           `json:"status"`
	Distance           *float64                     `json:"distance"`
	TotalWaterRequired *float64                     `json:"total_water_required"`
	Users              []*UserResponse              `json:"users"`
	TreeClusters       []*TreeClusterInListResponse `json:"treeclusters"`
	Transporter        *VehicleResponse             `json:"transporter"`
	Trailer            *VehicleResponse             `json:"trailer" validate:"optional"`
	CancellationNote   string                       `json:"cancellation_note"`
} // @Name WateringPlanInListResponse

type WateringPlanListResponse struct {
	Data       []*WateringPlanInListResponse `json:"data"`
	Pagination *Pagination                   `json:"pagination"`
} // @Name WateringPlanList

type WateringPlanCreateRequest struct {
	Date           time.Time `json:"date"`
	Description    string    `json:"description"`
	TreeClusterIDs []*int32  `json:"tree_cluster_ids"`
	TransporterID  *int32    `json:"transporter_id"`
	TrailerID      *int32    `json:"trailer_id"`
	Users          []*int32  `json:"users_ids"`
} // @Name WateringPlanCreate

type WateringPlanUpdateRequest struct {
	Date              time.Time            `json:"date"`
	Description       string               `json:"description"`
	TreeClusterIDs    []*int32             `json:"tree_cluster_ids"`
	TransporterID     *int32               `json:"transporter_id"`
	TrailerID         *int32               `json:"trailer_id"`
	Users             []*int32             `json:"users_ids"`
	CancellationNote  string               `json:"cancellation_note"`
	Status            WateringPlanStatus   `json:"status"`
	ConsumedWaterList []*ConsumedWaterItem `json:"consumed_water_list"`
} // @Name WateringPlanUpdate

type ConsumedWaterItem struct {
	WateringPlanID int32
	TreeClusterID  int32
	ConsumedWater  *float64
} // @Name ConsumedWaterItem
