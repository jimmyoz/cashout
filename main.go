package main

import (
	"cashout/config"
	_type "cashout/type"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	c := _type.Checkbook{}
	b := _type.Balance{}

	urlCashout := config.Url + config.Cashout
	urlWithDraw := config.Url + config.Withdraw
	urlBal := config.Url + config.Balance
	checkUrl := config.Url + config.Checkbook

	for {
		err := GetCheckbook(checkUrl, &c)
		//err := testGetCheckBook(&c)
		if err != nil {
			log.Printf("get checkbook failed with err: %v", err)
			time.Sleep(time.Hour * 6)
			continue
		}
		if has := hasCheckToCashout(&c); !has {
			time.Sleep(time.Hour * 6)
			continue
		} else {
			CashoutAllToBalance(urlCashout, &c)
			time.Sleep(time.Second * 120)
			suc, err := WithdrawBalance(urlWithDraw, urlBal, &b)
			if err != nil || !suc {
				log.Printf("withdraw failed with error: %v", err)
				time.Sleep(time.Hour * 6)
				continue
			} else {
				log.Printf("withdraw successfully")
				time.Sleep(time.Hour * 6)
				continue
			}
		}
	}
}

func GetCheckbook(url string, checkbook *_type.Checkbook) error {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("get checkbook from url %v failed with error: %v", url, err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &checkbook)
	if err != nil {
		return err
	}

	return nil
}

func CashoutAllToBalance(url string, checkbook *_type.Checkbook) {
	for _, lastCheck := range checkbook.Lastcheques {
		if lastCheck.Lastreceived.Payout == 0 {
			continue
		} else {
			postUrl := url + lastCheck.Peer
			err := postCashout(postUrl)
			if err != nil {
				continue
			} else {
				log.Printf("cashout with peer %v successfully", postUrl)
			}
		}
	}
}

func WithdrawBalance(urlBalance, urlBal string, bal *_type.Balance) (bool, error) {
	amt, err := getBalance(urlBal, bal)
	//amt, err := testGetBalance(bal)
	if err != nil {
		return false, err
	}
	if amt == 0 {
		return false, errors.New("there is no available balance to withdraw")
	}

	urlWithdraw := urlBalance + strconv.Itoa(int(amt))

	err = withdraw(urlWithdraw)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getBalance(url string, bal *_type.Balance) (int64, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("get balance from url %v failed with error: %v", url, err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &bal)

	if err != nil {
		return 0, err
	}

	log.Printf("ava amount: %d", bal.AvailableBalance)

	return int64(bal.AvailableBalance), nil
}

func withdraw(url string) error {
	res, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		log.Printf("failed to withdraw, an error occurred: %v", err)
		return err
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	log.Printf("withdraw trans: %v", string(result))

	return nil
}

func postCashout(url string) error {
	res, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		log.Printf("failed to cashout with peer %v, an error occurred: %v", url, err)
		return err
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	log.Printf("cashout trans: %v", string(result))
	return nil
}

func hasCheckToCashout(checkbook *_type.Checkbook) bool {
	if len(checkbook.Lastcheques) == 0 {
		return false
	}

	for _, lastCheck := range checkbook.Lastcheques {
		if lastCheck.Lastreceived.Payout != 0 {
			return true
		}
	}

	return false
}

func testGetCheckBook(checkbook *_type.Checkbook) error {
	jsonFile, err := os.Open("test.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	check, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(check, &checkbook)

	return err
}

func testGetBalance(bal *_type.Balance) (int64, error) {
	jsonFile, err := os.Open("testBal.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	check, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(check, &bal)

	if err != nil {
		return 0, err
	}

	return int64(bal.AvailableBalance), nil
}
