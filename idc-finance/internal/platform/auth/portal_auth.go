package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type portalLoginState struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DisplayName  string `json:"displayName"`
	CustomerID   int64  `json:"customerId"`
	CustomerName string `json:"customerName"`
}

var portalStateMu sync.Mutex

func RegisterPortalRoutes(router *gin.RouterGroup) {
	router.POST("/login", portalLogin)
}

func portalLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "登录参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	state := loadPortalState()
	if req.Username != state.Username || req.Password != state.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":      "AUTH_FAILED",
			"message":   "账号或密码错误",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(LoginResponse{
		Token:       PortalDemoToken,
		DisplayName: state.DisplayName,
		Roles:       []string{"portal-client"},
	}, middleware.GetRequestID(c)))
}

func ChangePortalPassword(currentPassword, newPassword string) error {
	newPassword = strings.TrimSpace(newPassword)
	if len(newPassword) < 8 {
		return errors.New("新密码至少需要 8 位")
	}

	portalStateMu.Lock()
	defer portalStateMu.Unlock()

	state := loadPortalStateLocked()
	if state.Password != currentPassword {
		return errors.New("当前密码不正确")
	}

	state.Password = newPassword
	return savePortalStateLocked(state)
}

func UpdatePortalIdentity(customerID int64, displayName, customerName string) error {
	portalStateMu.Lock()
	defer portalStateMu.Unlock()

	state := loadPortalStateLocked()
	if customerID > 0 {
		state.CustomerID = customerID
	}
	if value := strings.TrimSpace(displayName); value != "" {
		state.DisplayName = value
	}
	if value := strings.TrimSpace(customerName); value != "" {
		state.CustomerName = value
	}
	return savePortalStateLocked(state)
}

func currentPortalIdentity() (int64, string) {
	state := loadPortalState()
	return state.CustomerID, state.CustomerName
}

func loadPortalState() portalLoginState {
	portalStateMu.Lock()
	defer portalStateMu.Unlock()
	return loadPortalStateLocked()
}

func loadPortalStateLocked() portalLoginState {
	statePath := filepath.Join("data", "portal-auth-state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		state := defaultPortalState()
		_ = savePortalStateLocked(state)
		return state
	}

	var state portalLoginState
	if err := json.Unmarshal(data, &state); err != nil {
		state = defaultPortalState()
		_ = savePortalStateLocked(state)
		return state
	}

	if strings.TrimSpace(state.Username) == "" || strings.TrimSpace(state.Password) == "" {
		state = defaultPortalState()
		_ = savePortalStateLocked(state)
		return state
	}

	if strings.TrimSpace(state.DisplayName) == "" {
		state.DisplayName = "演示客户"
	}
	if strings.TrimSpace(state.CustomerName) == "" {
		state.CustomerName = state.DisplayName
	}
	if state.CustomerID <= 0 {
		state.CustomerID = 1
	}
	return state
}

func savePortalStateLocked(state portalLoginState) error {
	if err := os.MkdirAll("data", 0o755); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("data", "portal-auth-state.json"), payload, 0o644)
}

func defaultPortalState() portalLoginState {
	return portalLoginState{
		Username:     "portal",
		Password:     "Portal123!",
		DisplayName:  "演示客户",
		CustomerID:   1,
		CustomerName: "演示客户",
	}
}
