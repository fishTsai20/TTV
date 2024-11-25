package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"sync"
	"ttv-bot/cache"
	"ttv-bot/client"
	"ttv-bot/model"
)

type Service struct {
	bot                *api.BotAPI
	config             api.UpdateConfig
	commandFunctions   map[string]func(*api.Message)
	replyFunctions     map[string]func(message *api.Message)
	pageFunctions      map[string]func(upmsg *api.Message, page int, first bool)
	validAccountsCache *cache.ValidAccountsCache
	validJettonsCache  *cache.ValidJettonsCache
	validNFTsCache     *cache.ValidNFTsCache
	client             *client.Client
	chatMap            sync.Map
	stonFiPoolsCache   *stonFiPoolsCache
	topHoldersCache    *topHoldersCache
	nftCollectionCache *nftCollectionCache
	nftAssetsCache     *nftAssetsCache
	nftItemsCache      *nftItemsCache
}

func NewService(botToken string, debug bool, timeout int, apiKey string) *Service {
	bot, err := api.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = debug
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)

	config := api.NewUpdate(0)
	config.Timeout = timeout
	return &Service{
		bot:                bot,
		config:             config,
		commandFunctions:   make(map[string]func(*api.Message)),
		replyFunctions:     make(map[string]func(message *api.Message)),
		pageFunctions:      make(map[string]func(upmsg *api.Message, page int, first bool)),
		validAccountsCache: cache.NewValidAccountsCache(),
		validJettonsCache:  cache.NewValidJettonsCache(),
		validNFTsCache:     cache.NewValidNFTsCache(),
		client: &client.Client{
			ApiKey:  apiKey,
			BaseURL: "https://api.chainbase.com/api/v1",
		},
		stonFiPoolsCache:   newStonFiPoolsCache(),
		topHoldersCache:    newTopHoldersCache(),
		nftItemsCache:      newNFTItemsCache(),
		nftAssetsCache:     newNFTAssetsCache(),
		nftCollectionCache: newNFTCollectionCache(),
	}
}

