package main

import (
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 初期化確認
func validateInitialize() {
	voteSet := setupVotes(150, true)
	validateVote(voteSet)
	validateVoteError(voteSet)
	validateIndex(voteSet)
	validateCandidate(voteSet)
	validatePoliticalParty(voteSet)
}

func validateVote(voteSet []Vote) {
	for _, v := range voteSet {
		params := url.Values{}
		params.Add("name", v.Name)
		params.Add("address", v.Address)
		params.Add("mynumber", v.Mynumber)
		params.Add("candidate", v.Candidate)
		params.Add("keyword", v.Keyword)
		params.Add("vote_count", v.VoteCount)

		doc := httpRequestDoc("POST", "/vote", params)

		// 投票が成功したことの確認
		message := doc.Find(".text-danger").Text()
		if !strings.Contains(message, "投票に成功しました") {
			log.Print("正しい情報で投票ができません at POST /vote")
			os.Exit(1)
		}

		// DOMの構造確認
		if doc.Find("fieldset").Children().Size() != 14 {
			log.Print("DOM の構造が正しくありません at POST /vote")
			os.Exit(1)
		}
	}
}

func validateVoteError(voteSet []Vote) {
	// Case1: 個人情報に誤りがある場合
	v1 := voteSet[0]

	params := url.Values{}
	params.Add("name", "hoge")
	params.Add("address", v1.Address)
	params.Add("mynumber", v1.Mynumber)
	params.Add("candidate", v1.Candidate)
	params.Add("keyword", v1.Keyword)
	params.Add("vote_count", "0")

	doc := httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message := doc.Find(".text-danger").Text()
	if !strings.Contains(message, "個人情報に誤りがあります") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case2: 個人情報に誤りがある場合
	v2 := voteSet[1]

	params = url.Values{}
	params.Add("name", v2.Name)
	params.Add("address", "hoge")
	params.Add("mynumber", v2.Mynumber)
	params.Add("candidate", v2.Candidate)
	params.Add("keyword", v2.Keyword)
	params.Add("vote_count", "0")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "個人情報に誤りがあります") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case3: 個人情報に誤りがある場合
	v3 := voteSet[2]

	params = url.Values{}
	params.Add("name", v3.Name)
	params.Add("address", v3.Address)
	params.Add("mynumber", "1")
	params.Add("candidate", v3.Candidate)
	params.Add("keyword", v3.Keyword)
	params.Add("vote_count", "0")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "個人情報に誤りがあります") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case4: 投票数が上限を超えている場合
	v4 := voteSet[3]

	params = url.Values{}
	params.Add("name", v4.Name)
	params.Add("address", v4.Address)
	params.Add("mynumber", v4.Mynumber)
	params.Add("candidate", v4.Candidate)
	params.Add("keyword", v4.Keyword)
	params.Add("vote_count", "220")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "投票数が上限を超えています") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case5: 候補者が未記入の場合
	v5 := voteSet[4]

	params = url.Values{}
	params.Add("name", v5.Name)
	params.Add("address", v5.Address)
	params.Add("mynumber", v5.Mynumber)
	params.Add("candidate", "")
	params.Add("keyword", v5.Keyword)
	params.Add("vote_count", "0")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "候補者を記入してください") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case6: 候補者名が誤りの場合
	v6 := voteSet[5]

	params = url.Values{}
	params.Add("name", v6.Name)
	params.Add("address", v6.Address)
	params.Add("mynumber", v6.Mynumber)
	params.Add("candidate", "hoge")
	params.Add("keyword", v6.Keyword)
	params.Add("vote_count", "0")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "候補者を正しく記入してください") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}

	// Case7: 投票理由が空の場合
	v7 := voteSet[6]

	params = url.Values{}
	params.Add("name", v7.Name)
	params.Add("address", v7.Address)
	params.Add("mynumber", v7.Mynumber)
	params.Add("candidate", v7.Candidate)
	params.Add("keyword", "")
	params.Add("vote_count", "0")

	doc = httpRequestDoc("POST", "/vote", params)

	// 投票が成功したことの確認
	message = doc.Find(".text-danger").Text()
	if !strings.Contains(message, "投票理由を記入してください") {
		log.Print("エラーメッセージに誤りがあります at POST /vote")
		os.Exit(1)
	}
}

