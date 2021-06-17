package Server

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func InitTimestamp(n int, pNumber int) []int {
	stamp := make([]int, n+1)
	stamp[len(stamp)-1] = pNumber
	return stamp
}

func addTimestamp(stamp []int, msg string) string {
	stamp[pNumber] = stamp[pNumber] + 1
	toAdd := TimestampTostring(stamp) + "#"
	return toAdd + msg
}

func handleMsg(msg string, stamp []int, n int, holdback map[string]string, stayLong map[string]string) {
	result := strings.Split(msg, "#")
	length := len(result[0]) + 1
	runes := []rune(msg)
	realMsg := string(runes[length:])
	newstamp := stringToTimestamp(result[0])

	if checkStamp(newstamp, stamp, n) {
		fmt.Println(realMsg)
		whereUpdate := newstamp[n]
		stamp[whereUpdate] = newstamp[whereUpdate]

		UpdateHoldback(n, localStamp)

	} else {
		keystring := TimestampTostring(newstamp)
		holdback[keystring] = realMsg

		layout := "2000-15-01 15:00:00"
		currentTime := time.Now().Format(layout)
		stayLong[keystring] = currentTime

	}

}

func UpdateHoldback(n int, localStamp []int) {
	myFlag := false
	for key, value := range holdback {
		tempStamp := stringToTimestamp(key)
		if checkStamp(tempStamp, localStamp, n) {
			whereUpdate := tempStamp[n]
			localStamp[whereUpdate] = tempStamp[whereUpdate]
			myFlag = true

			_, ok1 := holdback[key]
			if ok1 {
				delete(holdback, key)
			}

			_, ok2 := stayLong[key]
			if ok2 {
				delete(stayLong, key)
			}

			fmt.Println(value)
		}
	}
	if myFlag {
		UpdateHoldback(n, localStamp)
	}

}

func checkStamp(newstamp []int, stamp []int, n int) bool {
	fromWhichP := newstamp[n]
	for i := range newstamp {
		if i == n {
			break
		}
		if i != fromWhichP {
			if newstamp[i] > stamp[i] {
				return false
			}

		} else {
			if newstamp[i] != stamp[i]+1 {
				return false
			}
		}
	}

	return true

}

func TimestampTostring(stamp []int) string {
	output := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stamp)), ","), "[]")
	return output
}

func stringToTimestamp(str string) []int {
	s := strings.Split(str, ",")
	var val = []int{}
	for _, i := range s {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		val = append(val, j)
	}
	return val

}
