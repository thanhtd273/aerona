package utils

import (
	"crypto"
	"crypto/md5"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/dto"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

func Hash(elements ...string) string {
	digester := crypto.MD5.New()
	for _, element := range elements {
		fmt.Fprint(digester, reflect.TypeOf(element))
		fmt.Fprint(digester, element)
	}
	return fmt.Sprintf("%x", digester.Sum(nil))
}

func FormatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func GenerateBrowserId(r *http.Request) string {
	// Get from cookie
	if cookie, err := r.Cookie("browserId"); err == nil {
		return cookie.Value
	}

	// If not in cookie, generate from User-Agent header
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		hash := md5.Sum([]byte(userAgent))
		return fmt.Sprintf("browser-%s", hash)
	}

	return uuid.New().String()
}

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return num
}

func ParseHourRanges(ranges []string) []dto.HourRange {
	result := make([]dto.HourRange, 0, len(ranges))
	for _, r := range ranges {
		parts := strings.Split(r, "-")
		if len(parts) == 2 {
			result = append(result, dto.HourRange{From: parts[0], To: parts[1]})
		}
	}
	return result
}

func GenerateUniqueId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}
	id := node.Generate()
	str := id.String()
	return str, nil
}
