package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type WebSearchOptions struct {
	SearchContextSize string `json:"search_context_size"`
}

type Request struct {
	Model           string           `json:"model"`
	Messages        []Message        `json:"messages"`
	WebSearchOptions WebSearchOptions `json:"web_search_options"`
}

type Usage struct {
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	SearchContextSize string `json:"search_context_size"`
}

type Response struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	Created   int64  `json:"created"`
	Usage     Usage  `json:"usage"`
	Citations json.RawMessage `json:"citations"`
}

func callPerplexityAPI(question string, apiKey string) (*Response, error) {
	url := "https://api.perplexity.ai/chat/completions"
	
	// リクエストの作成
	req := Request{
		Model: "sonar-pro",
		Messages: []Message{
			{
				Role:    "system",
				Content: "Be precise and concise.",
			},
			{
				Role:    "user",
				Content: question,
			},
		},
		WebSearchOptions: WebSearchOptions{
			SearchContextSize: "medium",
		},
	}
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	
	// HTTPリクエストの作成
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	
	// リクエストの送信
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// レスポンスの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// JSONのパース
	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	
	return &apiResp, nil
}

// 質問をファイルから読み込む関数
func loadQuestionsFromFile(filePath string) ([]string, error) {
	// ファイルを開く
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みエラー: %w", err)
	}

	// 区切り文字「---」で質問を分割
	text := string(content)
	questions := strings.Split(text, "---")
	
	// 各質問の前後の空白を削除
	for i, q := range questions {
		questions[i] = strings.TrimSpace(q)
	}
	
	// 空の質問を除外
	var filteredQuestions []string
	for _, q := range questions {
		if q != "" {
			filteredQuestions = append(filteredQuestions, q)
		}
	}
	
	return filteredQuestions, nil
}

func main() {
	// コマンドライン引数の処理
	questionFilePath := flag.String("file", "questions.txt", "質問が記載されたファイルのパス")
	flag.Parse()

	apiKey := os.Getenv("SONARAPI_KEY")
	if apiKey == "" {
		fmt.Println("環境変数SONARAPI_KEYが設定されていません")
		os.Exit(1)
	}
	
	// 質問をファイルから読み込む
	questions, err := loadQuestionsFromFile(*questionFilePath)
	if err != nil {
		fmt.Printf("質問の読み込みエラー: %v\n", err)
		os.Exit(1)
	}
	
	if len(questions) == 0 {
		fmt.Println("質問が見つかりませんでした")
		os.Exit(1)
	}
	
	var totalTokens int
	completionTokens := make([]int, 0, len(questions))
	
	for _, q := range questions {
		fmt.Printf("質問: %s\n", q)
		resp, err := callPerplexityAPI(q, apiKey)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			continue
		}
		
		fmt.Printf("Completion Tokens: %d\n\n", resp.Usage.CompletionTokens)
		completionTokens = append(completionTokens, resp.Usage.CompletionTokens)
		totalTokens += resp.Usage.CompletionTokens
	}
	
	// 平均値の計算
	var average float64
	if len(completionTokens) > 0 {
		average = float64(totalTokens) / float64(len(completionTokens))
	}
	
	// 結果の表示
	fmt.Println("==== 結果 ====")
	for i, tokens := range completionTokens {
		fmt.Printf("質問%d: %d tokens\n", i+1, tokens)
	}
	fmt.Printf("\n平均 Completion Tokens: %.2f\n", average)
} 
