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

	// There is conflicting information on whether or not this is a
	// requirement for POTA. It *should not be* because the implication is
	// that if I'm uploading a log, I'm probably the one who made the
	// contacts in that log. In any case, uncommenting this is a way to add
	// it if ever necessary. Next activation, I will try uploading without
	// it and see what happens. This is only really an issue because we have
	// no notion of a global log configuration right now, so no real place
	// to store this dynamically.
	//adif += "<OPERATOR:4>AI5A"

	tagAdifMap := map[string]string{
		"name":  "NAME",
		"pota":  "SIG_INFO",
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
