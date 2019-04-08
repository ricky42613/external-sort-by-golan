package main

import (
	"io"
	"os"
	"strconv"
	"strings"
)

func merge(keyword string, arr []string, start int, mid int, end int) {
	leftlen := mid - start + 1
	rightlen := end - mid
	arrleft := make([]string, leftlen)
	arrright := make([]string, rightlen)
	for i := 0; i < leftlen; i++ {
		arrleft[i] = arr[start+i]
	}
	for j := 0; j < rightlen; j++ {
		arrright[j] = arr[mid+1+j]
	}
	i, j, k := 0, 0, start
	var ikey string
	var jkey string
	var placei, placej, itail, jtail int
	for i < leftlen && j < rightlen && k <= end {
		if keyword == "" {
			if strings.Compare(arrleft[i], arrright[j]) < 0 {
				arr[k] = arrleft[i]
				i++
			} else {
				arr[k] = arrright[j]
				j++
			}
			k++
		} else {
			placei = strings.Index(arrleft[i], "@"+keyword)
			placej = strings.Index(arrright[j], "@"+keyword)
			itail = strings.Index(arrleft[i][placei+1:], "@")
			jtail = strings.Index(arrright[j][placej+1:], "@")
			ikey = arrleft[i][placei : itail+placei+1]
			jkey = arrright[j][placej : jtail+placej+1]
			if strings.Compare(ikey, jkey) < 0 {
				arr[k] = arrleft[i]
				i++
			} else {
				arr[k] = arrright[j]
				j++
			}
			k++
		}
	}
	for j < rightlen {
		arr[k] = arrright[j]
		j++
		k++
	}
	for i < leftlen {
		arr[k] = arrleft[i]
		i++
		k++
	}
}

func mergesort(key string, arr []string, start int, end int) {
	if start < end {
		mid := (start + end) / 2
		mergesort(key, arr, start, mid)
		mergesort(key, arr, mid+1, end)
		merge(key, arr, start, mid, end)
	}
}

func strcmp(a, b string) string {
	if strings.Compare(a, b) < 0 {
		return a
	}
	return b
}

func selection(a, b, key string) string {
	var akey, bkey string
	if key != "" {
		if a != "" {
			placea := strings.Index(a, "@"+key)
			atail := strings.Index(a[placea+1:], "@")
			akey = a[placea : atail+placea+1]
		} else {
			akey = ""
		}
		if b != "" {
			placeb := strings.Index(b, "@"+key)
			btail := strings.Index(b[placeb+1:], "@")
			bkey = b[placeb : btail+placeb+1]
		} else {
			bkey = ""
		}
	} else {
		akey = a
		bkey = b
	}
	if akey != "" && bkey != "" {
		if strcmp(akey, bkey) == akey {
			return a
		}
		return b
	} else if a == "" && b == "" {
		return ""
	} else if akey == "" {
		return b
	} else {
		return a
	}
}

func round(ptr []string, key string) []string {
	r := []string{}
	for i := 0; i < len(ptr); i += 2 {
		if i+1 < len(ptr) {
			r = append(r, selection(ptr[i], ptr[i+1], key))
		} else {
			if ptr[i] != "" {
				r = append(r, ptr[i])
			}
		}
	}
	return r
}

func choose(ptr []string, key string) string {
	r := ptr
	lenr := len(r)
	for lenr > 1 {
		r = round(r, key)
		lenr = len(r)
	}
	return r[0]
}

func checknum(arr []string, target string) int {
	l := len(arr)
	i := 0
	for i = 0; i < l; i++ {
		if arr[i] == target {
			break
		}
	}
	return i
}

func recordat(target []byte, begin string) []byte {
	length := len(target)
	if begin != "" {
		i := 1
		for {
			if string(target[length-len(begin)-i:length-i+1]) == "@"+begin {
				target = target[0 : length-len(begin)-i]
				break
			} else {
				i++
			}
		}
	} else {
		j := 1
		for {
			if target[length-j] == '\n' {
				target = target[:length-j+1]
				break
			} else {
				j++
			}
		}
	}
	return target
}

