package utils

func GetMasterDataCompare(provinceID, districtID, wardID int) int64 {
	return int64((provinceID*1000*2+districtID)*100000*2 + wardID)
}
