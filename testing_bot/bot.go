package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"pkg.nit.so/switchboard"
)

type testArgs struct {
	User     discordgo.User    `description:"Test user"`
	Channel  discordgo.Channel `description:"Test channel"`
	Optional *string           `description:"An optional argument, with no default"`
	Default  string            `description:"An optional argument, with a default" default:"testing"`
}

func testCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, args testArgs) {
	fmt.Printf("%#+v\n", args)
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Hello world!",
					Description: "Hello world from _Switchboard_!",
					Color:       0xFF55AA,
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error responding to command: %s", err)
	}
}

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("error creating Discord session: %s", err)
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	switchboardInstance := &switchboard.Switchboard{}
	_ = switchboardInstance.AddCommand(&switchboard.Command{
		Name:        "test",
		Description: "Hello world from Switchboard!",
		Handler:     testCommand,
		GuildID:     os.Getenv("DISCORD_GUILD_ID"),
	})
	session.AddHandler(switchboardInstance.HandleInteractionCreate)
	err = switchboardInstance.SyncCommands(session, os.Getenv("DISCORD_APP_ID"))
	if err != nil {
		log.Fatalf("error registering commands: %s", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening session: %s", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err = session.Close(); err != nil {
		log.Fatalf("error closing session: %s", err)
	}
}
