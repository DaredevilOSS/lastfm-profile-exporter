package common

import (
	"os"
)

type FileChannel chan string

const BaseURL = "https://ws.audioscrobbler.com/2.0/"
const MaxAPIRetries = 3

var ApiKey = os.Getenv("API_KEY")
