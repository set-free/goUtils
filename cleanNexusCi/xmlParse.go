package main

import (
	"encoding/xml"
	"log"
	"time"
)

const (
	dataFmt        string = "2006-01-02 15:4:5.999 UTC"
	defaultDataFmt string = "2006-01-02 15:04:05"
)

type OldArchives struct {
	date time.Time
	url  string
}

type XmlContent struct {
	XMLName xml.Name `xml:"content"`
	Data    XmlData  `xml:"data"`
}

type XmlData struct {
	XMLName     xml.Name `xml:"data"`
	ContentItem []XmlRrm `xml:"content-item"`
}

type XmlRrm struct {
	XMLName      xml.Name `xml:"content-item"`
	URI          string   `xml:"resourceURI"`
	Path         string   `xml:"relativePath"`
	Text         string   `xml:"text"`
	Leaf         string   `xml:"leaf"`
	LastModified string   `xml:"lastModified"`
	SizeOnDisk   string   `xml:"sizeOnDisk"`
}

func ParseXml(body []byte) XmlContent {
	var xmlNexus XmlContent
	err := xml.Unmarshal(body, &xmlNexus)
	if err != nil {
		log.Println(err)
		log.Fatal(err.Error())
	}

	return xmlNexus
}

func getDateBefore(days int) time.Time {
	curTime := time.Now()
	tHour := time.Duration(days * 24)
	curTimeMinusDays := curTime.Add(-tHour * time.Hour)
	return curTimeMinusDays
}

func FindOldArchives(xmlNexus XmlContent, minusDays time.Time) []OldArchives {
	var oldArchives []OldArchives

	for data := range xmlNexus.Data.ContentItem {
		dateTime, err := time.Parse(dataFmt, xmlNexus.Data.ContentItem[data].LastModified)
		if err != nil {
			log.Println(xmlNexus.Data.ContentItem[data].LastModified)
			log.Fatal(err.Error())
		}

		urlArchive := xmlNexus.Data.ContentItem[data].URI

		if minusDays.After(dateTime) {
			oldArchives = append(oldArchives, OldArchives{date: dateTime, url: urlArchive})
		}
	}

	for arch := range oldArchives {
		log.Printf("%s :: [%s] \n", oldArchives[arch].url, oldArchives[arch].date.Format(defaultDataFmt))
	}
	log.Println("Count old archives to delete: ", len(oldArchives))

	return oldArchives
}

// XML EXAMPLE:

//<content>
//	<data>
//		<content-item>
//			<resourceURI>https://nexus.com/nexus/service/local/repositories/repo/content/GEN/22219/GEN-22219-distrib.tgz</resourceURI>
//			<relativePath>/GEN/22219/GEN-22219-distrib.tgz</relativePath>
//			<text>RRM_GEN-DOTOOLS-22219-distrib.tgz</text>
//			<leaf>true</leaf>
//			<lastModified>2000-09-06 08:15:54.551 UTC</lastModified>
//			<sizeOnDisk>760832</sizeOnDisk>
//		</content-item>
//		<content-item>
//			<resourceURI>https://nexus.com/nexus/service/local/repositories/repo/content/GEN/22219/GEN-22219-distrib.tgz</resourceURI>
//			<relativePath>/GEN/22219/GEN-22219-distrib.tgz</relativePath>
//			<text>GEN-22219-distrib.tgz</text>
//			<leaf>true</leaf>
//			<lastModified>2000-09-06 08:15:54.551 UTC</lastModified>
//			<sizeOnDisk>760832</sizeOnDisk>
//		</content-item>
//	<data>
//<content>
