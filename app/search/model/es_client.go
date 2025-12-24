package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

// ESClient Elasticsearch 客户端封装
type ESClient struct {
	clients map[string]*elastic.Client
}

// NewESClient 创建 ES 客户端
func NewESClient(clusters map[string][]string) (*ESClient, error) {
	clients := make(map[string]*elastic.Client)

	for name, addresses := range clusters {
		client, err := elastic.NewClient(
			elastic.SetURL(addresses...),
			elastic.SetSniff(false),
		)
		if err != nil {
			logx.Errorf("failed to create ES client for cluster %s: %v", name, err)
			continue
		}
		clients[name] = client
		logx.Infof("ES cluster %s connected", name)
	}

	return &ESClient{clients: clients}, nil
}

// GetClient 获取指定集群的客户端
func (e *ESClient) GetClient(clusterName string) (*elastic.Client, bool) {
	client, ok := e.clients[clusterName]
	return client, ok
}

// Page 分页信息
type Page struct {
	Pn    int32 `json:"num"`
	Ps    int32 `json:"size"`
	Total int64 `json:"total"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Order  string            `json:"order"`
	Sort   string            `json:"sort"`
	Result []json.RawMessage `json:"result"`
	Page   *Page             `json:"page"`
	Debug  string            `json:"debug"`
}

// BasicSearchParams 基础搜索参数
type BasicSearchParams struct {
	AppID    string
	KW       string
	KwFields []string
	Order    []string
	Sort     []string
	Pn       int32
	Ps       int32
	Debug    bool
	Source   []string
}

// Search 执行搜索
func (e *ESClient) Search(ctx context.Context, clusterName, indexName string, query elastic.Query, params *BasicSearchParams) (*SearchResult, error) {
	result := &SearchResult{Debug: ""}

	// 调试模式
	if params.Debug {
		src, err := query.Source()
		if err == nil {
			data, _ := json.Marshal(src)
			result.Debug = string(data)
		}
	}

	// 获取客户端
	client, ok := e.clients[clusterName]
	if !ok {
		logx.Errorf("ES cluster not found: %s", clusterName)
		result.Debug = fmt.Sprintf("ES cluster not found: %s, %s", clusterName, result.Debug)
		return result, fmt.Errorf("ES cluster not found: %s", clusterName)
	}

	// 构建排序
	sorters := make([]elastic.Sorter, 0)
	if params.KW != "" {
		sorters = append(sorters, elastic.NewScoreSort().Desc())
	}
	for i, field := range params.Order {
		sortOrder := "desc"
		if i < len(params.Sort) {
			sortOrder = params.Sort[i]
		} else if len(params.Sort) > 0 {
			sortOrder = params.Sort[0]
		}

		if sortOrder == "desc" {
			sorters = append(sorters, elastic.NewFieldSort(field).Desc())
		} else {
			sorters = append(sorters, elastic.NewFieldSort(field).Asc())
		}
	}

	// 构建查询
	searchService := client.Search().
		Index(indexName).
		Query(query).
		From(int((params.Pn - 1) * params.Ps)).
		Size(int(params.Ps)).
		Pretty(true)

	// 添加排序
	for _, sorter := range sorters {
		searchService = searchService.SortBy(sorter)
	}

	// 指定返回字段
	if len(params.Source) > 0 {
		fsc := elastic.NewFetchSourceContext(true).Include(params.Source...)
		searchService = searchService.FetchSourceContext(fsc)
	}

	// 执行查询
	searchResp, err := searchService.Do(ctx)
	if err != nil {
		logx.Errorf("ES search failed: %v", err)
		result.Debug = result.Debug + " ES search failed"
		return result, err
	}

	// 解析结果
	var data []json.RawMessage
	for _, hit := range searchResp.Hits.Hits {
		var t json.RawMessage
		if err := json.Unmarshal(hit.Source, &t); err != nil {
			logx.Errorf("unmarshal hit source failed: %v", err)
			continue
		}
		data = append(data, t)
	}

	result.Order = strings.Join(params.Order, ",")
	result.Sort = strings.Join(params.Sort, ",")
	result.Result = data
	result.Page = &Page{
		Pn:    params.Pn,
		Ps:    params.Ps,
		Total: searchResp.Hits.TotalHits.Value,
	}

	return result, nil
}

// BulkUpdate 批量更新
func (e *ESClient) BulkUpdate(ctx context.Context, clusterName string, items []BulkUpdateItem) error {
	client, ok := e.clients[clusterName]
	if !ok {
		return fmt.Errorf("ES cluster not found: %s", clusterName)
	}

	bulkRequest := client.Bulk()
	for _, item := range items {
		request := elastic.NewBulkUpdateRequest().
			Index(item.IndexName).
			Id(item.IndexID).
			Doc(item.Fields).
			DocAsUpsert(true)
		bulkRequest.Add(request)
	}

	if bulkRequest.NumberOfActions() == 0 {
		return nil
	}

	_, err := bulkRequest.Do(ctx)
	if err != nil {
		logx.Errorf("ES bulk update failed: %v", err)
		return err
	}

	return nil
}

// BulkUpdateItem 批量更新项
type BulkUpdateItem struct {
	IndexName string
	IndexID   string
	Fields    map[string]interface{}
}

// Ping 健康检查
func (e *ESClient) Ping(ctx context.Context) error {
	for name, client := range e.clients {
		_, _, err := client.Ping(client.String()).Do(ctx)
		if err != nil {
			logx.Errorf("ES cluster %s ping failed: %v", name, err)
			return err
		}
	}
	return nil
}
