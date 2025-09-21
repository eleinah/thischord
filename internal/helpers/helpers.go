package helpers

import (
	"log/slog"

	"github.com/eleinah/thischord/internal/state"
)

func JoinUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID
	guild, err := interactionState.Session.State.Guild(guildID)
	if err != nil {
		slog.Error("Failed to get guild:", err.Error())
	}

	userID := interactionState.Interaction.Member.User.ID
	inChannel := false
	var channelID string

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			inChannel = true
			channelID = vs.ChannelID
			break
		}
	}

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.")
	}

	_, err = interactionState.Session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		slog.Error("Failed to join voice channel:", err.Error())
	}

	interactionState.Reply("Joined voice channel.")
}

func LeaveUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID
	guild, err := interactionState.Session.State.Guild(guildID)
	if err != nil {
		slog.Error("Failed to get guild:", err.Error())
	}

	userID := interactionState.Interaction.Member.User.ID
	inChannel := false

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			inChannel = true
			break
		}
	}

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.")
	}

	_, err = interactionState.Session.ChannelVoiceJoin(guildID, "", false, true)
	if err != nil {
		slog.Error("Failed to leave voice channel:", err.Error())
	}

	interactionState.Reply("Left voice channel.")
}
