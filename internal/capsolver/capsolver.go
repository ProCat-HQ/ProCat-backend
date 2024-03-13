// Package capsolver example:
// capSolver := capsolver.NewInstance(os.Getenv("CAP_SOLVER_API_KEY"),
//
//	"https://kgd.gov.kz/ru/app/culs-taxarrear-search-web",
//	"6LeVUdoUAAAAAL8OTsq8xN-1czb129z5zxWsXkMk")
//
// captcha, err := capSolver.GetGCaptcha()
//
//	if err != nil {
//		log.Fatalf("Can't get captcha: %s", err.Error())
//	}
package capsolver

import (
	capsolvergo "github.com/capsolver/capsolver-go"
	"github.com/sirupsen/logrus"
	"os"
)

type CapSolver struct {
	ApiKey     string
	URL        string
	WebsiteKey string
}

func NewInstance(apiKey, url, websiteKey string) *CapSolver {
	return &CapSolver{
		ApiKey:     apiKey,
		URL:        url,
		WebsiteKey: websiteKey,
	}
}

func (c *CapSolver) GetGCaptcha() (string, error) {
	capSolver := capsolvergo.CapSolver{ApiKey: c.ApiKey}
	s, err := capSolver.Solve(map[string]any{
		"type":       "ReCaptchaV2taskProxyLess",
		"websiteURL": c.URL,
		"websiteKey": c.WebsiteKey,
	})
	if err != nil {
		return "", err
	}
	return s.Solution.GRecaptchaResponse, nil
}

func writeCaptchaToFile(captcha, filename, path string) error {
	err := os.Mkdir(path, 0644)
	if err != nil {
		return err
	}
	f, err := os.Create(path + "/" + filename)
	if err != nil {
		f, err = os.Open(path + "/" + filename)
		if err != nil {
			logrus.Fatalf("Can't create or open file: %s", err.Error())
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Fatalf("Can't close file: %s", err.Error())
		}
	}(f)

	_, err = f.Write([]byte(captcha))
	if err != nil {
		return err
	}
	return nil
}
