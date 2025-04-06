package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// invitationPage
// gin func to handle web page for group invitaion
func (cfg *apiConfig) invitationPage(c *gin.Context) {
	tempGID := c.Param("группа")
	groupID, err := uuid.Parse(tempGID)
	if err != nil {
		c.HTML(400, "invalid_page.html", nil)
		return
	}

	tempUID := c.Param("пользователь")
	userID, err := uuid.Parse(tempUID)
	if err != nil {
		c.HTML(400, "invalid_page.html", nil)
		return
	}

	group, err := cfg.DB.GetGroupByID(c, groupID)
	if err != nil {
		c.HTML(400, "invalid_page.html", nil)
		return
	}

	log.Print(group.ImageUrl.String)

	c.HTML(200, "group_invitation.html", gin.H{
		"GroupID":  group.ID,
		"Name":     group.Name,
		"ImageUrl": group.ImageUrl.String,
		"MemberID": userID,
	})
}

// acceptInvitaion
// adds user to group
func (cfg *apiConfig) acceptInvitaion(c *gin.Context) {
	tempGID := c.Param("группа")
	groupID, err := uuid.Parse(tempGID)
	if err != nil {
		c.HTML(400, "invalid_page.html", nil)
		return
	}

	tempUID := c.Param("пользователь")
	userID, err := uuid.Parse(tempUID)
	if err != nil {
		c.HTML(400, "invalid_page.html", nil)
		return
	}

	err = cfg.DB.AddMember(c, db.AddMemberParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(201, utils.MessageObj("success"))
}

// sendGroupInvite
// sends email with group invitaion link
func (cfg *apiConfig) sendGroupInvite(c *gin.Context) {
	var reqDetails struct {
		Email string `json:"email" binding:"required"`
	}

	err := c.BindJSON(&reqDetails)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	tempGID := c.Param("group_id")
	groupID, err := uuid.Parse(tempGID)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	// retrieve id of user sending request
	tempUID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 401, utils.InvalidError, utils.MiddlewareError, nil)
		return
	}
	userID := tempUID.(uuid.UUID)

	fromUser, err := cfg.DB.GetUserByID(c, userID)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.DatabaseError, err)
		return
	}

	// get user from email to send mail to
	toUser, err := cfg.DB.GetUserByEmail(c, reqDetails.Email)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.DatabaseError, err)
		return
	}

	godotenv.Load()
	addr := os.Getenv("API")
	if addr == "" {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	link := fmt.Sprintf("http://%v:8080/group-invite/%v/%v", addr, groupID, toUser.ID)
	body := fmt.Sprintf("Hey %v, %v is sending you a group invite click here to join: %v", toUser.Name, fromUser.Name, link)
	err = utils.SendEmail(reqDetails.Email, "Expense Group Invite", body)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, "unable to send email", err)
		return
	}

	c.IndentedJSON(200, utils.MessageObj("email sent"))
}
