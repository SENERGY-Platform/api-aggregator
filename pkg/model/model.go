/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

type Characteristic struct {
	Id                 string           `json:"id"`
	Name               string           `json:"name"`
	DisplayUnit        string           `json:"display_unit"`
	Type               string           `json:"type"`
	MinValue           interface{}      `json:"min_value,omitempty"`
	MaxValue           interface{}      `json:"max_value,omitempty"`
	Value              interface{}      `json:"value,omitempty"`
	SubCharacteristics []Characteristic `json:"sub_characteristics"`
}

type Concept struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	BaseCharacteristicId string   `json:"base_characteristic_id"`
	CharacteristicIds    []string `json:"characteristic_ids"`
}

type ConceptInfo struct {
	Concept
	BaseCharacteristic Characteristic `json:"base_characteristic"`
}

type Function struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	ConceptId   string `json:"concept_id"`
}

type FunctionInfo struct {
	Function
	Concept ConceptInfo `json:"concept"`
}