func (s *Service) init() {
	// 设置命令列表

	commands := []api.BotCommand{
		{Command: model.StartCommand, Description: "Start the bot."},
		{Command: model.ListCommand, Description: "List accounts, nfts, jettons."},
		{Command: model.FakeCommand, Description: "Detect fake contracts with the same name."},
		{Command: model.HoneyPotCommand, Description: "Detect honeypot contracts."},
		{Command: model.NFTCommand, Description: "Query NFTs, including personal NFT assets, NFT collection, and NFT items."},
		{Command: model.JettonsCommand, Description: " jettons, including ston.fi new pools, jetton holders, 24-hour trading volume, jetton top 10 holders, jetton balance."},
		{Command: model.ListJettonsCommand, Description: "list Jettons"},
		{Command: model.ListJettonsByNameCommand, Description: "list Jettons with name filter"},
		{Command: model.ListJettonsDefaultCommand, Description: "list Jettons without any filters"},
		{Command: model.ListAccountsCommand, Description: "list accounts"},
		{Command: model.ListAccountsByNameCommand, Description: "list accounts with name filter"},
		{Command: model.ListAccountsDefaultCommand, Description: "list accounts without any filters"},
		{Command: model.ListNFTsCommand, Description: "list NFTs"},
		{Command: model.ListNFTsByNameCommand, Description: "list NFTs with name filter"},
		{Command: model.ListNFTsDefaultCommand, Description: "list NFTs without any filters"},

		{Command: model.JetttonNewPoolsCommand, Description: "Detect new jetton pools created in the last 24 hours."},
		{Command: model.FakeAccountCommand, Description: "Detect accounts that impersonate legitimate accounts with the same name."},
		{Command: model.FakeJettonCommand, Description: "Detect fake jettons that mimic real tokens."},
		{Command: model.FakeNFTCommand, Description: "Detect fake NFTs, including fake collections and items."},
		{Command: model.JettonHolderCommand, Description: "Monitor the number of jetton holders."},
		{Command: model.JettonChangesCommand, Description: "Track the 24-hour jetton changes and turnover rate."},
		{Command: model.JettonAmountCommand, Description: "Track the 24-hour jetton amounts (volume)."},
		{Command: model.JettonTopHoldersCommand, Description: "Monitor the top 10 jetton holders."},
		{Command: model.JettonBalanceCommand, Description: "Monitor the balance of jetton holders."},
		{Command: model.NFTCollectionCommand, Description: "Retrieve NFT collections and their items."},
		{Command: model.NFTItemCommand, Description: "Find the collection for a specific NFT item."},
		{Command: model.NFTAssetCommand, Description: "Query personal NFT assets."},
	}

	newCommands := api.NewSetMyCommands(commands...)
	if _, err := s.bot.Request(newCommands); err != nil {
		log.Println("Unable to set commands" + err.Error())
	}
	//command
	s.commandFunctions[model.StartCommand] = s.handleStartCommand
	s.commandFunctions[model.FakeCommand] = s.handleFakeCommand
	s.commandFunctions[model.ListCommand] = s.handleListCommand
	s.commandFunctions[model.FakeAccountCommand] = s.handleFakeAccountCommand
	s.commandFunctions[model.FakeJettonCommand] = s.handleFakeJettonCommand
	s.commandFunctions[model.FakeNFTCommand] = s.handleFakeNFTCommand
	s.commandFunctions[model.JettonHolderCommand] = s.handleJettonHoldersCommand
	s.commandFunctions[model.JettonsCommand] = s.handleJettonsCommand
	s.commandFunctions[model.ListJettonsCommand] = s.handleListJettonsCommand
	s.commandFunctions[model.ListJettonsByNameCommand] = s.handleListJettonsByNameCommand
	s.commandFunctions[model.ListJettonsDefaultCommand] = s.handleListJettonsDefaultCommand
	s.commandFunctions[model.ListAccountsCommand] = s.handleListAccountsCommand
	s.commandFunctions[model.ListAccountsByNameCommand] = s.handleListAccountsByNameCommand
	s.commandFunctions[model.ListAccountsDefaultCommand] = s.handleListAccountsDefaultCommand
	s.commandFunctions[model.JetttonNewPoolsCommand] = s.handle24HJettonNewPoolsCommand
	s.commandFunctions[model.JettonAmountCommand] = s.handleJettonAmountCommand
	s.commandFunctions[model.JettonChangesCommand] = s.handleJettonChangeCommand
	s.commandFunctions[model.JettonTopHoldersCommand] = s.handleJettonTopNCommand
	s.commandFunctions[model.JettonBalanceCommand] = s.handleJettonBalanceCommand
	s.commandFunctions[model.NFTCollectionCommand] = s.handleNFTCollectionCommand

	s.commandFunctions[model.ListNFTsCommand] = s.handleListNFTsCommand
	s.commandFunctions[model.ListNFTsByNameCommand] = s.handleListNFTsByNameCommand
	s.commandFunctions[model.ListNFTsDefaultCommand] = s.handleListNFTsDefaultCommand
	s.commandFunctions[model.NFTItemCommand] = s.handleNFTItemCommand
	s.commandFunctions[model.NFTAssetCommand] = s.handleNFTAssetsCommand
	s.commandFunctions[model.NFTCommand] = s.handleNFTsCommand
	s.commandFunctions[model.HoneyPotCommand] = s.handleHoneyPotCommand
	s.commandFunctions[model.HelpCommand] = s.handleHelpCommand

	//callback
	s.replyFunctions[model.ListNFTsDefaultCommand] = s.replyListNFTsByNameCommand
	s.replyFunctions[model.ListNFTsByNameCommand] = s.replyListNFTsByNameCommand
	s.replyFunctions[model.ListJettonsDefaultCommand] = s.replyListJettonsByNameCommand
	s.replyFunctions[model.ListJettonsByNameCommand] = s.replyListJettonsByNameCommand
	s.replyFunctions[model.ListAccountsByNameCommand] = s.replyListAccountsByNameCommand
	s.replyFunctions[model.ListAccountsByNameCommand] = s.replyListAccountsByNameCommand
	s.replyFunctions[model.FakeAccountCommand] = s.replyFakeAccountCommand
	s.replyFunctions[model.FakeJettonCommand] = s.replyFakeJettonCommand
	s.replyFunctions[model.FakeNFTCommand] = s.replyFakeNFTCommand
	s.replyFunctions[model.JettonHolderCommand] = s.replyJettonHoldersCommand
	s.replyFunctions[model.JettonAmountCommand] = s.replyJettonAmountCommand
	s.replyFunctions[model.JettonChangesCommand] = s.replyJettonChangeCommand
	s.replyFunctions[model.JettonTopHoldersCommand] = s.replyJettonTopNCommand
	s.replyFunctions[model.JettonBalanceCommand] = s.replyJettonBalanceCommand
	s.replyFunctions[model.NFTCollectionCommand] = s.replyNFTCollectionsCommand
	s.replyFunctions[model.NFTAssetCommand] = s.replyNFTAssetsCommand
	s.replyFunctions[model.NFTItemCommand] = s.replyNFTItemCommand
	s.replyFunctions[model.HoneyPotCommand] = s.replyHoneyPotCommand

	s.pageFunctions[model.NFTCollectionCommand] = s.replyNFTCollectionPage
	s.pageFunctions[model.NFTAssetCommand] = s.replyNFTAssetsPage
	s.pageFunctions[model.JettonTopHoldersCommand] = s.replyJettonTopNPage
	s.pageFunctions[model.JetttonNewPoolsCommand] = s.reply24HJettonNewPoolsCommand
	s.pageFunctions[model.ListJettonsCommand] = s.replyListJettonsPage
	s.pageFunctions[model.ListNFTsCommand] = s.replyListNFTsPage
	s.pageFunctions[model.ListAccountsCommand] = s.replyListAccountsPage
}
func (s *Service) Start() {
	s.init()
	updates := s.bot.GetUpdatesChan(s.config)
	for update := range updates {
		if update.CallbackQuery != nil {
			go s.processCallBack(&update)
			continue
		}
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		go s.processUpdate(&update)
	}
}

