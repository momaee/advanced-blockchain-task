package the_graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (g *theGraph) requestToTheGraph(ctx context.Context, query string, response interface{}) error {

	data := Payload{
		Query:     query,
		Variables: nil,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/graphprotocol/compound-v2", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code from the graphql. %v, %v", resp.StatusCode, err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, response)
	if err != nil {
		return err
	}

	return nil
}

func (g *theGraph) GetMintEvent(ctx context.Context, hash, logIndex string) (*Event, error) {
	lhash := strings.ToLower(hash)
	query := "{mintEvents(where:{id:\"" + lhash + "-" + logIndex + "\"}){id,amount,to,from,blockNumber,blockTime,underlyingAmount}}"
	response := MintEventResponse{}
	if err := g.requestToTheGraph(ctx, query, &response); err != nil {
		return nil, err
	}
	if len(response.Data.MintEvents) > 1 {
		return nil, fmt.Errorf("alarm. expected one result but found more than one. query: %v", query)
	}

	if len(response.Data.MintEvents) == 0 {
		return nil, nil
	}

	return response.Data.MintEvents[0], nil
}

func (g *theGraph) GetRedeemEvent(ctx context.Context, hash, logIndex string) (*Event, error) {
	lhash := strings.ToLower(hash)
	query := "{redeemEvents(where:{id:\"" + lhash + "-" + logIndex + "\"}){id,amount,to,from,blockNumber,blockTime,underlyingAmount}}"
	response := RedeemEventResponse{}
	if err := g.requestToTheGraph(ctx, query, &response); err != nil {
		return nil, err
	}

	if len(response.Data.RedeemEvents) > 1 {
		return nil, fmt.Errorf("alarm. expected one result but found more than one. query: %v", query)
	}

	if len(response.Data.RedeemEvents) == 0 {
		return nil, nil
	}

	return response.Data.RedeemEvents[0], nil
}

func (g *theGraph) GetLastDayMintEvents(ctx context.Context, contractAddress string) ([]*Event, error) {
	lContractAddress := strings.ToLower(contractAddress)
	time := time.Now().Unix() - time.Now().Unix()%60*60
	time -= 23 * 60 * 60
	query := "{mintEvents(where:{from:\"" + lContractAddress + "\",blockTime_gte:" + strconv.Itoa(int(time)) + "}){id,amount,to,from,blockNumber,blockTime,underlyingAmount}}"
	response := MintEventResponse{}
	if err := g.requestToTheGraph(ctx, query, &response); err != nil {
		return nil, err
	}
	return response.Data.MintEvents, nil
}

func (g *theGraph) GetLastDayRedeemEvents(ctx context.Context, contractAddress string) ([]*Event, error) {
	lContractAddress := strings.ToLower(contractAddress)
	time := time.Now().Unix() - time.Now().Unix()%60*60
	time -= 23 * 60 * 60
	query := "{redeemEvents(where:{to:\"" + lContractAddress + "\",blockTime_gte:" + strconv.Itoa(int(time)) + "}){id,amount,to,from,blockNumber,blockTime,underlyingAmount}}"
	response := RedeemEventResponse{}
	if err := g.requestToTheGraph(ctx, query, &response); err != nil {
		return nil, err
	}
	return response.Data.RedeemEvents, nil
}

func (g *theGraph) GetAllMarkets(ctx context.Context) ([]*Market, error) {
	query := "{markets{id}}"
	response := MarketResponse{}
	if err := g.requestToTheGraph(ctx, query, &response); err != nil {
		return nil, err
	}
	return response.Data.Markets, nil
}
