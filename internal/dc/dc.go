package dc

type DCConfig struct {
	Accounts   string
	Books      string
	Cliq       string
	CRM        string
	Expense    string
	Inventory  string
	Mail       string
	Projects   string
	Sheet      string
	WorkDrive  string
	Writer     string
	Download   string
}

var dcMap = map[string]DCConfig{
	"com": {
		Accounts:   "https://accounts.zoho.com",
		Books:      "https://www.zohoapis.com",
		Cliq:       "https://cliq.zoho.com",
		CRM:        "https://zohoapis.com",
		Expense:    "https://www.zohoapis.com",
		Inventory:  "https://www.zohoapis.com",
		Mail:       "https://mail.zoho.com",
		Projects:   "https://projectsapi.zoho.com",
		Sheet:      "https://sheet.zoho.com",
		WorkDrive:  "https://workdrive.zoho.com",
		Writer:     "https://www.zohoapis.com/writer",
		Download:   "https://download.zoho.com",
	},
	"eu": {
		Accounts:   "https://accounts.zoho.eu",
		Books:      "https://www.zohoapis.eu",
		Cliq:       "https://cliq.zoho.eu",
		CRM:        "https://zohoapis.eu",
		Expense:    "https://www.zohoapis.eu",
		Inventory:  "https://www.zohoapis.eu",
		Mail:       "https://mail.zoho.eu",
		Projects:   "https://projectsapi.zoho.eu",
		Sheet:      "https://sheet.zoho.eu",
		WorkDrive:  "https://workdrive.zoho.eu",
		Writer:     "https://www.zohoapis.eu/writer",
		Download:   "https://download.zoho.eu",
	},
	"in": {
		Accounts:   "https://accounts.zoho.in",
		Books:      "https://www.zohoapis.in",
		Cliq:       "https://cliq.zoho.in",
		CRM:        "https://zohoapis.in",
		Expense:    "https://www.zohoapis.in",
		Inventory:  "https://www.zohoapis.in",
		Mail:       "https://mail.zoho.in",
		Projects:   "https://projectsapi.zoho.in",
		Sheet:      "https://sheet.zoho.in",
		WorkDrive:  "https://workdrive.zoho.in",
		Writer:     "https://www.zohoapis.in/writer",
		Download:   "https://download.zoho.in",
	},
	"com.au": {
		Accounts:   "https://accounts.zoho.com.au",
		Books:      "https://www.zohoapis.com.au",
		Cliq:       "https://cliq.zoho.com.au",
		CRM:        "https://zohoapis.com.au",
		Expense:    "https://www.zohoapis.com.au",
		Inventory:  "https://www.zohoapis.com.au",
		Mail:       "https://mail.zoho.com.au",
		Projects:   "https://projectsapi.zoho.com.au",
		Sheet:      "https://sheet.zoho.com.au",
		WorkDrive:  "https://workdrive.zoho.com.au",
		Writer:     "https://www.zohoapis.com.au/writer",
		Download:   "https://download.zoho.com.au",
	},
	"jp": {
		Accounts:   "https://accounts.zoho.jp",
		Books:      "https://www.zohoapis.jp",
		Cliq:       "https://cliq.zoho.jp",
		CRM:        "https://zohoapis.jp",
		Expense:    "https://www.zohoapis.jp",
		Inventory:  "https://www.zohoapis.jp",
		Mail:       "https://mail.zoho.jp",
		Projects:   "https://projectsapi.zoho.jp",
		Sheet:      "https://sheet.zoho.jp",
		WorkDrive:  "https://workdrive.zoho.jp",
		Writer:     "https://www.zohoapis.jp/writer",
		Download:   "https://download.zoho.jp",
	},
	"ca": {
		Accounts:   "https://accounts.zohocloud.ca",
		Books:      "https://www.zohoapis.ca",
		Cliq:       "https://cliq.zohocloud.ca",
		CRM:        "https://zohoapis.ca",
		Expense:    "https://www.zohoapis.ca",
		Inventory:  "https://www.zohoapis.ca",
		Mail:       "https://mail.zohocloud.ca",
		Projects:   "https://projectsapi.zohocloud.ca",
		Sheet:      "https://sheet.zohocloud.ca",
		WorkDrive:  "https://workdrive.zohocloud.ca",
		Writer:     "https://www.zohoapis.ca/writer",
		Download:   "https://download.zohocloud.ca",
	},
	"sa": {
		Accounts:   "https://accounts.zoho.sa",
		Books:      "https://www.zohoapis.sa",
		Cliq:       "https://cliq.zoho.sa",
		CRM:        "https://zohoapis.sa",
		Expense:    "https://www.zohoapis.sa",
		Inventory:  "https://www.zohoapis.sa",
		Mail:       "https://mail.zoho.sa",
		Projects:   "https://projectsapi.zoho.sa",
		Sheet:      "https://sheet.zoho.sa",
		WorkDrive:  "https://workdrive.zoho.sa",
		Writer:     "https://www.zohoapis.sa/writer",
		Download:   "https://download.zoho.sa",
	},
	"uk": {
		Accounts:   "https://accounts.zoho.uk",
		Books:      "https://www.zohoapis.uk",
		Cliq:       "https://cliq.zoho.uk",
		CRM:        "https://zohoapis.uk",
		Expense:    "https://www.zohoapis.uk",
		Inventory:  "https://www.zohoapis.uk",
		Mail:       "https://mail.zoho.uk",
		Projects:   "https://projectsapi.zoho.uk",
		Sheet:      "https://sheet.zoho.uk",
		WorkDrive:  "https://workdrive.zoho.uk",
		Writer:     "https://www.zohoapis.uk/writer",
		Download:   "https://download.zoho.uk",
	},
	"com.cn": {
		Accounts:   "https://accounts.zoho.com.cn",
		Books:      "https://www.zohoapis.com.cn",
		Cliq:       "https://cliq.zoho.com.cn",
		CRM:        "https://zohoapis.com.cn",
		Expense:    "https://www.zohoapis.com.cn",
		Inventory:  "https://www.zohoapis.com.cn",
		Mail:       "https://mail.zoho.com.cn",
		Projects:   "https://projectsapi.zoho.com.cn",
		Sheet:      "https://sheet.zoho.com.cn",
		WorkDrive:  "https://workdrive.zoho.com.cn",
		Writer:     "https://www.zohoapis.com.cn/writer",
		Download:   "https://download.zoho.com.cn",
	},
}

var ValidDCs = []string{"com", "eu", "in", "com.au", "jp", "ca", "sa", "uk", "com.cn"}

func GetDC(dc string) DCConfig {
	if cfg, ok := dcMap[dc]; ok {
		return cfg
	}
	return dcMap["com"]
}

func AccountsURL(dc string) string    { return GetDC(dc).Accounts }
func BooksURL(dc string) string       { return GetDC(dc).Books }
func CliqURL(dc string) string        { return GetDC(dc).Cliq }
func CRMURL(dc string) string         { return GetDC(dc).CRM }
func ExpenseURL(dc string) string     { return GetDC(dc).Expense }
func InventoryURL(dc string) string   { return GetDC(dc).Inventory }
func MailURL(dc string) string        { return GetDC(dc).Mail }
func ProjectsURL(dc string) string    { return GetDC(dc).Projects }
func SheetURL(dc string) string       { return GetDC(dc).Sheet }
func WorkDriveURL(dc string) string   { return GetDC(dc).WorkDrive }
func WriterURL(dc string) string      { return GetDC(dc).Writer }
func DownloadURL(dc string) string    { return GetDC(dc).Download }
