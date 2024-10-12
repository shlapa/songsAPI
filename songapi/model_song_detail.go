/*
 * Music info
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// SongDetail struct for SongDetail
// swagger:model
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text string `json:"text"`
	Link string `json:"link"`
}

// NewSongDetail instantiates a new SongDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSongDetail(releaseDate string, text string, link string) *SongDetail {
	this := SongDetail{}
	this.ReleaseDate = releaseDate
	this.Text = text
	this.Link = link
	return &this
}

// NewSongDetailWithDefaults instantiates a new SongDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSongDetailWithDefaults() *SongDetail {
	this := SongDetail{}
	return &this
}

// GetReleaseDate returns the ReleaseDate field value
func (o *SongDetail) GetReleaseDate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ReleaseDate
}

// GetReleaseDateOk returns a tuple with the ReleaseDate field value
// and a boolean to check if the value has been set.
func (o *SongDetail) GetReleaseDateOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.ReleaseDate, true
}

// SetReleaseDate sets field value
func (o *SongDetail) SetReleaseDate(v string) {
	o.ReleaseDate = v
}

// GetText returns the Text field value
func (o *SongDetail) GetText() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Text
}

// GetTextOk returns a tuple with the Text field value
// and a boolean to check if the value has been set.
func (o *SongDetail) GetTextOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Text, true
}

// SetText sets field value
func (o *SongDetail) SetText(v string) {
	o.Text = v
}

// GetLink returns the Link field value
func (o *SongDetail) GetLink() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Link
}

// GetLinkOk returns a tuple with the Link field value
// and a boolean to check if the value has been set.
func (o *SongDetail) GetLinkOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Link, true
}

// SetLink sets field value
func (o *SongDetail) SetLink(v string) {
	o.Link = v
}

func (o SongDetail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["releaseDate"] = o.ReleaseDate
	}
	if true {
		toSerialize["text"] = o.Text
	}
	if true {
		toSerialize["link"] = o.Link
	}
	return json.Marshal(toSerialize)
}

type NullableSongDetail struct {
	value *SongDetail
	isSet bool
}

func (v NullableSongDetail) Get() *SongDetail {
	return v.value
}

func (v *NullableSongDetail) Set(val *SongDetail) {
	v.value = val
	v.isSet = true
}

func (v NullableSongDetail) IsSet() bool {
	return v.isSet
}

func (v *NullableSongDetail) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSongDetail(val *SongDetail) *NullableSongDetail {
	return &NullableSongDetail{value: val, isSet: true}
}

func (v NullableSongDetail) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSongDetail) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

