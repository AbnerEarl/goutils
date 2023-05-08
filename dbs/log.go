package dbs

import "time"

type LogInfoModel struct {
	BaseModel
	Method        string    `json:"method" gorm:"column:method;null;comment:'请求方法'"`
	ContentLength uint64    `json:"content_length" gorm:"column:content_length;null;comment:'内容长度'"`
	ExecuteTime   time.Time `json:"execute_time" gorm:"column:execute_time;null;comment:'执行时间'"`
	ContentType   string    `json:"content_type" gorm:"column:content_type;null;comment:'内容类型'"`
	CostTime      uint64    `json:"cost_time" gorm:"column:cost_time;null;comment:'花费时间'"`
	RequestUrl    string    `json:"request_url" gorm:"column:request_url;null;comment:'请求URL'"`
	RequestHost   string    `json:"request_host" gorm:"column:request_host;null;comment:'请求主机'"`
	UserAgent     string    `json:"user_agent" gorm:"column:user_agent;null;comment:'请求头'"`
	RemoteIp      string    `json:"remote_ip" gorm:"column:remote_ip;null;comment:'远程IP'"`
	RemoteAddr    string    `json:"remote_addr" gorm:"column:remote_addr;null;comment:'远程地址'"`
	ApiPath       string    `json:"api_path" gorm:"column:api_path;null;comment:'API路径'"`
	Referer       string    `json:"referer" gorm:"column:referer;null;comment:'网页关联'"`
	ApiDesc       string    `json:"api_desc" gorm:"column:api_desc;null;comment:'API功能'"`
	StatusCode    int       `json:"status_code" gorm:"column:status_code;null;comment:'状态码'"`
	ResponseData  string    `json:"response_data" gorm:"column:response_data;null;type:text;comment:'响应数据'"`
	AccountInfo   string    `json:"account_info" gorm:"column:account_info;null;type:text;comment:'用户信息'"`
	AccountName   string    `json:"account_name" gorm:"column:account_name;null;comment:'用户名称'"`
	AccountId     uint64    `json:"account_id" gorm:"column:account_id;null;comment:'用户ID'"`
	RequestParams string    `json:"request_params" gorm:"column:request_params;null;type:text;comment:'请求参数'"`
	RequestToken  string    `json:"request_token" gorm:"column:request_token;null;type:text;comment:'请求Token'"`
}

func (c *LogInfoModel) TableName() string {
	return "log_info"
}
