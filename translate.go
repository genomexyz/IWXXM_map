package main

import (
	"fmt"
	"os"
	"io"
	"regexp"
	"io/ioutil"
	"time"
	"strings"
	"strconv"
	"encoding/csv"
)

var wx_file = "ww.dat"
var file_stasiun = "stasiun.dat"

type DataSandi struct {
	stasiun string
	waktu time.Time
	sandi  string
}

type SandiTranslated struct {
	stasiun string
	waktu time.Time
	statmetar string
	autostat bool
	jenistrend string
	jeniswaktu string
	fmtrend time.Time
	tltrend time.Time
	arahangin int
	kecangin int
	gusty int
	anginvar1 int
	anginvar2 int
	vis int
	awanjumlah []string
	awantinggi []int
	wx []string
	suhu int
	dewpoint int
	tekanan int
	arahangintrend int
	kecangintrend int
	anginvartrend1 int
	anginvartrend2 int
	gustytrend int
	vistrend int
	wxtrend []string
	awanjumlahtrend []string
	awantinggitrend []int
	nswstat bool
	cavokstat bool
}

/*
func isitcavok(ds DataSandi) (bool):
	if 'CAVOK' in tafnowarray:
		return True
	else:
		return False
*/

func GenDataSandi(stsn string) (*DataSandi, error) {
	filename := "RAW/selected/"+stsn+".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return &DataSandi{stasiun: stsn, waktu: time.Now().UTC(), sandi: ""}, err
	}
	body_slice := strings.Split(string(body), "\n")
	waktu_parse, err := time.Parse("02/01/2006 15:04:05Z", body_slice[0])
	if err != nil {
		return &DataSandi{stasiun: stsn, waktu: time.Now().UTC(), sandi: ""}, err
	}
	
	sandi_str := body_slice[1]
	return &DataSandi{stasiun: stsn, waktu: waktu_parse, sandi: sandi_str}, err
}

