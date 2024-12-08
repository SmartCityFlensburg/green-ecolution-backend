/*
Green Space Management API

This is the API for the Green Ecolution Management System.

API version: develop
Contact: info@green-ecolution.de
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the TreeClusterUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TreeClusterUpdate{}

// TreeClusterUpdate struct for TreeClusterUpdate
type TreeClusterUpdate struct {
	Address string `json:"address"`
	Description string `json:"description"`
	Name string `json:"name"`
	SoilCondition SoilCondition `json:"soil_condition"`
	TreeIds []int32 `json:"tree_ids"`
}

type _TreeClusterUpdate TreeClusterUpdate

// NewTreeClusterUpdate instantiates a new TreeClusterUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTreeClusterUpdate(address string, description string, name string, soilCondition SoilCondition, treeIds []int32) *TreeClusterUpdate {
	this := TreeClusterUpdate{}
	this.Address = address
	this.Description = description
	this.Name = name
	this.SoilCondition = soilCondition
	this.TreeIds = treeIds
	return &this
}

// NewTreeClusterUpdateWithDefaults instantiates a new TreeClusterUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTreeClusterUpdateWithDefaults() *TreeClusterUpdate {
	this := TreeClusterUpdate{}
	return &this
}

// GetAddress returns the Address field value
func (o *TreeClusterUpdate) GetAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Address
}

// GetAddressOk returns a tuple with the Address field value
// and a boolean to check if the value has been set.
func (o *TreeClusterUpdate) GetAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Address, true
}

// SetAddress sets field value
func (o *TreeClusterUpdate) SetAddress(v string) {
	o.Address = v
}

// GetDescription returns the Description field value
func (o *TreeClusterUpdate) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *TreeClusterUpdate) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *TreeClusterUpdate) SetDescription(v string) {
	o.Description = v
}

// GetName returns the Name field value
func (o *TreeClusterUpdate) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *TreeClusterUpdate) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *TreeClusterUpdate) SetName(v string) {
	o.Name = v
}

// GetSoilCondition returns the SoilCondition field value
func (o *TreeClusterUpdate) GetSoilCondition() SoilCondition {
	if o == nil {
		var ret SoilCondition
		return ret
	}

	return o.SoilCondition
}

// GetSoilConditionOk returns a tuple with the SoilCondition field value
// and a boolean to check if the value has been set.
func (o *TreeClusterUpdate) GetSoilConditionOk() (*SoilCondition, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SoilCondition, true
}

// SetSoilCondition sets field value
func (o *TreeClusterUpdate) SetSoilCondition(v SoilCondition) {
	o.SoilCondition = v
}

// GetTreeIds returns the TreeIds field value
func (o *TreeClusterUpdate) GetTreeIds() []int32 {
	if o == nil {
		var ret []int32
		return ret
	}

	return o.TreeIds
}

// GetTreeIdsOk returns a tuple with the TreeIds field value
// and a boolean to check if the value has been set.
func (o *TreeClusterUpdate) GetTreeIdsOk() ([]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.TreeIds, true
}

// SetTreeIds sets field value
func (o *TreeClusterUpdate) SetTreeIds(v []int32) {
	o.TreeIds = v
}

func (o TreeClusterUpdate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TreeClusterUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["address"] = o.Address
	toSerialize["description"] = o.Description
	toSerialize["name"] = o.Name
	toSerialize["soil_condition"] = o.SoilCondition
	toSerialize["tree_ids"] = o.TreeIds
	return toSerialize, nil
}

func (o *TreeClusterUpdate) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"address",
		"description",
		"name",
		"soil_condition",
		"tree_ids",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varTreeClusterUpdate := _TreeClusterUpdate{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varTreeClusterUpdate)

	if err != nil {
		return err
	}

	*o = TreeClusterUpdate(varTreeClusterUpdate)

	return err
}

type NullableTreeClusterUpdate struct {
	value *TreeClusterUpdate
	isSet bool
}

func (v NullableTreeClusterUpdate) Get() *TreeClusterUpdate {
	return v.value
}

func (v *NullableTreeClusterUpdate) Set(val *TreeClusterUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableTreeClusterUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableTreeClusterUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTreeClusterUpdate(val *TreeClusterUpdate) *NullableTreeClusterUpdate {
	return &NullableTreeClusterUpdate{value: val, isSet: true}
}

func (v NullableTreeClusterUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTreeClusterUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


