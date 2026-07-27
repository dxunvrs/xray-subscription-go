package dto

type UserTraffic struct {
	UplinkFormatted   string `json:"uplinkFormatted"`
	DownlinkFormatted string `json:"downlinkFormatted"`
	TotalFormatted    string `json:"totalFormatted"`
	UplinkBytes       int64  `json:"uplinkBytes"`
	DownlinkBytes     int64  `json:"downlinkBytes"`
	TotalBytes        int64  `json:"totalBytes"`
}

type UserResponse struct {
	ID      uint        `json:"id"`
	Email   string      `json:"email"`
	UUID    string      `json:"uuid"`
	Traffic UserTraffic `json:"traffic"`
}