func GenSandiTranslated(ds *DataSandi) (*SandiTranslated) {
	st := SandiTranslated{stasiun: ds.stasiun, waktu: ds.waktu, statmetar: "NORMAL", autostat: false, jenistrend: "NONE", 
	jeniswaktu: "NONE", fmtrend: time.Now(), tltrend: time.Now(), arahangin: -9999, kecangin: -9999, gusty: -9999, 
	anginvar1: -9999, anginvar2: -9999, vis: -9999, wx: make([]string, 0), suhu: -9999, dewpoint: -9999, arahangintrend: -9999, kecangintrend: -9999, anginvartrend1: -9999, 
	anginvartrend2: -9999, gustytrend: -9999, vistrend: -9999, wxtrend: make([]string, 0), 
	awanjumlahtrend: make([]string, 0), awantinggitrend: make([]int, 0), nswstat: false, cavokstat: false}
	
	//make array from sandi
	sandi_trimmed := strings.TrimSpace(ds.sandi)
	space := regexp.MustCompile(`\s+`)
	sandi_trimmed = space.ReplaceAllString(sandi_trimmed, " ")
	
	//we dont need to translate sandi "="
	if string(sandi_trimmed[len(sandi_trimmed)-1]) == "=" {
		sandi_trimmed = string(sandi_trimmed[0:len(sandi_trimmed)-1])
	}
	
	sandi_slice := strings.Split(sandi_trimmed, " ")
	
	//split metar
	metarsplit := map[string]bool {
    "TEMPO": true,
    "BECMG": true,
    "NOSIG": true,
    "RMK": true }
    
    //cloud amount
    cloudall := map[string]bool {
    "FEW": true,
    "SCT": true,
    "BKN": true,
    "OVC": true }
    
    //list possible wx
	wx_map := make(map[string]string)  
    f, err := os.Open(wx_file)
	if err != nil {
		panic(err)
	}
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
		wx_map[record[0]] = record[1]
	}
	f.Close()
	
	boundary_idx := len(sandi_slice) - 1
	for i := range sandi_slice {
		if metarsplit[sandi_slice[i]] {
			boundary_idx = i
		}
	}
	
	metar_part := sandi_slice
	trend_part := make([]string, 0)
	if boundary_idx != len(sandi_slice) - 1 {
		metar_part = sandi_slice[:boundary_idx]
		trend_part = sandi_slice[boundary_idx+1:]
	}
	
	for i := range metar_part {
		if i == 1 {
			if metar_part[i] == "AMD" {
				st.statmetar = "AMD"
			} else if metar_part[i] == "COR" {
				st.statmetar = "COR"
			}
		} else if i == 3 {
			if metar_part[i] == "AUTO" {
				st.autostat = true
			} else if (len(metar_part[i])) == 7 {
				if string(metar_part[i][5:7]) == "KT" {
					arah_angin, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.arahangin = arah_angin
					}
					kec_angin, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.kecangin = kec_angin
					}
				}
			} else if (len(metar_part[i])) == 10 {
				if (string(metar_part[i][5]) == "G" && string(metar_part[i][8:10]) == "KT") {
					arah_angin, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.arahangin = arah_angin
					}
					kec_angin, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.kecangin = kec_angin
					}
					gusty, err := strconv.Atoi(string(metar_part[i][6:8]))
					if err == nil {
						st.gusty = gusty
					}
				}
			}
			
		} else if i == 4 {
			if metar_part[i] == "AUTO" {
				st.autostat = true
			} else if (len(metar_part[i])) == 7 {
				if string(metar_part[i][5:7]) == "KT" {
					arah_angin, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.arahangin = arah_angin
					}
					kec_angin, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.kecangin = kec_angin
					}
				} else if string(metar_part[i][3]) == "V" {
					anginvar1, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.anginvar1 = anginvar1
					}
					anginvar2, err := strconv.Atoi(string(metar_part[i][4:7]))
					if err == nil {
						st.anginvar2 = anginvar2
					}
				}
			} else if (len(metar_part[i])) == 10 {
				if (string(metar_part[i][5]) == "G" && string(metar_part[i][8:10]) == "KT") {
					arah_angin, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.arahangin = arah_angin
					}
					kec_angin, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.kecangin = kec_angin
					}
					gusty, err := strconv.Atoi(string(metar_part[i][6:8]))
					if err == nil {
						st.gusty = gusty
					}
				}
			} else if (len(metar_part[i])) == 4 {
				vis, err := strconv.Atoi(metar_part[i])
				if err == nil {
					st.vis = vis
				}
			}
		} else if i == 5 {
			if wx_map[metar_part[i]] != "" {
				st.wx = append(st.wx, metar_part[i])
			} else if (len(metar_part[i])) == 7 {
				if string(metar_part[i][3]) == "V" {
					anginvar1, err := strconv.Atoi(string(metar_part[i][:3]))
					if err == nil {
						st.anginvar1 = anginvar1
					}
					anginvar2, err := strconv.Atoi(string(metar_part[i][4:7]))
					if err == nil {
						st.anginvar2 = anginvar2
					}
				}
			} else if (len(metar_part[i])) == 6 {
				if cloudall[string(metar_part[i][:3])] {
					st.awanjumlah = append(st.awanjumlah, string(metar_part[i][:3]))
					tinggi_awan, err := strconv.Atoi(string(metar_part[i][3:6]))
					if err == nil {
						st.awantinggi = append(st.awantinggi, tinggi_awan)
					}
				}
			} else if (len(metar_part[i])) == 8 {
				if (cloudall[string(metar_part[i][:3])] && string(metar_part[i][6:8]) == "CB") {
					st.awanjumlah = append(st.awanjumlah, string(metar_part[i][:3]))
					tinggi_awan, err := strconv.Atoi(string(metar_part[i][3:6]))
					if err == nil {
						st.awantinggi = append(st.awantinggi, tinggi_awan)
					}
				}
			} else if (len(metar_part[i])) == 5 {
				if string(metar_part[i][2]) == "/" {
					suhu, err := strconv.Atoi(string(metar_part[i][:2]))
					if err == nil {
						st.suhu = suhu
					}
					dewpoint, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.dewpoint = dewpoint
					}
				}
			} else if (len(metar_part[i])) == 4 {
				vis, err := strconv.Atoi(metar_part[i])
				if err == nil {
					st.vis = vis
				}
			}
		} else if i > 5 {
			if wx_map[metar_part[i]] != "" {
				st.wx = append(st.wx, metar_part[i])
			} else if (len(metar_part[i])) == 6 {
				if cloudall[string(metar_part[i][:3])] {
					st.awanjumlah = append(st.awanjumlah, string(metar_part[i][:3]))
					tinggi_awan, err := strconv.Atoi(string(metar_part[i][3:6]))
					if err == nil {
						st.awantinggi = append(st.awantinggi, tinggi_awan)
					}
				}
			} else if (len(metar_part[i])) == 8 {
				if (cloudall[string(metar_part[i][:3])] && string(metar_part[i][6:8]) == "CB") {
					st.awanjumlah = append(st.awanjumlah, string(metar_part[i][:3]))
					tinggi_awan, err := strconv.Atoi(string(metar_part[i][3:6]))
					if err == nil {
						st.awantinggi = append(st.awantinggi, tinggi_awan)
					}
				}
			} else if (len(metar_part[i])) == 5 {
				if string(metar_part[i][0]) == "Q" {
					tekanan, err := strconv.Atoi(string(metar_part[i][1:5]))
					if err == nil {
						st.tekanan = tekanan
					}
				} else if string(metar_part[i][2]) == "/" {
					suhu, err := strconv.Atoi(string(metar_part[i][:2]))
					if err == nil {
						st.suhu = suhu
					}
					dewpoint, err := strconv.Atoi(string(metar_part[i][3:5]))
					if err == nil {
						st.dewpoint = dewpoint
					}
				}
			}
		}
		
	}
	
	
	//check NSW
	for i := range trend_part {
		if sandi_slice[i] == "NSW" {
			st.nswstat = true
		}
	}
	
	//check CAVOK
	for i := range trend_part {
		if sandi_slice[i] == "CAVOK" {
			st.cavokstat = true
		}
	}
	
	//for the rest of trend
