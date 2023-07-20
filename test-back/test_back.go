package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Collection struct {
	Address string `json : "collection"`
	Name    string `json : "name"`
	Symbol  string `json : "symbol"`
}

type NFT struct {
	Collection string `json : "collection"`
	Recipient  string `json : "recipient"`
	TokenID    uint64 `json : "tokenId"`
	TokenURI   string `json : "tokenUri"`
}

var collections []Collection
var nfts []NFT

func handleCreateCollection(w http.ResponseWriter, r *http.Request) {
	var collection Collection
	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := ethclient.Dial("YOUR_NODE_URL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contractABI, err := abi.JSON([]byte("YOUR_CONTRACT_ABI_JSON"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contractAddress := common.HexToAddress("YOUR_CONTRACT_ADDRESS")

	auth, err := bind.NewKeyedTransactorWithChainID(
		hexToKey("YOUR_CONTRACT_OWNER_PRIVATE_KEY"),
		big.NewInt(1337),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := contractABI.Pack("mintToken", common.HexToAddress(collection.Address), collection.Name, collection.Symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.SendTransaction(context.Background(), auth, &contractAddress, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collections = append(collections, collection)
	log.Println("Collection created:", collection.Address)

	w.WriteHeader(http.StatusCreated)
}

func handleMintNFT(w http.ResponseWriter, r *http.Request) {
	var nft NFT
	err := json.NewDecoder(r.Body).Decode(&nft)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := ethclient.Dial("YOUR_NODE_URL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contractABI, err := abi.JSON([]byte("YOUR_CONTRACT_ABI_JSON"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contractAddress := common.HexToAddress("YOUR_CONTRACT_ADDRESS")

	auth, err := bind.NewKeyedTransactorWithChainID(
		hexToKey("YOUR_CONTRACT_OWNER_PRIVATE_KEY"),
		big.NewInt(1337),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := contractABI.Pack("mintToken", common.HexToAddress(nft.Collection), common.HexToAddress(nft.Recipient), new(big.Int).SetUint64(nft.TokenID), nft.TokenURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.SendTransaction(context.Background(), auth, &contractAddress, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nfts = append(nfts, nft)
	log.Println("Token minted:", nft.TokenID)

	w.WriteHeader(http.StatusCreated)
}

func handleGetCollections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}

func handleGetNFTs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nfts)
}

func main() {
	// Добавление коллекции по умолчанию
	defaultCollection := Collection{
		Address: "0x123456789",
		Name:    "Default Collection",
		Symbol:  "DEF",
	}
	collections = append(collections, defaultCollection)

	// Добавление 5 NFT по умолчанию
	for i := 1; i <= 5; i++ {
		defaultNFT := NFT{
			Collection: defaultCollection.Address,
			Recipient:  "0x987654321",
			TokenID:    uint64(i),
			TokenURI:   "https://example.com/nft/" + strconv.Itoa(i),
		}
		nfts = append(nfts, defaultNFT)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/collections", handleGetCollections).Methods("GET")
	router.HandleFunc("/api/collections", handleCreateCollection).Methods("POST")
	router.HandleFunc("/api/nfts", handleGetNFTs).Methods("GET")
	router.HandleFunc("/api/nfts", handleMintNFT).Methods("POST")

	// Создание CORS-хэндлера с настройками по умолчанию
	corsHandler := cors.Default().Handler(router)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
