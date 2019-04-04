/*
 * Copyright 2019 InfAI (CC SES)
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

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJvd25lciIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.M33n6BgW1v-RcR0XaI4z288FwnctuijTuaHDIKBnKpI"

func testget(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + " " + string(b))
	}
	return nil
}

func testpost(url string, body interface{}) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + " " + string(b))
	}
	return nil
}

///filter/devices/state/:value/usertag/:tag/orderby/name/asc
func TestFilterDeviceStateUsertagOrderbyNameAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/usertag/tag/orderby/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&usertag=tag&sort=name.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

///filter/devices/state/:value/usertag/:tag/orderby/name/asc
func TestFilterDeviceStateUsertagOrderbyNameAscDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/usertag/tag/orderby/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&usertag=tag&sort=name")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

///filter/devices/state/:value/usertag/:tag/orderby/name/desc
func TestFilterDeviceStateUsertagOrderbyNameDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/usertag/tag/orderby/name/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&usertag=tag&sort=name.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value/tag/:tag/orderby/name/asc
func TestFilterDeviceStateTagOrderbyNameAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/tag/tag/orderby/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&tag=tag&sort=name.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value/tag/:tag/orderby/name/asc
func TestFilterDeviceStateTagOrderbyNameAscDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/tag/tag/orderby/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&tag=tag&sort=name")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

///filter/devices/state/:value/tag/:tag/orderby/name/desc
func TestFilterDeviceStateTagOrderbyNameDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/tag/tag/orderby/name/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&tag=tag&sort=name.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value/name/asc
func TestListAllDevicesAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&sort=name.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value/name/asc
func TestListAllDevicesDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/name/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&sort=name.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value/name/asc
func TestListAllDevices(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value&sort=name")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /filter/devices/state/:value
func TestListAllDevicesUnsorted(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/filter/devices/state/value")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?state=value")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset
func TestListDevicesUnsorted(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset
func TestListDevicesUnsortedDefaultsOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset
func TestListDevicesUnsortedDefaultsLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset
func TestSearchDevices(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset
func TestSearchDevicesDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset
func TestSearchDevicesDefaultOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset
func TestSearchDevicesDefaultLimitAndOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset/:orderfeature/asc
func TestListDevices(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset/:orderfeature/asc
func TestListDevicesAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset/:orderfeature/asc
func TestListDevicesAscLimitDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/devices/:limit/:offset/:orderfeature/desc
func TestListDevicesDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/devices/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset/:orderfeature/asc
func TestSearchDevicesAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature.asc&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset/:orderfeature/desc
func TestSearchDevicesDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature.desc&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset/:orderfeature/asc
func TestSearchDevicesAscDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?limit=100&offset=0&sort=orderfeature&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/devices/:query/:limit/:offset/:orderfeature/asc
func TestSearchDevicesAscDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/devices/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?offset=0&sort=orderfeature&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value
func TestSelectDeviceByTag(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDeviceByTagLimitedOrdered(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value&limit=100&offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDeviceByTagLimitedOrderedDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value&limit=100&offset=0&sort=orderfeature.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDeviceByTagLimitedOrderedDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value&limit=100&offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDeviceByTagLimitedOrderedDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value&offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/tag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDeviceByTagLimitedOrderedDefaultSort(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/tag/value/100/0/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?tag=value&offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value
func TestSelectDevicesByUsertag(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByUsertagWithLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value&limit=100&offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByUsertagWithLimitDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value&limit=100&offset=0&sort=orderfeature.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByUsertagWithLimitDefaultSortDir(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value&limit=100&offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByUsertagWithLimitDefaultSort(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value/100/0/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value&limit=100&offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByUsertagWithLimitDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/select/devices/usertag/value/100/0/name/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?usertag=value&offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/ids
func TestSelectDevicesByIds(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testpost(oldUrl+"/select/devices/ids", []string{"a", "b", "c"})
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?ids=a,b,c")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /select/devices/ids/:limit/:offset/:orderfeature/:direction
func TestSelectDevicesByIdsOrdered(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testpost(oldUrl+"/select/devices/ids/100/0/orderfeature/asc", []string{"a", "b", "c"})
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?ids=a,b,c&limit=100&offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /history/devices/:duration
func TestDeviceHistory(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/history/devices/4h")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/devices?log=4h")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /history/gateways/:duration
func TestGatewayHistory(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/history/gateways/4h")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?log=4h")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset
func TestGatewayList(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset
func TestGatewayListDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?offset=0")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset
func TestGatewayListDefaultOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset
func TestSearchGateway(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset
func TestSearchGatewayDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset
func TestSearchGatewayDefaultOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset
func TestSearchGatewayDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGatewayAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&offset=0&sort=orderfeature.asc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGatewayDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&offset=0&sort=orderfeature.desc")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGateway(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGatewayDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?offset=0&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGatewayDefaultOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?limit=100&sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /list/gateways/:limit/:offset/:orderfeature/asc
func TestListGatewayDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/list/gateways/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature")
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrdered(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature&limit=100&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrderedAsc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature.asc&limit=100&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrderedDesc(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/desc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature.desc&limit=100&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrderedDefaultLimit(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature&offset=0&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrderedDefaultOffset(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature&limit=100&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}

// /search/gateways/:query/:limit/:offset/:orderfeature/asc
func TestSearchGatewayOrderedDefault(t *testing.T) {
	oldUrl, newUrl, mockOld, mockNew, stop := newMock()
	defer stop()
	err := testget(oldUrl + "/search/gateways/" + url.QueryEscape("search#text") + "/100/0/orderfeature/asc")
	if err != nil {
		t.Error(err)
		return
	}
	err = testget(newUrl + "/hubs?sort=orderfeature&search=" + url.QueryEscape("search#text"))
	if err != nil {
		t.Error(err)
		return
	}
	if !mockOld.Compare(mockNew) {
		t.Error("\n", mockOld, "\n\n", mockNew)
		return
	}
}