func validateIndex(voteSet []Vote) {
	doc := httpRequestDoc("GET", "/", nil)

	// DOM の確認
	ppErrFlg := doc.Find("#people").Children().Size() != 11
	ptErrFlg := doc.Find("#parties").Children().Size() != 4
	sxErrFlg := doc.Find("#sex_ratio").Children().Size() != 2
	if ppErrFlg || ptErrFlg || sxErrFlg {
		log.Print("DOMの構造が正しくありません at GET /index")
		os.Exit(1)
	}

	// 個人の部の結果確認
	rank := map[string]int{}
	for _, v := range voteSet {
		cnt, _ := strconv.Atoi(v.VoteCount)
		rank[v.Candidate] = rank[v.Candidate] + cnt
	}
	l := List{}
	for k, v := range rank {
		e := Entry{k, v}
		l = append(l, e)
	}
	sort.Sort(l)

	doc.Find("#people").Children().Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			str := s.Text()
			// 得票数が同じ場合に順位がズレる可能性がある
			cand1 := "INIT STRING"
			if i > 0 {
				cand1 = strconv.Itoa(i+1) + ". " + l[l.Len()-i].name
			}
			cand2 := strconv.Itoa(i+1) + ". " + l[l.Len()-1-i].name
			cand3 := strconv.Itoa(i+1) + ". " + l[l.Len()-2-i].name
			if !strings.Contains(str, cand1) && !strings.Contains(str, cand2) && !strings.Contains(str, cand3) {
				log.Print("個人の部の選挙結果が正しくありません at GET /")
				os.Exit(1)
			}
		}
	})

	// 政党の部の結果確認
	partyRank := map[string]int{}
	for _, v := range voteSet {
		party := getCndInfo(v.Candidate).Party
		cnt, _ := strconv.Atoi(v.VoteCount)
		partyRank[party] = partyRank[party] + cnt
	}
	l = List{}
	for k, v := range partyRank {
		e := Entry{k, v}
		l = append(l, e)
	}
	sort.Sort(l)

	doc.Find("#parties").Children().Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			str := s.Text()
			// 得票数が同じ場合に順位がズレる可能性がある
			cand1 := "INIT STRING"
			if i > 0 {
				cand1 = strconv.Itoa(i+1) + ". " + l[l.Len()-i].name
			}
			cand2 := strconv.Itoa(i+1) + ". " + l[l.Len()-1-i].name
			cand3 := strconv.Itoa(i+1) + ". " + l[l.Len()-2-i].name
			if !strings.Contains(str, cand1) && !strings.Contains(str, cand2) && !strings.Contains(str, cand3) {
				log.Print("政党の部の選挙結果が正しくありません at GET /")
				os.Exit(1)
			}
		}
	})

	// 男女比率の結果確認
	sexRatio := map[string]int{}
	for _, v := range voteSet {
		sex := getCndInfo(v.Candidate).Sex
		cnt, _ := strconv.Atoi(v.VoteCount)
		sexRatio[sex] = sexRatio[sex] + cnt
	}

	man := strconv.Itoa(sexRatio["男"])
	women := strconv.Itoa(sexRatio["女"])
	doc.Find("#sex_ratio").Children().Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			str := s.Text()
			if !strings.Contains(str, man) && !strings.Contains(str, women) {
				log.Println(str)
				log.Println("man:" + man)
				log.Println("woman:" + women)
				log.Print("男女比率の選挙結果が正しくありません at GET /")
				os.Exit(1)
			}
		}
	})
}

