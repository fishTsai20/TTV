package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"ttv-bot/model"
)

type ValidAccountsCache struct {
	addressMap map[string]string
	nameMap    map[string]*model.TonAddr
	init       sync.Once
}

type ValidNFTsCache struct {
	addressMap map[string]string
	nameMap    map[string]*model.TonAddr
	init       sync.Once
}

type ValidJettonsCache struct {
	addressMap map[string]string
	nameMap    map[string]*model.TonAddr
	init       sync.Once
}

func NewValidAccountsCache() *ValidAccountsCache {
	return &ValidAccountsCache{
		nameMap:    make(map[string]*model.TonAddr),
		addressMap: make(map[string]string),
	}
}

func NewValidJettonsCache() *ValidJettonsCache {
	return &ValidJettonsCache{
		nameMap:    make(map[string]*model.TonAddr),
		addressMap: make(map[string]string),
	}
}

func NewValidNFTsCache() *ValidNFTsCache {
	return &ValidNFTsCache{
		nameMap:    make(map[string]*model.TonAddr),
		addressMap: make(map[string]string),
	}
}
func (v *ValidAccountsCache) initializeAccounts() {
	v.init.Do(func() {
		getAddressNameMap(v.nameMap, v.addressMap, "tonkeeper", "ton-assets", "accounts")
	})
}

func (v *ValidJettonsCache) initializeJettons() {
	v.init.Do(func() {
		getAddressNameMap(v.nameMap, v.addressMap, "tonkeeper", "ton-assets", "jettons")
	})
}

func (v *ValidNFTsCache) initializeNFTs() {
	v.init.Do(func() {
		getAddressNameMap(v.nameMap, v.addressMap, "tonkeeper", "ton-assets", "collections")
	})
}

func (v *ValidAccountsCache) GetAccountAddressByName(name string) *model.TonAddr {
	v.initializeAccounts()
	return v.nameMap[name]
}

func (v *ValidAccountsCache) GetAccountNameByAddress(addr *model.TonAddr) string {
	v.initializeAccounts()
	return v.addressMap[addr.Hex]
}

func (v *ValidJettonsCache) GetJettonAddressByName(name string) *model.TonAddr {
	v.initializeJettons()
	return v.nameMap[name]
}

func (v *ValidJettonsCache) GetJettonNameByAddress(addr *model.TonAddr) string {
	v.initializeJettons()
	return v.addressMap[addr.Hex]
}

func (v *ValidNFTsCache) GetNFTAddressByName(name string) *model.TonAddr {
	v.initializeNFTs()
	return v.nameMap[name]
}

func (v *ValidNFTsCache) GetNFTNameByAddress(addr *model.TonAddr) string {
	v.initializeNFTs()
	return v.addressMap[addr.Hex]
}

func (v *ValidAccountsCache) GetAll() map[string]*model.TonAddr {
	v.initializeAccounts()
	return v.nameMap
}

func (v *ValidNFTsCache) GetAll() map[string]*model.TonAddr {
	v.initializeNFTs()
	return v.nameMap
}

func (v *ValidJettonsCache) GetAll() map[string]*model.TonAddr {
	v.initializeJettons()
	return v.nameMap
}

func FuzzGetCache(name string, nameMap map[string]*model.TonAddr) map[string]*model.TonAddr {
	ret := make(map[string]*model.TonAddr)
	normalizedName := Normalize(name)
	for n, addr := range nameMap {
		if strings.Contains(Normalize(n), normalizedName) {
			ret[n] = addr
		}
	}
	return ret
}
func (v *ValidAccountsCache) FuzzyGET(name string) map[string]*model.TonAddr {
	v.initializeAccounts()
	return FuzzGetCache(name, v.nameMap)
}

func (v *ValidNFTsCache) FuzzyGET(name string) map[string]*model.TonAddr {
	v.initializeNFTs()
	return FuzzGetCache(name, v.nameMap)
}

func (v *ValidJettonsCache) FuzzyGET(name string) map[string]*model.TonAddr {
	v.initializeJettons()
	return FuzzGetCache(name, v.nameMap)
}

type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

func getFilesFromGitHub(repoOwner, repoName, path string) ([]GitHubContent, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=main", repoOwner, repoName, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contents []GitHubContent
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func fetchYAMLFile(url string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func parseYAMLData(data []byte, nameMap map[string]*model.TonAddr, addressMap map[string]string) error {
	var accounts []model.Account
	err := yaml.Unmarshal(data, &accounts)
	if err != nil {
		var account model.Account
		if strings.Contains(err.Error(), "cannot unmarshal !!map ") {
			err := yaml.Unmarshal(data, &account)
			if err != nil {
				return err
			}
			accounts = append(accounts, account)
		} else {
			return err
		}
	}

	for _, account := range accounts {
		err, parseAddr := model.ParseTonAddress(account.Address)
		if err == nil {
			addressMap[parseAddr.Hex] = account.Name
			nameMap[account.Name] = parseAddr
		}
	}

	return nil
}

func getAddressNameMap(nameMap map[string]*model.TonAddr, addressMap map[string]string, repoOwner string, repoName string, path string) {
	files, err := getFilesFromGitHub(repoOwner, repoName, path)
	if err != nil {
		log.Fatalf("Error fetching file list: %v", err)
	}
	// Iterate over each file and fetch/parse the YAML
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".yaml") { // Only process YAML files
			fmt.Printf("Fetching file: %s...\n", file.Name)

			// Fetch the YAML file from the download URL
			data, err := fetchYAMLFile(file.DownloadURL)
			if err != nil {
				log.Printf("Error fetching %s: %v\n", file.Name, err)
				continue
			}

			// Parse the YAML data
			err = parseYAMLData(data, nameMap, addressMap)
			if err != nil {
				log.Printf("Error parsing YAML from %s: %v\n", file.Name, err)
				continue
			}
		}
	}
}

func Normalize(input string) string {
	// 正则表达式：匹配所有非字母和数字的字符
	re := regexp.MustCompile("[^a-zA-Z0-9]")
	// 转为小写并去除非字母和数字字符
	return re.ReplaceAllString(strings.ToLower(input), "")
}
