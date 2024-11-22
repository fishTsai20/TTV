package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"ttv-bot/client"
	"ttv-bot/model"
)

func (s *Service) query(id string, params map[string]string, msg api.Chattable) (*client.Response, error) {
	res, err := s.client.Query(id, params, msg, s.SendMessage)
	return res, err
}
func (s *Service) queryJettonHolderCountAmount(addr *model.TonAddr, msg api.Chattable) (float64, error) {
	res, err := s.query("690285", map[string]string{"contract_address": addr.MainnetBounceable}, msg)
	if err != nil {
		return 0, err
	}
	if len(res.Data.Data) == 0 {
		return -1, nil
	}
	ret, ok := res.Data.Data[0][0].(float64)
	if !ok {
		return 0, fmt.Errorf("turn query res into float64 failed, %+v", zap.Any("res", res.Data.Data[0]))
	}
	return ret, nil
}

func (s *Service) queryJettonTotalAmount(addr *model.TonAddr, msg api.Chattable) (float64, error) {
	res, err := s.query("690295", map[string]string{"jetton_master": addr.MainnetBounceable}, msg)
	if err != nil {
		return 0, err
	}
	if len(res.Data.Data) == 0 {
		return -1, nil
	}
	ret, ok := res.Data.Data[0][0].(float64)
	if !ok {
		return 0, fmt.Errorf("turn query res into float64 failed, %+v", zap.Any("res", res.Data.Data[0]))
	}
	return ret, nil
}

func (s *Service) queryJetton24hVol(addr *model.TonAddr, msg api.Chattable) (float64, error) {
	res, err := s.query("690294", map[string]string{"jetton_master": addr.MainnetBounceable}, msg)
	if err != nil {
		return 0, err
	}
	if len(res.Data.Data) == 0 {
		return -1, nil
	}
	ret, ok := res.Data.Data[0][0].(float64)
	if !ok {
		return 0, fmt.Errorf("turn query res into float64 failed, %+v", zap.Any("res", res.Data.Data[0]))
	}
	return ret, nil
}

func (s *Service) query24HJettonNewPools(msg api.Chattable) ([]model.Pool, error) {
	res, err := s.query("690291", map[string]string{"interval": "1"}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get ston.fi new pools")
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}
		var pools []model.Pool
		for _, row := range res.Data.Data {
			if len(row) < 3 {
				return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
			} else {
				txHash, ok1 := row[0].(string)
				pool, ok2 := row[1].(string)
				createdAt, ok3 := row[2].(string)
				if !ok1 || !ok2 || !ok3 {
					return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
				} else {
					pools = append(pools, model.Pool{
						TxHash:    txHash,
						Pool:      pool,
						CreatedAt: createdAt,
					})
				}
			}
		}
		return pools, nil
	}
}

func (s *Service) queryJettonTopN(addr *model.TonAddr, n int, msg api.Chattable) ([]model.Jetton, error) {
	res, err := s.query("690286", map[string]string{"jetton_master": addr.MainnetBounceable, "n": fmt.Sprintf("%+v", n)}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get top%d holders", n)
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}
		var jettons []model.Jetton
		for _, row := range res.Data.Data {
			if len(row) < 6 {
				return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
			} else {
				jettonWalletAddress, ok1 := row[0].(string)
				walletAddress, ok2 := row[1].(string)
				contractAddress, ok3 := row[2].(string)
				jettonBalance, ok4 := row[3].(float64)
				rnn, ok5 := row[4].(float64)
				totalAmount, ok6 := row[5].(float64)

				if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
					return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
				} else {
					jettons = append(jettons, model.Jetton{
						JettonWalletAddress: jettonWalletAddress,
						WalletAddress:       walletAddress,
						JettonMasterAddress: contractAddress,
						JettonBalance:       model.FormatNumber(jettonBalance),
						Percent:             model.FormatNumber(jettonBalance / totalAmount * 100),
						Rnn:                 int(rnn),
					})
				}
			}
		}
		sort.Slice(jettons, func(i, j int) bool {
			return jettons[i].Rnn < jettons[j].Rnn
		})
		return jettons, nil
	}
}

func (s *Service) queryJettonBalance(msg api.Chattable, jettonMaster *model.TonAddr, wallet *model.TonAddr) (*model.Jetton, error) {
	res, err := s.query("690288", map[string]string{"jetton_master": jettonMaster.MainnetBounceable, "wallet_address": wallet.MainnetBounceable}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get jetton current balance")
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}
		row := res.Data.Data[0]
		if len(row) < 4 {
			return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
		}
		jettonWalletAddress, ok1 := row[0].(string)
		walletAddress, ok2 := row[1].(string)
		jettonMasterAddress, ok3 := row[2].(string)
		jettonBalance, ok4 := row[3].(string)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
		} else {
			return &model.Jetton{
				JettonWalletAddress: jettonWalletAddress,
				WalletAddress:       walletAddress,
				JettonMasterAddress: jettonMasterAddress,
				JettonBalance:       jettonBalance,
			}, nil
		}
	}
}

