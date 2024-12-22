package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nedik/spp-lobby/controllers"
	"github.com/nedik/spp-lobby/routes"
	"github.com/nedik/spp-lobby/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type InvalidInput struct {
	SomeValue string `json:"some_value"`
}

type TestEnvironment struct {
	HTTPRecorder *httptest.ResponseRecorder
	Context      *gin.Context
	Router       *gin.Engine
}

func setupRouter() TestEnvironment {
	var testEnvironment TestEnvironment
	testEnvironment.HTTPRecorder = httptest.NewRecorder()
	testEnvironment.Context, testEnvironment.Router = gin.CreateTestContext(testEnvironment.HTTPRecorder)

	serverController := controllers.NewServerController()
	serverRouteController := routes.NewServerRouteController(serverController)
	routerGroup := testEnvironment.Router.Group("/")

	serverRouteController.InitServerRoutes(routerGroup)
	return testEnvironment
}

func TestGetServersEmptyList(t *testing.T) {
	testEnvironment := setupRouter()
	getServersAndAssert(t, testEnvironment, []types.Server{})
}

func TestRegisterNewServer(t *testing.T) {
	testEnvironment := setupRouter()

	registeredServer := registerServerAndAssert(t, testEnvironment, 23073, "Test Server")

	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer})
}

func TestRegisterNewServerTwice(t *testing.T) {
	testEnvironment := setupRouter()

	registerServerAndAssert(t, testEnvironment, 23073, "Test Server")
	registeredServer := registerServerAndAssert(t, testEnvironment, 23073, "Test Server")

	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer})
}

func TestRegisterTwoNewServers(t *testing.T) {
	testEnvironment := setupRouter()

	registeredServer1 := registerServerAndAssert(t, testEnvironment, 23073, "Test Server")
	registeredServer2 := registerServerAndAssert(t, testEnvironment, 23074, "Test Server 2")

	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer1, registeredServer2})
}

func TestRegisterTwoNewServersTwice(t *testing.T) {
	testEnvironment := setupRouter()

	registeredServer1 := registerServerAndAssert(t, testEnvironment, 23073, "Test Server")
	registeredServer2 := registerServerAndAssert(t, testEnvironment, 23074, "Test Server 2")
	registerServerAndAssert(t, testEnvironment, 23073, "Test Server")
	registerServerAndAssert(t, testEnvironment, 23074, "Test Server 2")

	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer1, registeredServer2})
}

func TestRegisterServerAndThenUpdateIt(t *testing.T) {
	testEnvironment := setupRouter()

	registerServerInput := createRegisterServerInput(23073, "Test server")
	registeredServer := registerServerInputAndAssert(t, testEnvironment, registerServerInput)
	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer})

	advanced := !*registerServerInput.Advanced
	registerServerInput.Advanced = &advanced
	registeredServer = registerServerInputAndAssert(t, testEnvironment, registerServerInput)
	getServersAndAssert(t, testEnvironment, []types.Server{registeredServer})
}

func TestFieldsWithDefaultValuesOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()

	advanced := new(bool)
	antiCheatOn := new(bool)
	bonusFrequency := new(uint16)
	country := new(string)
	info := new(string)
	numBots := new(uint8)
	private := new(bool)
	realistic := new(bool)
	respawn := new(uint32)
	survival := new(bool)
	wm := new(bool)
	*advanced = false
	*antiCheatOn = false
	*bonusFrequency = 0
	*country = ""
	*info = ""
	*numBots = 0
	*private = false
	*realistic = false
	*respawn = 0
	*survival = false
	*wm = false

	registerServerInput := types.RegisterServerInput{
		Advanced:       advanced,
		AntiCheatOn:    antiCheatOn,
		BonusFrequency: bonusFrequency,
		Country:        country,
		CurrentMap:     "ctf_Ash",
		GameStyle:      "CTF",
		Info:           info,
		MaxPlayers:     10,
		Name:           "Name",
		NumBots:        numBots,
		OS:             "Linux",
		Players:        []string{},
		Port:           23073,
		Private:        private,
		Realistic:      realistic,
		Respawn:        respawn,
		Survival:       survival,
		Version:        "1.0",
		WM:             wm,
	}

	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusCreated, returnedCode)
}

func TestMissingFieldsOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()

	advanced := new(bool)
	antiCheatOn := new(bool)
	bonusFrequency := new(uint16)
	info := new(string)
	numBots := new(uint8)
	realistic := new(bool)
	respawn := new(uint32)
	survival := new(bool)
	wm := new(bool)
	*advanced = true
	*antiCheatOn = true
	*bonusFrequency = 10
	*info = "Test Server Info"
	*numBots = 1
	*realistic = true
	*respawn = 1
	*survival = true
	*wm = true

	newServerInputWithMissingFields := types.RegisterServerInput{
		Advanced:       advanced,
		AntiCheatOn:    antiCheatOn,
		BonusFrequency: bonusFrequency,
		CurrentMap:     "ctf_Ash",
		Info:           info,
		MaxPlayers:     32,
		NumBots:        numBots,
		OS:             "Linux",
		Port:           23073,
		Realistic:      realistic,
		Respawn:        respawn,
		Survival:       survival,
		Version:        "1.0",
		WM:             wm,
	}

	invalidInputJson, _ := json.Marshal(newServerInputWithMissingFields)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestInvalidInputOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()

	invalidInput := InvalidInput{
		SomeValue: "test",
	}

	invalidInputJson, _ := json.Marshal(invalidInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongCountryOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	country := new(string)
	*country = createLongString(types.MaxCountryLength + 1)
	registerServerInput.Country = country
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongMapNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.CurrentMap = createLongString(types.MaxMapSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongGameStyleOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.GameStyle = createLongString(types.MaxGameStyleSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongInfoOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	info := new(string)
	*info = createLongString(types.MaxInfoSize + 1)
	registerServerInput.Info = info
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, createLongString(types.MaxNameSize+1))
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongOSOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.OS = createLongString(types.MaxOSSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongPlayerNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Players[1] = createLongString(types.MaxPlayerNameSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongVersionOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Version = createLongString(types.MaxVersionSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func createLongString(size int) string {
	var stringBuilder strings.Builder
	stringBuilder.Grow(size)
	for i := 0; i < size; i++ {
		stringBuilder.WriteByte(0)
	}
	result := stringBuilder.String()
	return result
}

func registerServerAndAssert(t *testing.T, env TestEnvironment, port uint16, name string) types.Server {
	registerServerInput := createRegisterServerInput(port, name)
	serverJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(env, "/servers", serverJson)
	assert.Equal(t, http.StatusCreated, returnedCode)

	registeredServer := types.ConvertRegisterServerInputToServer(registerServerInput)
	return registeredServer
}

func registerServerInputAndAssert(t *testing.T, env TestEnvironment, registerServerInput types.RegisterServerInput) types.Server {
	serverJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(env, "/servers", serverJson)
	assert.Equal(t, http.StatusCreated, returnedCode)

	registeredServer := types.ConvertRegisterServerInputToServer(registerServerInput)
	return registeredServer
}

func createRegisterServerInput(port uint16, name string) types.RegisterServerInput {
	advanced := new(bool)
	antiCheatOn := new(bool)
	bonusFrequency := new(uint16)
	country := new(string)
	info := new(string)
	numBots := new(uint8)
	private := new(bool)
	realistic := new(bool)
	respawn := new(uint32)
	survival := new(bool)
	wm := new(bool)
	*advanced = true
	*antiCheatOn = true
	*bonusFrequency = 10
	*country = "PL"
	*info = "Test Server Info"
	*numBots = 1
	*private = true
	*realistic = true
	*respawn = 1
	*survival = true
	*wm = true

	registerServerInput := types.RegisterServerInput{
		Advanced:       advanced,
		AntiCheatOn:    antiCheatOn,
		BonusFrequency: bonusFrequency,
		Country:        country,
		CurrentMap:     "ctf_Ash",
		GameStyle:      "CTF",
		Info:           info,
		MaxPlayers:     32,
		Name:           name,
		NumBots:        numBots,
		OS:             "Linux",
		Players:        []string{"test_player_1", "another player"},
		Port:           port,
		Private:        private,
		Realistic:      realistic,
		Respawn:        respawn,
		Survival:       survival,
		Version:        "1.0",
		WM:             wm,
	}

	return registerServerInput
}

func getServersAndAssert(t *testing.T, env TestEnvironment, expected_servers []types.Server) {
	code, body := getServers(env)
	assert.Equal(t, http.StatusOK, code)
	expected_servers_json, _ := json.Marshal(expected_servers)
	assert.Equal(t, string(expected_servers_json), body)
}

func getServers(env TestEnvironment) (int, string) {
	env.Context.Request, _ = http.NewRequest("GET", "/servers", nil)
	httpRecorder := httptest.NewRecorder()
	env.Router.ServeHTTP(httpRecorder, env.Context.Request)
	return httpRecorder.Code, httpRecorder.Body.String()
}

func sendJsonToPostEndpoint(env TestEnvironment, endpoint string, dataJson []byte) (int, string) {
	env.Context.Request, _ = http.NewRequest("POST", endpoint, strings.NewReader(string(dataJson)))
	httpRecorder := httptest.NewRecorder()
	env.Router.ServeHTTP(httpRecorder, env.Context.Request)
	return httpRecorder.Code, httpRecorder.Body.String()
}
