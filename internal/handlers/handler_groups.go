package handlers

import (
	"database/sql"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReqDetails struct{
	Name		string	`json:"name"`
	Description	string	`json:"description"`
}


// userCreateGroup
// Creates a expense group in db
// Returns created group details 
func(cfg *apiConfig) userCreateGroup(c *gin.Context){
	var reqDetails ReqDetails

	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return
	}

	userID := tempID.(uuid.UUID)

	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return 
	}

	descriptionval := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}

	createdGroup, err := cfg.DB.CreateGroup(c, db.CreateGroupParams{
		Name: reqDetails.Name,
		Description: descriptionval,
		CreatedBy: userID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	err = cfg.DB.AddMember(c, db.AddMemberParams{
		UserID: userID,
		GroupID: createdGroup.ID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(201, createdGroup)
}


// getGroupByID
// Retrieves a group's details
// Specified by ID
func(cfg *apiConfig) getGroupByID(c *gin.Context){
	tempID := c.Param("group_id")
	groupID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, err)
		return 
	}

	group, err := cfg.DB.GetGroupByID(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(200, group)
}


// updateGroupDetails
// Updates a group details if authenticated and creator of group requests
func(cfg *apiConfig) updateGroupDetails(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InternalError, utils.RequestBodyError, nil)
		return
	}

	userID := tempID.(uuid.UUID)

	var reqDetails ReqDetails
	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	// get groupID from params
	// get group details for default values if field empty in json
	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)
	
	group, err := cfg.DB.GetGroupByID(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	if group.CreatedBy != userID{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return
	}

	if(reqDetails.Name == ""){
		reqDetails.Name = group.Name
	}
	if(reqDetails.Description == ""){
		reqDetails.Description = group.Description.String
	}

	descriptionVal := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}

	err = cfg.DB.UpdateGroup(c, db.UpdateGroupParams{
		Name: reqDetails.Name,
		Description: descriptionVal,
		ID: groupID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(204, utils.MessageObj("updated"))
}


// deleteGroup
// deletes a group created by user
func(cfg *apiConfig) deleteUserGroup(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return 
	}
	userID := tempID.(uuid.UUID)

	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	group, err := cfg.DB.GetGroupByID(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	if userID != group.CreatedBy{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return
	}

	err = cfg.DB.DeleteGroup(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(204, utils.MessageObj("deleted successfully"))
}


// addGroupMember
// Adds a member to a particular group
func(cfg *apiConfig) addGroupMember(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return 
	}
	creatorID := tempID.(uuid.UUID)

	tempUID := c.Param("user_id")
	userID := uuid.MustParse(tempUID)

	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	group, _ := cfg.DB.GetGroupByID(c, groupID)
	if group.CreatedBy != creatorID{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return 
	} 

	err := cfg.DB.AddMember(c, db.AddMemberParams{
		UserID: userID,
		GroupID: groupID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(201, utils.MessageObj("added successfully"))
}


// getGroupMembers
// Retrieves a particular group's members
func(cfg *apiConfig) getGroupMembers(c *gin.Context){
	tempID := c.Param("group_id")
	groupID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, err)
		return 
	}

	tempUID, _ := c.Get("userID")
	userID := tempUID.(uuid.UUID)

	_, err = cfg.DB.CheckMemeber(c, db.CheckMemeberParams{
		GroupID: groupID,
		UserID: userID,
	})
	if err == sql.ErrNoRows{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return
	}

	members, err := cfg.DB.GetGroupMembers(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}
	

	c.IndentedJSON(200, members)
}


// deleteGroupMember
func(cfg *apiConfig) deleteGroupMember(c *gin.Context){
	tempCID, _ := c.Get("userID")
	creatorID := tempCID.(uuid.UUID)

	tempID := c.Param("group_id")
	groupID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, err)
		return 
	}

	tempID = c.Param("user_id")
	userID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, err)
		return 
	}

	// Checks who sent the deletion request
	// Continues only if sent by group creator or the user deleting himself
	group, _ := cfg.DB.GetGroupByID(c, groupID)
	if group.CreatedBy != creatorID && creatorID != userID{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return 
	} 

	// If delete is for creator then request shoudl also be sent by creator
	if userID == group.CreatedBy && creatorID != group.CreatedBy{
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return 
	}


	err = cfg.DB.DeleteGroupMember(c, db.DeleteGroupMemberParams{
		UserID: userID,
		GroupID: groupID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(204, utils.MessageObj("deleted successfully"))
}


// getUserAllGroups
// Retrieves all the groups a user is a part of or has created
func(cfg *apiConfig) getUserAllGroups(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return 
	}
	userID := tempID.(uuid.UUID)

	groups, err := cfg.DB.GetUserAllGroups(c, userID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(200, groups)
}


// getUserGroups
// Retrives all the groups createdby user
func(cfg *apiConfig) getUserGroups(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return 
	}
	userID := tempID.(uuid.UUID)

	groups, err := cfg.DB.GetUserGroups(c, userID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(200, groups)
}


// deleteLoggedInUser
// deletes the user sending request itself.
// For group leave option
func(cfg *apiConfig) deleteLoggedInUser(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.MiddlewareError, nil)
		return
	}
	userID := tempID.(uuid.UUID)

	tempGID := c.Param("group_id")
	if tempGID == ""{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, nil)
		return
	}
	groupID, err := uuid.Parse(tempGID)
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.IDParseError, err)
		return 
	}

	err = cfg.DB.DeleteGroupMember(c, db.DeleteGroupMemberParams{
		UserID: userID,
		GroupID: groupID,
	})
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(204, utils.MessageObj("user deleted success"))

}