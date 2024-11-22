package model

import (
	"strconv"
)

type NFT struct {
	NFTAddress           string `json:"nft_address"`
	NFTIndex             int    `json:"nft_index"`
	WalletAddress        string `json:"wallet_address"`
	NFTCollectionAddress string `json:"nft_collection_address"`
	ContentUri           string `json:"content_uri"`
	ContentName          string `json:"content_name"`
	ContentDescription   string `json:"content_description"`
	ContentImage         string `json:"content_image"`
	ContentImageData     string `json:"content_image_data"`
	Count                int    `json:"items_count"`
}

func (n NFT) ToTgText() string {
	res := "\n"
	res += "💎\n"
	if n.ContentName != "" {
		res += "*" + n.ContentName + "*\n"
	}
	if n.ContentDescription != "" {
		res += "--" + n.ContentDescription + "\n"
	}
	if n.ContentUri != "" {
		res += "[uri ↗️](" + n.ContentUri + ")\n"
	}
	res += "*nft: *[" + n.NFTAddress + "](https://tonscan.org/nft/" + n.NFTAddress + ")\n"
	res += "----------\n*index: *" + strconv.Itoa(n.NFTIndex) + "\n"
	res += "*collection: *[" + n.NFTCollectionAddress + "](https://tonscan.org/nft/" + n.NFTCollectionAddress + ")\n"
	if n.WalletAddress != "" {
		res += "*wallet: *[" + n.WalletAddress + "](https://tonscan.org/address/" + n.WalletAddress + ")\n"
	}
	if n.ContentImage != "" {
		res += "*image: *" + n.ContentImage + "\n"
	}
	if n.ContentImageData != "" {
		res += "*imageData: *" + n.ContentImageData + "\n"
	}
	return res
}
