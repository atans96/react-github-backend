package service

import (
	"backend/src/types"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (f *Fetch) FetchRateLimit(token string) (types.ApiRateLimitRule, error) {
	var ja types.ApiRateLimitRule
	for {
		makeNew, err := http.NewRequest("GET", "https://api.github.com/rate_limit", nil)
		if err != nil {
			panic(err)
		}
		makeNew.Header.Set("Authorization", "Bearer "+token)
		client := http.Client{}
		b, err := client.Do(makeNew)
		if err != nil {
			return types.ApiRateLimitRule{}, err
		}
		bo, err := ioutil.ReadAll(b.Body)
		go b.Body.Close()
		if err != nil {
			return types.ApiRateLimitRule{}, err
		}
		var apiRateLimit types.ApiRateLimit
		err = json.Unmarshal(bo, &apiRateLimit)
		core := apiRateLimit.Resources.Core
		if core.Remaining < 50 {
			exp := time.Unix(core.Reset, 0)
			dur := time.Until(exp)
			fmt.Printf("RATELIMIT REACHED - SLEEPING %v\n\n", dur)
			time.Sleep(dur)
			continue
		}
		//fmt.Printf("NEW RATELIMIT LOADED - %d remaining, exp %v", core.Remaining, time.Unix(core.Reset, 0))
		ja = core
		break
	}

	return ja, nil
}
