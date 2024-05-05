package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/valyala/fastjson"
)

// Strip out similar URLs by unique hostname-path-paramName and some other noise pattern
// cat urls.txt | durl
// only grep url have parameter
// cat urls.txt | durl -p

var (
	excludeStatic bool
	excludeNoise  bool
	haveParam     bool
	handleJson    bool
	limit         int
	ext           string
	targetScope   string
	jsonField     string
)

func main() {
	// cli aguments
	flag.BoolVar(&excludeStatic, "s", true, "Exclude static files extensions")
	flag.BoolVar(&excludeNoise, "n", true, "Exclude noise content pattern like blogspot, calender, etc")
	flag.BoolVar(&haveParam, "p", false, "Enable check if input have parameter")
	flag.IntVar(&limit, "l", 100, "Limit length of path item (default 100)")
	flag.StringVar(&ext, "e", "", "Blacklist regex string")
	flag.StringVar(&targetScope, "t", "", "Target scope")
	flag.StringVar(&jsonField, "f", "", "Field to select in JSON data (only apply for JSON input)")

	flag.Parse()
	var p fastjson.Parser
	if jsonField != "" {
		handleJson = true
	}

	data := make(map[string]string)
	hostMapping := make(map[string]string)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		raw := strings.TrimSpace(sc.Text())
		if sc.Err() != nil && raw == "" {
			continue
		}

		var original string

		// handle json
		// check if the input is JSON or not
		if jsonField != "" {
			v, err := p.Parse(raw)
			if err != nil {
				continue
			}
			original = raw
			raw = string(v.GetStringBytes(jsonField))
		}

		if excludeStatic {
			if IsStaticPattern(raw) {
				continue
			}
		}

		// parsing the URL
		u, err := url.Parse(raw)
		if err != nil || u.Hostname() == "" {
			continue
		}

		// check if the url host is in scope or not
		if targetScope != "" {
			if !strings.Contains(u.Hostname(), targetScope) {
				continue
			}
		}

		hash := hashUrl(u)
		if hash == "" {
			continue
		}

		_, exist := data[hash]
		if !exist {
			if excludeNoise {
				if IsBlackList(raw) {
					_, notSeenYet := hostMapping[u.Hostname()]
					if !notSeenYet {
						hostMapping[u.Hostname()] = raw
						fmt.Println(raw)
					}
					continue
				}
			}

			if ext != "" {
				if !RegexCheck(ext, raw) {
					continue
				}
			}

			if handleJson {
				data[hash] = original
				fmt.Println(original)
			} else {
				data[hash] = raw
				fmt.Println(data[hash])
			}

		}
	}
}

// IsBlackList check if url is blacklisted or not
func IsBlackList(raw string) bool {
	calenderPattern := `(\d{2,4})(-|/)(\d{1,2})(-|/)(\d{1,2})`
	if RegexCheck(calenderPattern, raw) {
		return true
	}

	noiseContent := `/(articles|about|blog|event|events|shop|post|posts|product|products|docs|support|pages|media|careers|jobs|video|videos|resource|resources)/.*`
	if RegexCheck(noiseContent, raw) {
		return true
	}

	// e.g: /abc/1234
	idContentNoExt := `.*\/[0-9]+$`
	if RegexCheck(idContentNoExt, raw) {
		return true
	}
	// e.g: /abc/1234.html
	idContent := `.*\/[0-9]+\.[a-z]+`
	return RegexCheck(idContent, raw)
}

func RegexCheck(pattern string, raw string) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return r.MatchString(raw)
}

// IsStaticPattern check if url is blacklisted or not
func IsStaticPattern(raw string) bool {
	staticPattern := `(?i)\.(png|apng|bmp|gif|ico|cur|jpg|jpeg|jfif|pjp|pjpeg|svg|tif|tiff|webp|xbm|3gp|aac|flac|mpg|mpeg|mp3|mp4|m4a|m4v|m4p|oga|ogg|ogv|mov|wav|webm|eot|woff|woff2|ttf|otf|css)(?:\?|#|$)`

	if RegexCheck(staticPattern, raw) {
		return true
	}

	// check if have param
	if haveParam {
		return !RegexCheck(`\?.*\=`, raw)
	}

	return false
}

// hashUrl gen unique hash base on url
func hashUrl(u *url.URL) string {
	// length check for path element or seeing too much "-"
	if strings.Count(u.Path, "/") >= 1 {
		paths := strings.Split(u.Path, "/")
		for _, item := range paths {
			if len(item) > limit || strings.Count(item, "-") > 3 {
				return ""
			}
		}
	}

	var queries []string
	for k := range u.Query() {
		queries = append(queries, k)
	}
	sort.Strings(queries)
	query := strings.Join(queries, "-")

	data := fmt.Sprintf("%v-%v-%v", u.Hostname(), u.Path, query)
	return genHash(data)
}

// genHash gen SHA1 hash from string
func genHash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%v", hashed)
}