func (s *Service) processUpdate(update *api.Update) {
	upmsg := update.Message
	if upmsg.IsCommand() {
		go s.processCommond(update)
	} else {
		go s.processMessage(update)
	}
}

func (s *Service) processCommond(update *api.Update) {
	upmsg := update.Message
	if fn, exists := s.commandFunctions[upmsg.Command()]; exists {
		fn(upmsg)
	}
}

func (s *Service) processCallBack(update *api.Update) {
	if callback := update.CallbackQuery; callback != nil {
		if fn, exists := s.commandFunctions[callback.Data]; exists {
			if _, existsRe := s.replyFunctions[callback.Data]; existsRe {
				s.chatMap.Store(update.CallbackQuery.Message.Chat.ID, callback.Data)
			}
			fn(callback.Message)
		} else if f := s.getPageFunction(callback.Data); f != nil {
			split := strings.Split(callback.Data, "-")
			if len(split) > 1 {
				page, err := strconv.Atoi(split[1])
				if err != nil {
					log.Println(err)
					return
				}
				f(callback.Message, page, false)
			}
		}
	}
}

func (s *Service) processMessage(update *api.Update) {
	if command, exists := s.chatMap.Load(update.Message.Chat.ID); exists {
		if fn, existsFunc := s.replyFunctions[command.(string)]; existsFunc {
			fn(update.Message)
		}
	}
}

func (s *Service) getPageFunction(data string) func(upmsg *api.Message, page int, first bool) {
	split := strings.Split(data, "-")
	if len(split) > 1 {
		f := s.pageFunctions[split[0]]
		return f
	}
	return nil

}
