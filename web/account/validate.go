package account

//
//func ValidateHandler(c *gin.Context) {
//	serviceQuery, _ := c.GetQuery("service")
//	serviceTicketQuery, ok := c.GetQuery("ticket")
//	if !ok {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// First, we need check the Server Ticket and delete it.
//	stDataString := db.Redis.Get(serviceTicketQuery).Val()
//	db.Redis.Del(serviceTicketQuery)
//	// Service Ticket not found.
//	if stDataString == "" {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// Get Service Ticket data.
//	stData := strings.Split(stDataString, "|")
//	if len(stData) != 2 {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//	userIDStr := stData[0]
//	userID, err := strconv.Atoi(userIDStr)
//	if err != nil {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//	serviceIDStr := stData[1]
//	serviceID, err := strconv.Atoi(serviceIDStr)
//	if err != nil {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// Check service token.
//	serviceData := new(db.Service)
//	db.MySQL.Model(&db.Service{}).Where("secret = ?", serviceQuery).Find(&serviceData)
//	if serviceData.ID != uint(serviceID) {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// Check user existed.
//	userData := new(db.User)
//	db.MySQL.Model(&db.User{}).Where(&db.User{Model: gorm.Model{ID: uint(userID)}}).Find(&userData)
//	if userData.ID == 0 {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// Get the user service auth token.
//	auth := tool.GetServiceAuth(serviceData.ID, userData.ID)
//	if auth.ID == 0 {
//		tool.ValidateResponse(c, false, "")
//		return
//	}
//
//	// All done!
//	tool.ValidateResponse(c, true, userData.Name)
//}
