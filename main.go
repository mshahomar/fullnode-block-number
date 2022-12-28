package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"

	"github.com/fatih/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	tronOne      string
	tronTwo      string
	tronUsername string
	tronPassword string
	bnbOne       string
	bnbUsername  string
	bnbPassword  string
	ethOne       string
	ethUsername  string
	ethPassword  string
	tronGridUrl  string
	bscScanUrl   string
	etherscanUrl string
)

type TronResponse struct {
	BlockID      string        `json:"blockID"`
	BlockHeader  BlockHeader   `json:"block_header"`
	Transactions []Transaction `json:"transactions"`
}

type BlockHeader struct {
	Data             RawData `json:"raw_data"`
	WitnessSignature string  `json:"witness_signature"`
}

type RawData struct {
	Number         int    `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int    `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type Transaction struct{}

type EthBasedResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func scrapeEtherscan(url string, c *http.Client) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error requesting EtherScan: %s\n", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("Error getting reqponse:", err)
	}
	defer resp.Body.Close()

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error creating new goquery")
	}

	if resp.StatusCode == 403 {
		log.Printf("Could not make request to EtherScan. Getting HTTP response %d\n", resp.StatusCode)
	}

	s := dom.Find(".text-size-1").Last().Text()
	return s

}

func scrapeBscScan(url string, c *http.Client) string {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Error creating new request to BSC", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("Error getting response from BSC:", err)
	}

	defer resp.Body.Close()

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("error creating new goquery")
	}

	if resp.StatusCode == 403 {
		log.Fatal("Could not make request to bscscan.com. Getting HTTP response 403 from Cloudflare")
	}

	// fmt.Println("Response Status Code:", resp.StatusCode) // for testing, cause sometime getting response 403 - Cloudflare

	bscLatestBlock := dom.Find("#lastblock").Text()
	return bscLatestBlock

}

func queryTronGrid(url string, c *http.Client) TronResponse {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal("Error querying TronGrid:", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("Error sending request to TronGrid:", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body from TronGrid:", err)
	}

	var response TronResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Error unmarshalling json response from TronGrid:", err)
	}

	return response
}

func queryTron(tronUrl, tronUser, tronPass string, c *http.Client) TronResponse {
	req, err := http.NewRequest(http.MethodPost, tronUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Encoded to Base64 username and password
	tronEncodedAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", tronUser, tronPass)))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", tronEncodedAuth))

	resp, err := c.Do(req)

	if err != nil {
		log.Fatal("Error sending request to hosted TRX fullnode:", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading response body from", tronUrl, " ERR:", err)
	}

	var response TronResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Error unmarshalling json response from", tronUrl, " ERR:", err)
	}

	return response

}

func queryEthBased(url, user, pass string, c *http.Client) EthBasedResponse {
	body := []byte(`{
		"method": "eth_blockNumber",
		"params": [],
		"id": 1,
		"jsonrpc": "2.0"
	}`)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	// Encoded to Base64 username and password
	bnbEncodedAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", bnbEncodedAuth))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)

	if err != nil {
		log.Fatal("Error getting response from", url, " ERR:", err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading response body from", url, " ERR:", err)
	}

	var response EthBasedResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Error unmarshalling json response from", url, " ERR:", err)
	}

	return response
}

func HexToInt64(resultStr string) int64 {
	res, err := strconv.ParseInt(resultStr, 0, 64)
	if err != nil {
		log.Fatal("Error Convering Hex to Int64:", err)
	}
	return res
}

func formatDecimalToString(num interface{}) string {
	p := message.NewPrinter(language.English)

	if s, ok := num.(string); ok {
		i, _ := strconv.ParseInt(s, 10, 0)
		return p.Sprintf("%d\n", i)
	} else {
		switch num.(type) {
		case int:
			return p.Sprintf("%d\n", num)
		case int64:
			return p.Sprintf("%d\n", num)
		case float64:
			return p.Sprintf("%f\n", num)
		default:
			return p.Sprintf("%v\n", num)
		}
	}

}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Unable to load .env file.")
		return
	}
}

