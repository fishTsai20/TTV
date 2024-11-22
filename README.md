# README

## Overview

TTV (Ton TokenVigil Analytics) is an AI-powered sophisticated token analysis platform specialized in TON coin security and risk assessment. The platform combines real-time monitoring of smart contract behaviors, holder distributions, and market activities with advanced machine learning algorithms for predictive risk analysis. TTV provides comprehensive tools for detecting honeypot contracts, identifying fake contracts, and analyzing jettons and NFTs to enhance security and transparency within the TON ecosystem.

## Bot

[https://t.me/tonTokenVigil_bot](https://t.me/tonTokenVigil_bot)

## Features

### HoneyPot
Honey pot contract detection identifies suspicious smart contracts that are designed to deceive users by trapping them into making irreversible transactions. These contracts are commonly known as "honeypots," where users are unable to withdraw their funds after interacting with the contract.

### Fake
Fake contract detection looks for contracts with the same name or identifiers as legitimate ones, potentially deceiving users into interacting with fraudulent contracts. It includes:

- **fakeAccount**: Detects fake accounts with the same name as legitimate ones.
- **fakeJetton**: Detects fake jettons, tokens designed to impersonate legitimate ones.
- **fakeNFT**: Detects fake NFTs, including fake NFT collections and items.

### Jettons
Jetton-related features focus on analyzing and monitoring jetton contracts and holders. Key functionalities include:

- **jettons**: Query basic information about jetton contracts.
- **24HJetttonNewPools**: Monitors new jetton pools created within the last 24 hours.
- **jettonHolder**: Tracks the number of jetton holders.
- **24HJettonChanges**: Tracks 24-hour jetton changes and trading volume.
- **24HJettonAmounts**: Tracks the 24-hour jetton amounts (volume).
- **jettonTop10**: Lists the top 10 jetton holders.
- **jettonBalance**: Tracks the balance of jetton holders.

### NFTs
NFT-related commands allow users to query and explore NFT collections, items, and assets:

- **nftCollection**: Retrieves NFT collections and their contents.
- **nftItem**: Queries the specific items in a given NFT collection.
- **nftAsset**: Queries an individual's personal NFT assets.

## TG Commands List

| Command                | Description                                                               |
|------------------------|---------------------------------------------------------------------------|
| **start**              | Start the bot.                                                            |
| **fake**               | Detect fake contracts with the same name.                               |
| **honeyPot**           | Detect honeypot contracts.                                              |
| **nfts**               | Query NFTs, including personal NFT assets, NFT collection, and NFT items. |
| **jettons**            | Query jettons, including ston.fi new pools, jetton holders, 24-hour trading volume, jetton top 10 holders, jetton balance. |
| **24HJetttonNewPools** | Detect new jetton pools created in the last 24 hours.                     |
| **fakeAccount**        | Detect fake accounts with the same name as legitimate ones.               |
| **fakeJetton**         | Detect fake jettons that mimic real tokens.                               |
| **fakeNFT**            | Detect fake NFTs, including fake collections and items.                   |
| **jettonHolder**       | Monitor the number of jetton holders.                                     |
| **24HJettonChanges**   | Track the 24-hour jetton changes and turnover rate.                       |
| **24HJettonAmounts**   | Track the 24-hour jetton amounts (volume).                                |
| **jettonTop10**        | Monitor the top 10 jetton holders.                                        |
| **jettonBalance**      | Monitor the balance of jetton holders.                                    |
| **nftCollection**      | Retrieve the NFT collection and its items.                                |
| **nftItem**            | Find the collection for a specific NFT item.                              |
| **nftAsset**           | Query an individual's personal NFT assets.                                |

## Usage

### Example Commands
To use any of the commands, simply input the desired command in the format below:
- /start
- /fake

## Installation
1. Clone the repository
   ```
   git clone https://github.com/fishTsai20/TTV.git
   cd ./TTV
    ``` 
2. Install dependencies
   ```
    go mod download
   ```
3. Build the application
   ```
   go build -o ttv-bot
    ```
4. Run the bot with a token and API key:
    ```
   ./ttv-bot ttv-bot --bot-token <your-bot-token> --api-key <your-api-key>
   ```
## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Powered By

This project is powered by [Chainbase](https://chainbase.com/), providing robust infrastructure and blockchain data analytics solutions.