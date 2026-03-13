package handlers

import (
	"encoding/json"
	"llm-ops/models"
	"llm-ops/services"
	"llm-ops/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	userService *services.UserService
)

func InitUserService(us *services.UserService) {
	userService = us
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("User attempting to list users - Username: %s, Role: %s", claims.Username, claims.Role)

	if claims.Role != "admin" {
		log.Printf("Unauthorized access attempt by user %s with role %s", claims.Username, claims.Role)
		utils.ErrorResponse(w, http.StatusForbidden, "需要管理员权限")
		return
	}

	users := userService.GetAllUsers()
	log.Printf("Retrieved users: %+v", users)

	utils.SuccessResponse(w, map[string]interface{}{
		"data": users,
	}, "")
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		log.Printf("No user found in context")
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	log.Printf("Add user attempt by: %s (role: %s)", claims.Username, claims.Role)

	if claims.Role != "admin" {
		log.Printf("Unauthorized access attempt by user %s with role %s", claims.Username, claims.Role)
		utils.ErrorResponse(w, http.StatusForbidden, "需要管理员权限")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.SetPassword(user.Password); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := userService.AddUser(&user); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, nil, "用户添加成功")
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := utils.GetUserFromContext(r.Context())
	if !ok || currentUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "需要管理员权限")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if !models.IsValidRole(req.Role) {
		utils.ErrorResponse(w, http.StatusBadRequest, "无效的角色")
		return
	}

	if username == currentUser.Username {
		utils.ErrorResponse(w, http.StatusForbidden, "不能修改自己的角色")
		return
	}

	err := userService.UpdateUserRole(username, req.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, nil, "角色更新成功")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := utils.GetUserFromContext(r.Context())
	if !ok || currentUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "需要管理员权限")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if username == currentUser.Username {
		utils.ErrorResponse(w, http.StatusForbidden, "不能删除自己的账号")
		return
	}

	err := userService.DeleteUser(username)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, nil, "用户删除成功")
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "未找到用户信息")
		return
	}

	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	err := userService.UpdatePassword(claims.UserID, req.NewPassword)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, nil, "密码修改成功")
}
