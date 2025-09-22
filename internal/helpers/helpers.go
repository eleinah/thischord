package helpers

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/eleinah/thischord/internal/state"
)

func ListCommandOptions(options []*discordgo.ApplicationCommandOption) string {
	strOptions := "["
	for i, option := range options {
		if i < len(options)-1 {
			strOptions += fmt.Sprintf("%s,", option.Name)
		} else {
			strOptions += fmt.Sprintf("%s", option.Name)
		}
	}
	strOptions += "]"

	return strOptions
}

func JoinUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID

	inChannel, channelID := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.", true)
	}

	if state.VoiceConnected {
		interactionState.Reply("I am already in a voice channel.", true)
	} else {
		_, err := interactionState.Session.ChannelVoiceJoin(guildID, channelID, false, true)
		if err != nil {
			slog.Error("Failed to join voice channel:", "error", err.Error())
		}

		state.VoiceConnected = true
		interactionState.Reply("Joined voice channel.", false)
	}
}

func LeaveUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID

	inChannel, _ := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.Reply("You are not in a voice channel.", true)
		return
	}

	if !state.VoiceConnected {
		interactionState.Reply("I am not in a voice channel.", true)
	} else {
		_, err := interactionState.Session.ChannelVoiceJoin(guildID, "", false, true)
		if err != nil {
			slog.Error("Failed to leave voice channel:", "error", err.Error())
		}
		state.VoiceConnected = false
		interactionState.Reply("Left voice channel.", false)
	}
}

func IsUserInVoiceChannel(interactionState *state.InteractionState) (inChannel bool, channelID string) {
	guildID := interactionState.Interaction.GuildID
	guild, err := interactionState.Session.State.Guild(guildID)
	if err != nil {
		slog.Error("Failed to get guild:", "error", err.Error())
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == interactionState.Interaction.Member.User.ID {
			return true, vs.ChannelID
		}
	}
	return false, ""
}
