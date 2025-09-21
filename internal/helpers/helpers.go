package helpers

import (
	"log/slog"

	"github.com/eleinah/thischord/internal/state"
)

func JoinUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID

	inChannel, channelID := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.", true)
	}

	_, err := interactionState.Session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		slog.Error("Failed to join voice channel:", err.Error())
	}

	interactionState.Reply("Joined voice channel.", false)
}

func LeaveUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID

	inChannel, _ := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.", true)
		return
	}

	_, err := interactionState.Session.ChannelVoiceJoin(guildID, "", false, true)
	if err != nil {
		slog.Error("Failed to leave voice channel:", err.Error())
	}

	interactionState.Reply("Left voice channel.", false)
}

func IsUserInVoiceChannel(interactionState *state.InteractionState) (inChannel bool, channelID string) {
	guildID := interactionState.Interaction.GuildID
	guild, err := interactionState.Session.State.Guild(guildID)
	if err != nil {
		slog.Error("Failed to get guild:", err.Error())
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == interactionState.Interaction.Member.User.ID {
			return true, vs.ChannelID
		}
	}
	return false, ""
}

// TODO: Add a function to get whether the bot is in a voice channel.
