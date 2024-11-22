package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
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
	validCache         *cache.ValidCache
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
		bot:              bot,
		config:           config,
		commandFunctions: make(map[string]func(*api.Message)),
		replyFunctions:   make(map[string]func(message *api.Message)),
		pageFunctions:    make(map[string]func(upmsg *api.Message, page int, first bool)),
		validCache:       cache.NewValidCache(),
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
	s.commandFunctions[model.StartCommand] = s.handleStartCommand
	s.commandFunctions[model.FakeCommand] = s.handleFakeCommand
	s.commandFunctions[model.FakeAccountCommand] = s.handleFakeTypeCommand
	s.commandFunctions[model.FakeJettonCommand] = s.handleFakeTypeCommand
	s.commandFunctions[model.FakeNFTCommand] = s.handleFakeTypeCommand
	s.commandFunctions[model.JettonHolderCommand] = s.handleJettonHoldersCommand
	s.commandFunctions[model.JettonsCommand] = s.handleJettonsCommand
	s.commandFunctions[model.JetttonNewPoolsCommand] = s.handle24HJettonNewPoolsCommand
	s.commandFunctions[model.JettonAmountCommand] = s.handleJettonAmountCommand
	s.commandFunctions[model.JettonChangesCommand] = s.handleJettonChangeCommand
	s.commandFunctions[model.JettonTopHoldersCommand] = s.handleJettonTopNCommand
	s.commandFunctions[model.JettonBalanceCommand] = s.handleJettonBalanceCommand
	s.commandFunctions[model.NFTCollectionCommand] = s.handleNFTCollectionCommand
	s.commandFunctions[model.NFTItemCommand] = s.handleNFTItemCommand
	s.commandFunctions[model.NFTAssetCommand] = s.handleNFTAssetsCommand
	s.commandFunctions[model.NFTCommand] = s.handleNFTsCommand

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

	s.pageFunctions[model.NFTCollectionCommand] = s.replyNFTCollectionPage
	s.pageFunctions[model.NFTAssetCommand] = s.replyNFTAssetsPage
	s.pageFunctions[model.JettonTopHoldersCommand] = s.replyJettonTopNPage
	s.pageFunctions[model.JetttonNewPoolsCommand] = s.reply24HJettonNewPoolsCommand
}
func (s *Service) Start() {
	s.init()
	updates, err := s.bot.GetUpdatesChan(s.config)
	if err != nil {
		panic("Can't get Updates")
	}
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
