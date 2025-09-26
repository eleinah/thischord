package helpers

import (
	"errors"
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

func JoinUserVoiceChannel(interactionState *state.InteractionState) (vc *discordgo.VoiceConnection, err error) {
	guildID := interactionState.Interaction.GuildID

	inChannel, channelID := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.InvisibleReply("You are not in a voice channel.")
		return &discordgo.VoiceConnection{}, errors.New("user not in voice channel")
	}

	if state.VoiceConnected {
		interactionState.InvisibleReply("I am already in a voice channel.")
		return &discordgo.VoiceConnection{}, errors.New("bot already in voice channel")
	} else {
		vc, err := interactionState.Session.ChannelVoiceJoin(guildID, channelID, false, true)
		if err != nil {
			slog.Error("Failed to join voice channel:", "error", err.Error())
		}

		state.VoiceConnected = true
		interactionState.Reply("Joined voice channel.")
		return vc, nil
	}
}

func LeaveUserVoiceChannel(interactionState *state.InteractionState) {
	guildID := interactionState.Interaction.GuildID

	inChannel, _ := IsUserInVoiceChannel(interactionState)

	if !inChannel {
		interactionState.InvisibleReply("You are not in a voice channel.")
		return
	}

	if !state.VoiceConnected {
		interactionState.InvisibleReply("I am not in a voice channel.")
	} else {
		_, err := interactionState.Session.ChannelVoiceJoin(guildID, "", false, true)
		if err != nil {
			slog.Error("Failed to leave voice channel:", "error", err.Error())
		}
		state.VoiceConnected = false
		interactionState.Reply("Left voice channel.")
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
