package postprocess

import (
	"math/big"
	"math/rand"
	"time"
)

func GenerateNextId() int64 {
	randProvider := rand.NewSource(time.Now().Unix())
	nextId := big.NewInt(0).Sqrt(big.NewInt(10000))
	random := big.NewInt(randProvider.Int63())
	nextId = nextId.Mul(nextId, random)
	return nextId.Int64()
}

func GenerateTimestampForContent(translateContent string) int64 {
	return generateTimestamp(CountAlphaI(translateContent))
}

func generateTimestamp(iCount int64) int64 {
	ts := time.Now().UnixMilli()
	if iCount != 0 {
		iCount = iCount + 1
		return ts - ts%iCount + iCount
	} else {
		return ts
	}
}

func GenerateRandomNumber() int64 {
	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return myRand.Int63n(100000) + 83000000000
}
