package main

import (
	"fmt"
)

func hzToBand(frequency int) string {
	switch {
	case frequency >= 1356700 && frequency <= 1378000:
		return "2190m"
	case frequency >= 4720000 && frequency <= 4790000:
		return "630m"
	case frequency >= 501000 && frequency <= 504000:
		return "560m"
	case frequency >= 1800000 && frequency <= 2000000:
		return "160m"
	case frequency >= 3500000 && frequency <= 4000000:
		return "80m"
	case frequency >= 5060000 && frequency <= 5450000:
		return "60m"
	case frequency >= 7000000 && frequency <= 7300000:
		return "40m"
	case frequency >= 10100000 && frequency <= 10150000:
		return "30m"
	case frequency >= 14000000 && frequency <= 14350000:
		return "20m"
	case frequency >= 18068000 && frequency <= 18168000:
		return "17m"
	case frequency >= 21000000 && frequency <= 21400000:
		return "15m"
	case frequency >= 24890000 && frequency <= 24990000:
		return "12m"
	case frequency >= 28000000 && frequency <= 29700000:
		return "10m"
	case frequency >= 40000000 && frequency <= 45000000:
		return "8m"
	case frequency >= 50000000 && frequency <= 54000000:
		return "6m"
	case frequency >= 54000001 && frequency <= 6990000:
		return "5m"
	case frequency >= 7000000 && frequency <= 7100000:
		return "4m"
	case frequency >= 144000000 && frequency <= 148000000:
		return "2m"
	case frequency >= 222000000 && frequency <= 225000000:
		return "1.25m"
	case frequency >= 430000000 && frequency <= 450000000:
		return "70cm"
	case frequency >= 902000000 && frequency <= 928000000:
		return "33cm"
	case frequency >= 1240000000 && frequency <= 1300000000:
		return "23cm"
	case frequency >= 2300000000 && frequency <= 2450000000:
		return "13cm"
	case frequency >= 3300000000 && frequency <= 3500000000:
		return "9cm"
	case frequency >= 5650000000 && frequency <= 5925000000:
		return "6cm"
	case frequency >= 10000000000 && frequency <= 10500000000:
		return "3cm"
	case frequency >= 2400000000 && frequency <= 2425000000:
		return "1.25cm"
	case frequency >= 4700000000 && frequency <= 4720000000:
		return "6mm"
	case frequency >= 7550000000 && frequency <= 8100000000:
		return "4mm"
	case frequency >= 1199800000 && frequency <= 1230000000:
		return "2.5mm"
	case frequency >= 1340000000 && frequency <= 1490000000:
		return "2mm"
	case frequency >= 2410000000 && frequency <= 2500000000:
		return "1mm"
	case frequency >= 300000000 && frequency <= 7500000000:
		return "submm"
	default:
		return "unknown"
	}
}

func renderHz(hz int) string {
	if hz < 1000 {
		return fmt.Sprintf("%d Hz", hz)
	} else if hz < 1000000 {
		return fmt.Sprintf("%.3f KHz", float64(hz)/1000)
	} else {
		return fmt.Sprintf("%.3f MHz", float64(hz)/1000000)
	}
}
