package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/pomaretta/mcpingserver"
)

func main() {
	bindAddr := flag.String("bindAddr", "0.0.0.0:25565", "Ip:port combination to bind to")
	pingMotdFlg := flag.String("motd", "A Golang inplace server", "Motd, with newlines and & colors")
	kickMsgFlg := flag.String("kickMsg", "&4This is not a joinable server!", "Kickmsg, with newlines and & colors")
	playerCnt := flag.Int("players", 0, "Player count")
	playerCap := flag.Int("cap", 20, "Player cap")
	faviconLocation := flag.String("favicon", "favicon.png", "Location of the favicon to display. Must be 64x64 PNG")
	serverVersionName := flag.String("serverVersionName", "1.12.2", "Version of server (string version)")
	serverVersionNum := flag.Int("serverVersionNumber", 340, "Version of server (number version)")

	flag.Parse()

	kickMsg := mcpingserver.TranslateColorCodes(*kickMsgFlg)
	pingMotd := mcpingserver.TranslateColorCodes(*pingMotdFlg)

	kickJson := mcpingserver.ConvertMCChat(kickMsg)
	motdJson := mcpingserver.ConvertMCChat(pingMotd)

	faviconb64 := readFavicon(*faviconLocation)

	pingResponse := mcpingserver.PingResponse{
		mcpingserver.VersionEntry{*serverVersionName, uint(*serverVersionNum)},
		mcpingserver.PlayersEntry{*playerCap, *playerCnt, []mcpingserver.PlayerEntry{}},
		motdJson,
		faviconb64,
	}

	legacyPing := mcpingserver.LegacyPingResponse{
		*playerCnt, *playerCap, *serverVersionNum, *serverVersionName, pingMotd}

	responder := mcpingserver.CreateSimpleResponder(&pingResponse, kickJson, &legacyPing)

	pingServer := mcpingserver.CreatePingServer(*bindAddr, responder)

	fmt.Println("Binding to", *bindAddr)
	err := pingServer.Bind()
	if err != nil {
		panic(err)
	}
	err = pingServer.AcceptConnections(handleError)
	if err != nil {
		panic(err)
	}
}

func handleError(err error) {
	fmt.Println("Error occurred: ", err)
}

func readFavicon(loc string) string {
	faviconData, err := ioutil.ReadFile(loc)
	if err != nil {
		fmt.Println("WARNING: Failed to load favicon!", err)
		fmt.Println("WARNING: Server will not respond with a favicon!")
		return ""
	} else {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(faviconData)
	}
}