/*	for i := range trend_part {
		if i == 0 {
			if trend_part[i] == "TEMPO" {
				st.jenistrend = "TEMPORARY_FLUCTUATIONS"
			} else if trend_part[i] == "BECMG" {
				st.jenistrend = "BECOMING"
			}
		} else {
			if wx_map[trend_part[i]] != "" {
				st.wxtrend = append(st.wxtrend, trend_part[i])
			} else if (len(trend_part[i])) == 6 {
				if string(trend_part[i][:2]) == "FM" {
					st.jeniswaktu = "FROM"
				}
			} 
		}
	}
*/	
	return &st
}

func (st *SandiTranslated) save() string {
	awan_jumlah := "Cerah"
	if len(st.awanjumlah) > 0 {
		if st.awanjumlah[0] == "FEW" {
			awan_jumlah = "Sedikit Berawan"
		} else if st.awanjumlah[0] == "SCT" {
			awan_jumlah = "Berawan"
		} else if st.awanjumlah[0] == "BKN" {
			awan_jumlah = "Sangat Berawan"
		} else if st.awanjumlah[0] == "OVC" {
			awan_jumlah = "Berawan Penuh"
		}
	}
	
	awan_tinggi := 0
	if len(st.awanjumlah) > 0 {
		awan_tinggi = st.awantinggi[0] * 100
	}
	
	cuaca := "NOSIG"
	if len(st.wx) > 0 {
		cuaca = st.wx[0]
	}
	
	strfile := fmt.Sprintf("%s,%d,%d,%d,%d,%d,%s,%d,%s,%d,%d,%d", st.waktu.Format("02-01-2006 15:04:05Z"), st.arahangin, st.kecangin,
	st.anginvar1, st.anginvar2, st.vis, awan_jumlah, awan_tinggi, cuaca, st.suhu, st.dewpoint, st.tekanan)
	fmt.Println(strfile)
	
//	return ioutil.WriteFile("translated/"+st.stasiun+".txt", []byte(strfile), 0600)
	return strfile
}

func main() {
	var stasiun_list []string
	var stasiun_list_coordinate [][]string
	
	//read file, extract ICAO code
	f, err := os.Open(file_stasiun)
	if err != nil {
		return
	}
	defer f.Close()
	
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
		stasiun_list_coordinate = append(stasiun_list_coordinate, []string{record[2], record[3]})
	}
	fmt.Println(stasiun_list)
	f.Close()
	
	str_translated := "icao_code,lat,lon,time,arah_angin,kecepatan_angin,"+
	"anginvar1, anginvar2,visibility,jumlah_awan,tinggi_awan,cuaca,suhu,dewpoint,tekanan"
	for i := range stasiun_list {
		ds, err := GenDataSandi(stasiun_list[i])
		if err != nil {
			fmt.Println(stasiun_list[i], "Broken File or file not exist, continue...")
			continue
		}
		
		translated_sandi := GenSandiTranslated(ds)
		translated_segment := translated_sandi.save()
		
		str_translated += stasiun_list[i]+","+stasiun_list_coordinate[i][0]+","+
		stasiun_list_coordinate[i][1]+","+translated_segment+"\n"
		
		
		//write translated sandi in file
		fmt.Println(translated_sandi, translated_sandi.wx)
	}
	
	err = ioutil.WriteFile("static/data/recent_weather.txt", []byte(str_translated), 0600)
	if err != nil {
		fmt.Println("Saving file error")
	}
}
