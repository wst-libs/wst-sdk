package im

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/wst-libs/wst-sdk/errors"
	"github.com/wst-libs/wst-sdk/im/sdk"
	"github.com/wst-libs/wst-sdk/utils"
)

func SetOutPutHeader(ctx *context.Context) {
	ctx.Output.Header("Connection", "close")
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.Header("Server", utils.Version)
}

func ParseInputHander(ctx *context.Context) error {
	if !strings.Contains(ctx.Input.Header("Content-Type"), "application/json") {
		log.Println("not application/json.")
		// return error.Error()("Content-Type not application/json")
		// return errors.New("Content-Type is application/json")
	}
	return nil
}

func checkRegisterCommon(info RegisteredUsers) (bool, error) {
	if len(info.Id) > 32 || len(info.Id) <= 0 {
		return false, errors.New("user id too long or short")
	}

	if len(info.Name) > 32 || len(info.Name) <= 0 {
		return false, errors.New("user name too long or short")
	}

	if len(info.Portrait) > 32 || len(info.Portrait) <= 0 {
		return false, errors.New("user portrait too long or short")
	}
	return true, nil
}

// Processing user registration
func RegisterUsersHandler(ctx *context.Context) []byte {
	// Parse request body to json
	var request RequestRegisteredUsers
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	// Check request common param
	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	// Parse request uid
	// uid := ctx.Input.Param(":uid")
	// err = errors.VerifyUid(uid)
	// if err != nil {
	// 	ret := utils.ResponseCommon{
	// 		Code:    errors.UidErr,
	// 		Message: err.Error(),
	// 		Version: request.RequestCommon.Version,
	// 		SeqNum:  request.RequestCommon.SeqNum,
	// 		From:    request.RequestCommon.From,
	// 		To:      request.RequestCommon.To,
	// 		Type:    request.RequestCommon.Type,
	// 		Number:  request.RequestCommon.Number,
	// 	}
	// 	out, _ := json.Marshal(ret)
	// 	return out
	// }
	_, err = checkRegisterCommon(request.RegisteredUsers)
	if err != nil {
		return errors.ImplementErr(errors.UserInfoErr, request.RequestCommon, err.Error())
	}

	s := sdk.NewRCServer()
	token, err := s.GetTokenFromUser(request.RegisteredUsers.Id, request.RegisteredUsers.Name, request.RegisteredUsers.Portrait)
	if err != nil {
		return errors.ImplementErr(errors.UserTokenErr, request.RequestCommon, err.Error())
	}

	v := ResponseRegisteredUsers{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  request.RequestCommon.SeqNum,
			From:    request.RequestCommon.From,
			To:      request.RequestCommon.To,
			Type:    request.RequestCommon.Type,
			Number:  request.RequestCommon.Number,

			Code:    0,
			Message: "success",
		},
		utils.TOKEN{
			Token: token,
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing create session
func CreateSessionHandler(ctx *context.Context) []byte {
	// Parse request body to json
	var request RequestCreateSession
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	// Check request common param
	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	// Get session id by uri
	sid := ctx.Input.Param(":sid")
	err = errors.VerifyUid(sid)
	if err != nil {
		ret := utils.ResponseCommon{
			Code:    errors.UidErr,
			Message: err.Error(),
			Version: request.RequestCommon.Version,
			SeqNum:  request.RequestCommon.SeqNum,
			From:    request.RequestCommon.From,
			To:      request.RequestCommon.To,
			Type:    request.RequestCommon.Type,
			Number:  request.RequestCommon.Number,
		}
		out, _ := json.Marshal(ret)
		return out
	}

	s := sdk.NewRCServer()
	err = s.CreateChatRoom(sid, request.Name)
	if err != nil {
		return errors.ImplementErr(errors.CreateSessionErr, request.RequestCommon, err.Error())
	}

	v := ResponseCreateSession{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  request.RequestCommon.SeqNum,
			From:    request.RequestCommon.From,
			To:      request.RequestCommon.To,
			Type:    request.RequestCommon.Type,
			Number:  request.RequestCommon.Number,

			Message: "success",
			Code:    0,
		},
		utils.ID{
			Id: sid,
		},
	}
	out, err := json.Marshal(v)
	if err != nil {
	}
	return out
}

// Processing delete session
func DeleteSessionHandler(ctx *context.Context) []byte {
	// Parse request body to json
	var request RequestDelSession
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}
	// Check request common params
	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}
	// Get session id by request params
	ids := make([]string, request.Size)
	for i, v := range request.List {
		ids[i] = v.Id
	}
	s := sdk.NewRCServer()
	err = s.DeleteChatRoom(ids)
	if err != nil {
		return errors.ImplementErr(errors.DeleteSessionErr, request.RequestCommon, err.Error())
	}

	v := ResponseDelSession{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  request.RequestCommon.SeqNum,
			From:    request.RequestCommon.From,
			To:      request.RequestCommon.To,
			Type:    request.RequestCommon.Type,
			Number:  request.RequestCommon.Number,

			Message: "success",
			Code:    0,
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing delete session by session id
func DeleteSessionByIDHandler(ctx *context.Context) []byte {
	// Parse request body to json
	var request RequestDelSession
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}
	// Check request common
	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	sid := ctx.Input.Param(":sid")

	s := sdk.NewRCServer()
	err = s.DeleteChatRoom([]string{sid})
	if err != nil {
		return errors.ImplementErr(errors.DeleteSessionErr, request.RequestCommon, err.Error())
	}

	v := utils.ResponseCommon{
		Version: utils.Version,
		SeqNum:  request.RequestCommon.SeqNum,
		From:    request.RequestCommon.From,
		To:      request.RequestCommon.To,
		Type:    request.RequestCommon.Type,
		Number:  request.RequestCommon.Number,

		Message: "success",
		Code:    0,
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing delete session by user id (deprecated)
func DeleteSessionByUIDHandler(ctx *context.Context) []byte {
	var request RequestDelSession
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	// s := sdk.NewRCServer()
	v := ResponseDelSession{
		utils.ResponseCommon{
			Version: "V1.0",
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing get session (deprecated)
func GetSessionHandler(ctx *context.Context) []byte {

	v := ResponseGetSession{
		utils.ResponseCommon{
			Version: "V1.0",
		},
		GetSession{},
	}
	// s := sdk.NewRCServer()
	// s.GetAllChatRoom()
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing get session infomation by session id
func GetSessionByIDHandler(ctx *context.Context) []byte {

	sid := ctx.Input.Param(":sid")
	s := sdk.NewRCServer()
	info, err := s.GetChatRoomById(sid)
	if err != nil {
		return errors.ImplementErr(errors.DeleteSessionErr, utils.RequestCommon{}, err.Error())
	}

	v := ResponseGetSession{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  1,
			From:    "",
			To:      "",
			Type:    "",
			Number:  "",
			Code:    0,
		},
		GetSession{
			List: []utils.RoomInfo{
				utils.RoomInfo{
					Id:         info.Id,
					Name:       info.Name,
					CreateTime: info.Time,
				},
			},
			Size: 1,
		},
	}

	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing get session by user id
func GetSessionByUIDHandler(ctx *context.Context) []byte {
	v := ResponseGetSession{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  1,
			From:    "",
			To:      "",
			Type:    "",
			Number:  "",
			Code:    0,
		},
		GetSession{},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

// Processing get users by session id
func GetUsersBySessionIDHandler(ctx *context.Context) []byte {

	uid := ctx.Input.Param(":sid")

	s := sdk.NewRCServer()
	users, err := s.GetUsersByRoomId(uid)
	if err != nil {
		return errors.ImplementErr(errors.DeleteSessionErr, utils.RequestCommon{}, err.Error())
	}

	var userIds []utils.ID
	for i, v := range users {
		userIds[i].Id = v
	}
	v := ResponseGetSessionUsers{
		utils.ResponseCommon{
			Version: utils.Version,
			SeqNum:  1,
			From:    "",
			To:      "",
			Type:    "",
			Number:  "",
			Code:    0,
		},
		GetSessionUsers{
			List: userIds,
			Size: len(users),
		},
	}

	out, err := json.Marshal(v)
	if err != nil {
		log.Println("GetUsersBySessionId Error: ", err.Error())
	}
	return out
}

func PutSessionByUIDHandler(ctx *context.Context) []byte {
	var request RequestJoinSession
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	uid := ctx.Input.Param(":uid")
	sid := ctx.Input.Param(":sid")
	s := sdk.NewRCServer()
	s.JoinRoomByUserId(uid, sid)

	v := ResponseJoinSession{
		utils.ResponseCommon{
			Version: "V1.0",
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

func PostMessageToUserByIDHandler(ctx *context.Context) []byte {
	var request RequestSendMessage
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	sourceId := request.Uid
	targetId := ctx.Input.Param(":uid")
	content := request.Content

	log.Println("source id: ", sourceId)
	log.Println("target id: ", targetId)
	log.Println("content is: ", content)

	s := sdk.NewRCServer()
	s.SendMsgUserToUsers(sourceId, targetId, content)

	v := ResponseSendMessage{
		utils.ResponseCommon{
			Version: utils.Version,
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

func PostMessageToSessionByIDHandler(ctx *context.Context) []byte {
	var request RequestSendMessage
	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	com := errors.IsCommonErr(request.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	sourceId := request.Uid
	targetId := ctx.Input.Param(":sid")
	content := request.Content

	s := sdk.NewRCServer()
	s.SendMsgUserToSession(sourceId, targetId, content)

	v := ResponseSendMessage{
		utils.ResponseCommon{
			Version: utils.Version,
		},
	}
	out, err := json.Marshal(v)
	if err != nil {

	}
	return out
}

func ReqMsgHandler(ctx *context.Context) []byte {
	var req ReqMsg
	err := json.Unmarshal(ctx.Input.RequestBody, &req)
	if err != nil {
		log.Println("Error: ", err.Error())
		return errors.ParseJsonFailed()
	}

	com := errors.IsCommonErr(req.RequestCommon)
	if com.Code != 0 {
		outerr, _ := json.Marshal(com)
		return outerr
	}

	s := sdk.NewRCServer()
	switch req.Data.Type {
	case "private":
		s.SendMsgPrivate(req.Data.SourceId, req.Data.TargetId, req.Data.Content)
	case "group":
		s.SendMsgGroup(req.Data.SourceId, req.Data.TargetId, req.Data.Content)
	case "chatroom":
		s.SendMsgRoom(req.Data.SourceId, req.Data.TargetId, req.Data.Content)
	default:
		return errors.ReqTypeErr(req.RequestCommon)
	}

	v := utils.ResponseCommon{
		Version: utils.Version,
		SeqNum:  req.SeqNum,
		From:    req.To,
		To:      req.From,
		Type:    req.Type,
		Number:  req.Number,
		Message: "success",
		Code:    0,
	}
	out, _ := json.Marshal(v)
	return out
}
