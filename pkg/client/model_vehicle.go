/*
Green Space Management API

This is the API for the Green Ecolution Management System.

API version: develop
Contact: info@green-ecolution.de
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Vehicle type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Vehicle{}

// Vehicle struct for Vehicle
type Vehicle struct {
	AdditionalInformation map[string]interface{} `json:"additional_information"`
	CreatedAt             string                 `json:"created_at"`
	Description           string                 `json:"description"`
	DrivingLicense        DrivingLicense         `json:"driving_license"`
	Height                float32                `json:"height"`
	Id                    int32                  `json:"id"`
	Length                float32                `json:"length"`
	Model                 string                 `json:"model"`
	NumberPlate           string                 `json:"number_plate"`
	Provider              string                 `json:"provider"`
	Status                VehicleStatus          `json:"status"`
	Type                  VehicleType            `json:"type"`
	UpdatedAt             string                 `json:"updated_at"`
	WaterCapacity         float32                `json:"water_capacity"`
	Weight                float32                `json:"weight"`
	Width                 float32                `json:"width"`
}

type _Vehicle Vehicle

// NewVehicle instantiates a new Vehicle object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVehicle(additionalInformation map[string]interface{}, createdAt string, description string, drivingLicense DrivingLicense, height float32, id int32, length float32, model string, numberPlate string, provider string, status VehicleStatus, type_ VehicleType, updatedAt string, waterCapacity float32, weight float32, width float32) *Vehicle {
	this := Vehicle{}
	this.AdditionalInformation = additionalInformation
	this.CreatedAt = createdAt
	this.Description = description
	this.DrivingLicense = drivingLicense
	this.Height = height
	this.Id = id
	this.Length = length
	this.Model = model
	this.NumberPlate = numberPlate
	this.Provider = provider
	this.Status = status
	this.Type = type_
	this.UpdatedAt = updatedAt
	this.WaterCapacity = waterCapacity
	this.Weight = weight
	this.Width = width
	return &this
}

// NewVehicleWithDefaults instantiates a new Vehicle object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVehicleWithDefaults() *Vehicle {
	this := Vehicle{}
	return &this
}

// GetAdditionalInformation returns the AdditionalInformation field value
func (o *Vehicle) GetAdditionalInformation() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.AdditionalInformation
}

// GetAdditionalInformationOk returns a tuple with the AdditionalInformation field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetAdditionalInformationOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.AdditionalInformation, true
}

// SetAdditionalInformation sets field value
func (o *Vehicle) SetAdditionalInformation(v map[string]interface{}) {
	o.AdditionalInformation = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Vehicle) GetCreatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetCreatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Vehicle) SetCreatedAt(v string) {
	o.CreatedAt = v
}

// GetDescription returns the Description field value
func (o *Vehicle) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *Vehicle) SetDescription(v string) {
	o.Description = v
}

// GetDrivingLicense returns the DrivingLicense field value
func (o *Vehicle) GetDrivingLicense() DrivingLicense {
	if o == nil {
		var ret DrivingLicense
		return ret
	}

	return o.DrivingLicense
}

// GetDrivingLicenseOk returns a tuple with the DrivingLicense field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetDrivingLicenseOk() (*DrivingLicense, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DrivingLicense, true
}

// SetDrivingLicense sets field value
func (o *Vehicle) SetDrivingLicense(v DrivingLicense) {
	o.DrivingLicense = v
}

// GetHeight returns the Height field value
func (o *Vehicle) GetHeight() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Height
}

// GetHeightOk returns a tuple with the Height field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetHeightOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Height, true
}

// SetHeight sets field value
func (o *Vehicle) SetHeight(v float32) {
	o.Height = v
}

// GetId returns the Id field value
func (o *Vehicle) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Vehicle) SetId(v int32) {
	o.Id = v
}

// GetLength returns the Length field value
func (o *Vehicle) GetLength() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Length
}

// GetLengthOk returns a tuple with the Length field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetLengthOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Length, true
}

// SetLength sets field value
func (o *Vehicle) SetLength(v float32) {
	o.Length = v
}

// GetModel returns the Model field value
func (o *Vehicle) GetModel() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Model
}

// GetModelOk returns a tuple with the Model field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetModelOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Model, true
}

// SetModel sets field value
func (o *Vehicle) SetModel(v string) {
	o.Model = v
}

// GetNumberPlate returns the NumberPlate field value
func (o *Vehicle) GetNumberPlate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NumberPlate
}

// GetNumberPlateOk returns a tuple with the NumberPlate field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetNumberPlateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumberPlate, true
}

// SetNumberPlate sets field value
func (o *Vehicle) SetNumberPlate(v string) {
	o.NumberPlate = v
}

// GetProvider returns the Provider field value
func (o *Vehicle) GetProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Provider
}

// GetProviderOk returns a tuple with the Provider field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Provider, true
}

// SetProvider sets field value
func (o *Vehicle) SetProvider(v string) {
	o.Provider = v
}

// GetStatus returns the Status field value
func (o *Vehicle) GetStatus() VehicleStatus {
	if o == nil {
		var ret VehicleStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetStatusOk() (*VehicleStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *Vehicle) SetStatus(v VehicleStatus) {
	o.Status = v
}

// GetType returns the Type field value
func (o *Vehicle) GetType() VehicleType {
	if o == nil {
		var ret VehicleType
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetTypeOk() (*VehicleType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *Vehicle) SetType(v VehicleType) {
	o.Type = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *Vehicle) GetUpdatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetUpdatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *Vehicle) SetUpdatedAt(v string) {
	o.UpdatedAt = v
}

// GetWaterCapacity returns the WaterCapacity field value
func (o *Vehicle) GetWaterCapacity() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.WaterCapacity
}

// GetWaterCapacityOk returns a tuple with the WaterCapacity field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetWaterCapacityOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WaterCapacity, true
}

// SetWaterCapacity sets field value
func (o *Vehicle) SetWaterCapacity(v float32) {
	o.WaterCapacity = v
}

// GetWeight returns the Weight field value
func (o *Vehicle) GetWeight() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Weight
}

// GetWeightOk returns a tuple with the Weight field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetWeightOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Weight, true
}

// SetWeight sets field value
func (o *Vehicle) SetWeight(v float32) {
	o.Weight = v
}

// GetWidth returns the Width field value
func (o *Vehicle) GetWidth() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Width
}

// GetWidthOk returns a tuple with the Width field value
// and a boolean to check if the value has been set.
func (o *Vehicle) GetWidthOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Width, true
}

// SetWidth sets field value
func (o *Vehicle) SetWidth(v float32) {
	o.Width = v
}

func (o Vehicle) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Vehicle) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["additional_information"] = o.AdditionalInformation
	toSerialize["created_at"] = o.CreatedAt
	toSerialize["description"] = o.Description
	toSerialize["driving_license"] = o.DrivingLicense
	toSerialize["height"] = o.Height
	toSerialize["id"] = o.Id
	toSerialize["length"] = o.Length
	toSerialize["model"] = o.Model
	toSerialize["number_plate"] = o.NumberPlate
	toSerialize["provider"] = o.Provider
	toSerialize["status"] = o.Status
	toSerialize["type"] = o.Type
	toSerialize["updated_at"] = o.UpdatedAt
	toSerialize["water_capacity"] = o.WaterCapacity
	toSerialize["weight"] = o.Weight
	toSerialize["width"] = o.Width
	return toSerialize, nil
}

func (o *Vehicle) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"additional_information",
		"created_at",
		"description",
		"driving_license",
		"height",
		"id",
		"length",
		"model",
		"number_plate",
		"provider",
		"status",
		"type",
		"updated_at",
		"water_capacity",
		"weight",
		"width",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varVehicle := _Vehicle{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVehicle)

	if err != nil {
		return err
	}

	*o = Vehicle(varVehicle)

	return err
}

type NullableVehicle struct {
	value *Vehicle
	isSet bool
}

func (v NullableVehicle) Get() *Vehicle {
	return v.value
}

func (v *NullableVehicle) Set(val *Vehicle) {
	v.value = val
	v.isSet = true
}

func (v NullableVehicle) IsSet() bool {
	return v.isSet
}

func (v *NullableVehicle) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVehicle(val *Vehicle) *NullableVehicle {
	return &NullableVehicle{value: val, isSet: true}
}

func (v NullableVehicle) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVehicle) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
