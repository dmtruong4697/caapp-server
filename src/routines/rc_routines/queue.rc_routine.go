package rcroutines

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	"time"
)

func PairUser() {
	for {
		var userA db_models.GeneralQueueUser
		var userB db_models.GeneralQueueUser

		var allUserInQueue []db_models.GeneralQueueUser
		database.DB.Order("join_at ASC").Find(&allUserInQueue)

		if len(allUserInQueue) < 2 {
			time.Sleep(1 * time.Second)
			continue
		}

		firstUserInQueue := allUserInQueue[0]
		for i := 1; i < len(allUserInQueue); i++ {
			if IsMatchingUser(firstUserInQueue, allUserInQueue[i]) {
				userB = allUserInQueue[i]
				userA = firstUserInQueue
				break
			}
		}

	}
}

func IsMatchingUser(userA db_models.GeneralQueueUser, userB db_models.GeneralQueueUser) bool {
	if (userA.TargetGender != userB.Gender) || (userB.TargetGender != userA.Gender) {
		return false
	}
	return true
}
