/*
 * Copyright 2021 InfAI (CC SES)
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

package pkg

import (
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	importRepo "github.com/SENERGY-Platform/import-repository/lib/client"
	"github.com/SENERGY-Platform/import-repository/lib/model"
)

type ImportTypeWithCriteria struct {
	model.ImportType
	Criteria []ImportTypeCriteria
}

type ImportTypeCriteria struct {
	FunctionId string
	AspectId   string
}

func (this *Lib) GetImportTypesWithAspect(token auth.Token, aspectIds []string) (importTypes []ImportTypeWithCriteria, err error, code int) {
	temp, _, err, code := this.importRepo.ListImportTypes(token, importRepo.ImportTypeListOptions{
		Limit:    9999,
		Offset:   0,
		SortBy:   "name.asc",
		Criteria: []model.ImportTypeFilterCriteria{{AspectIds: aspectIds}},
	})
	if err != nil {
		return nil, err, code
	}
	for _, t := range temp {
		importTypes = append(importTypes, ImportTypeWithCriteria{
			ImportType: t,
			Criteria:   importTypeContentVariableToCertList(t.Output),
		})
	}
	return importTypes, err, code
}

func importTypeContentVariableToCertList(cv model.ContentVariable) []ImportTypeCriteria {
	result := []ImportTypeCriteria{{
		FunctionId: cv.FunctionId,
		AspectId:   cv.AspectId,
	}}
	for _, sub := range cv.SubContentVariables {
		result = append(result, importTypeContentVariableToCertList(sub)...)
	}
	return result
}

func (this *Lib) GetImportTypes(token auth.Token) (importTypes []ImportTypeWithCriteria, err error, code int) {
	temp, _, err, code := this.importRepo.ListImportTypes(token, importRepo.ImportTypeListOptions{
		Limit:  9999,
		Offset: 0,
		SortBy: "name.asc",
	})
	if err != nil {
		return nil, err, code
	}
	for _, t := range temp {
		importTypes = append(importTypes, ImportTypeWithCriteria{
			ImportType: t,
			Criteria:   importTypeContentVariableToCertList(t.Output),
		})
	}
	return importTypes, err, code
}