func sep(result []string, target []byte, recbeg string) []string {
	result = result[:0]
	if recbeg != "" {
		result = strings.SplitAfter(string(target), "@"+recbeg)
		result = result[1:]
		for i := 0; i < len(result); i++ {
			result[i] = "@" + recbeg + result[i]
			if i != len(result)-1 {
				result[i] = result[i][:len(result[i])-len(recbeg)-1]
			}
		}
	} else {
		result = strings.SplitAfter(string(target), "\n")
	}
	return result
}

func main() {
	var record int64
	argc := len(os.Args)
	argi := 1
	inputfile := ""
	outputfile := ""
	key := ""
	recbeg := ""
	for argi = 1; argi < argc; argi += 2 {
		if os.Args[argi] == "-f" {
			inputfile += os.Args[argi+1]
		} else if os.Args[argi] == "-o" {
			outputfile += os.Args[argi+1]
		} else if os.Args[argi] == "-k" {
			key += os.Args[argi+1]
		} else if os.Args[argi] == "-rb" {
			recbeg += os.Args[argi+1]
		}
	}
	filein, err := os.Open(inputfile)
	f, err := os.Stat(inputfile)
	size := f.Size()
	if err != nil {
		panic(err)
	}
	defer filein.Close()
	cnt := 0
	var tmp = make([]byte, 1000000000)
	tmpfilename := []string{}
	datastruct := []string{}
	for {
		space, errspace := filein.ReadAt(tmp, record)
		if space == 0 {
			break
		}
		if errspace == io.EOF {
			size = size - record
			var laststr = make([]byte, size)
			last, lasterr := filein.ReadAt(laststr, record)
			if lasterr != nil {
				panic(lasterr)
			}
			if laststr[0] == '\n' {
				laststr = laststr[1:]
			}
			if last != 0 {
				tmpfilename = append(tmpfilename, strconv.Itoa(cnt)+".txt")
				datastruct = sep(datastruct, laststr, recbeg)
				tmplast, tmplasterr := os.OpenFile(tmpfilename[cnt], os.O_WRONLY|os.O_CREATE, 0666)
				if tmplasterr != nil {
					panic(tmplasterr)
				}
				os.Truncate(tmpfilename[cnt], 0)
				defer tmplast.Close()
				mergesort(key, datastruct, 0, len(datastruct)-1)
				for t := 0; t < len(datastruct); t++ {
					tmplast.WriteString(datastruct[t])
				}
			}
			cnt++
			break
		} else if errspace != nil {
			panic(errspace)
		}
		tmpfilename = append(tmpfilename, strconv.Itoa(cnt)+".txt")
		start := 0
		if tmp[0] == '\n' {
			tmp = tmp[1:]
			start++
		}
		tmp = recordat(tmp, recbeg)
		record = record + int64(len(tmp)) + int64(start)
		tmpfile, tmperr := os.OpenFile(tmpfilename[cnt], os.O_WRONLY|os.O_CREATE, 0666) //memory grow up
		if tmperr != nil {
			panic(tmperr)
		}
		datastruct = sep(datastruct, tmp, recbeg)
		mergesort(key, datastruct, 0, len(datastruct)-1)
		os.Truncate(tmpfilename[cnt], 0)
		for t := 0; t < len(datastruct); t++ {
			tmpfile.WriteString(datastruct[t])
		}
		cnt++
		datastruct = datastruct[:0]
		defer tmpfile.Close()
	}
	// 臨時檔案產生完成
	k := 0
	i := 0
	//slicearr := [11][]byte{}
	slicearr := make([][]byte, len(tmpfilename))
	//var slicerecord = [11]int{}
	slicerecord := make([]int, len(tmpfilename))
	//datalist := [11][]string{}
	datalist := make([][]string, len(tmpfilename))
	flag := make([]int, len(tmpfilename))
	for k = 0; k < cnt; k++ {
		slicearr[k] = make([]byte, 100000000)
	}
	for {
		slice, sliceerr := os.Open(tmpfilename[i])
		if sliceerr != nil {
			panic(sliceerr)
		}
		sliceinfo, sliceinfoerr := os.Stat(tmpfilename[i])
		if sliceinfoerr != nil {
			panic(sliceinfoerr)
		}
		slicesize := sliceinfo.Size()
		slicespace, slicespaceerr := slice.ReadAt(slicearr[i], int64(slicerecord[i]))
		if slicespace == 0 {
			break
		}
		if slicespaceerr == io.EOF {
			flag[i] = 1
			lastslice := make([]byte, slicesize-int64(slicerecord[i]))
			lastslicespace, lastslicespaceerr := slice.ReadAt(lastslice, int64(slicerecord[i]))
			if lastslicespace != 0 {
				if lastslicespaceerr != nil {
					panic(lastslicespaceerr)
				}
				laststart := 0
				if lastslice[0] == '\n' {
					lastslice = lastslice[1:]
					laststart++
				}
				datalist[i] = sep(datalist[i], lastslice, recbeg)
				slicerecord[i] = slicerecord[i] + len(lastslice) + laststart
			}
			break
		} else if slicespaceerr != nil {
			panic(slicespaceerr)
		}
		slicestart := 0
		if slicearr[i][0] == '\n' {
			slicearr[i] = slicearr[i][1:]
			slicestart++
		}
		tmpslice := make([]byte, 100000000)
		slicearr[i] = recordat(slicearr[i], recbeg)
		tmpslice = slicearr[i]
		datalist[i] = sep(datalist[i], tmpslice, recbeg)
		slicerecord[i] = slicerecord[i] + len(slicearr[i]) + slicestart
		i++
	}
	//從各個臨時檔中切1/10
	a := make([]int, len(datalist))
	export := []string{}
	exp, experr := os.OpenFile(outputfile, os.O_WRONLY|os.O_CREATE, 0666)
	if experr != nil {
		panic(experr)
	}
	os.Truncate(outputfile, 0)
	ptr := []string{datalist[0][a[0]], datalist[1][a[1]], datalist[2][a[2]], datalist[3][a[3]], datalist[4][a[4]], datalist[5][a[5]], datalist[6][a[6]], datalist[7][a[7]], datalist[8][a[8]], datalist[9][a[9]], datalist[10][a[10]]}
	n := 0
	c := 0
	for a[0] < len(datalist[0]) || a[1] < len(datalist[1]) || a[2] < len(datalist[2]) || a[3] < len(datalist[3]) || a[4] < len(datalist[4]) || a[5] < len(datalist[5]) || a[6] < len(datalist[6]) || a[7] < len(datalist[7]) || a[8] < len(datalist[8]) || a[9] < len(datalist[9]) || a[10] < len(datalist[10]) {
		export = append(export, choose(ptr, key))
		n = checknum(ptr, export[len(export)-1])
		a[n]++
		if len(export) == 10000 {
			for t := 0; t < len(export); t++ {
				exp.WriteString(export[t])
			}
			export = export[:0]
		}
		if a[n] == len(datalist[n]) {
			if flag[n] == 1 {
				c++
				if c == len(datalist) {
					break
				} else {
					datalist[n] = append(datalist[n], "")
				}
			} else {
				fp, ferr := os.Open(tmpfilename[n])
				fpinfo, fpinfoerr := os.Stat(tmpfilename[n])
				defer fp.Close()
				if fpinfoerr != nil {
					panic(fpinfoerr)
				}
				if ferr != nil {
					panic(ferr)
				}
				next := make([]byte, 100000000)
				fr, frerr := fp.ReadAt(next, int64(slicerecord[n]))
				if frerr == io.EOF {
					flag[n] = 1
					lastone := make([]byte, fpinfo.Size()-int64(slicerecord[n]))
					lastfr, lastfrerr := fp.ReadAt(lastone, int64(slicerecord[n]))
					if lastfrerr != nil {
						panic(lastfrerr)
					}
					if lastfr != 0 {
						if lastone[0] == '\n' {
							lastone = lastone[1:]
							slicerecord[n]++
						}
						datalist[n] = sep(datalist[n], lastone, recbeg)
						slicerecord[n] += len(lastone)
						a[n] = 0
					}
				} else if frerr != nil {
					panic(frerr)
				} else if fr != 0 {
					if next[0] == '\n' {
						next = next[1:]
						slicerecord[n]++
					}
					next = recordat(next, recbeg)
					datalist[n] = sep(datalist[n], next, recbeg)
					slicerecord[n] += len(next)
					a[n] = 0
				}
			}
		}
		ptr[n] = datalist[n][a[n]]
	}
	if len(export) != 0 {
		for t := 0; t < len(export); t++ {
			exp.WriteString(export[t])
		}
	}
	defer exp.Close()
}
