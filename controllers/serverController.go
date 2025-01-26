package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/nedik/spp-lobby/types"

	"github.com/gin-gonic/gin"
)

type TServersByUpdatedTime = *types.SafeTreeMap[int64, []types.Server]
type TServersByIPPort = *types.SafeMap[string, types.Server]

const SERVER_EXPIRY_TIME_IN_SECONDS = 5 * 60 // 5 minutes

type ServerController struct {
    serversByUpdatedTime TServersByUpdatedTime
    serversByIPPort      TServersByIPPort
}

func NewServerController() ServerController {
    return ServerController{
        serversByUpdatedTime: types.NewSafeTreeMap[int64, []types.Server](),
        serversByIPPort:      types.NewSafeMap[string, types.Server](),
    }
}

func treeToList(serversTree TServersByUpdatedTime) []types.Server {
    serversList := []types.Server{}

    serversTree.Mutex.Lock()
    defer serversTree.Mutex.Unlock()

    for it := serversTree.UnsafeTreeMap.Iterator(); it.Valid(); it.Next() {
        for _, currentServer := range it.Value() {
            serversList = append(serversList, currentServer)
        }
    }
    return serversList
}

func (self *ServerController) ListAllServers(c *gin.Context) {
    removeExpiredServers(self.serversByUpdatedTime, self.serversByIPPort)
    c.JSON(http.StatusOK, treeToList(self.serversByUpdatedTime))
}

func (self *ServerController) RegisterServer(c *gin.Context) {
    var registerServerInput types.RegisterServerInput

    if err := c.ShouldBindJSON(&registerServerInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    if !types.ValidateRegisterServerInput(registerServerInput) {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    incomingServerIP := c.ClientIP()
    incomingServerPort := registerServerInput.Port
    incomingServer := types.ConvertRegisterServerInputToServer(registerServerInput)
    incomingServer.IP = incomingServerIP
    incomingServerIPPort := convertToIPPort(incomingServerIP, incomingServerPort)
    incomingServer.UpdatedAt = time.Now().Unix()

    // Find and update duplicate if exists
    foundServer, serverFound := self.serversByIPPort.Get(incomingServerIPPort)
    if serverFound {
        self.updateInServersTree(foundServer, incomingServer)
        self.serversByIPPort.Set(incomingServerIPPort, incomingServer)

        c.JSON(http.StatusCreated, gin.H{})
        return
    }

    // If doesn't exist, then add a new one
    appendToServersTree(self.serversByUpdatedTime, incomingServer)
    self.serversByIPPort.Set(incomingServerIPPort, incomingServer)
    c.JSON(http.StatusCreated, gin.H{})
}

func (self *ServerController) GetSpecificServer(c *gin.Context) {
    removeExpiredServers(self.serversByUpdatedTime, self.serversByIPPort)
    ip := c.Param("ip")
    port, err := getPortFromParams(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    foundServer, err := self.findServer(ip, port)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusNotFound, foundServer)
}

func (self *ServerController) GetPlayersOfServer(c *gin.Context) {
    removeExpiredServers(self.serversByUpdatedTime, self.serversByIPPort)
    ip := c.Param("ip")
    port, err := getPortFromParams(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    foundServer, err := self.findServer(ip, port)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, foundServer.Players)
}

func getPortFromParams(c *gin.Context) (uint16, error) {
    port64, err := strconv.ParseUint(c.Param("port"), 10, 16)
    if err != nil {
        return 0, errors.New("Invalid port")
    }
    port := uint16(port64)

    return port, nil
}

func (self *ServerController) findServer(ip string, port uint16) (*types.Server, error) {
    serverIPPort := convertToIPPort(ip, port)
    candidateServer, serverFound := self.serversByIPPort.Get(serverIPPort)
    if serverFound {
        return &candidateServer, nil
    }

    return nil, errors.New("server not found")
}

func convertToIPPort(ip string, port uint16) string {
    return ip + ":" + strconv.FormatUint(uint64(port), 10)
}

func appendToServersTree(serversTree TServersByUpdatedTime, newServer types.Server) {
    serverListNow, found := serversTree.Get(newServer.UpdatedAt)
    if found {
        serverListNowUpdated := append(serverListNow, newServer)
        serversTree.Del(newServer.UpdatedAt)
        serversTree.Set(newServer.UpdatedAt, serverListNowUpdated)
    } else {
        serversTree.Set(newServer.UpdatedAt, []types.Server{newServer})
    }
}

func (self *ServerController) updateInServersTree(existingServer types.Server, updatedServer types.Server) {
    lastUpdatedAt := existingServer.UpdatedAt
    serversListAtTime, found := self.serversByUpdatedTime.Get(lastUpdatedAt)

    // Find server and remove it from the tree
    if found {
        // Find the server in the list
        var serverIndex *int = nil
        for index, currentServer := range serversListAtTime {
            if currentServer.IP == existingServer.IP && currentServer.Port == existingServer.Port {
                serverIndex = new(int)
                *serverIndex = index
                break
            }
        }

        if serverIndex != nil {
            // Remove server from the list
            lastIndex := len(serversListAtTime) - 1
            serversListAtTime[*serverIndex] = serversListAtTime[lastIndex]
            serversListAtTime = serversListAtTime[:lastIndex]

            // Remove current list and put back  on the tree without the server
            self.serversByUpdatedTime.Del(lastUpdatedAt)
            self.serversByUpdatedTime.Set(lastUpdatedAt, serversListAtTime)
        }
    }

    // Add the server again to the tree
    appendToServersTree(self.serversByUpdatedTime, updatedServer)
}

func removeExpiredServers(serversTree TServersByUpdatedTime, serversMap TServersByIPPort) {
    timeNow := time.Now().Unix()
    serversTreeKeysToDelete := []int64{}

    serversTree.Mutex.Lock()
    defer serversTree.Mutex.Unlock()

    for it := serversTree.UnsafeTreeMap.Iterator(); it.Valid(); it.Next() {
        if timeNow - it.Key() > SERVER_EXPIRY_TIME_IN_SECONDS {
            for _, currentServer := range it.Value() {
                serverIPPort := convertToIPPort(currentServer.IP, currentServer.Port)
                serversMap.Delete(serverIPPort)
            }
            serversTreeKeysToDelete = append(serversTreeKeysToDelete, it.Key())
        }
    }

    for _, serversTreeKeyToDelete := range serversTreeKeysToDelete {
        serversTree.UnsafeTreeMap.Del(serversTreeKeyToDelete)
    }
}
