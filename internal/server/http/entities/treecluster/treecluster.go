package treecluster

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/pagination"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
)

type TreeClusterWateringStatus string // @Name WateringStatus

const (
	TreeClusterWateringStatusGood     TreeClusterWateringStatus = "good"
	TreeClusterWateringStatusModerate TreeClusterWateringStatus = "moderate"
	TreeClusterWateringStatusBad      TreeClusterWateringStatus = "bad"
	TreeClusterWateringStatusUnknown  TreeClusterWateringStatus = "unknown"
)

type TreeSoilCondition string // @Name SoilCondition

const (
	TreeSoilConditionSchluffig TreeSoilCondition = "schluffig"
)

type TreeClusterResponse struct {
	ID             int32                     `json:"id,omitempty"`
	CreatedAt      time.Time                 `json:"created_at,omitempty"`
	UpdatedAt      time.Time                 `json:"updated_at,omitempty"`
	WateringStatus TreeClusterWateringStatus `json:"watering_status,omitempty"`
	LastWatered    time.Time                 `json:"last_watered,omitempty"`
	MoistureLevel  float64                   `json:"moisture_level,omitempty"`
	Region         string                    `json:"region,omitempty"`
	Address        string                    `json:"address,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Archived       bool                      `json:"archived,omitempty"`
	Latitude       float64                   `json:"latitude,omitempty"`
	Longitude      float64                   `json:"longitude,omitempty"`
	Trees          []*tree.TreeResponse      `json:"trees,omitempty"`
	SoilCondition  TreeSoilCondition         `json:"soil_condition,omitempty"`
} // @Name TreeCluster

type TreeClusterListResponse struct {
	Data       []*TreeClusterResponse `json:"data,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
} // @Name TreeClusterList

type TreeClusterCreateRequest struct {
	WateringStatus TreeClusterWateringStatus `json:"watering_status,omitempty"`
	LastWatered    time.Time                 `json:"last_watered,omitempty"`
	MoistureLevel  float64                   `json:"moisture_level,omitempty"`
	Region         string                    `json:"region,omitempty"`
	Address        string                    `json:"address,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Archived       bool                      `json:"archived,omitempty"`
	Latitude       float64                   `json:"latitude,omitempty"`
	Longitude      float64                   `json:"longitude,omitempty"`
	TreeIDs        []*int32                  `json:"tree_ids,omitempty"`
	SoilCondition  TreeSoilCondition         `json:"soil_condition,omitempty"`
} // @Name TreeClusterCreate

type TreeClusterUpdateRequest struct {
	WateringStatus TreeClusterWateringStatus `json:"watering_status,omitempty"`
	LastWatered    time.Time                 `json:"last_watered,omitempty"`
	MoistureLevel  float64                   `json:"moisture_level,omitempty"`
	Region         string                    `json:"region,omitempty"`
	Address        string                    `json:"address,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Archived       bool                      `json:"archived,omitempty"`
	Latitude       float64                   `json:"latitude,omitempty"`
	Longitude      float64                   `json:"longitude,omitempty"`
	TreeIDs        []*int32                  `json:"tree_ids,omitempty"`
	SoilCondition  TreeSoilCondition         `json:"soil_condition,omitempty"`
} // @Name TreeClusterUpdate

type TreeClusterAddTreesRequest struct {
	TreeIDs []*int32 `json:"tree_ids,omitempty"`
} // @Name TreeClusterAddTrees
