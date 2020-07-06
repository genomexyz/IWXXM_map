package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"time"
	"strings"
	"encoding/csv"
	curl "github.com/andelf/go-curl"
)

type DataSandi struct {
	stasiun string
	waktu []byte
	sandi  []byte
}

func (ds *DataSandi) save() error {
	filename := "RAW/"+ds.stasiun + ".txt"
	separator := []byte("\n")
	isi := append(ds.waktu, separator...)
	isi = append(isi, ds.sandi...)
	return ioutil.WriteFile(filename, isi, 0600)
}


func main() {
	file_stasiun := "stasiun.dat"
	var stasiun_list []string
	
	//read file, extract ICAO code
	f, err := os.Open(file_stasiun)
	if err != nil {
		return
	}
	defer f.Close()
	
	fmt.Println(f)
	r := csv.NewReader(f)
	
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		stasiun_list = append(stasiun_list, record[0])
	}
	fmt.Println(stasiun_list)
	
	//declare time now
	time_now := time.Now().UTC()
//	fmt.Println(time_now.String().Format("01/Jan/2006:15:04:05 -0700"))
//	fmt.Println(time.Time(time_now).Format("01/Jan/2006:15:04:05 -0700"))
	tahun := time_now.Year()
	bulan := int(time_now.Month())
//	hari := int(time_now.Day())
	fmt.Println(tahun, bulan)
	
	//iter all ICAO code
	var link string
	var body_html string
	easy := curl.EasyInit()
	defer easy.Cleanup()
	for i := range stasiun_list {
		body_html = ""
		link = fmt.Sprintf("http://aviation.bmkg.go.id/latest/metar.php?i=%s&y=%d&m=%d", stasiun_list[i], tahun, bulan)
		fmt.Println(link)
		
		easy.Setopt(curl.OPT_URL, link)
		easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, _ interface{}) bool {
			body_html += string(ptr)
			return true
		})
		
		if err := easy.Perform(); err != nil {
			fmt.Println("ERROR %v: continue...\n", err)
			continue
		}
		
		if body_html == "" {
			fmt.Println("Empty link: continue...\n")
			continue
		}
		
		fmt.Println(body_html)
		
		//get last sandi
		sandi_slice := strings.Split(body_html, "\n")
		last_segment := sandi_slice[len(sandi_slice)-1]
		if last_segment == "" {
			last_segment = sandi_slice[len(sandi_slice)-2]
		}
		
		last_segment_slice := strings.Split(last_segment, "\t")
		
		if len(last_segment_slice) != 3 {
			fmt.Println("HTML body broken: continue...\n")
			continue
		}
		
		last_sandi := last_segment_slice[2]
		last_waktu := last_segment_slice[0]

		dataset := &DataSandi{stasiun: stasiun_list[i], waktu: []byte(last_waktu), sandi: []byte(last_sandi)}
		dataset.save()
	}
}
