package query

var CheckFriendChannelQuery = `
	SELECT channels.*
	FROM channels
	JOIN channel_members cm1 ON channels.id = cm1.channel_id
	JOIN channel_members cm2 ON channels.id = cm2.channel_id
	WHERE channels.type = 'friend'
	AND cm1.user_id = ?
	AND cm2.user_id = ?
	AND cm1.channel_id = cm2.channel_id
`

func GetCheckFriendChannelQuery() string {
	return CheckFriendChannelQuery
}
