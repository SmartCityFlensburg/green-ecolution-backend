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

// checks if the RoleList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RoleList{}

// RoleList struct for RoleList
type RoleList struct {
	Data []Role `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type _RoleList RoleList

// NewRoleList instantiates a new RoleList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRoleList(data []Role, pagination Pagination) *RoleList {
	this := RoleList{}
	this.Data = data
	this.Pagination = pagination
	return &this
}

// NewRoleListWithDefaults instantiates a new RoleList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRoleListWithDefaults() *RoleList {
	this := RoleList{}
	return &this
}

// GetData returns the Data field value
func (o *RoleList) GetData() []Role {
	if o == nil {
		var ret []Role
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *RoleList) GetDataOk() ([]Role, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *RoleList) SetData(v []Role) {
	o.Data = v
}

// GetPagination returns the Pagination field value
func (o *RoleList) GetPagination() Pagination {
	if o == nil {
		var ret Pagination
		return ret
	}

	return o.Pagination
}

// GetPaginationOk returns a tuple with the Pagination field value
// and a boolean to check if the value has been set.
func (o *RoleList) GetPaginationOk() (*Pagination, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Pagination, true
}

// SetPagination sets field value
func (o *RoleList) SetPagination(v Pagination) {
	o.Pagination = v
}

func (o RoleList) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RoleList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	toSerialize["pagination"] = o.Pagination
	return toSerialize, nil
}

func (o *RoleList) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"data",
		"pagination",
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

	varRoleList := _RoleList{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRoleList)

	if err != nil {
		return err
	}

	*o = RoleList(varRoleList)

	return err
}

type NullableRoleList struct {
	value *RoleList
	isSet bool
}

func (v NullableRoleList) Get() *RoleList {
	return v.value
}

func (v *NullableRoleList) Set(val *RoleList) {
	v.value = val
	v.isSet = true
}

func (v NullableRoleList) IsSet() bool {
	return v.isSet
}

func (v *NullableRoleList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRoleList(val *RoleList) *NullableRoleList {
	return &NullableRoleList{value: val, isSet: true}
}

func (v NullableRoleList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRoleList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


