package commands

import (
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils/rinha"
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"dungeon", "dg", "boss"},
		Run:       runDungeon,
		Available: true,
		Cooldown:  5,
		Usage:     "j!dungeon",
		Category:  1,
		Help:      "Adentre na dungeon",
	})
}

func runDungeon(session disgord.Session, msg *disgord.Message, args []string) {
	galo, _ := rinha.GetGaloDB(msg.Author.ID)
	if galo.Type == 0 {
		msg.Reply(context.Background(), session, msg.Author.Mention()+", Voce nao tem um galo, use j!galo para criar um")
		return
	}
	if len(args) == 0 {
		authorAvatar, _ := msg.Author.AvatarURL(512, true)
		msg.Reply(context.Background(), session, &disgord.Embed{
			Title: "Dungeon",
			Footer: &disgord.EmbedFooter{
				Text:    msg.Author.Username,
				IconURL: authorAvatar,
			},
			Color:       65535,
			Description: fmt.Sprintf("Voce esta no andar **%d**\nUse j!dungeon battle para batalhar contra o chefe", galo.Dungeon),
		})
		return
	}
	battleMutex.RLock()
	if currentBattles[msg.Author.ID] != "" {
		battleMutex.RUnlock()
		msg.Reply(context.Background(), session, "Voce ja esta lutando com o "+currentBattles[msg.Author.ID])
		return
	}
	battleMutex.RUnlock()
	if len(rinha.Dungeon) == galo.Dungeon {
		return
	}
	dungeon := rinha.Dungeon[galo.Dungeon]
	galoAdv := dungeon.Boss
	LockEvent(msg.Author.ID, "Boss "+rinha.Classes[galoAdv.Type].Name)
	defer UnlockEvent(msg.Author.ID)
	winner, _ := ExecuteRinha(msg, session, rinhaOptions{
		galoAuthor:  &galo,
		galoAdv:     &galoAdv,
		authorName:  rinha.GetName(msg.Author.Username, galo),
		advName:     "Boss " + rinha.Classes[galoAdv.Type].Name,
		authorLevel: rinha.CalcLevel(galo.Xp),
		advLevel:    rinha.CalcLevel(galoAdv.Xp),
		noItems: true,
	})
	if winner == 0 {
		diffGalo, endMsg := rinha.DungeonWin(dungeon.Level, galo)
		update := rinha.Diff(galo, diffGalo)
		update["dungeon"] = galo.Dungeon + 1
		rinha.UpdateGaloDB(msg.Author.ID, update)
		tag := msg.Author.Username + "#" + msg.Author.Discriminator.String()
		telemetry.Debug(fmt.Sprintf("%s wins %s", tag, endMsg), map[string]string{
			"user":         strconv.FormatUint(uint64(msg.Author.ID), 10),
			"dungeonLevel": string(galo.Dungeon),
		})
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16776960,
			Title:       "Dungeon",
			Description: fmt.Sprintf("Parabens %s voce consegiu derrotar o boss e avançar para o andar **%d** %s", msg.Author.Username, galo.Dungeon+1, endMsg),
		})
	} else {
		msg.Reply(context.Background(), session, &disgord.Embed{
			Color:       16711680,
			Title:       "Dungeon",
			Description: fmt.Sprintf("Parabens %s, voce perdeu. Use j!dungeon battle para tentar novamente", msg.Author.Username),
		})
	}
}
