package model

// NavConfig 存储导航栏全局配置
type NavConfig struct {
    ID                uint   `json:"id" gorm:"primaryKey"`
    
    // Service & Support 下拉菜单
    ServiceTitle1     string `json:"service_title1"`
    ServiceDesc1      string `json:"service_desc1"`
    ServiceLink1      string `json:"service_link1"`
    
    ServiceTitle2     string `json:"service_title2"`
    ServiceDesc2      string `json:"service_desc2"`
    ServiceLink2      string `json:"service_link2"`

    ServiceTel        string `json:"service_tel"`
    ServiceEmail      string `json:"service_email"`

    // About 下拉菜单
    AboutTitle1       string `json:"about_title1"`
    AboutLink1        string `json:"about_link1"`
    AboutTitle2       string `json:"about_title2"`
    AboutLink2        string `json:"about_link2"`
    AboutTitle3       string `json:"about_title3"`
    AboutLink3        string `json:"about_link3"`
    AboutTitle4       string `json:"about_title4"`
    AboutLink4        string `json:"about_link4"`

    // 统计数据
    TickerSymbol      string `json:"ticker_symbol"` // 300990.SZ
    EmployeeCount     string `json:"employee_count"` // 2200+
    RdEngineerCount   string `json:"rd_engineer_count"` // 300+
}