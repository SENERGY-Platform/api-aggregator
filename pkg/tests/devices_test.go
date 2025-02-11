package tests

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/api-aggregator/pkg"
	"github.com/SENERGY-Platform/api-aggregator/pkg/api"
	"github.com/SENERGY-Platform/api-aggregator/pkg/tests/environment"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDevicesEndpoint(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repoUrl, err := environment.New(ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	c := client.NewClient(repoUrl, nil)

	serverPort, err := getFreePortStr()
	if err != nil {
		t.Error(err)
		return
	}

	go api.Start(pkg.New(pkg.Config{
		ServerPort: serverPort,
		IotUrl:     repoUrl,
	}))

	devices := []models.Device{
		{
			Id:      "urn:ses:device:d1",
			Name:    "foo 1",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d2",
			Name:    "foo 2",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d3",
			Name:    "foo 3",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d4",
			Name:    "bar 1",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d5",
			Name:    "bar 2",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d6",
			Name:    "bar 3",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d7",
			Name:    "batz 1",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d8",
			Name:    "batz 2",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d9",
			Name:    "batz 3",
			OwnerId: userId,
		},

		{
			Id:      "urn:ses:device:d99",
			Name:    "plug Kühlschrank Backofen",
			OwnerId: userId,
		},
	}

	location := models.Location{
		Id:        "urn:ses:location:l1",
		Name:      "test",
		DeviceIds: []string{"urn:ses:device:d2", "urn:ses:device:d3", "urn:ses:device:d4"},
	}

	t.Run("create devices", createDevices(c, devices))
	t.Run("create locations", createLocations(c, []models.Location{location}))

	time.Sleep(1 * time.Second)

	t.Run(testDeviceQuery(``, serverPort, devices))
	t.Run(testDeviceQuery(`?limit=100&offset=0`, serverPort, devices))
	t.Run(testDeviceQuery(`?limit=2&offset=1`, serverPort, []models.Device{
		{
			Id:      "urn:ses:device:d5",
			Name:    "bar 2",
			OwnerId: userId,
		},
		{
			Id:      "urn:ses:device:d6",
			Name:    "bar 3",
			OwnerId: userId,
		},
	}))
}

func testDeviceQuery(query string, port string, expected []models.Device) (string, func(t *testing.T)) {
	return query, func(t *testing.T) {
		result := []models.Device{}
		err := pkg.GetJson(testjwt, "http://localhost:"+port+"/devices"+strings.ReplaceAll(query, ":", "%3A"), &result)
		if err != nil {
			t.Error(err)
			return
		}
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].Name < expected[j].Name
		})
		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Error("\n", string(resultJson), "\n", string(expectedJson))
		}
	}
}

func testDeviceQueryRaw(query string, port string, expected []models.Device) (string, func(t *testing.T)) {
	return query, func(t *testing.T) {
		result := []models.Device{}
		err := pkg.GetJson(testjwt, "http://localhost:"+port+"/devices"+query, &result)
		if err != nil {
			t.Error(err)
			return
		}
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].Name > expected[j].Name
		})
		sort.Slice(result, func(i, j int) bool {
			return result[i].Name > result[j].Name
		})
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Error(string(resultJson), "\n", string(expectedJson))
		}
	}
}

func createDevices(c client.Interface, devices []models.Device) func(t *testing.T) {
	return func(t *testing.T) {
		for _, d := range devices {
			_, err, _ := c.SetDevice(testjwt, d, client.DeviceUpdateOptions{})
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}

func createLocations(c client.Interface, locations []models.Location) func(t *testing.T) {
	return func(t *testing.T) {
		for _, l := range locations {
			_, err, _ := c.SetLocation(testjwt, l)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}

var testjwt = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzaUtabW9aUHpsMmRtQnBJdS1vSkY4ZVVUZHh4OUFIckVOcG5CcHM5SjYwIn0.eyJleHAiOjE2MTEyMTkwMzAsImlhdCI6MTYxMTIxNTQzMCwiYXV0aF90aW1lIjoxNjExMjE1NDI5LCJqdGkiOiJiZjY0NGI3Yy04YTZjLTQyYmMtOWNkNS0wODQ1NGU3ZmY1NDkiLCJpc3MiOiJodHRwOi8vZmdzZWl0c3JhbmNoZXIud2lmYS5pbnRlcm4udW5pLWxlaXB6aWcuZGU6ODA4Ny9hdXRoL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNjIxOWRjNDItYjhkMC00YjQyLTg1MWEtMWM1OTU2MTQ5OTQ0IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiZnJvbnRlbmQiLCJub25jZSI6IjkxZjdkYWViLTE2YzEtNDcwNy04YjNkLTFkY2NjOWI1NzMzMyIsInNlc3Npb25fc3RhdGUiOiI4ODAxMTBiOC1jYmM1LTQ2YWItYjRmYS1lM2Q4OGEwYjdhYzEiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGV2ZWxvcGVyIiwidW1hX2F1dGhvcml6YXRpb24iLCJ1c2VyIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQiLCJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRldmVsb3BlciIsInVtYV9hdXRob3JpemF0aW9uIiwidXNlciJdLCJuYW1lIjoiRGVtbyBVc2VyIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiZGVtby51c2VyIiwiZ2l2ZW5fbmFtZSI6IkRlbW8iLCJmYW1pbHlfbmFtZSI6IlVzZXIiLCJlbWFpbCI6IiJ9.H4YEIT3zJP67xjORJNiD6zr208SkcXi4LmN-sosZs5lPKLOJzEIxqn6gFphGRrcBfWHa5QZftFdpuw_yItvQSVMr5mNnsf8sNEFEsGu4zVfipnQdg3crY5TBK7dKiczhQarBPALFXP5Q2Q8uuuX2wAta76j3gl6X5RIbcwlNqKC5BG3LIoFbYVIeeKqhgNFEON5H530klJBzZ2pvLAXZxptQZUMydWTik7DrJrYSx_sPPCJtrz_d5UVT0ppkZ5h_MZGa4fJM0aVXW0hs7gxGEIQGpSY5-wma9EpP_C-mfY53jDOLn0etRfNEgjZo4116yLqamt-3MsY7_GB9fydkuw`

const userId = "6219dc42-b8d0-4b42-851a-1c5956149944"

func getFreePortStr() (string, error) {
	intPort, err := getFreePort()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(intPort), nil
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
