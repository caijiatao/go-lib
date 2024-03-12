package chat

import (
	"time"
)

// UserService 提供用户相关的服务
type UserService struct {
	usersByID map[int]*User
}

// AddUser 添加用户
func (userService *UserService) AddUser(userID int, name string, passHash string) {
	// 实现添加用户的逻辑
}

// RemoveUser 移除用户
func (userService *UserService) RemoveUser(userID int) {
	// 实现移除用户的逻辑
}

// AddFriendRequest 添加好友请求
func (userService *UserService) AddFriendRequest(fromUserID, toUserID int) {
	// 实现添加好友请求的逻辑
}

// ApproveFriendRequest 批准好友请求
func (userService *UserService) ApproveFriendRequest(fromUserID, toUserID int) {
	// 实现批准好友请求的逻辑
}

// RejectFriendRequest 拒绝好友请求
func (userService *UserService) RejectFriendRequest(fromUserID, toUserID int) {
	// 实现拒绝好友请求的逻辑
}

// User 表示用户
type User struct {
	UserID                         int
	Name                           string
	PassHash                       string
	FriendsByID                    map[int]*User
	FriendIDsToPrivateChats        map[int]*PrivateChat
	GroupChatsByID                 map[int]*GroupChat
	ReceivedFriendRequestsByUserID map[int]*AddRequest
	SentFriendRequestsByUserID     map[int]*AddRequest
}

// MessageUser 向用户发送消息
func (user *User) MessageUser(friendID int, message string) {
	// 实现向用户发送消息的逻辑
}

// MessageGroup 向群组发送消息
func (user *User) MessageGroup(groupID int, message string) {
	// 实现向群组发送消息的逻辑
}

// SendFriendRequest 发送好友请求
func (user *User) SendFriendRequest(friendID int) {
	// 实现发送好友请求的逻辑
}

// ReceiveFriendRequest 接收好友请求
func (user *User) ReceiveFriendRequest(friendID int) {
	// 实现接收好友请求的逻辑
}

// ApproveFriendRequest 批准好友请求
func (user *User) ApproveFriendRequest(friendID int) {
	// 实现批准好友请求的逻辑
}

// RejectFriendRequest 拒绝好友请求
func (user *User) RejectFriendRequest(friendID int) {
	// 实现拒绝好友请求的逻辑
}

// Chat 表示聊天
type Chat struct {
	ChatID   int
	Users    []*User
	Messages []Message
}

// PrivateChat 表示私聊
type PrivateChat struct {
	Chat
}

// NewPrivateChat 创建私聊
func NewPrivateChat(firstUser, secondUser *User) *PrivateChat {
	privateChat := &PrivateChat{}
	privateChat.Users = append(privateChat.Users, firstUser, secondUser)
	return privateChat
}

// GroupChat 表示群聊
type GroupChat struct {
	Chat
}

// AddUser 添加用户到群聊
func (groupChat *GroupChat) AddUser(user *User) {
	// 实现添加用户到群聊的逻辑
}

// RemoveUser 从群聊移除用户
func (groupChat *GroupChat) RemoveUser(user *User) {
	// 实现从群聊移除用户的逻辑
}

// Message 表示消息
type Message struct {
	MessageID int
	Message   string
	Timestamp time.Time
}

// AddRequest 表示添加好友请求
type AddRequest struct {
	FromUserID    int
	ToUserID      int
	RequestStatus RequestStatus
	Timestamp     time.Time
}

// RequestStatus 表示请求状态
type RequestStatus int

const (
	// UNREAD 未读
	UNREAD RequestStatus = iota
	// READ 已读
	READ
	// ACCEPTED 已接受
	ACCEPTED
	// REJECTED 已拒绝
	REJECTED
)
