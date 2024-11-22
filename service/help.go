package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Service) handleHelpCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "\n\n**Welcome to TTV - Ton TokenVigil Analytics Bot!**\n\n" +
		"Here are the available commands you can use:\n\n" +
		"### **General Commands:**\n" +
		"- **/start**: Start the bot.\n" +
		"- **/start**\n\n" +
		"### **Fake Contract Detection:**\n" +
		"- **/fake**: Detect fake contracts with the same name.\n" +
		"- **/fakeAccount**: Detect accounts that impersonate legitimate accounts with the same name.\n" +
		"- **/fakeJetton**: Detect fake jettons that mimic real tokens.\n" +
		"- **/fakeNFT**: Detect fake NFTs, including fake collections and items.\n\n" +
		"### **Honeypot Contract Detection:**\n" +
		"- **/honeyPot**: Detect honeypot contracts designed to trap users into irreversible transactions.\n\n" +
		"### **Jetton Queries:**\n" +
		"- **/jettons**: Query jettons, including ston.fi new pools, jetton holders, 24-hour trading volume, jetton top 10 holders, jetton balance.\n" +
		"- **/24HJetttonNewPools**: Detect new jetton pools created in the last 24 hours.\n" +
		"- **/jettonHolder**: Monitor the number of jetton holders.\n" +
		"- **/24HJettonChanges**: Track the 24-hour jetton changes and turnover rate.\n" +
		"- **/24HJettonAmounts**: Track the 24-hour jetton amounts (volume).\n" +
		"- **/jettonTop10**: Monitor the top 10 jetton holders.\n" +
		"- **/jettonBalance**: Monitor the balance of jetton holders.\n\n" +
		"### **NFT Queries:**\n" +
		"- **/nfts**: Query NFTs, including personal NFT assets, NFT collection, and NFT items.\n" +
		"- **/nftCollection**: Retrieve NFT collections and their items.\n" +
		"- **/nftItem**: Find the collection for a specific NFT item.\n" +
		"- **/nftAsset**: Query personal NFT assets.\n\n" +
		"---\n\nPowered by [Chainbase](https://chainbase.com/)"
	msg.ParseMode = api.ModeMarkdown
	b.SendMessage(msg)
	b.chatMap.Delete(upmsg.Chat.ID)
}