func validateCandidate(voteSet []Vote) {
	rank := map[string]int{}
	for _, v := range voteSet {
		cnt, _ := strconv.Atoi(v.VoteCount)
		rank[v.Candidate] = rank[v.Candidate] + cnt
	}
	l := List{}
	for k, v := range rank {
		e := Entry{k, v}
		l = append(l, e)
	}
	sort.Sort(l)

	// 上位2人の個人ページを確認する
	for i, cnd := range l {
		if i >= l.Len()-2 {
			cndInfo := getCndInfo(cnd.name)
			doc := httpRequestDoc("GET", "/candidates/"+cndInfo.ID, nil)
			doc.Find("#info p").Each(func(i int, s *goquery.Selection) {
				str := s.Text()
				if i == 0 {
					// 得票数の確認
					if !strings.Contains(str, strconv.Itoa(cnd.value)) {
						log.Print("得票数の情報が正しくありません at GET /candidates/:id")
						os.Exit(1)
					}
				} else if i == 1 {
					// 政党名の確認
					if !strings.Contains(str, cndInfo.Party) {
						log.Print("政党の情報が正しくありません at GET /candidates/:id")
						os.Exit(1)
					}
				} else if i == 2 {
					// 性別の確認
					if !strings.Contains(str, cndInfo.Sex) {
						log.Print("性別の情報が正しくありません at GET /candidates/:id")
						os.Exit(1)
					}
				}
			})

			// キーワードの確認
			keyRank := map[string]int{}
			for _, v := range voteSet {
				if v.Candidate == cnd.name {
					cnt, _ := strconv.Atoi(v.VoteCount)
					keyRank[v.Keyword] = keyRank[v.Keyword] + cnt
				}
			}
			keyList := List{}
			for k, v := range keyRank {
				e := Entry{k, v}
				keyList = append(keyList, e)
			}
			sort.Sort(keyList)

			doc.Find("#info ul").Children().Each(func(i int, s *goquery.Selection) {
				if i < 2 {
					str := s.Text()
					// 得票数の確認
					key1 := "INIT STRING"
					if i != 0 {
						key1 = keyList[keyList.Len()-i].name
					}
					key2 := keyList[keyList.Len()-1-i].name
					key3 := keyList[keyList.Len()-2-i].name
					key4 := keyList[keyList.Len()-3-i].name
					if !strings.Contains(str, key1) && !strings.Contains(str, key2) && !strings.Contains(str, key3) && !strings.Contains(str, key4) {
						log.Print("支持者の声が正しくありません at GET /candidates/:id")
						os.Exit(1)
					}
				}
			})
		}
	}
}

func validatePoliticalParty(voteSet []Vote) {
	doc := httpRequestDoc("GET", "/political_parties/国民元気党", nil)

	var votes int
	keyRank := map[string]int{}
	for _, v := range voteSet {
		if getPatryInfo(v.Candidate) == "国民元気党" {
			cnt, _ := strconv.Atoi(v.VoteCount)
			votes = votes + cnt
			keyRank[v.Keyword] = keyRank[v.Keyword] + cnt
		}
	}

	// 得票数の確認
	docVotesTxt := doc.Find("#votes").Text()
	docVotes, _ := strconv.Atoi(docVotesTxt)
	if docVotes != votes {
		log.Print("得票数が正しくありません at GET /political_parties/:name")
		os.Exit(1)
	}

	// 党員の確認
	memberSet := membersOf("国民元気党")
	doc.Find("#members").Children().Each(func(i int, s *goquery.Selection) {
		flg := false
		str := s.Text()
		for _, member := range memberSet {
			if strings.Contains(str, member) {
				flg = true
			}
		}
		if !flg {
			log.Print("候補者が正しくありません at GET /political_parties/:name")
			os.Exit(1)
		}
	})

	// 支持者の声の確認
	keyList := List{}
	for k, v := range keyRank {
		e := Entry{k, v}
		keyList = append(keyList, e)
	}
	sort.Sort(keyList)
	doc.Find("#voices").Children().Each(func(i int, s *goquery.Selection) {
		if i < 2 {
			str := s.Text()
			key1 := "INIT STRING"
			if i != 0 {
				key1 = keyList[keyList.Len()-i].name
			}
			key2 := keyList[keyList.Len()-1-i].name
			key3 := keyList[keyList.Len()-2-i].name
			key4 := keyList[keyList.Len()-3-i].name
			if !strings.Contains(str, key1) && !strings.Contains(str, key2) && !strings.Contains(str, key3) && !strings.Contains(str, key4) {
				log.Print("支持者の声が正しくありません at GET /political_parties/:name")
				os.Exit(1)
			}
		}
	})
}

// follows for sort

// Entry for sort
type Entry struct {
	name  string
	value int
}

// List for sort
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	}
	return (l[i].value < l[j].value)
}
