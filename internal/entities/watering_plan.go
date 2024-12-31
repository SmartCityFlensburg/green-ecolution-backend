package entities

import (
	"time"

	"github.com/google/uuid"
)

type WateringPlanStatus string

const (
	WateringPlanStatusPlanned     WateringPlanStatus = "planned"
	WateringPlanStatusActive      WateringPlanStatus = "active"
	WateringPlanStatusCanceled    WateringPlanStatus = "canceled"
	WateringPlanStatusFinished    WateringPlanStatus = "finished"
	WateringPlanStatusNotCompeted WateringPlanStatus = "not competed"
	WateringPlanStatusUnknown     WateringPlanStatus = "unknown"
)

type WateringPlan struct {
	ID                 int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Date               time.Time
	Description        string
	Status             WateringPlanStatus
	Distance           *float64
	TotalWaterRequired *float64
	UserIDs            []*uuid.UUID
	TreeClusters       []*TreeCluster
	Transporter        *Vehicle
	Trailer            *Vehicle
	CancellationNote   string
	Evaluation         []*EvaluationValue
}

type WateringPlanCreate struct {
	Date           time.Time `validate:"required"`
	Description    string
	TreeClusterIDs []*int32 `validate:"required,min=1,dive,required"`
	TransporterID  *int32   `validate:"required"`
	TrailerID      *int32
	// Users           []*int32
}

type WateringPlanUpdate struct {
	Date             time.Time `validate:"required"`
	Description      string
	TreeClusterIDs   []*int32 `validate:"required,min=1,dive,required"`
	TransporterID    *int32   `validate:"required"`
	TrailerID        *int32
	CancellationNote string
	Status           WateringPlanStatus `validate:"oneof=planned active canceled finished 'not competed' unknown"`
	Evaluation       []*EvaluationValue
	// Users           []*int32
}

type EvaluationValue struct {
	WateringPlanID int32
	TreeClusterID  int32
	ConsumedWater  *float64
}