func (s *Service) queryNFTitemsByCollection(addr *model.TonAddr, limit int, msg api.Chattable) ([]model.NFT, error) {
	res, err := s.query("690301", map[string]string{"nft_collection": addr.MainnetBounceable, "limit": strconv.Itoa(limit)}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT items by collection")
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}
		var NFTs []model.NFT
		for _, row := range res.Data.Data {
			if len(row) < 10 {
				return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
			} else {
				NFTAddress, ok1 := row[0].(string)
				NFTIndex, ok2 := row[1].(string)
				walletAddress, ok3 := row[2].(string)
				NFTCollectionAddress, ok4 := row[3].(string)
				contentUri, ok5 := row[4].(string)
				contentName, ok6 := row[5].(string)
				contentDescription, ok7 := row[6].(string)
				contentImage, ok8 := row[7].(string)
				contentImageData, ok9 := row[8].(string)
				count, ok10 := row[9].(float64)

				if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 {
					return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
				} else {
					NFTIndexI, err := strconv.Atoi(NFTIndex)
					if err != nil {
						continue
					}
					NFTs = append(NFTs, model.NFT{
						NFTAddress:           NFTAddress,
						NFTIndex:             NFTIndexI,
						WalletAddress:        walletAddress,
						NFTCollectionAddress: NFTCollectionAddress,
						ContentUri:           contentUri,
						ContentName:          contentName,
						ContentDescription:   contentDescription,
						ContentImage:         contentImage,
						ContentImageData:     contentImageData,
						Count:                int(count),
					})
				}
			}
		}
		sort.Slice(NFTs, func(i, j int) bool {
			return NFTs[i].NFTIndex < NFTs[j].NFTIndex
		})
		return NFTs, nil
	}
}

func (s *Service) queryNFTCollectionByItem(addr *model.TonAddr, msg api.Chattable) (*model.NFT, error) {
	res, err := s.query("690302", map[string]string{"nft_item": addr.MainnetBounceable}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT collection by item")
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}

		row := res.Data.Data[0]
		if len(row) < 9 {
			return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
		}
		NFTAddress, ok1 := row[0].(string)
		NFTIndex, ok2 := row[1].(string)
		walletAddress, ok3 := row[2].(string)
		NFTCollectionAddress, ok4 := row[3].(string)
		contentUri, ok5 := row[4].(string)
		contentName, ok6 := row[5].(string)
		contentDescription, ok7 := row[6].(string)
		contentImage, ok8 := row[7].(string)
		contentImageData, ok9 := row[8].(string)

		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
			return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
		}
		NFTIndexI, err := strconv.Atoi(NFTIndex)
		if err != nil {
			return nil, fmt.Errorf("invalid NFT index: %+v", zap.Any("NFT index", NFTIndex))
		}
		return &model.NFT{
			NFTAddress:           NFTAddress,
			NFTIndex:             NFTIndexI,
			WalletAddress:        walletAddress,
			NFTCollectionAddress: NFTCollectionAddress,
			ContentUri:           contentUri,
			ContentName:          contentName,
			ContentDescription:   contentDescription,
			ContentImage:         contentImage,
			ContentImageData:     contentImageData,
		}, nil

	}
}

func (s *Service) queryNFTAssetsByWallet(addr *model.TonAddr, limit int, msg api.Chattable) ([]model.NFT, error) {
	res, err := s.query("690303", map[string]string{"wallet_address": addr.MainnetBounceable, "limit": strconv.Itoa(limit)}, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT assets by wallet")
	} else {
		if len(res.Data.Data) == 0 {
			return nil, nil
		}
		var NFTs []model.NFT
		for _, row := range res.Data.Data {
			if len(row) < 10 {
				return nil, fmt.Errorf("invalid row length: %+v", zap.Any("row", row))
			} else {
				NFTAddress, ok1 := row[0].(string)
				NFTIndex, ok2 := row[1].(string)
				walletAddress, ok3 := row[2].(string)
				NFTCollectionAddress, ok4 := row[3].(string)
				contentUri, ok5 := row[4].(string)
				contentName, ok6 := row[5].(string)
				contentDescription, ok7 := row[6].(string)
				contentImage, ok8 := row[7].(string)
				contentImageData, ok9 := row[8].(string)
				count, ok10 := row[9].(float64)

				if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 {
					return nil, fmt.Errorf("invalid row data: %+v", zap.Any("row", row))
				} else {
					NFTIndexI, err := strconv.Atoi(NFTIndex)
					if err != nil {
						continue
					}
					NFTs = append(NFTs, model.NFT{
						NFTAddress:           NFTAddress,
						NFTIndex:             NFTIndexI,
						WalletAddress:        walletAddress,
						NFTCollectionAddress: NFTCollectionAddress,
						ContentUri:           contentUri,
						ContentName:          contentName,
						ContentDescription:   contentDescription,
						ContentImage:         contentImage,
						ContentImageData:     contentImageData,
						Count:                int(count),
					})
				}
			}
		}
		return NFTs, nil
	}
}

func (s *Service) queryIsContractDeployed(addr *model.TonAddr, msg api.Chattable) (bool, error) {
	res, err := s.query("690304", map[string]string{"account": addr.MainnetBounceable}, msg)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	if len(res.Data.Data) == 0 {
		return false, nil
	}
	ret, ok := res.Data.Data[0][0].(float64)
	if !ok {
		return false, fmt.Errorf("turn query res into bool failed, %+v", zap.Any("res", res.Data.Data[0]))
	}
	if ret == 1 {
		return true, nil
	}
	return false, nil
}