func main() {

	c := &http.Client{}
	c.Timeout = time.Duration(5) * time.Second

	// get sensitive info from .env
	tronOne = os.Getenv("TRONONE")
	tronTwo = os.Getenv("TRONTWO")
	tronUsername = os.Getenv("TRONUSER")
	tronPassword = os.Getenv("TRONPASS")
	bnbOne = os.Getenv("BNBONE")
	bnbUsername = os.Getenv("BNBUSER")
	bnbPassword = os.Getenv("BNBPASS")
	ethOne = os.Getenv("ETHONE")
	ethUsername = os.Getenv("ETHUSER")
	ethPassword = os.Getenv("ETHPASS")
	tronGridUrl = os.Getenv("TRONGRID")
	bscScanUrl = os.Getenv("BSCSCAN")
	etherscanUrl = os.Getenv("ETHERSCAN")

	// fmt.Printf("tronOne = %s \t tronTwo = %s\n", tronOne, tronTwo)
	// fmt.Printf("tronPassword = %s\n", tronPassword)
	// fmt.Printf("bnbPassword = %s\n", bnbPassword)
	// fmt.Printf("ethPassword = %s\n", ethPassword)
	// fmt.Printf("trongrid url = %s\n", tronGridUrl)

	respTronOne := queryTron(tronOne, tronUsername, tronPassword, c)
	respTronTwo := queryTron(tronTwo, tronUsername, tronPassword, c)
	respTronGrid := queryTronGrid(tronGridUrl, c)
	respBnbOne := queryEthBased(bnbOne, bnbUsername, bnbPassword, c)
	respEthOne := queryEthBased(ethOne, ethUsername, ethPassword, c)
	bscLatestBlock := scrapeBscScan(bscScanUrl, c)
	bscLastBlock, _ := strconv.ParseInt(bscLatestBlock, 10, 0)
	etherscanLatestBlock := scrapeEtherscan(etherscanUrl, c)
	ethLastBlock, _ := strconv.ParseInt(etherscanLatestBlock, 10, 0)

	tronGridBlockNum := respTronGrid.BlockHeader.Data.Number
	fmt.Print("TRONGRID: ")
	color.HiGreen(fmt.Sprint(formatDecimalToString(tronGridBlockNum)))

	fmt.Print("BscScan: ")
	color.HiGreen(fmt.Sprint(formatDecimalToString(bscLastBlock)))
	fmt.Print("EtherScan: ")
	color.HiGreen(fmt.Sprint(formatDecimalToString(etherscanLatestBlock)))
	color.Green("----------------------")

	color.Magenta("::: Tron01 :::")
	tronOneBlockNum := respTronOne.BlockHeader.Data.Number
	tronOneDiff := tronGridBlockNum - tronOneBlockNum
	fmt.Print("Latest Block: ")
	color.HiMagenta(fmt.Sprintf(formatDecimalToString(tronOneBlockNum)))
	fmt.Print("Difference with TronGrid: ")
	color.HiMagenta(fmt.Sprint(formatDecimalToString(tronOneDiff)))

	color.Red("\n::: Tron02 :::")
	tronTwoBlockNum := respTronTwo.BlockHeader.Data.Number
	tronTwoDiff := tronGridBlockNum - tronTwoBlockNum
	fmt.Print("Latest Block: ")
	color.HiRed(fmt.Sprintf(formatDecimalToString(tronTwoBlockNum)))
	fmt.Print("Difference with TronGrid: ")
	color.HiRed(fmt.Sprint(tronTwoDiff))

	color.Yellow("\n::: BNB01 :::")

	bnbBlockNumber := HexToInt64(respBnbOne.Result)
	bnbOneDiff := bscLastBlock - bnbBlockNumber
	fmt.Print("Latest Block: ")
	color.HiYellow(fmt.Sprintf(formatDecimalToString(bnbBlockNumber)))
	fmt.Print("Difference with BscScan: ")
	color.HiYellow(fmt.Sprint(formatDecimalToString(int(bnbOneDiff))))

	color.Cyan("\n::: ETH01 :::")
	ethOneBlockNumber := HexToInt64(respEthOne.Result)
	ethOneDiff := ethLastBlock - ethOneBlockNumber
	fmt.Print("Latest Block: ")
	color.HiCyan(fmt.Sprint(formatDecimalToString(ethOneBlockNumber)))
	fmt.Print("Difference with EtherScan: ")
	color.HiCyan(fmt.Sprint(formatDecimalToString(ethOneDiff)))

	scrapeEtherscan(etherscanUrl, c)

}
