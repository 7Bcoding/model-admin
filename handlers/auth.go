package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"llm-ops/models"
	"llm-ops/utils"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var cert = `-----BEGIN CERTIFICATE-----
MIIE2TCCAsGgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCYxDjAMBgNVBAoTBWFk
bWluMRQwEgYDVQQDDAtjZXJ0X2hjb3NqdjAeFw0yMjA3MDMwMzEzMjdaFw00MjA3
MDMwMzEzMjdaMCYxDjAMBgNVBAoTBWFkbWluMRQwEgYDVQQDDAtjZXJ0X2hjb3Nq
djCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAOQg1mnj+ORVR/TbjqPP
eqmwBtOD0hLjlO3mNpt9zgq7zLIuFzt6FNhxyJb8XD8xca9Pu/djjT38m+DeZKlw
iGlo7M2k1fh6tC0DLi6vstVHpgnDCiV8NQHedEOygdc17MMiGqDvVr9VrbX31KAO
0vrtboRWwuRCBXzkYMblh3lWf/D4MUhRQTQyNKlDd058VZ54J/zSRgrcqrK0Z4Ll
KPT6ynRc6M0lGYKgFB18N2pLQh97pFiHQTg5JkvnmSKWJ2Xtm+oPhKi3VPlzyS9m
er3S39Pkati5T0Pq3lorjUQHcXUMzXQYt1jOMplsPGWp1ULvncpYGXPHyCuEcR4u
eVlNKmMLKeG71n8LESx0ntepUrg2K79B//MjtA/aV/60tZNqnze18rUAHK6FDt6q
RvKQkMLqYmK5ynj5IHXCFzlSZ89cQYM3dMRWwJuxwHZJc7Ssah916QBfu6ilCuV6
NV43QCHjeIWHVslvToSojtaBlcIwB0I50l5OLUf1kmNZm+LXPqw2P4yw47KnlDOh
BVTAEBFmEM+iGcQmwIXHICq5w48Biix6QnJjqi43RDxRWJ/xlYnmtTEtIm+SxGf/
60UkguQNysnlYtoQWJuDuTi7/ykZWm93Hpl5tQjiX3s6Xxe2ckocBU8Rl/VwRdWA
o1RE5hopdWgq25ujQTraKAnPAgMBAAGjEDAOMAwGA1UdEwEB/wQCMAAwDQYJKoZI
hvcNAQELBQADggIBAB+ZqZ5wJGAy6gwRbr4Op5lEb1xM+Tgvp/PePzZfZWQ4yjcE
odcTrJ0xaXvgfcypL3sdqZkW5FtDh6mWV3USuhyDXhYsNMufND7HC1+gWQpd4+NI
24PZMQ74Ho38k3qxs+7rXGqtHlAsNmKWNu4XyfrdIuBRcx0cFadyIQM3PBDF8jHY
3KzKRpzr12UI1mTKn1K1BO0ptNyATqVd8sqOnkf8iKrbhB8F78p5Tc6JuhcECz1O
bdBfqDJdNMFATL8IvLlPOMQtT3/+TmUuCheBpQXkKtuGxhdldLKF+GTK3uW9wbZE
Pt6Aa2i9CHlsJNd+kYg0ICswk9JZ7RtWUI0vZ7Ao4dGHb6BrDqs9Wm6lfJcASW/W
ZkJVviXJqVV3I69jBFwOv9P7OD9o/TumasVxvw0Kt68s9UIT6ar1LbMpfiCzPlAV
zwa8YopO4k6z4BCUxnMcNNYR9N4w1Kc0H5dCZY8PSTODIvW8dQWF73qelRDprP89
Cp4M0g9yWFk4xpg/FUlpR7hSIj1iDb7ziuuY/qKncKjV0W7OhawHaynBSJkn8uov
/1iG06quMJlZ3asqg4YAmjnrrFqJJXJhxxh3yQ+oxOizfQT+XY1dizoDADzHC+ah
RpIVUCWOHx7+WTJyxiZPqjA7yvBRqWc5eCQsQJ2g+rc3+xHARgzgyJ/SGbcZ
-----END CERTIFICATE-----`

func init() {
	casdoorsdk.InitConfig("https://iam.beta.cloud", "f3c9bbc3ec62ed982b4a", "52ac3134b3ba16b3361754972cd870db51caa77c", cert, "beta", "faas-ops")
}

type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`
	// contains filtered or unexported fields

}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Handling Login Request ===")

	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	log.Printf("Login attempt for user: %s", loginReq.Username)

	// 获取用户信息
	user, err := userService.GetUserByUsername(loginReq.Username)
	if err != nil {
		log.Printf("User not found: %v", err)
		utils.ErrorResponse(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 验证密码
	if !user.CheckPassword(loginReq.Password) {
		log.Printf("Invalid password for user: %s", loginReq.Username)
		utils.ErrorResponse(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(int(user.ID), user.Username, string(user.Role))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 返回用户信息和 token
	utils.SuccessResponse(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}, "Login successful")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// 如果是 OPTIONS 请求，直接返回
	if r.Method == "OPTIONS" {
		return
	}

	// 返回成功响应
	utils.SuccessResponse(w, nil, "Logged out successfully")
}

func HandleFeishuLogin(w http.ResponseWriter, r *http.Request) {
	// get code and state from json payload
	var loginReq struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	code := loginReq.Code
	if code == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "code is required")
		return
	}

	state := loginReq.State
	if state == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "state is required")
		return
	}

	client := casdoorsdk.NewClient("https://iam.beta.cloud", "f3c9bbc3ec62ed982b4a", "52ac3134b3ba16b3361754972cd870db51caa77c", cert, "beta", "faas-ops")

	log.Printf("Feishu login with code: %s, state: %s", code, state)

	token, err := client.GetOAuthToken(code, state)
	if err != nil {
		log.Printf("Failed to get OAuth token: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get OAuth token")
		return
	}

	// 获取用户信息
	claims, err := client.ParseJwtToken(token.AccessToken)
	if err != nil {
		log.Printf("Failed to parse JWT token: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to parse JWT token")
		return
	}

	log.Printf("Feishu user: %s, %s", claims.User.Name, claims.User.Email)

	exists, err := userService.UserExists(claims.User.Name)
	if err != nil {
		log.Printf("Failed to check if user exists: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	var opsUser *models.User
	if !exists {
		log.Printf("User %s does not exist, creating user", claims.User.Name)

		utils.ErrorResponse(w, http.StatusInternalServerError, "user not found")
		return

		// opsUser = &models.User{
		// 	Username:    claims.User.Name,
		// 	AccountName: claims.User.Name,
		// 	Role:        string(models.RoleUser),
		// }
		// err = userService.AddUser(opsUser)
		// if err != nil {
		// 	log.Printf("Failed to create user: %v", err)
		// 	utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		// 	return
		// }
	}

	opsUser, err = userService.GetUserByUsername(claims.User.Name)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	// 生成 JWT token
	opsToken, err := utils.GenerateToken(int(opsUser.ID), opsUser.Username, string(opsUser.Role))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 返回用户信息和 token
	utils.SuccessResponse(w, map[string]interface{}{
		"token": opsToken,
		"user": map[string]interface{}{
			"id":       opsUser.ID,
			"username": opsUser.Username,
			"role":     opsUser.Role,
		},
	}, "Login successful")

	log.Printf("Feishu login successful")
}
