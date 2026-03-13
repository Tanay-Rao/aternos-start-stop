package aternos_discord_bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"github.com/sleeyax/aternos-discord-bot/database"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"github.com/sleeyax/aternos-discord-bot/message"
	"strings"
)

const (
	limitedCommandRoleID = "1480411247810842634"
	fullAccessRoleID     = "1480411813534502993"
)

var limitedRoleAllowedCommands = map[string]bool{
	StartCommand:   true,
	StatusCommand:  true,
	PlayersCommand: true,
	PingCommand:    true,
	InfoCommand:    true,
}

func hasAnyRole(member *discordgo.Member, roleIDs ...string) bool {
	if member == nil {
		return false
	}

	for _, memberRoleID := range member.Roles {
		for _, allowedRoleID := range roleIDs {
			if memberRoleID == allowedRoleID {
				return true
			}
		}
	}

	return false
}

// handleCommands responds to incoming interactive commands on discord.
func (ab *Bot) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()

	// wrap functions around our utilities to make life easier
	sendText := func(content string) {
		respondWithText(s, i, content)
	}
	sendHiddenText := func(content string) {
		respondWithHiddenText(s, i, content)
	}
	sendErrorText := func(content string, err error) {
		respondWithError(s, i, content, err)
	}

	hasFullAccess := hasAnyRole(i.Member, fullAccessRoleID)
	hasLimitedAccess := hasAnyRole(i.Member, limitedCommandRoleID)

	if !hasFullAccess {
		if !hasLimitedAccess || !limitedRoleAllowedCommands[command.Name] {
			sendHiddenText(message.FormatWarning("You don't have permission to use this command."))
			return
		}
	}

	switch command.Name {
	case HelpCommand:
		sendHiddenText(message.FormatDefault(faq))
	case PingCommand:
		sendHiddenText(message.FormatDefault("Pong!"))
	case ConfigureCommand:
		options := optionsToMap(command.Options)

		err := ab.Database.UpdateServerSettings(&models.ServerSettings{
			GuildID:       i.GuildID,
			SessionCookie: options[SessionOption].StringValue(),
			ServerCookie:  options[ServerOption].StringValue(),
		})
		if err != nil {
			sendErrorText("Failed to save configuration.", err)
			break
		}

		sendText(message.FormatSuccess("Configuration changed successfully."))
	case StatusCommand:
		fallthrough
	case InfoCommand:
		fallthrough
	case PlayersCommand:
		fallthrough
commands.go
commands.go
+6
-7

package aternos_discord_bot

import "github.com/bwmarrin/discordgo"

const (
	HelpCommand      = "help"
	PingCommand      = "ping"
	ConfigureCommand = "configure"
	StartCommand     = "start"
	StopCommand      = "stop"
	StatusCommand    = "status"
	InfoCommand      = "info"
	PlayersCommand   = "players"
	SessionOption    = "session"
	ServerOption     = "server"
)

var (
	adminPermissions int64 = discordgo.PermissionManageServer
	userPermissions  int64 = discordgo.PermissionUseSlashCommands
	dmPermission           = false
	userPermissions int64 = discordgo.PermissionUseSlashCommands
	dmPermission          = false
)

// List of available discord commands.
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        ConfigureCommand,
		Description: "Save configuration settings",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         SessionOption,
				Description:  "Set the ATERNOS_SESSION cookie value",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
			{
				Name:         ServerOption,
				Description:  "Set the ATERNOS_SERVER cookie value",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
			},
		},
		DefaultMemberPermissions: &adminPermissions,
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     StartCommand,
		Description:              "Start the minecraft server",
		DefaultMemberPermissions: &adminPermissions,
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     StopCommand,
		Description:              "Stop the minecraft server",
		DefaultMemberPermissions: &adminPermissions,
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     PingCommand,
		Description:              "Check if the discord bot is still alive",
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     StatusCommand,
		Description:              "Get the minecraft server status",
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     InfoCommand,
		Description:              "Get detailed information about the minecraft server status",
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     PlayersCommand,
		Description:              "List active players",
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
	{
		Name:                     HelpCommand,
		Description:              "Get help",
		DefaultMemberPermissions: &adminPermissions,
		DefaultMemberPermissions: &userPermissions,
		DMPermission:             &dmPermission,
	},
}
