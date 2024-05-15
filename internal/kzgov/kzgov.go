package kzgov

import (
	"encoding/json"
	"errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/capsolver"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	kzGovWebsiteKey = "6LeVUdoUAAAAAL8OTsq8xN-1czb129z5zxWsXkMk"
	kzGovUrl        = "https://kgd.gov.kz/ru/app/culs-taxarrear-search-web"
	kzGovRestUrl    = "https://kgd.gov.kz/apps/services/culs-taxarrear-search-web/rest/search"
)

type Request struct {
	GRecaptchaResponse string `json:"g-recaptcha-response"`
	IinBin             string `json:"iinBin"`
}

func GetArrear(iinBin string) (model.KZGovResponse, error) {
	var kzGovResponse model.KZGovResponse

	if !CheckIinFormat(iinBin) {
		return kzGovResponse, errors.New("wrong iinBin format")
	}

	capSolver := capsolver.NewInstance(os.Getenv("CAP_SOLVER_API_KEY"), kzGovUrl, kzGovWebsiteKey)
	captcha, err := capSolver.GetGCaptcha()
	if err != nil {
		return kzGovResponse, err
	}

	requestBody := Request{
		GRecaptchaResponse: captcha,
		IinBin:             iinBin,
	}

	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return kzGovResponse, err
	}

	res, err := http.Post(kzGovRestUrl, "application/json", strings.NewReader(string(requestJson)))
	if err != nil {
		return kzGovResponse, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	err = json.Unmarshal(body, &kzGovResponse)
	return kzGovResponse, err
}

func CompareNames(kzGovNameKk, kzGovNameRu, name string) bool {
	if name == "" {
		return false
	}
	nameParts := strings.Split(name, " ")
	for i := range len(nameParts) {
		nameParts[i] = strings.ToUpper(nameParts[i])
	}
	if kzGovNameKk != "" {
		kkParts := strings.Split(kzGovNameKk, " ")
		if len(kkParts) != len(nameParts) {
			return false
		}
		for i := range len(nameParts) {
			if nameParts[i] != kkParts[i] {
				return false
			}
		}
		return true
	}
	if kzGovNameRu != "" {
		ruParts := strings.Split(kzGovNameRu, " ")
		if len(ruParts) != len(nameParts) {
			return false
		}
		for i := range len(nameParts) {
			if nameParts[i] != ruParts[i] {
				return false
			}
		}
		return true
	}

	return false
}

func CheckIinFormat(iinBin string) bool {
	iinBinLenPattern := "\\d{12}"
	if iinBin == "" {
		return false
	}
	matchLen, _ := regexp.MatchString(iinBinLenPattern, iinBin)
	if !matchLen {
		return false
	}

	date := iinBin[0:6]
	gender := iinBin[6:7]
	return isValidBirthday(date) && isValidGender(gender) && isValidControl(iinBin)
}

func isValidBirthday(date string) bool {
	birthdayPattern := "\\d{2}(0[1-9]|1[012])(0[1-9]|[12][0-9]|3[01])$"
	if date == "" {
		return false
	}
	match, _ := regexp.MatchString(birthdayPattern, date)
	return match
}

func isValidGender(gender string) bool {
	genderPattern := "[0-6]"
	if gender == "" {
		return false
	}
	match, _ := regexp.MatchString(genderPattern, gender)
	return match
}

func isValidControl(iinBin string) bool {
	b1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	b2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}
	iinDigits := make([]int, 0)
	control := 0
	for i := range 12 {
		d, err := strconv.Atoi(iinBin[i : i+1])
		if err != nil {
			return false
		}
		iinDigits = append(iinDigits, d)
		if i != 11 {
			control += d * b1[i]
		}
	}
	control %= 11
	if control == 10 {
		control = 0
		for i := range 11 {
			control += iinDigits[i] * b2[i]
		}
		control %= 11
	}
	return control == iinDigits[11]
}
