package helper

import (
	"gitlab.com/prestrafe/prestrafe-bot/gsiclient"
	"gitlab.com/prestrafe/prestrafe-bot/smclient"
)

// Check for data received from sourcemod and gsi backends and see if they match
func CompareData(smData *smclient.FullPlayerInfo, gsiData *gsiclient.GameState) bool {

	/* Team timeout check is to make sure that the server is doesn't send inaccurate player data.
	These values will only match if the player is in the server, and change at frequent interval
	to make sure the player is always in the server. */

	// This can get out of sync if the server/client was slow to update, so we will also compare to previous data.
	timeoutDelayed := (smData.TimeoutsCTPrev == *gsiData.Map.TeamCT.Timeouts) && (smData.TimeoutsTPrev == *gsiData.Map.TeamT.Timeouts)

	timeoutAhead := gsiData.PreviousState != nil
	// Make sure we don't have memory access violation
	if timeoutAhead {
		timeoutAhead = gsiData.PreviousState.Map != nil
	}

	if timeoutAhead {
		timeoutAhead = (gsiData.PreviousState.Map.TeamCT != nil) && (gsiData.PreviousState.Map.TeamT != nil)
	}

	if timeoutAhead {
		timeoutAhead = (gsiData.PreviousState.Map.TeamCT.Timeouts != nil) && (gsiData.PreviousState.Map.TeamT.Timeouts != nil)
	}

	if timeoutAhead {
		timeoutAhead = (smData.TimeoutsCT == *gsiData.PreviousState.Map.TeamCT.Timeouts) && (smData.TimeoutsT == *gsiData.PreviousState.Map.TeamT.Timeouts)
	}

	timeoutInSync := (smData.TimeoutsCT == *gsiData.Map.TeamCT.Timeouts) && (smData.TimeoutsT == *gsiData.Map.TeamT.Timeouts)

	timeoutCheck := timeoutDelayed || timeoutInSync || timeoutAhead

	/* SM exclusive commands don't return anything if the player is spectating for now, because the three steamids don't match.
	This is because currently GSI commands will return data about the spectated player and the SM backend might not have that data.
	Therefore, it will be confusing if SM commands return the information of the original player.
	However, it is likely possible to obtain and send the data of the spectated player instead. */

	steamIDCheck := (smData.SteamId == gsiData.Player.SteamId) && (smData.SteamId == gsiData.Provider.SteamId)

	// Server update is sent every 2 second. GSI update is sent every 2.5 seconds. A 6 seconds gap should be sufficient in case of packet loss.
	timestampCheck := (smData.TimeStamp <= int(gsiData.Provider.Timestamp)+3) && (smData.TimeStamp >= int(gsiData.Provider.Timestamp)-3)

	return timeoutCheck && steamIDCheck && timestampCheck
}
