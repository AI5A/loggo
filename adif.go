package main

import (
	"fmt"
	"strconv"
)

type ADIFExportable interface {
	ToADIF() string
}

func (qso QSO) ToADIF() string {
	var adif string

	// QSO date/time will be local time, convert to UTC
	qsoDate := qso.CreatedAt.UTC().Format("20060102")
	qsoTime := qso.CreatedAt.UTC().Format("1504")
	qsoBand := hzToBand(qso.Frequency)
	qsoFreq := fmt.Sprintf("%.4f", float64(qso.Frequency)/1000000)
	qsoSent := fmt.Sprintf("%d", qso.Sent)
	qsoRcvd := fmt.Sprintf("%d", qso.Received)

	// Pull out interesting tags from comment field
	tags := commentToTags(qso.Comment)

	adif += "<CALL:" + strconv.Itoa(len(qso.Callsign)) + ">" + qso.Callsign
	adif += "<FREQ:" + strconv.Itoa(len(qsoFreq)) + ">" + qsoFreq
	adif += "<MODE:" + strconv.Itoa(len(qso.Mode)) + ">" + qso.Mode
	adif += "<RST_RCVD:" + strconv.Itoa(len(qsoRcvd)) + ">" + qsoRcvd
	adif += "<RST_SENT:" + strconv.Itoa(len(qsoSent)) + ">" + qsoSent
	adif += "<QSO_DATE:" + strconv.Itoa(len(qsoDate)) + ">" + qsoDate
	adif += "<TIME_ON:" + strconv.Itoa(len(qsoTime)) + ">" + qsoTime
	adif += "<BAND:" + strconv.Itoa(len(qsoBand)) + ">" + qsoBand

	tagAdifMap := map[string]string{
		"name":  "NAME",
		"pota":  "POTA_REF",
		"qth":   "QTH",
		"grid":  "GRIDSQUARE",
		"skcc":  "SKCC",
		"fists": "FISTS",
	}

	// Iterate the tag/ADIF map and add any tags that are present
	for tag, adifTag := range tagAdifMap {
		value, ok := tags[tag]
		if ok {
			adif += "<" + adifTag + ":" + strconv.Itoa(len(value)) + ">" + value
		}
	}

	adif += "<COMMENT:" + strconv.Itoa(len(qso.Comment)) + ">" + qso.Comment
	adif += "<EOR>"
	return adif
}
